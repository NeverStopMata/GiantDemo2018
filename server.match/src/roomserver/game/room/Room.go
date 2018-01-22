package room

// 房间类

import (
	"runtime/debug"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"base/env"
	"base/glog"
	"common"
	"roomserver/client/dbclient"
	"roomserver/client/rcenterclient"
	"roomserver/client/voiceclient"
	"roomserver/conf"
	"roomserver/game/bll"
	"roomserver/game/consts"
	"roomserver/game/room/internal"
	"roomserver/game/scn"
	"roomserver/game/scn/playermgr"
	"roomserver/game/scn/plr"
	tm "roomserver/game/team"
	ri "roomserver/interfaces"
	"roomserver/redismgr"
	"roomserver/udp"
	"usercmd"
)

const (
	MaxTopPlayer int = 9 //排行人数
)

const (
	ROOM_CONTROL_END  = 1 // 结束
	ROOM_CONTROL_STOP = 2 // 停止
)

// 麦上人信息
type SpeakerData struct {
	speakerId  uint64 // 麦上玩家id
	speakEndTm uint32 // 麦上结束时间
}

// XXX(Jinq): Scene从Room中分离出来
type Room struct {
	scn.Scene                                              // 继承场景信息
	roomType              uint32                           // 房间类型
	sceneId               uint32                           // 地图id
	id                    uint32                           // 房间id
	name                  string                           // 房间名
	isclosed              int32                            // 是否关闭标识
	endTime               int64                            // 结束时间
	startTime             int64                            // 开始时间
	IsCustom              bool                             // 是否自建房间
	playernum             int32                            // 在线人数
	freeRanks             common.FreeRankList              // 自由排行榜
	teamRanks             common.TeamRankList              // 组队排行榜
	speakermap            map[uint32]*SpeakerData          // 麦上玩家
	teamlist              map[uint32]*tm.Team              // 队伍列表
	chan_Control          chan int                         // 房间事件
	Chan_AddPlayer        chan ri.IPlayerTask              // 添加玩家
	chan_RemovePlayerById chan uint64                      // 删除玩家
	chan_PlayerCmd        chan *internal.PlayerCmd         // 玩家输入
	Chan_GetPlayerNum     chan chan []int                  // 其他协程获取房间内玩家数、机器人数
	chan_BindUdpSession   chan map[uint64]*udp.UdpSess     // 通知绑定udpsession
	isInited              int32                            // 是否有人进去了，初始完
	topplayers            [MaxTopPlayer]*usercmd.MsgPlayer // 排行榜
	firstSpeak            bool
	newLogin              map[uint64]bool
	newLoginMutex         sync.Mutex
}

// 创建房间
func NewRoom(sceneId uint32, rtype, rid uint32, iscustom bool, name string) *Room {
	room := Room{
		roomType:              rtype,
		sceneId:               sceneId,
		id:                    rid,
		name:                  name,
		isclosed:              -1,
		IsCustom:              iscustom,
		speakermap:            make(map[uint32]*SpeakerData),
		teamlist:              make(map[uint32]*tm.Team),
		Chan_AddPlayer:        make(chan ri.IPlayerTask, 32),
		chan_RemovePlayerById: make(chan uint64, 32),
		chan_PlayerCmd:        make(chan *internal.PlayerCmd, 1024),
		chan_Control:          make(chan int, 1),
		Chan_GetPlayerNum:     make(chan chan []int),
		chan_BindUdpSession:   make(chan map[uint64]*udp.UdpSess, 32),
		newLogin:              make(map[uint64]bool),
		firstSpeak:            true,
	}
	room.endTime = int64(room.GetRoomTypeConfig().LastTm) + time.Now().Unix()
	room.Pool = scn.NewPool()
	for index := 0; index < MaxTopPlayer; index++ {
		room.topplayers[index] = &usercmd.MsgPlayer{}
	}
	room.Start()
	return &room
}

//地图配置
func (this *Room) GetMapConfig() *conf.XmlRoom {
	return conf.ConfigMgr_GetMe().GetRoomById(this.sceneId)
}
func (this *Room) GetRoomTypeConfig() *conf.XmlRoomModel {
	return conf.ConfigMgr_GetMe().GetRoomType(int(this.roomType))
}

func (this *Room) GetNewLoginUsers() map[uint64]bool {
	this.newLoginMutex.Lock()
	defer this.newLoginMutex.Unlock()
	if this.firstSpeak {
		this.firstSpeak = false
		return nil
	}
	m := this.newLogin
	this.newLogin = make(map[uint64]bool)
	return m
}

