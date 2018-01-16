package redismgr

import (
	"strconv"

	"base/env"
	"base/glog"
	"base/redis"
	"common"
	"roomserver/client/rcenterclient"
	"usercmd"
)

var (
	MAX_RELATION_TIME = 600 // 与我相关过期时间
	MAX_COMMENT_TIME  = 600 // 新手评价时间
)

type RedisMgr struct {
}

type RoomTypeID struct {
	Type uint32
	ID   uint32
}

var redismgr *RedisMgr

func GetMe() *RedisMgr {
	if redismgr == nil {
		redismgr = &RedisMgr{}
	}
	return redismgr
}

func (this *RedisMgr) Init() bool {
	r0addr := env.Get("global", "redis0")
	err := redis.RegRedisKey(common.TokenRedis, r0addr)
	if err != nil {
		glog.Error("[REDIS] 连接失败 ", r0addr, ",", err)
		return false
	}
	r1addr := env.Get("global", "redis1")
	err = redis.RegRedisKey(common.PlayerRedis, r1addr)
	if err != nil {
		glog.Error("[REDIS] 连接失败 ", r1addr, ",", err)
		return false
	}
	r2addr := env.Get("global", "redis2")
	err = redis.RegRedisKey(common.CacheRedis, r2addr)
	if err != nil {
		glog.Error("[REDIS] 连接失败 ", r2addr, ",", err)
		return false
	}
	return true
}

func Final() {
	redis.CloseAll()
}

// 获取token数据
func (this *RedisMgr) LoadFromRedis(key string, obj interface{}) bool {
	op := redis.NewRedisOp(common.TokenRedis)
	if op == nil {
		glog.Error("[REDIS] 打开redis失败", common.TokenRedis)
		return false
	}
	if !op.Exist(key) {
		glog.Error("[REDIS] key不存在", key)
		return false
	}
	err := op.GetObject(key, obj)
	if err != nil {
		glog.Error("[REDIS] 获取数据失败", key, ",", err)
		return false
	}
	if !op.Del(key) {
		glog.Error("[REDIS] 删除失败", key)
		return false
	}
	return true
}

// 获取自建房间时间
func (this *RedisMgr) GetTempRoomTime(roomtype, roomid uint32) int32 {
	op := redis.NewRedisOp(common.PlayerRedis)
	if op == nil {
		glog.Error("[REDIS] 获取opt失败", common.PlayerRedis)
		return 0
	}
	tmproom := getTmpRoomKey(roomtype, roomid)
	if !op.Exist(tmproom) {
		return 0
	}
	lasttime, _ := redis.Int32(op.Do("TTL", tmproom))
	return lasttime
}

// 玩家上线
func (this *RedisMgr) UserOnline(key string, userid uint64, serverid uint16, roomtype, roomid, model, sceneID uint32, iscustom bool) bool {
	op := redis.NewRedisOp(common.PlayerRedis)
	if op == nil {
		glog.Error("[REDIS] 获取opt失败", common.PlayerRedis)
		return false
	}
	player := "player:" + strconv.FormatInt(int64(userid), 10)
	userstate, _ := redis.Int(op.GetField(player, "State"))
	switch userstate {
	case common.UserStateFPlaying, common.UserStateTPlaying, common.UserStateQPlaying:
		oldserverid, _ := redis.Int(op.GetField(player, "RServerId"))
		if oldserverid != 0 && uint16(oldserverid) != rcenterclient.GetMe().Id {
			oldkey, err := redis.String(op.GetField(player, "Key"))
			if err != nil {
				glog.Error("[REDIS] 设置失败", key, ",", userid, ",", roomid, ",", err)
				return false
			}
			retCmd := &usercmd.ReqCheckRelogin{
				Key: oldkey,
				Id:  userid,
			}
			rcenterclient.GetMe().SendCmdToServer(uint16(oldserverid), usercmd.CmdType_ChkReLogin, retCmd)
		}
	default:
		switch roomtype {
		case common.RoomTypeTeam:
			op.SetFields(player, "State", common.UserStateTPlaying)
		case common.RoomTypeQuick:
			op.SetFields(player, "State", common.UserStateQPlaying)
		}
	}
	if !iscustom {
		if model != common.UserModelWatch {
			switch roomtype {
			case common.RoomTypeTeam:
				op.SetFields(player, "RServerId", serverid, "Model", common.UserModelTeam)
			case common.RoomTypeQuick:
				op.SetFields(player, "RServerId", serverid, "Model", common.UserModelQuick)
			}
		} else {
			op.SetFields(player, "RServerId", serverid, "Model", common.UserModelWatch)
		}
	} else {
		op.SetFields(player, "RServerId", serverid, "Model", common.UserModelCustom)
	}
	op.SetFields(player, "Key", key)
	op.SetExpire(player, 3600*7)
	//glog.Info("[REDIS] 玩家上线 ", key, ",", userid, ",", serverid, ",", roomid, ",", model)
	return true
}

