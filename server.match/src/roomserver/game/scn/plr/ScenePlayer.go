// 包 plr 处理玩家类相关功能。
package plr

// 玩家类

import (
	"base/glog"
	"common"
	"math"
	"roomserver/conf"
	"roomserver/game/bll"
	"roomserver/game/cll"
	"roomserver/game/consts"
	"roomserver/game/interfaces"
	"roomserver/game/scn/plr/internal"
	ri "roomserver/interfaces"
	"roomserver/util"
	"time"
	"usercmd"
)

type ScenePlayer struct {
	MoveHelper                                        // 检查移动消息包的辅助类
	ScenePlayerViewHelper                             // 玩家视野相关辅助类
	ScenePlayerNetMsgHelper                           // 房间玩家协议处理辅助类
	ScenePlayerPool                                   // 对象池
	Sess                    ri.IPlayerTask            // 网络会话（只有真实玩家有、robot没有）。 Scene.AddPlayer() 中设置
	room                    IRoom                     // 所在房间
	ID                      uint64                    // 玩家id
	BallId                  uint32                    // 玩家球id（一次定义，后面不变）
	Name                    string                    // 玩家昵称
	Key                     string                    // 登录的Key
	SelfAnimal              *bll.BallPlayer           // 玩家球
	KillNum                 uint32                    // 击杀数量
	udata                   *common.UserData          // 玩家主数据
	Rank                    uint32                    // 结算排名
	StartTime               time.Time                 // 进入房间时间
	IsLive                  bool                      // 生死
	Skill                   interfaces.ISkillPlayer   // 技能信息
	AI                      interfaces.IScenePlayerAI // AI信息
	IsRobot                 bool                      // 是否机器人
	isMoved                 bool                      // 是否移动过
	isRunning               bool                      // 当前是否在奔跑
	IsActClose              bool                      // 主动断开socket
	IsTimeout               bool                      // 是否是超时
	deadTime                int64                     // 死亡时间
	OldAnimalID             int32                     // 旧的动物编号（重连时需要用到）
	Reconn                  bool                      // 重连
	msgPool                 *internal.MsgPool
}

// NewIScenePlayerAI = ai.NewIScenePlayerAI
var NewIScenePlayerAI func(player *ScenePlayer) interfaces.IScenePlayerAI

// NewISkillPlayer = skill.NewISkillPlayer
var NewISkillPlayer func(player *ScenePlayer) interfaces.ISkillPlayer

// NewISkillBall = skill.NewISkillBall
var NewISkillBall func(player *ScenePlayer, ball *bll.BallSkill) interfaces.ISkillBall

func NewScenePlayer(udata *common.UserData, name string, room IRoom, isRobot bool) *ScenePlayer {
	p := &ScenePlayer{
		room:      room,
		udata:     udata,
		ID:        udata.Id,
		Name:      name,
		StartTime: time.Now(),
		IsLive:    true,
		IsRobot:   isRobot,
		msgPool:   internal.NewMsgPool(),
	}
	p.AI = NewIScenePlayerAI(p)
	p.Init()
	return p
}
func (this *ScenePlayer) Init() {
	this.ScenePlayerPool.Init()
	this.ScenePlayerNetMsgHelper.Init(this)
	this.ScenePlayerViewHelper.Init()
	this.BallId = this.GetScene().NewBallPlayerId()
	this.Skill = NewISkillPlayer(this)
	this.SelfAnimal = bll.NewBallPlayer(this, this.BallId)
	this.GetScene().AddBall(this.SelfAnimal)
	if this.IsRobot || conf.ConfigMgr_GetMe().Global.Pystress != 0 {
		this.AI.InitAI()
	}
	this.SelfAnimal.SetHP(consts.DefaultMaxHP)
	this.SelfAnimal.SetMP(consts.DefaultMaxMP)

}

func (this *ScenePlayer) SendChat(str string) {
	op := &usercmd.MsgSceneChat{Id: this.ID, Msg: str}
	this.BroadCastMsg(usercmd.MsgTypeCmd_SceneChat, op)
}

// 释放技能
func (this *ScenePlayer) CastSkill(op *usercmd.MsgCastSkill) {
	this.Skill.CastSkill(op.Skillid, this.Face)
}