func (this *Room) AddLoginUser(UID uint64) {
	this.newLoginMutex.Lock()
	defer this.newLoginMutex.Unlock()
	if this.firstSpeak {
		return
	}
	this.newLogin[UID] = true
}

// 开启房间
func (this *Room) Start() bool {
	if !atomic.CompareAndSwapInt32(&this.isclosed, -1, 0) {
		return false
	}

	// 按大小初始化房间大小
	this.Scene.Init(this)

	// 设置初始数据 (时间，战队)
	lasttime := this.GetRoomTypeConfig().LastTm

	if this.IsCustom {
		tlasttime := redismgr.GetMe().GetTempRoomTime(this.roomType, this.id)
		if tlasttime != 0 {
			lasttime = uint32(tlasttime)
		}
	}
	this.startTime = time.Now().Unix()
	this.endTime = time.Now().Unix() + int64(lasttime)

	// 开启逻辑处理协程
	go this.Loop()
	glog.Info("[房间] 创建房间 createroom[sid:", this.sceneId, ",type:", this.roomType, ",id:", this.id, "] ", this.IsCustom, ",", this.endTime, ",", lasttime, ",", this.name)
	return true
}

// 停止房间
func (this *Room) Stop() bool {
	if !atomic.CompareAndSwapInt32(&this.isclosed, 0, 1) {
		return false
	}
	close(this.Chan_AddPlayer)
	close(this.chan_RemovePlayerById)
	close(this.chan_Control)
	close(this.chan_PlayerCmd)
	close(this.Chan_GetPlayerNum)
	close(this.chan_BindUdpSession)
	glog.Info("[房间] 销毁房间 [", this.roomType, ",", this.id, "] ", len(this.Players), ",", this.startTime)
	return true
}

func (this *Room) IsClosed() bool {
	return atomic.LoadInt32(&this.isclosed) != 0
}

//主循环
func (this *Room) Loop() {
	var (
		timeTicker = time.NewTicker(time.Millisecond * consts.FrameTimeMS)
		tloop      uint64
	)

	defer func() {
		this.Stop()
		timeTicker.Stop()
		if err := recover(); err != nil {
			glog.Error("[异常] 房间线程出错 [", this.roomType, ",", this.id, "] ", err, "\n", string(debug.Stack()))
		}
	}()

	this.ResetFrame()

	for {
		select {
		case <-timeTicker.C:
			this.IncreaseFrame()
			this.Scene.Render()
			//100ms
			if tloop%consts.FrameCountBy100MS == 0 {
				this.Render5()
			}
			//400ms
			if tloop%(consts.FrameCountBy100MS*4) == 0 {
				this.SendTeamMemPos()
			}
			//1s
			if tloop%(consts.FrameCountBy100MS*10) == 0 {
				this.TimeAction()
			}
			if tloop%(consts.FrameCountBy100MS*40) == 0 {
				this.PushTopList(nil)
			}
			tloop += 1
		case op := <-this.chan_PlayerCmd:
			if !this.IsClosed() {
				player, ok := this.Players[op.PlayerID]
				if ok {
					player.OnRecvPlayerCmd(op.Cmd, op.Data, op.Flag)
				} else {
					glog.Info("chan_PlayerCmd:no player,", op.PlayerID, " cmd:", op.Cmd)
				}
			}
		case player := <-this.Chan_AddPlayer:
			this.AddPlayer(player)
		case playerId := <-this.chan_RemovePlayerById:
			this.RemovePlayerById(playerId)
		case ctrl := <-this.chan_Control:
			switch ctrl {
			case ROOM_CONTROL_END:
			case ROOM_CONTROL_STOP:
				glog.Info("[ctrl]", this.name, " ctrl: ", ctrl)
				this.destory()
			}
			return
		case c := <-this.Chan_GetPlayerNum:
			rn := int(this.GetRobotNum())
			pn := len(this.Players) - rn
			c <- []int{rn, pn}
		case datas := <-this.chan_BindUdpSession:
			for key, val := range datas {
				player, ok := this.Players[key]
				if ok && player.Sess != nil {
					player.Sess.BindUdpSession(val)
				}
			}
		}
	}
}

func (this *Room) PostBindUdpSession(data map[uint64]*udp.UdpSess) {
	this.chan_BindUdpSession <- data
}

