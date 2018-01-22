package ai

// 玩家AI

import (
	b3core "base/behavior3go/core"
	"base/glog"
	"common"
	"math/rand"
	"roomserver/conf"
	"roomserver/game/interfaces"
	"roomserver/game/scn"
	"roomserver/game/scn/plr"
	"roomserver/util"
)

type AIData struct {
	blackboard  *b3core.Blackboard
	actionTimer float64
	enemyId     uint32
	enemyTimer  float64

	lastPos       util.Vector2
	checkStopTime float64
	expireTime    int64 // 机器回收时间ms
	speedRate     float64

	bevData *conf.XmlAIBehaveData //ai配置
	bevTree *b3core.BehaviorTree
}

type ScenePlayerAI struct {
	AiData     *AIData
	SelfPlayer *plr.ScenePlayer
}

func NewIScenePlayerAI(player *plr.ScenePlayer) interfaces.IScenePlayerAI {
	return &ScenePlayerAI{AiData: &AIData{}, SelfPlayer: player}
}

func NewSceneAIPlayer(room scn.IRoom) *plr.ScenePlayer {
	var udata common.UserData
	udata.Id = room.NewRobotUID()
	//name := conf.ConfigMgr_GetMe().RobotNams.Names[rand.Intn(len(conf.ConfigMgr_GetMe().RobotNams.Names))]
	name := conf.ConfigMgr_GetMe().GetRobotRandName(156)
	data := conf.ConfigMgr_GetMe().RobotDatas.Players[rand.Intn(len(conf.ConfigMgr_GetMe().RobotDatas.Players))]
	udata.Icon = data.Icon
	udata.Account = data.Acc
	udata.PassIcon = data.PassIcon
	udata.Location = data.Location
	udata.Sex = data.Sex
	udata.ShowPos = 1

	p := room.NewScenePlayer(&udata, name, true)
	ai := p.AI.(*ScenePlayerAI)

	cfg := conf.ConfigMgr_GetMe().GetAISrcDataByScore(0)
	rand := rand.Int() % cfg.PowerSum
	for i := 0; i < len(cfg.Powers); i++ {
		if rand < cfg.Powers[i] {
			ai.AiData.bevData = conf.ConfigMgr_GetMe().GetAiBevDataByRoomType(cfg.Behaves[i])
			break
		}
	}
	if ai.AiData.bevData == nil {
		glog.Error("bevData nil:", cfg.Behaves, "  room:", room.RoomType())
	}

	if ai.AiData.bevTree != nil {
		ai.AiData.bevTree.Close(nil, ai.AiData.blackboard)
	}
	ai.AiData.bevTree = GetBevTree(ai.AiData.bevData.Aifile)
	ai.AiData.speedRate = cfg.SpeedRate
	return p
}

func NewCopyPlayerAI(room scn.IRoom, copy_player *plr.ScenePlayer) *plr.ScenePlayer {
	var udata common.UserData
	udata.Id = room.NewRobotUID()
	name := copy_player.Name
	udata.Icon = copy_player.GetUserData().Icon
	udata.Account = copy_player.GetUserData().Account
	udata.PassIcon = copy_player.GetUserData().PassIcon
	udata.Location = copy_player.GetUserData().Location
	udata.Sex = copy_player.GetUserData().Sex
	udata.ShowPos = copy_player.GetUserData().ShowPos

	p := room.NewScenePlayer(&udata, name, true)
	if p != nil {
		// 设置复制体的回收时间
		ai := p.AI.(*ScenePlayerAI)
		ai.AiData.expireTime = copy_player.AI.GetExpireTime()
	}
	return p
}

func (this *ScenePlayerAI) AddAICtrl() {
	glog.Info("AddAICtrl ", this.SelfPlayer.GetId())
	cfg := conf.ConfigMgr_GetMe().GetAISrcDataByScore(0)
	rand := rand.Int() % cfg.PowerSum
	for i := 0; i < len(cfg.Powers); i++ {
		if rand < cfg.Powers[i] {
			this.AiData.bevData = conf.ConfigMgr_GetMe().GetAiBevDataByRoomType(cfg.Behaves[i])
			break
		}
	}
	if this.AiData.bevData == nil {
		glog.Error("bevData nil:", cfg.Behaves, "  room:", this.SelfPlayer.GetScene().RoomType())
	}

	if this.AiData.bevTree != nil {
		this.AiData.bevTree.Close(nil, this.AiData.blackboard)
	}
	this.AiData.bevTree = GetBevTree(this.AiData.bevData.Aifile)
	this.AiData.speedRate = cfg.SpeedRate
}

func (this *ScenePlayerAI) RemoveAICtrl() {
	glog.Info("RemoveAICtrl ", this.SelfPlayer.GetId())
	if this.AiData.bevTree != nil {
		this.AiData.bevTree.Close(nil, this.AiData.blackboard)
		this.AiData.bevTree = nil
	}
	this.AiData.bevData = nil

	// 确保不在移动中
	this.SelfPlayer.Power = 0
}