func (this *ScenePlayer) Run(op *usercmd.MsgRun) {
	if this.isRunning {
		return
	}
	if this.Power == 0 {
		return
	}
	if this.Skill.GetCurSkillId() != 0 {
		return
	}
	this.isRunning = true
}

// 移动
func (this *ScenePlayer) Move(power, angle float64, face uint32) {
	if power != 0 {
		power = 1 // power恒为1,减少移动同步影响因素
	}
	this.Power = power
	this.Face = face
	if power != 0 {
		this.Angle = angle
	}
	if power == 0 {
		this.isRunning = false
	}
}

func (this *ScenePlayer) ClientCloseSocket() {
	if this.IsActClose == true {
		return
	}
	this.IsActClose = true
	this.room.DeleteActOffline(this.udata.Account)
	this.IsLive = false
	cmd := &usercmd.ClientHeartBeat{}
	this.SendCmd(usercmd.MsgTypeCmd_ActCloseSocket, cmd)
}

//死了
func (this *ScenePlayer) Dead(killer *ScenePlayer) {
	this.RealDead(killer)
}

//clearExp 是否清理经验(机器人使用)，默认false
func (this *ScenePlayer) RealDead(killer *ScenePlayer) {
	if false != this.IsRobot {
		this.AI.OnDeadEvent(killer)
	}

	if killer != nil {
		killer.UpdateExp(consts.DefaultBallPlayerExp)
	}

	msg := &this.msgPool.MsgDeath
	msg.MaxScore = uint32(this.GetExp())
	msg.Id = this.ID
	msg.Animalid = uint32(this.SelfAnimal.GetAnimalId())
	if killer == nil {
		msg.KillId = 0
		msg.KillName = ""
		msg.Killanimalid = 0
	} else {
		killer.KillNum++
		msg.KillId = killer.ID
		msg.KillName = killer.Name
		msg.Killanimalid = uint32(killer.SelfAnimal.GetAnimalId())
	}
	this.BroadCastMsg(usercmd.MsgTypeCmd_Death, msg)
	this.OnDead()
	if false == this.IsRobot {
		this.room.SaveRoomData(this)
	}
	if false == this.IsActClose {
		this.GetScene().AddOffline(this)
	} else {
		this.room.DeleteActOffline(this.udata.Account)
	}
}

func (this *ScenePlayer) OnDead() {
	this.SelfAnimal.OnDead()
	this.IsLive = false
	this.GetScene().RemoveBall(this.SelfAnimal) //移除
	this.GetScene().RemoveAnimalPhysic(this.SelfAnimal.PhysicObj)
}

func (this *ScenePlayer) GetRelifeMsg() *usercmd.MsgS2CRelife {
	msg := &this.msgPool.MsgRelife
	msg.Name = this.Name
	msg.Frame = this.GetScene().Frame()
	msg.SnapInfo = this.GetSnapInfo()
	msg.Curmp = uint32(this.SelfAnimal.GetMP())
	msg.Animalid = uint32(this.SelfAnimal.GetAnimalId())
	msg.Curhp = uint32(this.SelfAnimal.GetHP())
	return msg
}