// 获取麦上玩家
func (this *Room) getSpeakerId(tid uint32) uint64 {
	speaker, ok := this.speakermap[tid]
	if !ok {
		return 0
	}
	if speaker.speakEndTm < uint32(time.Now().Unix()) {
		speaker.speakerId = 0
	}
	return speaker.speakerId
}

// 设置麦上玩家
func (this *Room) setSpeakerId(tid uint32, userid uint64) {
	speaker := &SpeakerData{
		speakerId:  userid,
		speakEndTm: uint32(time.Now().Unix()) + 3600,
	}
	this.speakermap[tid] = speaker
	glog.Info("[房间] 设置麦上玩家speakEndTm:", speaker.speakEndTm)
}

// 获取剩余时间
func (this *Room) getSpeakRemainTm(tid uint32) uint32 {
	speaker, ok := this.speakermap[tid]
	if !ok {
		return 0
	}
	if speaker.speakEndTm < uint32(time.Now().Unix()) {
		return 0
	}
	return speaker.speakEndTm - uint32(time.Now().Unix())
}

// 语间服务器地址
func (this *Room) SendVoiceInfo(user ri.IPlayerTask) {
	retCmd := &usercmd.RetVoiceInfo{
		Address: env.Get("room", "voiceudp"),
		RoomId:  this.id,
	}
	user.SendCmd(usercmd.MsgTypeCmd_VoiceInfo, retCmd)
	glog.Info("[登录]返回语音服务器地址: ", user.Name(), " 地址是: ", retCmd.Address, "房间是:", retCmd.RoomId)
}

// 抢麦
func (this *Room) ReqToSpeak(playerId uint64, banVoice bool, newLogin map[uint64]bool) bool {
	player, ok := this.Players[playerId]
	if !ok || player.Sess == nil {
		return false
	}

	if this.getSpeakerId(player.UData().TeamId) != 0 {
		player.Sess.RetErrorMsg(common.ErrorCodeMicBusy)
		return false
	}

	go voiceclient.GetMe().ToSpeak(this.id, playerId, player.UData().TeamId, 3600, banVoice, newLogin)

	this.setSpeakerId(player.UData().TeamId, playerId)

	// 发送给说话者
	toSpeakCmd := &this.Pool.MsgToSpeak
	toSpeakCmd.Time = this.getSpeakRemainTm(player.UData().TeamId)
	player.SendCmd(usercmd.MsgTypeCmd_ToSpeak, toSpeakCmd)

	// 发送给所有人
	SpeakUserCmd := &this.Pool.MsgSpeakUser
	SpeakUserCmd.Id = player.ID
	SpeakUserCmd.Name = player.Name
	SpeakUserCmd.NeedTm = this.getSpeakRemainTm(player.UData().TeamId)
	if banVoice {
		SpeakUserCmd.IsBaned = 1
	} else {
		SpeakUserCmd.IsBaned = 0
	}
	for _, p := range this.Players {
		if p.UData().TeamId != player.UData().TeamId {
			continue
		}
		p.SendCmd(usercmd.MsgTypeCmd_SpeakUser, SpeakUserCmd)
	}

	glog.Info("[语音] 请求抢麦 [", this.roomType, ",", this.id, "] ", player.UData().Id, ",", player.UData().Account, ",", SpeakUserCmd.IsBaned)
	return true
}

//添加玩家
func (this *Room) AddPlayer(player ri.IPlayerTask) {

	p, ok := this.Players[player.ID()]
	if ok && p.Sess != nil {
		p.Sess.RetErrorMsg(common.ErrorCodeReLogin)
		p.Sess.Close()
		glog.Info("AddPlayer.....", p.Sess.RemoteAddrStr())
		p.Sess = nil
	}

	var team *tm.Team
	switch this.roomType {
	case common.RoomTypeTeam:
		if player.UData().TeamId == 0 {
			player.RetErrorMsg(common.ErrorCodeOutTeam)
			player.Close()
			glog.Error("[登录] 不在任何队伍 [", this.roomType, ",", this.id, "] ", player.UData().Id, ",", player.UData().Account)
			return
		}
		team, ok = this.teamlist[player.UData().TeamId]
		if !ok {
			team = &tm.Team{
				IsNewbie: player.UData().IsNewbie,
				MemList:  make(map[uint64]bool),
			}
			if player.UData().IsLeader {
				team.LeaderID = player.UData().Id
			}
			team.MemList[player.ID()] = player.UData().IsNewbie
			this.teamlist[player.UData().TeamId] = team
		} else {
			if player.UData().IsNewbie {
				team.IsNewbie = true
			}
			if player.UData().IsLeader {
				team.LeaderID = player.UData().Id
			}
			team.MemList[player.ID()] = player.UData().IsNewbie
		}
	}

	//根据房间人数关闭/开启机器人
	if this.roomType == common.RoomTypeTeam {
		this.Scene.DeleteTeamRoomRobot(player.ID(), player.UData().TeamId)
	} else {
		this.Scene.CheckRobotNum()
	}
	// 添加到场景
	this.Scene.AddPlayer(player, team, false, nil, 0, 0)

	// 语音服务器信息
	this.SendVoiceInfo(player)
	this.SetInited()
	glog.Info("[登录] 进入房间成功 [", this.roomType, ",", this.id, "],", player.ID(), ",", player.UData().Account, ",", this.GetPlayerNum(), ",", player.UData().Robot, ",set:", player.UData().Sex, " sceneId= ", this.sceneId)
}

