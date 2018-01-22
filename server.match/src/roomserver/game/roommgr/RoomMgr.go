package roommgr

// 房间管理类

import (
	"runtime/debug"
	"sync"
	"time"

	"base/glog"
	"common"
	"roomserver/client/rcenterclient"
	"roomserver/client/voiceclient"
	rm "roomserver/game/room"
	"roomserver/interfaces"
	"roomserver/redismgr"
)

type RoomMgr struct {
	mutex sync.RWMutex // protect rooms
	rooms map[uint32]*rm.Room

	frame int64

	// 匹配相关
	nextRoomId      uint32
	playerId2roomId map[uint64]uint32
	roomId2playerId map[uint32][]uint64
	suspensePlayers []interfaces.IPlayerTask
	mutexMatch      sync.Mutex
}

var (
	roommgr      *RoomMgr
	roommgr_once sync.Once
)

func GetMe() *RoomMgr {
	if roommgr == nil {
		roommgr_once.Do(func() {
			roommgr = &RoomMgr{
				rooms:           make(map[uint32]*rm.Room),
				nextRoomId:      1,
				playerId2roomId: make(map[uint64]uint32),
				roomId2playerId: make(map[uint32][]uint64),
			}
			roommgr.init()
		})
	}
	return roommgr
}

func (this *RoomMgr) init() {
	go func() {
		tick := time.NewTicker(time.Second * 1)
		defer func() {
			if err := recover(); err != nil {
				glog.Error("[异常] 报错 ", err, "\n", string(debug.Stack()))
			}
			tick.Stop()
		}()
		for {
			select {
			case <-tick.C:
				this.timeAction()
			}
		}
	}()
}

// 定时事件
func (this *RoomMgr) timeAction() {
	var prooms []*rm.Room
	nowTime := time.Now().Unix()

	this.mutex.RLock()
	rms := make([]*rm.Room, 0)
	for _, room := range this.rooms {
		rms = append(rms, room)
	}
	this.mutex.RUnlock()

	this.frame++
	for _, room := range rms {
		if !room.CheckInited() {
			continue
		}
		channum := make(chan []int, 1)
		room.Chan_GetPlayerNum <- channum
		num := <-channum
		playernum := num[1]
		allpalyernum := num[0] + num[1]
		if room.EndTime() > nowTime && playernum != 0 {
			continue
		}
		if this.frame%50 == 0 && room.RoomType() == common.RoomTypeQuick {
			if allpalyernum != 0 && room.EndTime() > nowTime {
				continue
			}
		}
		prooms = append(prooms, room)
	}
	for _, room := range prooms {
		this.RemoveRoom(room.ID())
		if !room.IsCustom {
			rcenterclient.GetMe().RemoveRoom(room.RoomType(), room.ID())
		}
		rcenterclient.GetMe().UpdateServer()
		room.Control(rm.ROOM_CONTROL_STOP)
		glog.Info("[房间]删除房间 deleteroom [", room.RoomType(), ",", room.ID(), "],", this.GetNum(),
			",time:", nowTime-room.EndTime(), " now:", nowTime, " end:", room.EndTime())
	}
}

// 添加房间
func (this *RoomMgr) AddRoom(room *rm.Room) *rm.Room {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	oldroom, ok := this.rooms[room.ID()]
	if ok {
		room.Control(rm.ROOM_CONTROL_END)
		return oldroom
	}
	this.rooms[room.ID()] = room
	return room
}

// 删除房间
func (this *RoomMgr) RemoveRoom(rid uint32) {

	this.mutexMatch.Lock()
	if ids, ok := this.roomId2playerId[rid]; ok {
		for _, id := range ids {
			delete(this.playerId2roomId, id)
		}
		delete(this.roomId2playerId, rid)
	}
	this.mutexMatch.Unlock()

	this.mutex.Lock()
	delete(this.rooms, rid)
	this.mutex.Unlock()
}

// 获取房间列表
func (this *RoomMgr) GetRoomIds() (froomids []uint32, troomids []uint32, qroomids []uint32) {
	this.mutex.RLock()
	for _, room := range this.rooms {
		switch room.RoomType() {
		case common.RoomTypeTeam:
			troomids = append(troomids, room.ID())
		case common.RoomTypeQuick:
			qroomids = append(qroomids, room.ID())
		}
	}
	this.mutex.RUnlock()
	return
}

// 获取玩家列表
func (this *RoomMgr) getPlayerIDs() (tuids, quids, cuids []uint64) {
	this.mutex.RLock()
	for _, room := range this.rooms {
		switch room.RoomType() {
		case common.RoomTypeTeam:
			for _, player := range room.Players {
				tuids = append(tuids, player.ID)
				if room.IsCustom {
					cuids = append(cuids, player.ID)
				}
			}
		case common.RoomTypeQuick:
			for _, player := range room.Players {
				quids = append(quids, player.ID)
				if room.IsCustom {
					cuids = append(cuids, player.ID)
				}
			}
		}
	}
	this.mutex.RUnlock()
	return
}