// 复活
func (this *ScenePlayer) Relife() {
	// 分身机器人不会复活
	if true == this.IsRobot {
		if this.AI.GetExpireTime() != 0 {
			return
		}
	}

	if true == this.IsLive {
		return
	}
	if common.RoomTypeTeam == this.room.RoomType() && this.IsRobot == true {
		if !this.room.IsTeamMemberLessThan(this.udata.TeamId, 6) {
			this.room.RemovePlayerById(this.ID)
			return
		}
	}

	this.CleanPower()
	this.IsLive = true
	this.deadTime = 0
	this.IsActClose = false

	scene := this.GetScene()
	// 添加一个新的玩家球
	exp := this.GetExp()
	animalId := this.SelfAnimal.GetAnimalId()

	ball := bll.NewBallPlayer(this, this.BallId)
	this.SelfAnimal = ball
	this.SetExp(exp)
	this.SelfAnimal.SetAnimalId(animalId)

	scene.AddBall(this.SelfAnimal) //添加一个新的

	this.SendRoundMsg(usercmd.MsgTypeCmd_ReLife, this.GetRelifeMsg()) //通知复活
	if this.IsRobot {
		return
	}

	// 清除视野大小，设置新视野
	this.UpdateView(scene)
	this.UpdateViewPlayers(scene)
	this.ResetMsg()

	// 玩家视野中的所有球，发送给自己
	newMsg := &this.msgPool.MsgSceneTCP
	newMsg.Reset()
	this.LookFeeds = make(map[uint32]*bll.BallFeed)
	addfeeds, _ := this.UpdateVeiwFeeds()
	newMsg.Adds = append(newMsg.Adds, addfeeds...)

	this.LookBallSkill = make(map[uint32]*bll.BallSkill)
	adds, _ := this.UpdateVeiwBallSkill()
	newMsg.Adds = append(newMsg.Adds, adds...)

	this.LookBallFoods = make(map[uint32]*bll.BallFood)
	addfoods, _ := this.UpdateVeiwFoods()
	newMsg.Adds = append(newMsg.Adds, addfoods...)

	newMsg.AddPlayers = append(newMsg.AddPlayers, bll.PlayerBallToMsgBall(this.SelfAnimal))

	for _, other := range this.Others {
		newMsg.AddPlayers = append(newMsg.AddPlayers, bll.PlayerBallToMsgBall(other.SelfAnimal))
	}

	this.SendCmd(usercmd.MsgTypeCmd_SceneTCP, newMsg)

	this.RefreshPlayer()
}

// 获取位置
func (this *ScenePlayer) GetLocation() uint32 {
	if this.udata.ShowPos == 0 {
		return 0
	}
	return this.udata.Location
}

func (this *ScenePlayer) ResetMsg() {
	this.ScenePlayerPool.ResetMsg()
	this.ScenePlayerViewHelper.ResetMsg()
	this.isMoved = false
}