func (this *Room) CheckInited() bool {
	return atomic.LoadInt32(&this.isInited) != 0
}

func (this *Room) SetInited() {
	atomic.CompareAndSwapInt32(&this.isInited, 0, 1)
}

// 从房间里删除玩家
func (this *Room) RemovePlayerById(playerId uint64) {
	scenePlayer, ok := this.Players[playerId]
	if !ok {
		return
	}

	playerTask := scenePlayer.Sess
	if playerTask != nil {
		if playerTask.ID() == playerId {
			scenePlayer.Sess = nil
		}
		if playerTask.IsTimeout() {
			scenePlayer.IsTimeout = true
		}
	}

	//退出房间处理
	this.Scene.RemovePlayer(scenePlayer.ID)
	playermgr.GetMe().Remove(scenePlayer)
	// 通知其它人删除玩家
	rmCmd := &usercmd.MsgRemovePlayer{
		Id: scenePlayer.ID,
	}
	this.BroadcastMsg(usercmd.MsgTypeCmd_RemovePlayer, rmCmd)

	glog.Info("[房间] 删除玩家 [timeout:", scenePlayer.IsTimeout, ",", this.roomType, ",", this.id, ",", this.name, "] ", playerId, ",", scenePlayer.UData().Account, ",", len(this.Players), ", roomID:", scenePlayer.UData().RoomId)
}

// 向房间协程发送指令
func (this *Room) Control(ctrl int) bool {
	if this.IsClosed() {
		return false
	}
	this.chan_Control <- ctrl
	return true
}

func (this *Room) IncPlayerNum(player ri.IPlayerTask) {
	atomic.AddInt32(&this.playernum, 1)
	glog.Info("[房间] 房间人数增加 [", this.roomType, ",", this.id, "] ", "当前:", this.GetPlayerNum())
}

func (this *Room) DecPlayerNum() {
	atomic.AddInt32(&this.playernum, -1)
	glog.Info("[房间] 房间人数减少 [", this.roomType, ",", this.id, "] ", "当前:", this.GetPlayerNum())
}

func (this *Room) GetPlayerNum() int32 {
	return atomic.LoadInt32(&this.playernum)
}

// 发送队友位置
func (this *Room) SendTeamMemPos() {
	if this.roomType != common.RoomTypeTeam {
		return
	}
	for tid, _ := range this.teamlist {
		info := this.GetTeamInfo(tid)
		this.BroadcastTeamMsg(tid, usercmd.MsgTypeCmd_UpdateTeamInfo, info)
	}
}

//GetTeam 获取队伍信息
func (this *Room) GetTeam(teamid uint32) *tm.Team {
	team, ok := this.teamlist[teamid]
	if !ok {
		return nil
	}
	return team
}

//队伍信息
func (this *Room) GetTeamInfo(teamid uint32) *usercmd.UpdateTeamInfoMsg {
	_, ok := this.teamlist[teamid]
	if !ok {
		return nil
	}

	data := usercmd.UpdateTeamInfoMsg{}
	var topPlayer *plr.ScenePlayer
	for pid, player := range this.Players {
		//player, _ := this.GetPlayer(pid)
		//队友
		if player.UData().TeamId == teamid {
			m := &usercmd.TeamInfoMsg{Playerid: pid, X: float32(player.SelfAnimal.Pos.X), Y: float32(player.SelfAnimal.Pos.Y)}
			data.Members = append(data.Members, m)
		}

		//top
		//glog.Info("loop top:", player.name, "  exp:", player.roomExp, "  rank:", player.rank)
		if this.roomType == common.RoomTypeTeam {
			if player.Rank == 1 {

				if topPlayer == nil || player.GetExp() > topPlayer.GetExp() {
					topPlayer = player
				}

			}
		}
	}

	//top
	if topPlayer != nil {
		top := &usercmd.TeamInfoMsg{Playerid: topPlayer.ID, X: float32(topPlayer.SelfAnimal.Pos.X), Y: float32(topPlayer.SelfAnimal.Pos.Y)}
		data.TopPlayers = append(data.TopPlayers, top)
		//glog.Info("update top:", topPlayer.name, "  exp:", topPlayer.roomExp)
	}

	return &data
}

