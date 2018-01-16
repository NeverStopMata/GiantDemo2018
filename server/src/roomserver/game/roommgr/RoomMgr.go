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
}

var (
	roommgr      *RoomMgr
	roommgr_once sync.Once
)

func GetMe() *RoomMgr {
	if roommgr == nil {
		roommgr_once.Do(func() {
			roommgr = &RoomMgr{
				rooms: make(map[uint32]*rm.Room),
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

func (this *RoomMgr) AddPlayer(player interfaces.IPlayerTask) bool {
	room := this.getRoomById(player.UData().RoomId)
	switch player.UData().Model {
	case common.UserModelJoin:
		if room == nil {
			player.RetErrorMsg(common.ErrorCodeRoom)
			glog.Error("[房间] 加入的房间不存在 ", player.UData().Id, ",", player.UData().Account, ",", player.UData().RoomId)
			return false
		}
		glog.Info("[房间] 加入模式 [", room.RoomType(), ",", room.ID(), "],[", player.UData().Id, ",", player.UData().Account, "]")
	case common.UserModelTeam:
		if room == nil {
			room = this.NewRoom(player.UData().SceneId, common.RoomTypeTeam, player.UData().RoomId, player, 0)
		}
		glog.Info("[房间] 组队模式 [", room.RoomType(), ",", room.ID(), "],[", player.UData().Id, ",", player.UData().Account, "],[", player.UData().TeamId, ",", player.UData().TeamName, "],", room.GetPlayerNum())
	case common.UserModelQuick:
		if room == nil {
			room = this.NewRoom(player.UData().SceneId, common.RoomTypeQuick, player.UData().RoomId, player, 0)
		}
		glog.Info("[房间] 闪电模式 [", room.RoomType(), ",", room.ID(), "],[", player.UData().Id, ",", player.UData().Account, "],[", player.UData().TeamId, ",", player.UData().TeamName, "],", room.GetPlayerNum())
	default:
		if room == nil {
			room = this.NewRoom(player.UData().SceneId, common.RoomTypeQuick, player.UData().RoomId, player, 0)
		}
		glog.Info("[房间] 闪电模式 [", room.RoomType(), ",", room.ID(), "],[", player.UData().Id, ",", player.UData().Account, "],[", player.UData().TeamId, ",", player.UData().TeamName, "],", room.GetPlayerNum())
	}

	if !room.IsClosed() {
		player.SetRoom(room)
		room.IncPlayerNum(player)
		room.Chan_AddPlayer <- player
		if !room.IsCustom {
			rcenterclient.GetMe().UpdateRoom(room.RoomType(), room.ID(), room.GetPlayerNum(), common.UserOnline, false, 0, 0)
		}
	}

	glog.Info("[房间] 分配成功 [", room.RoomType(), ",", room.ID(), "],", player.ID(), ",", player.UData().Account, ",", player.UData().Model, ",", player.UData().IsCustom)
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