func (this *ScenePlayer) SendSceneMsg() {

	var (
		Eats          []*usercmd.BallEat
		Adds          []*usercmd.MsgBall
		AddPlayers    []*usercmd.MsgPlayerBall
		Moves         []*usercmd.BallMove
		Hits          []*usercmd.HitMsg
		Removes       []uint32
		RemovePlayers []uint32
	)

	//feed的添加删除消息单独处理
	addfeeds, delfeeds := this.UpdateVeiwFeeds()
	Adds = append(Adds, addfeeds...)
	Removes = append(Removes, delfeeds...)

	adds, dels := this.UpdateVeiwBallSkill()
	Adds = append(Adds, adds...)
	Removes = append(Removes, dels...)

	addfoods, delfoods := this.UpdateVeiwFoods()
	Adds = append(Adds, addfoods...)
	Removes = append(Removes, delfoods...)

	addplayers, delplayers := this.updateViewBallPlayer()
	AddPlayers = append(AddPlayers, addplayers...)
	RemovePlayers = append(RemovePlayers, delplayers...)

	Eats = append(Eats, this.ScenePlayerPool.MsgEats...)
	Hits = append(Hits, this.ScenePlayerPool.MsgHits...)

	ball := this.SelfAnimal
	if this.isMoved {
		ballmove := this.ScenePlayerPool.MsgBallMove
		ballmove.Id = ball.GetID()
		ballmove.X = int32(ball.Pos.X * bll.MsgPosScaleRate)
		ballmove.Y = int32(ball.Pos.Y * bll.MsgPosScaleRate)

		// angle && face
		if (this.SelfAnimal.HasForce() == false || this.Power == 0) && this.Face != 0 {
			ballmove.Face = uint32(this.Face)
			ballmove.Angle = 0
		} else {
			ballmove.Face = 0
			ballmove.Angle = int32(this.Angle)
		}

		ballmove.State = 0
		if this.isRunning {
			ballmove.State = 2
		}
		if skillid := this.Skill.GetCurSkillId(); skillid != 0 {
			ballmove.State = skillid
		}

		Moves = append(Moves, &ballmove)
	}

	//玩家广播
	for _, other := range this.Others {
		Eats = append(Eats, other.ScenePlayerPool.MsgEats...)
		Hits = append(Hits, other.ScenePlayerPool.MsgHits...)
		if other.isMoved {
			ball = other.SelfAnimal
			ballmove := other.ScenePlayerPool.MsgBallMove
			ballmove.Id = ball.GetID()
			ballmove.X = int32(ball.Pos.X * bll.MsgPosScaleRate)
			ballmove.Y = int32(ball.Pos.Y * bll.MsgPosScaleRate)

			// angle && face
			if (other.SelfAnimal.HasForce() == false || other.Power == 0) && other.Face != 0 {
				ballmove.Face = uint32(other.Face)
				ballmove.Angle = 0
			} else {
				ballmove.Face = 0
				ballmove.Angle = int32(other.Angle)
			}

			ballmove.State = 0
			if other.isRunning {
				ballmove.State = 2
			}
			if skillid := other.Skill.GetCurSkillId(); skillid != 0 {
				ballmove.State = skillid
			}

			if other != this {
				Moves = append(Moves, &ballmove)
			}
		}
	}

	// 玩家视野中的所有消息，发送给自己
	for _, cell := range this.LookCells {
		Moves = append(Moves, cell.MsgMoves...)
	}

	if len(Adds) != 0 || len(Removes) != 0 {
		//剔除自己
		if len(Adds) != 0 {
			for k, v := range Adds {
				if v.Id == this.SelfAnimal.GetID() {
					Adds = append(Adds[:k], Adds[k+1:]...)
					break
				}
			}
		}

		if len(Removes) != 0 {
			for k, v := range Removes {
				if v == this.SelfAnimal.GetID() {
					Removes = append(Removes[:k], Removes[k+1:]...)
					break
				}
			}
		}
	}

	if len(Eats) != 0 || len(Adds) != 0 || len(AddPlayers) != 0 || len(Hits) != 0 || len(Removes) != 0 || len(RemovePlayers) != 0 {
		msg := &this.msgPool.MsgSceneTCP
		msg.Eats = Eats
		msg.Adds = Adds
		msg.AddPlayers = AddPlayers
		msg.Hits = Hits
		msg.Removes = Removes
		msg.RemovePlayers = RemovePlayers
		this.SendCmd(usercmd.MsgTypeCmd_SceneTCP, msg)
	}
	//	glog.Info("fuck u!")
	//glog.Info(len(this.GetScene().GetMovingCubes()))

	if len(Moves) != 0 || len(this.GetScene().GetMovingCubes()) != 0 {
		msg := &this.msgPool.MsgSceneUDP
		msg.ChangingInf = []*usercmd.CubeReDst{}
		for i, v := range this.GetScene().GetMovingCubes() {
			if (this.GetScene().GetCubeState(i) == 1 || this.GetScene().GetCubeState(i) == 2) && this.GetScene().GetCubeMoveDrct(i) >= 0 {
				msg.ChangingInf = append(msg.ChangingInf, &usercmd.CubeReDst{
					CubeIndex:      i,
					RemainDistance: 1000 - v,
				})
			} else if (this.GetScene().GetCubeState(i) == 1 || this.GetScene().GetCubeState(i) == 0) && this.GetScene().GetCubeMoveDrct(i) <= 0 {
				msg.ChangingInf = append(msg.ChangingInf, &usercmd.CubeReDst{
					CubeIndex:      i,
					RemainDistance: -v,
				})
			} else if (this.GetScene().GetCubeState(i) == -1 || this.GetScene().GetCubeState(i) == 0) && this.GetScene().GetCubeMoveDrct(i) >= 0 {
				msg.ChangingInf = append(msg.ChangingInf, &usercmd.CubeReDst{
					CubeIndex:      i,
					RemainDistance: -v,
				})
			} else if (this.GetScene().GetCubeState(i) == -2 || this.GetScene().GetCubeState(i) == -1) && this.GetScene().GetCubeMoveDrct(i) <= 0 {
				msg.ChangingInf = append(msg.ChangingInf, &usercmd.CubeReDst{
					CubeIndex:      i,
					RemainDistance: -1000 - v,
				})
			}
			if v == 0 {
				//删除这个
				this.GetScene().RemoveMovingCube(i)
			}

		}
		for _, v := range msg.ChangingInf {
			glog.Info(v)
		}

		msg.Moves = Moves
		msg.Frame = this.GetScene().Frame()

		if this.Sess != nil {
			// 优先采用可丢包的原生UDP信道发送
			this.Sess.SendUDPCmd(usercmd.MsgTypeCmd_SceneUDP, msg)
		}
	}
}