// 广播结束消息
func (this *Room) BroadcastEndMsg() {
	retCmd := &usercmd.MsgEndRoom{}
	size := this.GetRoomTypeConfig().PlayerNum
	toids := make([]uint64, 0)
	teamids := make(map[uint32]uint32)
	leaderids := make(map[uint32]uint64)
	switch this.roomType {
	case common.RoomTypeQuick, common.RoomTypeTeam:
		tmpList := common.FreeRankList{}
		for _, p := range this.Players {
			exp := p.GetExp()
			if common.RoomTypeTeam != this.roomType {
				tmpList = append(tmpList, common.FreeRank{Id: p.ID, Score: float64(exp)})
			}
			if p.IsRobot == false {
				toids = append(toids, p.ID)
			}
			if common.RoomTypeTeam == this.roomType {
				teamids[p.UData().TeamId] = teamids[p.UData().TeamId] + uint32(exp)
				if p.UData().IsLeader {
					leaderids[p.UData().TeamId] = p.ID
				}
			}
			//glog.Info("BroadcastEndMsg ", p.ID, ",", p.UData().Account, ",", p.UData().TeamId, ",", p.UData().IsLeader)
		}
		for key, val := range teamids {
			tmpList = append(tmpList, common.FreeRank{Id: uint64(key), Score: float64(val)})
		}
		sort.Sort(tmpList)
		for index, v := range tmpList {
			// 团战模式
			if this.roomType == common.RoomTypeTeam {
				for _, p := range this.Players {
					if p.UData().TeamId == uint32(v.Id) {
						retCmd.Players = append(retCmd.Players, &usercmd.EndPlayer{
							Id:        p.ID,
							UName:     p.UData().Account,
							Name:      p.Name,
							Score:     uint64(p.GetExp()),
							KillNum:   p.KillNum,
							Sex:       uint32(p.UData().Sex),
							Icon:      p.UData().Icon,
							Location:  p.GetLocation(),
							Rank:      p.Rank,
							PassIcon:  p.UData().PassIcon,
							IsFollow:  false,
							AnimalId:  uint32(p.SelfAnimal.GetAnimalId()),
							TeamScore: teamids[p.UData().TeamId],
							TeamName:  p.UData().TeamName,
							LeaderId:  leaderids[p.UData().TeamId],
						})
						//glog.Info("BroadcastEndMsg2222 ", p.ID, ",", p.UData().Account, ",", p.UData().TeamId, ",", p.UData().IsLeader)
					}
				}
				continue
			}

			// 非团战模式
			p, ok := this.Players[v.Id]
			if !ok {
				continue
			}
			if this.roomType != common.RoomTypeTeam {
				p.Rank = uint32(index) + 1
			}
			if len(retCmd.Players) < int(size) {
				retCmd.Players = append(retCmd.Players, &usercmd.EndPlayer{
					Id:        p.ID,
					UName:     p.UData().Account,
					Name:      p.Name,
					Score:     uint64(p.GetExp()),
					KillNum:   p.KillNum,
					Sex:       uint32(p.UData().Sex),
					Icon:      p.UData().Icon,
					Location:  p.GetLocation(),
					Rank:      p.Rank,
					PassIcon:  p.UData().PassIcon,
					IsFollow:  false,
					AnimalId:  uint32(p.SelfAnimal.GetAnimalId()),
					TeamScore: teamids[p.UData().TeamId],
					TeamName:  p.UData().TeamName,
					LeaderId:  leaderids[p.UData().TeamId],
				})
			}
		}
	}
	for _, p := range this.Players {
		if p.Sess == nil {
			continue
		}
		retCmd.UserSelf = &usercmd.EndPlayer{
			Id:        p.ID,
			UName:     p.UData().Account,
			Name:      p.Name,
			Score:     uint64(p.GetExp()),
			KillNum:   p.KillNum,
			Sex:       uint32(p.UData().Sex),
			Icon:      p.UData().Icon,
			Location:  p.GetLocation(),
			Rank:      p.Rank,
			PassIcon:  p.UData().PassIcon,
			IsFollow:  true,
			AddMoney:  0,
			AnimalId:  uint32(p.SelfAnimal.GetAnimalId()),
			TeamScore: teamids[p.UData().TeamId],
			TeamName:  p.UData().TeamName,
			LeaderId:  leaderids[p.UData().TeamId],
		}
		//glog.Info("BroadcastEndMsg_333 ", p.ID, ",", p.UData().Account, ",", p.UData().TeamId, ",", p.UData().IsLeader)
		IsFollow := redismgr.GetMe().HasFollowed(p.ID, toids)
		for _, vp := range retCmd.Players {
			vp.IsFollow = IsFollow[vp.Id]
		}
		retCmd.Level = p.UData().Level
		p.SendCmd(usercmd.MsgTypeCmd_EndRoom, retCmd)
	}
}