// 获取房间列表
func (this *RoomMgr) GetRooms() (rooms []*rm.Room) {
	this.mutex.RLock()
	for _, room := range this.rooms {
		rooms = append(rooms, room)
	}
	this.mutex.RUnlock()
	return
}

// 根据id获取有效房间
func (this *RoomMgr) getRoomById(rid uint32) *rm.Room {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	room, ok := this.rooms[rid]
	if !ok {
		return nil
	}
	return room
}

func (this *RoomMgr) GetNum() int {
	this.mutex.RLock()
	roomnum := len(this.rooms)
	this.mutex.RUnlock()
	return roomnum
}

func (this *RoomMgr) NewRoom(sceneId uint32, rtype, rid uint32,
	player interfaces.IPlayerTask, robot int) *rm.Room {

	isnew := false
	if player.UData().PlayNum <= 5 {
		isnew = true
	}

	room := this.AddRoom(rm.NewRoom(sceneId, rtype, rid, player.UData().IsCustom, player.UData().RoomName))
	voiceclient.GetMe().NewRoom(rid, uint32(room.EndTime()-time.Now().Unix()))
	if !room.IsCustom {
		rcenterclient.GetMe().AddRoom(rtype, rid, uint32(room.EndTime()), sceneId, false, player.UData().HideScore, player.UData().UnCoop, robot, isnew, player.UData().Level)
	}
	glog.Info("[房间]create new room:", sceneId, "  rid:", rid, " isnew:", isnew, " playnum:", player.UData().PlayNum)
	return room
}

const DEFAULT_MATCH_NUM = 1

func (this *RoomMgr) AddPlayer(player interfaces.IPlayerTask) bool {

	var roomId uint32 = 0
	this.mutexMatch.Lock()
	if val, ok := this.playerId2roomId[player.ID()]; ok {
		roomId = val
	}
	this.mutexMatch.Unlock()

	var room *rm.Room = nil
	if roomId != 0 {
		room = this.getRoomById(roomId)
	}

	if room == nil {
		// 开始游戏，还没有房间。
		this.mutexMatch.Lock()
		this.suspensePlayers = append(this.suspensePlayers, player)
		if len(this.suspensePlayers) >= DEFAULT_MATCH_NUM {
			roomId = this.nextRoomId
			this.nextRoomId++
			player0 := this.suspensePlayers[0]
			room = this.NewRoom(player0.UData().SceneId, common.RoomTypeQuick, roomId, player0, 0)
			this.roomId2playerId[roomId] = make([]uint64, 0)
			for i := 0; i < DEFAULT_MATCH_NUM; i++ {
				playeri := this.suspensePlayers[i]
				this.playerId2roomId[playeri.ID()] = roomId
				this.roomId2playerId[roomId] = append(this.roomId2playerId[roomId], playeri.ID())
				playeri.UData().RoomId = roomId
				playeri.SetRoom(room)
				room.IncPlayerNum(playeri)
				room.Chan_AddPlayer <- playeri
				if !room.IsCustom {
					rcenterclient.GetMe().UpdateRoom(room.RoomType(), room.ID(), room.GetPlayerNum(), common.UserOnline, false, 0, 0)
				}
			}
			this.suspensePlayers = this.suspensePlayers[DEFAULT_MATCH_NUM:]
		} else {

			// TODO: 可以做个协议通知客户端，哪些玩家正在匹配等待中

		}
		this.mutexMatch.Unlock()

	} else {
		// 已经有房间了（如断线后，重新进入）
		if !room.IsClosed() {
			player.UData().RoomId = roomId
			player.SetRoom(room)
			room.IncPlayerNum(player)
			room.Chan_AddPlayer <- player
			if !room.IsCustom {
				rcenterclient.GetMe().UpdateRoom(room.RoomType(), room.ID(), room.GetPlayerNum(), common.UserOnline, false, 0, 0)
			}
		}
	}
	//glog.Info("[房间] 分配成功 [", room.RoomType(), ",", room.ID(), "],", player.ID(), ",", player.UData().Account, ",", player.UData().Model, ",", player.UData().IsCustom)
	return true
}

// 通知语音服创建房间
func (this *RoomMgr) SyncVoiceRoom() {
	timenow := time.Now().Unix()
	this.mutex.RLock()
	for _, v := range this.rooms {
		voiceclient.GetMe().NewRoom(v.ID(), uint32(v.EndTime()-timenow))
	}
	this.mutex.RUnlock()
}

func (this *RoomMgr) Final() {
	tuids, quids, cuids := this.getPlayerIDs()
	var rms []*redismgr.RoomTypeID
	rooms := GetMe().GetRooms()
	for _, room := range rooms {
		rms = append(rms, &redismgr.RoomTypeID{
			Type: room.RoomType(), ID: room.ID()})
	}

	redismgr.GetMe().AllUserOffline(tuids, quids, cuids, rms)
}