// 检查能否被吃
func (this *ScenePlayer) CanBeEat() bool {
	if this.IsLive {
		return true
	}
	return false
}

// 处理房间结束
func (this *ScenePlayer) DoEndRoom(room IRoom) bool {
	if room == nil || true == this.IsRobot {
		return false
	}
	if this.udata.Level == 0 {
		this.udata.Level = 1
	}
	if this.udata.Scores == 0 {
		this.udata.Scores = 1
	}
	glog.Info("[玩家] 房间结算 [", room.RoomType(), ",", room.ID(), "],[", this.ID, ",", this.udata.Account, "] rank= ", this.Rank)
	return true
}

// 发送普通消息
func (this *ScenePlayer) SendCmd(cmd usercmd.MsgTypeCmd, msg common.Message) bool {
	if this.Sess == nil {
		return false
	}
	data, flag, err := common.EncodeGoCmd(uint16(cmd), msg)
	if err != nil {
		return false
	}

	this.AsyncSend(data, flag)
	return true
}

// 广播消息
func (this *ScenePlayer) BroadCastMsg(cmd usercmd.MsgTypeCmd, msg common.Message) bool {
	this.room.BroadcastMsg(cmd, msg)
	return true
}

// 给周围发送消息
func (this *ScenePlayer) SendRoundMsg(cmd usercmd.MsgTypeCmd, msg common.Message) bool {
	data, flag, err := common.EncodeGoCmd(uint16(cmd), msg)
	if err != nil {
		return false
	}

	return this.AsyncRoundMsg(data, flag)
}

func (this *ScenePlayer) AsyncRoundMsg(data []byte, flag byte) bool {
	this.AsyncSend(data, flag)
	for _, player := range this.RoundPlayers {
		player.AsyncSend(data, flag)
	}
	return true
}

// 发送二进制消息
func (this *ScenePlayer) AsyncSend(buffer []byte, flag byte) {
	if this.Sess != nil {
		this.Sess.AsyncSend(buffer, flag)
	}
}
func (this *ScenePlayer) UpdateMove(perTime float64, frameRate float64) {
	if !this.IsLive {
		return
	}

	// 玩家球移动
	ball := this.SelfAnimal
	ball.UpdateForce(perTime)
	if ball.Move(perTime, frameRate) {
		ball.FixMapEdge() //修边
		this.isMoved = true
		ball.ResetRect()

		if this.isRunning {
			cost := frameRate * float64(consts.FrameTimeMS) * consts.DefaultRunCostMP
			cost = 0 //mata:infinite running
			diff := ball.GetMP() - cost
			if diff <= 0 {
				this.isRunning = false
			} else {
				ball.SetMP(diff)
			}
		}
	}
}

// 场景帧驱动
func (this *ScenePlayer) Update(perTime float64, now int64, scene IScene) {
	if this.IsRobot == false && this.Sess == nil {
		return
	}

	curmp := this.SelfAnimal.GetMP()
	curexp := this.GetExp()
	curanimal := this.SelfAnimal.GetAnimalId()

	// 有在释放技能，恢复转向
	this.Skill.TryTurn(&this.Angle, &this.Face)

	// 角色朝向，每帧只算一次。避免多次计算，因此代码挪至开头
	this.SelfAnimal.SetAngleVelAndNormalize(
		math.Cos(math.Pi*this.Angle/180),
		-math.Sin(math.Pi*this.Angle/180))

	this.Skill.Update()

	if this.AI.IsOK() {
		this.AI.UpdateAI(perTime)
	}

	var frameRate float64 = 2

	// 更新球
	ball := this.SelfAnimal

	// 玩家球移动
	this.UpdateMove(perTime, frameRate)

	this.UpdateView(scene)

	if this.IsLive {
		var rect util.Square
		rect.CopyFrom(this.GetViewRect())
		rect.SetRadius(this.SelfAnimal.GetEatRange())
		cells := this.GetScene().GetAreaCells(&rect)
		for _, newcell := range cells {
			newcell.EatByPlayer(ball, this)
		}
	}

	// 更新视野中的玩家
	this.UpdateViewPlayers(scene)

	if curanimal != this.SelfAnimal.GetAnimalId() || curexp != this.GetExp() || curmp != this.SelfAnimal.GetMP() {
		this.RefreshPlayer()
	}
	if curanimal != this.SelfAnimal.GetAnimalId() {
		this.AsyncAnimal()
	}
}