// 删除
func (this *Room) destory() {

	// 关闭房间
	this.Stop()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				glog.Error("[异常] 报错 ", err, "\n", string(debug.Stack()))
			}
		}()
		// 更新排行榜
		this.PushTopList(nil)
		// 同步玩家数据/玩家处理房间结算
		if !this.IsCustom {
			this.SyncDatas()
		}
		// 广播房间结束
		this.BroadcastEndMsg()
		// 玩家数据结算保存
		this.SaveRoomDatas()
		// 重置玩家状态
		redismgr.GetMe().ClearUserRoom(this.roomType, this.id, this.getPlayerIDs(), this.IsCustom)
	}()
	// 清除scenePlayer
	playermgr.GetMe().Removes(this.Players)

	//glog.Info("[房间] 结算完成 [", this.roomType, ",", this.id, "] ", this.iscustom, ",", len(this.Players), ",", this.maxplayernum, ",", this.GetPlayerNum())
}

// getPlayerIDs 获取玩家ID列表
func (this *Room) getPlayerIDs() (ids []uint64) {
	for _, p := range this.Players {
		ids = append(ids, p.ID)
	}
	return ids
}

// 同步玩家数据
func (this *Room) SyncDatas() bool {
	for _, p := range this.Players {
		if false == p.IsRobot {
			p.DoEndRoom(this)
		}
	}
	glog.Info("[SyncDatas] 房间同步 roomid= ", this.id, " roomType= ", this.roomType)
	return true
}

func (this *Room) SaveRoomData(p *plr.ScenePlayer) {
	nowTime := time.Now()
	datas := &common.ReqRoomInc{}
	urdata := &common.RoomIncData{
		Id:        p.UData().Id,
		Location:  p.UData().Location,
		Icon:      p.UData().Icon,
		NickName:  p.Name,
		Rank:      p.Rank,
		KillNum:   p.KillNum,
		AddExp:    uint32(p.GetExp()),
		GameTime:  uint32((nowTime.Unix() - p.StartTime.Unix())),
		IsEndRoom: false,
	}

	datas.UData = append(datas.UData, urdata)
	datas.RoomType = this.roomType
	datas.RoomId = this.id
	dbclient.GetMe().RefreshRoomInc(datas)
}

// 保存结算数据 (自由/团队)
func (this *Room) SaveRoomDatas() bool {
	if this.IsCustom {
		return false
	}
	nowTime := time.Now()
	datas := &common.ReqRoomInc{}
	for _, p := range this.Players {
		if p.IsRobot {
			continue
		}
		urdata := &common.RoomIncData{
			Id:        p.UData().Id,
			Location:  p.UData().Location,
			IsEndRoom: true,
			Icon:      p.UData().Icon,
			NickName:  p.Name,
			Rank:      p.Rank,
			Level:     p.UData().Level,
			Scores:    p.UData().Scores,
			KillNum:   p.KillNum,
			AddExp:    uint32(p.GetExp()),
			GameTime:  uint32((nowTime.Unix() - p.StartTime.Unix())),
			AnimalID:  uint32(p.SelfAnimal.GetAnimalId()),
			SceneId:   this.sceneId,
		}
		datas.UData = append(datas.UData, urdata)
	}

	if this.roomType == common.RoomTypeTeam {
		for _, team := range this.teamlist {
			var memlist []uint64
			for memid, _ := range team.MemList {
				memlist = append(memlist, memid)
			}
			datas.TData = append(datas.TData, memlist)
		}
	}
	datas.RoomType = this.roomType
	datas.RoomId = this.id
	dbclient.GetMe().RefreshRoomInc(datas)
	return true
}