//InitAI 初始化AI
func (this *ScenePlayerAI) InitAI() {
	this.AiData.actionTimer = 0
	this.AiData.blackboard = CreateBlackBoard()

	this.AiData.enemyId = 0
	this.AiData.enemyTimer = 0
}

//被攻击
func (this *ScenePlayerAI) OnBeHit(enemyBallId uint32) {
	if this.AiData != nil {
		var evoTime float64 = 3
		cfg := this.AiData.bevData
		if cfg != nil && cfg.EvoTime > 0 {
			evoTime = cfg.EvoTime
		}
		this.AiData.enemyId = enemyBallId
		this.AiData.enemyTimer = evoTime
	}

}

//死亡事件
func (this *ScenePlayerAI) OnDeadEvent(enemy interface{}) {
	if this.AiData != nil {
		this.AiData.actionTimer = 3 //3秒后复活
	}
}
func (this *ScenePlayerAI) FindNearAttackAnimal(disNeed float64) uint32 {
	var minAni uint32
	var minDis float64

	sqrDisNeed := disNeed * disNeed * 2

	for _, viewPlayer := range this.SelfPlayer.Others {
		if viewPlayer == this.SelfPlayer {
			continue
		}
		if this.SelfPlayer.SelfAnimal.PreTryHit(viewPlayer.SelfAnimal) {
			dis := viewPlayer.SelfAnimal.GetPosV().SqrMagnitudeTo(this.SelfPlayer.SelfAnimal.GetPosV())
			if disNeed > 0 && dis > sqrDisNeed {
				continue //0不确认距离
			}
			if minAni == 0 || dis < minDis {
				minAni = viewPlayer.SelfAnimal.GetID()
				minDis = dis
			}
		}
	}
	return minAni
}

//UpdateAI 更新AI
func (this *ScenePlayerAI) UpdateAI(perTime float64) {
	//if !this.IsLive {
	//	this.Relife()
	//}

	var robotFrameTime float64 = 0.2
	this.AiData.actionTimer -= perTime
	if this.AiData.actionTimer < 0 {
		if this.SelfPlayer.IsLive {
			//glog.Info("name:", this.name, " level:", this.SelfAnimal.GetAnimalId(), " mp:", this.SelfAnimal.GetMP())
			this.AiData.actionTimer = robotFrameTime
			//scene.room.Move(&MoveOp{this.id, 0.5, 360 * rand.Float64()})
			this.PreUpdate(robotFrameTime)
			this.UpdateBevTree()
			this.UpdateRequst()

			ai := conf.ConfigMgr_GetMe().GetRoomAI()
			//glog.Info("ai exp:", this.room.RoomType(), "  exp:", ai.MaxExp, "  ai:", ai)
			if ai.MaxExp > 0 && this.SelfPlayer.GetExp() >= ai.MaxExp {
				this.SelfPlayer.RealDead(this.SelfPlayer)
			}
		} else {
			this.SelfPlayer.Relife()
		}

	}
}

func (this *ScenePlayerAI) PreUpdate(perTime float64) {
	if this.AiData.enemyId > 0 {
		this.AiData.enemyTimer -= perTime
		//glog.Info("AI.sub enemy:", this.id, " time:", this.AiData.enemyTimer)

		if this.AiData.enemyTimer <= 0 {
			this.AiData.enemyId = 0
			this.AiData.enemyTimer = 0
			//glog.Info("AIHIT.no enemy:", this.id)
		}
	}
	this.AiData.blackboard.Set("enemy", this.AiData.enemyId, "", "")
	//可能不动了
	this.AiData.checkStopTime -= perTime
	if this.AiData.checkStopTime <= 0 {
		this.AiData.checkStopTime = 3

		dis := this.AiData.lastPos.SqrMagnitudeTo(this.SelfPlayer.SelfAnimal.GetPosV())
		this.AiData.lastPos = *this.SelfPlayer.SelfAnimal.GetPosV()
		//glog.Info("AI.walk:", this.name, " dis:", this.SelfAnimal.timeWalkDistance, " d:", dis)
		isStop := false
		if dis < 0.2*0.2 {
			//glog.Info("!!stop:", this.name, " d:", dis, " at:", this.SelfAnimal.pos)
			isStop = true
		}
		this.AiData.blackboard.Set("isStop", isStop, "", "")

	}
}

func (this *ScenePlayerAI) UpdateBevTree() {
	if this.SelfPlayer != nil {
		this.AiData.bevTree.Tick(this.SelfPlayer, this.AiData.blackboard)
	} else {
		panic("unknow error!")
	}
}

func (this *ScenePlayerAI) UpdateRequst() {

}

func (this *ScenePlayerAI) GetBevData() *conf.XmlAIBehaveData {
	return this.AiData.bevData
}

func (this *ScenePlayerAI) GetExpireTime() int64 {
	return this.AiData.expireTime
}

func (this *ScenePlayerAI) IsOK() bool {
	return this.AiData.bevTree != nil
}