//增加经验，不会降级
func (this *ScenePlayer) UpdateExp(addexp int32) bool {
	isChangeLevel := false
	if 0 == addexp {
		return false
	}
	exp := this.GetExp()
	if addexp > 0 {
		exp += uint32(addexp)
	} else {
		if exp > uint32(util.AbsInt(int(addexp))) {
			exp -= uint32(util.AbsInt(int(addexp)))
		} else {
			exp = 0
		}
	}
	this.SetExp(exp)
	OldAnimalID := int32(this.SelfAnimal.GetAnimalId())
	if addexp > 0 {
		newanimalid := int32(conf.ConfigMgr_GetMe().GetAnimalIdByExp(this.SceneID(), exp))
		if newanimalid > OldAnimalID {
			isChangeLevel = true
			this.SelfAnimal.SetAnimalId(newanimalid)
			this.SelfAnimal.UpdateSize()
			this.SelfAnimal.SetMP(consts.DefaultMaxMP)
			this.SelfAnimal.SetHP(consts.DefaultMaxHP)
			this.SelfAnimal.SetHpMax(consts.DefaultMaxHP)
		}
	}
	this.SelfAnimal.UpdateSize()
	return isChangeLevel
}

func (this *ScenePlayer) GetSnapInfo() *usercmd.MsgPlayerSnap {
	msg := &this.msgPool.MsgPlayerSnap
	msg.Snapx = float32(this.SelfAnimal.Pos.X)
	msg.Snapy = float32(this.SelfAnimal.Pos.Y)
	msg.Angle = float32(this.Angle)
	msg.Id = this.ID
	return msg
}

// 玩家定时器 (一秒一次)
func (this *ScenePlayer) TimeAction(room IRoom, timenow time.Time) bool {
	if false == this.IsLive {
		return true
	}

	nowsec := timenow.Unix()
	// 定时器
	this.SelfAnimal.TimeAction(nowsec)
	mp := int32(this.SelfAnimal.GetMP())
	maxmp := consts.DefaultMaxMP
	curhp := this.SelfAnimal.GetHP()
	animalid := this.SelfAnimal.GetAnimalId()
	maxhp := consts.DefaultMaxHP
	addmp := consts.DefaultMpRecover
	addhp := consts.DefaultHpRecover
	if addhp <= 0 {
		addhp = 1
	}
	if addmp <= 0 {
		addmp = 1
	}

	if 0 != addmp {
		if mp+int32(addmp) > int32(maxmp) {
			this.SelfAnimal.SetMP(float64(maxmp))
		} else {
			this.SelfAnimal.SetMP(float64(mp + int32(addmp)))
		}
	}
	if uint32(curhp) < uint32(maxhp) {
		if uint32(uint32(curhp)+uint32(addhp)) > uint32(maxhp) {
			this.SelfAnimal.SetHP(int32(maxhp))
		} else {
			this.SelfAnimal.SetHP(int32(curhp) + int32(addhp))
		}
	}

	this.RefreshPlayer()
	if animalid != this.SelfAnimal.GetAnimalId() {
		this.AsyncAnimal()
	}

	return true
}

func (this *ScenePlayer) RefreshPlayer() {
	if this.Sess == nil {
		return
	}
	msg := &this.msgPool.MsgRefreshPlayer
	msg.Player.Id = this.ID
	msg.Player.Name = this.Name
	msg.Player.Local = this.GetLocation()
	msg.Player.IsLive = this.IsLive
	msg.Player.SnapInfo = this.GetSnapInfo()
	msg.Player.Curexp = this.GetExp()
	msg.Player.BallId = this.SelfAnimal.GetID()
	msg.Player.Curmp = uint32(this.SelfAnimal.GetMP())
	msg.Player.Curhp = uint32(this.SelfAnimal.GetHP())
	msg.Player.Animalid = uint32(this.SelfAnimal.GetAnimalId())
	msg.Player.TeamName = this.udata.TeamName

	msg.Player.BombNum = int32(this.SelfAnimal.GetAttr(bll.AttrBombNum))
	msg.Player.HammerNum = int32(this.SelfAnimal.GetAttr(bll.AttrHammerNum))

	this.Sess.SendCmd(usercmd.MsgTypeCmd_RefreshPlayer, msg)
}