// 检查玩家是否上排行榜
func (this *Room) inTopList(uids ...uint64) bool {
	for _, p := range this.freeRanks {
		for _, u := range uids {
			if u == p.Id {
				return true
			}
		}
	}
	return false
}

// 发送排行榜
func (this *Room) PushTopList(player ri.IPlayerTask) {
	switch this.roomType {
	case common.RoomTypeQuick:
		ltime := this.endTime - time.Now().Unix()
		if ltime > 0 {
			this.Pool.MsgTopRank.EndTime = uint32(ltime)
		} else {
			this.Pool.MsgTopRank.EndTime = 0
		}
		this.Pool.MsgTopRank.Players = []*usercmd.MsgPlayer{}
		tmpList := common.FreeRankList{}
		for _, p := range this.Players {
			tmpList = append(tmpList, common.FreeRank{Id: p.ID, Score: float64(p.GetExp())})
		}
		sort.Sort(tmpList)
		for index, v := range tmpList {
			p, ok := this.Players[v.Id]
			if !ok {
				continue
			}
			if len(this.Pool.MsgTopRank.Players) < MaxTopPlayer {
				this.topplayers[index].Id = v.Id
				this.topplayers[index].Name = p.Name
				this.topplayers[index].Curexp = uint32(v.Score)
				this.topplayers[index].Local = p.GetLocation()
				this.Pool.MsgTopRank.Players = append(this.Pool.MsgTopRank.Players, this.topplayers[index])
			}
		}
		for k, v := range tmpList {
			p, ok := this.Players[v.Id]
			if !ok {
				continue
			}
			p.Rank = uint32(k + 1)
			this.Pool.MsgTopRank.Rank = p.Rank
			this.Pool.MsgTopRank.KillNum = uint32(v.Score)
			p.SendCmd(usercmd.MsgTypeCmd_Top, &this.Pool.MsgTopRank)
		}
	case common.RoomTypeTeam:
		teamids := make(map[uint32]uint32)
		teamexps := make(map[uint32]uint32)
		m_team_name := make(map[uint32]uint32)
		for _, p := range this.Players {
			teamexps[p.UData().TeamId] = teamexps[p.UData().TeamId] + p.GetExp()
			teamids[p.UData().TeamId] = teamids[p.UData().TeamId] + 1
			m_team_name[p.UData().TeamId] = p.UData().TeamName
		}
		tmpList := common.FreeRankList{}
		retCmd := &usercmd.RetTeamRankList{}
		for key, _ := range teamids {
			tmpList = append(tmpList, common.FreeRank{Id: uint64(key), Score: float64(teamexps[key])})
		}
		sort.Sort(tmpList)
		for _, v := range tmpList {
			retCmd.Teams = append(retCmd.Teams, &usercmd.RetTeamRankList_TeamRank{
				Tname: m_team_name[uint32(v.Id)],
				Num:   teamids[uint32(v.Id)],
				Score: float64(v.Score),
			})
		}
		for _, p := range this.Players {
			for index, _ := range tmpList {
				if p.UData().TeamId == uint32(tmpList[index].Id) {
					p.Rank = uint32(index) + 1
				}
			}
		}
		ltime := this.endTime - time.Now().Unix()
		if ltime > 0 {
			retCmd.EndTime = uint32(ltime)
		} else {
			retCmd.EndTime = 0
		}
		for _, c := range this.Players {
			retCmd.KillNum = c.GetExp()
			c.SendCmd(usercmd.MsgTypeCmd_TeamRankList, retCmd)
		}
	default:
		glog.Error("[排名] 未知房间类型 ", this.roomType, ",", this.id)
	}
}

// 广播 msg
func (this *Room) BroadcastMsg(msgNo usercmd.MsgTypeCmd, msg common.Message) {
	data, flag, err := common.EncodeGoCmd(uint16(msgNo), msg)
	if err != nil {
		glog.Error("[广播] 发送消息 ", msgNo, " 失败!")
		return
	}
	for _, c := range this.Players {
		c.AsyncSend(data, flag)
	}
}