// 玩家下线
func (this *RedisMgr) UserOffline(userid uint64) bool {
	op := redis.NewRedisOp(common.PlayerRedis)
	if op == nil {
		glog.Error("[REDIS] 获取opt失败", common.PlayerRedis)
		return false
	}
	player := "player:" + strconv.FormatInt(int64(userid), 10)
	if !op.Exist(player) {
		glog.Error("[REDIS] 玩家数据不存在 ", userid)
		return false
	}
	fields, err := redis.Ints(op.GetFields(player, "Model", "State"))
	if err != nil {
		glog.Error("[REDIS] 失败", userid, ",", err)
		return false
	}
	if len(fields) != 2 {
		glog.Error("[REDIS] 数量错误 ", userid, ",", fields)
		return false
	}
	if fields[0] == int(common.UserModelWatch) {
		op.SetField(player, "Model", 0)
	}
	switch fields[1] {
	case common.UserStateFPlaying, common.UserStateTPlaying, common.UserStateQPlaying:
		op.SetField(player, "State", common.UserStateOnline)
	}
	return true
}

// 清空房间中玩家
func (this *RedisMgr) ClearUserRoom(roomtype, roomid uint32, playerids []uint64, iscustom bool) bool {
	op := redis.NewRedisOp(common.PlayerRedis)
	if op == nil {
		glog.Error("[REDIS] 获取opt失败", common.PlayerRedis)
		return false
	}
	for _, pid := range playerids {
		if !iscustom {
			player := "player:" + strconv.FormatInt(int64(pid), 10)
			switch roomtype {
			case common.RoomTypeQuick:
				{
					op.SetFields(player, "QRoomId", 0, "State", common.UserStateOnline)
				}
			}
		} else {
			op.Del("cplayer:" + strconv.FormatInt(int64(pid), 10))
			op.SetField("player:"+strconv.FormatInt(int64(pid), 10), "State", common.UserStateOnline)
		}
	}
	glog.Info("[REDIS] 清除玩家房间状态 ", len(playerids))
	return true
}

func (this *RedisMgr) SetFRoomId(userid uint64, roomid uint32, expir int) {
	op := redis.NewRedisOp(common.PlayerRedis)
	if op == nil {
		glog.Error("[REDIS] 获取opt失败", common.PlayerRedis)
		return
	}
	player := "player:" + strconv.FormatInt(int64(userid), 10)
	op.SetField(player, "RoomId", roomid)
	//op.SetExpire(player, expir)
}

// 所有玩家下线
func (this *RedisMgr) AllUserOffline(
	tuids, quids, cuids []uint64,
	rooms []*RoomTypeID) bool {

	op := redis.NewRedisOp(common.PlayerRedis)
	if op == nil {
		glog.Error("[REDIS] 获取opt失败", common.PlayerRedis)
		return false
	}
	// DEL
	//	for _, uid := range fuids {
	//		keyplayer := "player:" + strconv.FormatInt(int64(uid), 10)
	//		op.SetFields(keyplayer, "State", common.UserStateOnline, "Model", 0, "RServerId", 0, "RoomId", 0, "RoomOwner", 0)
	//	}
	for _, uid := range tuids {
		keyplayer := "player:" + strconv.FormatInt(int64(uid), 10)
		op.SetFields(keyplayer, "State", common.UserStateOnline, "Model", 0, "RServerId", 0, "TRoomId", 0, "TeamId", 0, "TeamName", 0, "LeaderId", 0, "IsLeader", false, "IsNewbie", false, "RoomOwner", 0)
	}
	for _, uid := range quids {
		keyplayer := "player:" + strconv.FormatInt(int64(uid), 10)
		op.SetFields(keyplayer, "State", common.UserStateOnline, "Model", 0, "RServerId", 0, "QRoomId", 0, "RoomOwner", 0)
	}
	for _, uid := range cuids {
		op.Del("cplayer:" + strconv.FormatInt(int64(uid), 10))
	}
	for _, room := range rooms {
		op.Del(getTmpRoomKey(room.Type, room.ID))
	}
	glog.Info("[REDIS] 所有玩家下线完成 ")
	return true
}

// 获取屏蔽语音时间
func (this *RedisMgr) IsBanVoice(toid uint64) (bool, bool) {
	op := redis.NewRedisOp(common.CacheRedis)
	if op == nil {
		glog.Error("[REDIS] 打开redis失败")
		return false, false
	}
	banvoiceday := "banvoiceday:" + strconv.FormatUint(toid, 10)
	if op.Exist(banvoiceday) {
		return true, false
	}
	banvoicemin := "banvoicemin:" + strconv.FormatUint(toid, 10)
	return false, op.Exist(banvoicemin)
}

func (this *RedisMgr) HasFollowed(userid uint64, toids []uint64) (isfollow map[uint64]bool) {
	op := redis.NewRedisOp(common.PlayerRedis)
	if op == nil {
		glog.Error("[REDIS] 获取opt失败")
		return
	}
	isfollow = make(map[uint64]bool)
	following := "following:" + strconv.FormatUint(userid, 10)
	for _, toid := range toids {
		isfriend, _ := op.Do("ZSCORE", following, toid)
		if isfriend != nil {
			isfollow[toid] = true
		}
	}
	return
}

func getTmpRoomKey(roomtype, roomid uint32) string {
	return "tmproom:" + strconv.FormatInt(int64(roomtype), 10) + ":" + strconv.FormatInt(int64(roomid), 10)
}