func (this *ScenePlayer) AsyncAnimal() {
	msg := &this.msgPool.MsgAsyncPlayerAnimal
	msg.Id = this.ID
	msg.Animalid = uint32(this.SelfAnimal.GetAnimalId())
	this.SendRoundMsg(usercmd.MsgTypeCmd_AsyncPlayerAnimal, msg)
}

//重置摇杆力
func (this *ScenePlayer) CleanPower() {
	this.Power = 0
	this.Angle = 0
}

func (this *ScenePlayer) SetIsRunning(v bool) {
	this.isRunning = v
}

func (this *ScenePlayer) GetId() uint64 {
	return this.ID
}

func (this *ScenePlayer) GetFrame() uint32 {
	return this.room.Frame()
}

func (this *ScenePlayer) GetExp() uint32 {
	return uint32(this.SelfAnimal.GetAttr(bll.AttrExp))
}

func (this *ScenePlayer) SetExp(exp uint32) {
	this.SelfAnimal.SetAttr(bll.AttrExp, float64(exp))
}

func (this *ScenePlayer) GetScene() IScene {
	return this.room.GetPlayerIScene()
}

func (this *ScenePlayer) GetBallScene() bll.IScene {
	return this.room.GetBallIScene()
}

func (this *ScenePlayer) SceneID() uint32 {
	return this.room.SceneID()
}

func (this *ScenePlayer) GetUserData() *common.UserData {
	return this.udata
}

func (this *ScenePlayer) FindNearBallByKind(kind interfaces.BallKind, dir *util.Vector2, cells []*cll.Cell, ballType uint32) (interfaces.IBall, float64) {
	return this.ScenePlayerViewHelper.FindNearBallByKind(this.SelfAnimal, kind, dir, cells, ballType)
}

func (this *ScenePlayer) UpdateView(scene IScene) {
	if !this.IsLive {
		return
	}
	this.ScenePlayerViewHelper.UpdateView(scene, this.SelfAnimal, this.room.RoomSize(), this.room.CellNumX(), this.room.CellNumY())
}

func (this *ScenePlayer) UpdateViewPlayers(scene IScene) {
	this.ScenePlayerViewHelper.UpdateViewPlayers(scene, this.SelfAnimal)
}

func (this *ScenePlayer) GetID() uint64 {
	return this.ID
}

// 当前摇杆力度（目前恒为0或者1，来简化同步计算）
func (this *ScenePlayer) GetPower() float64 {
	return this.Power
}

func (this *ScenePlayer) IsRunning() bool {
	return this.isRunning
}

func (this *ScenePlayer) GetIsLive() bool {
	return this.IsLive
}

func (this *ScenePlayer) RoomType() uint32 {
	return this.room.RoomType()
}

// 玩家主数据
func (this *ScenePlayer) UData() *common.UserData {
	return this.udata
}

func (this *ScenePlayer) SetUData(udata *common.UserData) {
	this.udata = udata
}

func (this *ScenePlayer) KilledByPlayer(killer bll.IScenePlayer) {
	this.Dead(killer.(*ScenePlayer))
}

func (this *ScenePlayer) NewSkillBall(sb *bll.BallSkill) interfaces.ISkillBall {
	return NewISkillBall(this, sb) // skill.NewISkillBall
}

func (this *ScenePlayer) GetAngle() float64 {
	return this.Angle
}

func (this *ScenePlayer) GetFace() uint32 {
	return this.Face
}

// XXX Rename this.IsRobot to this.isRobo
func (this *ScenePlayer) GetIsRobot() bool {
	return this.IsRobot
}

func (this *ScenePlayer) DeadTime() int64 {
	return this.deadTime
}

func (this *ScenePlayer) SetDeadTime(deadTime int64) {
	this.deadTime = deadTime
}