// 队伍广播 msg
func (this *Room) BroadcastTeamMsg(teamid uint32, msgNo usercmd.MsgTypeCmd, msg common.Message) {
	data, flag, err := common.EncodeGoCmd(uint16(msgNo), msg)
	if err != nil {
		glog.Error("[广播] 发送team消息 ", msgNo, " 失败!")
		return
	}

	for _, c := range this.Players {
		if c.UData().TeamId == teamid {
			c.AsyncSend(data, flag)
		}
	}
}

//广播(剔除特定ID)
func (this *Room) BroadcastMsgExcept(msgNo usercmd.MsgTypeCmd, msg common.Message, uid uint64) {
	data, flag, err := common.EncodeGoCmd(uint16(msgNo), msg)
	if err != nil {
		glog.Error("[广播] 发送消息 ", msgNo, " 失败!")
		return
	}
	for _, c := range this.Players {
		if c.ID == uid {
			continue
		}
		c.AsyncSend(data, flag)
	}
}

// 房间定时器 (一秒一次)
func (this *Room) TimeAction() bool {
	timeNow := time.Now()
	nowmsec := time.Now().UnixNano() / 1000000
	var flag = false
	// glog.Info("-----------------roomtype----------", this.roomType)
	if this.roomType == common.RoomTypeTeam {
		this.Scene.AddTeamRoomRobotNum()
	} else {
		this.Scene.CheckRobotNum()
	}
	for _, player := range this.Players {
		player.TimeAction(this, timeNow)
	}
	if this.roomType == common.RoomTypeQuick {
		return true
	}
	for _, player := range this.Players {
		if false == player.IsLive && false == player.IsRobot && 0 != player.DeadTime() && player.DeadTime()+60 < timeNow.Unix() {
			glog.Info("TimeAction _ ", player.UData().Account, ", ", player.ID, " , ", player.IsLive)
			this.Scene.DeleteActOffline(player.UData().Account)
			redismgr.GetMe().SetFRoomId(player.ID, 0, 0)
			delete(this.Players, player.ID)
			flag = true
		}
		if true == player.IsRobot && player.AI.GetExpireTime() != 0 && player.AI.GetExpireTime() < nowmsec {
			glog.Info("[房间] 分身机器人移除")
			this.RemovePlayerById(player.ID)
		}
	}
	if true == flag {
		rcenterclient.GetMe().UpdateRoom(this.roomType, this.id, this.GetPlayerNum(), common.UserOffline, false, 0, 0)
	}
	return true
}

// PostPlayerCmd 发送玩家命令到 chan_PlayerCmd. 命令在房间协程中执行。
func (this *Room) PostPlayerCmd(playerID uint64, cmd usercmd.MsgTypeCmd,
	data []byte, flag byte) {

	playerCmd := &internal.PlayerCmd{PlayerID: playerID, Cmd: cmd, Flag: flag}
	// Must copy data.
	playerCmd.Data = make([]byte, len(data))
	copy(playerCmd.Data, data)
	this.chan_PlayerCmd <- playerCmd
} // PostPlayerCmd()

// ResetPlayerTask 在房间中查找玩家，并将玩家的PlayerTask置空。
func (this *Room) ResetPlayerTask(playerTaskID uint64) {
	player, _ := this.GetPlayer(playerTaskID)
	if player != nil {
		player.Sess = nil
	}
}

func (this *Room) RoomType() uint32 {
	return this.roomType
}

func (this *Room) ID() uint32 {
	return this.id
}

// PostToRemovePlayerById 计划移除玩家，将在房间协程中执行动作。
func (this *Room) PostToRemovePlayerById(playerID uint64) {
	this.chan_RemovePlayerById <- playerID
}

// 判断队伍成员个数是否小于某数
func (this *Room) IsTeamMemberLessThan(teamID uint32, count int) bool {
	if count <= 0 {
		return false
	}
	for _, p := range this.Players {
		if p.UData().TeamId != teamID {
			continue
		}
		count--
		if count <= 0 {
			return false // 可提早中断遍历
		}
	}
	return true
}

func (this *Room) GetPlayerIScene() plr.IScene {
	return &this.Scene
}

func (this *Room) GetBallIScene() bll.IScene {
	return &this.Scene
}

func (this *Room) NewScenePlayer(udata *common.UserData, name string, isRobot bool) *plr.ScenePlayer {
	return plr.NewScenePlayer(udata, name, this, isRobot)
}

func (this *Room) StartTime() int64 {
	return this.startTime
}

func (this *Room) EndTime() int64 {
	return this.endTime
}

func (this *Room) Name() string {
	return this.name
}

func (this *Room) SceneID() uint32 {
	return this.sceneId
}
