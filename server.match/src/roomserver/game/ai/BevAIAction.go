package ai

import (
	b3 "base/behavior3go"
	b3config "base/behavior3go/config"
	b3core "base/behavior3go/core"
	"base/glog"
	bmath "base/math"
	"math/rand"
	"roomserver/game/bll"
	"roomserver/game/interfaces"
	"roomserver/game/scn/plr"
	"roomserver/util"
	"time"
)

//---------------------------------------condition------------------------------------------------
//CheckBall
type CheckBall struct {
	b3core.Condition
	index string
}

func (this *CheckBall) Initialize(setting *b3config.BTNodeCfg) {
	this.Condition.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
}

func (this *CheckBall) OnTick(tick *b3core.Tick) b3.Status {
	id := tick.Blackboard.GetUInt32(this.index, "", "")
	if id < 1 {
		return b3.FAILURE
	}
	player := tick.GetTarget().(*plr.ScenePlayer)
	tball := player.FindNearBall(uint32(id)) //GetScene().GetBall(uint32(id))
	if tball == nil {
		//glog.Info("CheckBall:miss ball", id)
		tick.Blackboard.Set(this.index, uint32(0), "", "")
		return b3.FAILURE
	}
	return b3.SUCCESS
}

//CheckNearPlayer
type CheckNearPlayer struct {
	b3core.Condition
	index string
}

func (this *CheckNearPlayer) Initialize(setting *b3config.BTNodeCfg) {
	this.Condition.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
}

func (this *CheckNearPlayer) OnTick(tick *b3core.Tick) b3.Status {
	id := tick.Blackboard.GetUInt32(this.index, "", "")
	if id < 1 {
		return b3.FAILURE
	}
	player := tick.GetTarget().(*plr.ScenePlayer)
	tball := player.FindVeiwAnimal(uint32(id))
	if tball == nil {
		//glog.Info("CheckNearPlayer:miss ball", id)
		return b3.FAILURE
	}

	if !tball.GetPlayer().GetIsLive() {
		return b3.FAILURE
	}
	return b3.SUCCESS
}

//CheckNearAttackPlayer
type CheckNearAttackPlayer struct {
	b3core.Condition
	index string
}

func (this *CheckNearAttackPlayer) Initialize(setting *b3config.BTNodeCfg) {
	this.Condition.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
}

func (this *CheckNearAttackPlayer) OnTick(tick *b3core.Tick) b3.Status {
	id := tick.Blackboard.GetUInt32(this.index, "", "")
	if id < 1 {
		return b3.FAILURE
	}
	player := tick.GetTarget().(*plr.ScenePlayer)
	tball := player.FindVeiwAnimal(uint32(id))
	if tball == nil {
		// glog.Info("CheckNearAttackPlayer:miss ball", id)
		tick.Blackboard.SetTree(this.index, uint32(0), "")
		player.Face = 0
		return b3.FAILURE
	}

	if /*player.SelfAnimal.GetAnimalId() <= tball.GetAnimalId() tmp1 <= tmp2 || */ !player.SelfAnimal.PreTryHit(tball) {
		// glog.Info("CheckNearAttackPlayer:PreTryHit fail", id)
		tick.Blackboard.SetTree(this.index, uint32(0), "")
		player.Face = 0
		return b3.FAILURE
	}

	//cfg := conf.ConfigMgr_GetMe().GetAISrcDataByScore(player.room.HScore, player.room.roomType)
	cfg := player.AI.GetBevData()
	var followTime int64 = 8
	var followRange float64 = 6
	if cfg != nil && cfg.FollowTime > 0 && cfg.FollowRange > 0 {
		followTime = cfg.FollowTime
		followRange = cfg.FollowRange
	}
	//超时保护
	var startTime = tick.Blackboard.GetInt64("attackTime", "", "")
	var currTime int64 = time.Now().UnixNano() / 1000000
	if currTime-startTime > followTime && player.SelfAnimal.GetPosV().SqrMagnitudeTo(tball.GetPosV()) > followRange*followRange {
		tick.Blackboard.Set("attackTime", int64(0), "", "")
		tick.Blackboard.SetTree(this.index, uint32(0), "")
		// glog.Info("attachTime failed")
		player.Face = 0
		return b3.FAILURE
	}

	return b3.SUCCESS
}

//AttrLimit
type AttrLimit struct {
	b3core.Condition
	limit int32
	attr  bll.AttrType
}

func (this *AttrLimit) Initialize(setting *b3config.BTNodeCfg) {
	this.Condition.Initialize(setting)
	this.limit = int32(setting.GetPropertyAsInt("limit"))
	this.attr = bll.AttrType(setting.GetPropertyAsInt("attr"))
}

func (this *AttrLimit) OnTick(tick *b3core.Tick) b3.Status {
	player := tick.GetTarget().(*plr.ScenePlayer)
	if int(player.SelfAnimal.GetAttr(this.attr)) < int(this.limit) {
		return b3.SUCCESS
	}
	return b3.FAILURE
}

//TargetAttrLess
type TargetAttrLess struct {
	b3core.Condition
	index string
	attr  bll.AttrType
}

func (this *TargetAttrLess) Initialize(setting *b3config.BTNodeCfg) {
	this.Condition.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
	this.attr = bll.AttrType(setting.GetPropertyAsInt("attr"))
}

func (this *TargetAttrLess) OnTick(tick *b3core.Tick) b3.Status {
	player := tick.GetTarget().(*plr.ScenePlayer)

	id := tick.Blackboard.GetUInt32(this.index, "", "")
	if id < 1 {
		return b3.FAILURE
	}
	tball := player.FindVeiwAnimal(uint32(id))
	if tball == nil {
		return b3.FAILURE
	}
	if tball.GetAttr(this.attr) < player.SelfAnimal.GetAttr(this.attr) {
		glog.Info("TargetAttrLess ok:", tball.GetAttr(this.attr), player.SelfAnimal.GetAttr(this.attr))
		return b3.SUCCESS
	}
	if !tball.GetPlayer().GetIsRobot() {
		glog.Info("TargetAttrLess fail:", tball.GetAttr(this.attr), player.SelfAnimal.GetAttr(this.attr))
	}

	return b3.FAILURE
}

//CheckBool
type CheckBool struct {
	b3core.Condition
	keyname string
}

func (this *CheckBool) Initialize(setting *b3config.BTNodeCfg) {
	this.Condition.Initialize(setting)
	this.keyname = setting.GetPropertyAsString("keyname")
}

func (this *CheckBool) OnTick(tick *b3core.Tick) b3.Status {
	var b = tick.Blackboard.GetBool(this.keyname, "", "")
	if b {
		//glog.Info("CheckBool ok:", this.keyname)
		return b3.SUCCESS
	}
	return b3.FAILURE
}

//---------------------------------------actions------------------------------------------------
//RandAction
type RandAction struct {
	b3core.Action
	index string
	min   float64
	max   float64
}

func (this *RandAction) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
	this.min = setting.GetProperty("min")
	this.max = setting.GetProperty("max")
}

func (this *RandAction) OnTick(tick *b3core.Tick) b3.Status {
	val := this.min + rand.Float64()*(this.max-this.min)
	tick.Blackboard.Set(this.index, val, "", "")
	return b3.SUCCESS
}

//TurnTarget
type TurnTarget struct {
	b3core.Action
	index string
}

func (this *TurnTarget) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
}

func (this *TurnTarget) OnTick(tick *b3core.Tick) b3.Status {
	id := tick.Blackboard.GetUInt32(this.index, "", "")
	if id < 1 {
		return b3.FAILURE
	}
	player := tick.GetTarget().(*plr.ScenePlayer)
	tball := player.FindNearBall(uint32(id))
	if tball == nil {
		//glog.Info("TurnTarget:miss ball", id, " cell:", len(player.lookCells))
		tick.Blackboard.Set(this.index, uint32(0), "", "")
		return b3.FAILURE
	}
	var currTime int64 = time.Now().UnixNano() / 1000000
	starttime := tick.Blackboard.GetInt64("targetTime", "", "")
	if starttime > 0 && currTime-starttime > 6000 {
		tick.Blackboard.Set(this.index, uint32(0), "", "")
		//glog.Info("TurnTarget:overtime", id, "  time:", currTime-starttime, " start:", starttime)
		return b3.FAILURE
	}
	bx, by := tball.GetPos()
	bv := util.Vector2{bx, by}
	v := bv.SubMethod(player.SelfAnimal.GetPosV())
	angle := 360 - bmath.Vector2DToAngle(&bmath.Vector2{float32(v.X), float32(v.Y)})

	//	player.GetScene().room.Move(&MoveOp{player.id, 0, float64(angle)})
	player.Angle = float64(angle)
	//glog.Info("TurnTarget:", tball.GetID(), " x=", bx, " y=", by, " angle=", angle, " dis=", bv.DistanceTo(&player.SelfAnimal.pos))
	return b3.SUCCESS
}

//TurnTargetPlayer
type TurnTargetPlayer struct {
	b3core.Action
	index string
}

func (this *TurnTargetPlayer) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
}

func (this *TurnTargetPlayer) OnTick(tick *b3core.Tick) b3.Status {
	id := tick.Blackboard.GetUInt32(this.index, "", "")
	if id < 1 {
		return b3.FAILURE
	}
	player := tick.GetTarget().(*plr.ScenePlayer)
	tball := player.FindVeiwAnimal(uint32(id))
	if tball == nil {
		//glog.Info("TurnTargetPlayer:miss ball", id, " cell:", len(player.lookCells))
		tick.Blackboard.Set(this.index, uint32(0), "", "")
		return b3.FAILURE
	}

	bx, by := tball.GetPos()
	bv := util.Vector2{bx, by}
	v := bv.SubMethod(player.SelfAnimal.GetPosV())
	angle := 360 - bmath.Vector2DToAngle(&bmath.Vector2{float32(v.X), float32(v.Y)})

	//	player.GetScene().room.Move(&MoveOp{player.id, 0, float64(angle)})
	player.Angle = float64(angle)
	//glog.Info("TurnTargetPlayer:", tball.GetID(), " x=", bx, " y=", by, " angle=", angle)
	return b3.SUCCESS
}

//TurnAwayTargetPlayer
type TurnAwayTargetPlayer struct {
	b3core.Action
	index string
}

func (this *TurnAwayTargetPlayer) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
}

func (this *TurnAwayTargetPlayer) OnTick(tick *b3core.Tick) b3.Status {
	id := tick.Blackboard.GetUInt32(this.index, "", "")
	if id < 1 {
		return b3.FAILURE
	}
	player := tick.GetTarget().(*plr.ScenePlayer)
	tball := player.FindVeiwAnimal(uint32(id))
	if tball == nil {
		//glog.Info("AIHIT.TurnAwayTarget:miss ball", id, " cell:", len(player.lookCells))
		tick.Blackboard.Set(this.index, uint32(0), "", "")
		return b3.FAILURE
	}

	bx, by := tball.GetPos()
	bv := util.Vector2{bx, by}
	v := bv.SubMethod(player.SelfAnimal.GetPosV())
	v = v.MultiMethod(-1) //取反

	//修复方向
	x, y := player.SelfAnimal.GetPos()
	var boardSize float64 = 3
	if x < boardSize {
		v = &util.Vector2{1, 0}
	} else if x > (player.GetScene().RoomSize() - boardSize) {
		v = &util.Vector2{-1, 0}
	} else if y > (player.GetScene().RoomSize() - boardSize) {
		v = &util.Vector2{0, -1}
	} else if y < boardSize {
		v = &util.Vector2{0, 1}
	}

	//if bmath.AbsF64(v.x) > bmath.AbsF64(v.y) {
	//	v.y = 0
	//} else {
	//	v.x = 0
	//}

	angle := 360 - bmath.Vector2DToAngle(&bmath.Vector2{float32(v.X), float32(v.Y)})
	//	player.GetScene().room.Move(&MoveOp{player.id, 0, float64(angle)})
	player.Angle = float64(angle)
	//glog.Info("TurnAwayTarget:", tball.GetID(), " x=", bx, " y=", by, " angle=", angle)

	//glog.Info("AIHIT.TurnAwayTarget:", tball.GetID(), " x=", int(bx), " y=", int(by), " angle=", int(angle), " ev:", player.aiData.evoAngle)
	return b3.SUCCESS
}

//TurnIndex
type TurnIndex struct {
	b3core.Action
	index string
}

func (this *TurnIndex) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
}

func (this *TurnIndex) OnTick(tick *b3core.Tick) b3.Status {
	angle := tick.Blackboard.GetFloat64(this.index, "", "")
	player := tick.GetTarget().(*plr.ScenePlayer)
	//	player.GetScene().room.Move(&MoveOp{player.id, 0, float64(angle)})
	player.Angle = float64(angle)
	return b3.SUCCESS
}

//FindNearUnit
type FindNearUnit struct {
	b3core.Action
	index    string
	unitKind uint32
}

func (this *FindNearUnit) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
	this.unitKind = uint32(setting.GetPropertyAsInt("unitKind"))
}

func (this *FindNearUnit) OnTick(tick *b3core.Tick) b3.Status {
	player := tick.GetTarget().(*plr.ScenePlayer)
	ball, _ := player.FindNearBallByKind(interfaces.BallKind_None, nil, nil, this.unitKind)
	var id uint32
	if ball != nil {
		id = ball.GetID()
		tick.Blackboard.Set(this.index, id, "", "")

		//bx, by := ball.GetPos()
		//bv := util.Vector2{bx, by}
		//glog.Info("FindNearUnit:", this.unitKind, " id:", id, " near=", player.GetNearFoodCount(), " type=", ball.GetType(), " dis=", bv.DistanceTo(&player.SelfAnimal.pos))
		//		glog.Info("FindNearUnit ", ball.(*BallFood).ballType)
		var currTime int64 = time.Now().UnixNano() / 1000000
		tick.Blackboard.Set("targetTime", currTime, "", "")
		return b3.SUCCESS
	}
	tick.Blackboard.Set(this.index, id, "", "")

	//glog.Info("nofindNearUnit:", this.unitKind, " id:", id)
	return b3.FAILURE

}

//FindNearUnit2
type FindNearUnit2 struct {
	b3core.Action
	index    string
	unitKind uint32
	dis      float64
}

func (this *FindNearUnit2) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
	this.unitKind = uint32(setting.GetPropertyAsInt("unitKind"))
	this.dis = setting.GetProperty("dis")
}

func (this *FindNearUnit2) OnTick(tick *b3core.Tick) b3.Status {
	player := tick.GetTarget().(*plr.ScenePlayer)
	tick.Blackboard.Set(this.index, uint32(0), "", "")

	var rect util.Square
	rect.CopyFrom(player.GetViewRect())
	rect.SetRadius(this.dis)
	cells := player.GetScene().GetAreaCells(&rect)

	ball, _ := player.FindNearBallByKind(interfaces.BallKind_None, nil, cells, this.unitKind)
	if nil == ball {
		return b3.FAILURE
	}

	fx, fy := ball.GetPos()
	pos := util.Vector2{fx, fy}
	dis := pos.SqrMagnitudeTo(player.SelfAnimal.GetPosV())
	if dis > this.dis {
		return b3.FAILURE
	}
	id := ball.GetID()
	tick.Blackboard.Set(this.index, id, "", "")

	var currTime int64 = time.Now().UnixNano() / 1000000
	tick.Blackboard.Set("targetTime", currTime, "", "")
	return b3.SUCCESS
}

//FindAttackTarget
type FindAttackTarget struct {
	b3core.Action
	index  string
	fRange float64
}

func (this *FindAttackTarget) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
	this.fRange = setting.GetProperty("range")
}

func (this *FindAttackTarget) OnTick(tick *b3core.Tick) b3.Status {
	player := tick.GetTarget().(*plr.ScenePlayer)
	//	if player.GetScene().room.isNew {
	//		glog.Info("what is it?")
	//		return b3.FAILURE
	//	}

	//cfg := conf.ConfigMgr_GetMe().GetAISrcDataByScore(player.room.HScore, player.room.roomType)
	cfg := player.AI.GetBevData()
	var frange float64 = this.fRange
	if cfg != nil && cfg.AttackRange > 0 {
		frange = cfg.AttackRange
	}

	anid := player.AI.FindNearAttackAnimal(frange)
	if anid != 0 {
		id := anid
		tick.Blackboard.Set(this.index, id, "", "")
		var currTime int64 = time.Now().UnixNano() / 1000000
		tick.Blackboard.Set("attackTime", currTime, "", "")
		//		glog.Info("FindAttackTarget:", player.id, " id:", id)
		player.Face = id
		return b3.SUCCESS
	}
	tick.Blackboard.Set(this.index, uint32(0), "", "")

	//	glog.Info("nofindNearUnit:", player.id, " range:", this.fRange, " near:", len(player.Others))
	return b3.FAILURE

}

//SubTree
type SubTreeNode struct {
	b3core.Action
	sTree    *b3core.BehaviorTree
	treeName string
}

func (this *SubTreeNode) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.treeName = setting.GetPropertyAsString("treeName")
	this.sTree = GetBevTree(this.treeName)
	if nil == this.sTree {
		glog.Error("SubTreeNode Get SubTree Failed, treeName: ", this.treeName)
	}
	glog.Info("SubTreeNode::Initialize ", this, " treeName ", this.treeName)
}

func (this *SubTreeNode) OnTick(tick *b3core.Tick) b3.Status {
	if nil == this.sTree {
		return b3.ERROR
	}
	if tick.GetTarget() == nil {
		panic("unknow error!")
	}
	player := tick.GetTarget().(*plr.ScenePlayer)
	//	glog.Info("subtree: ", this.treeName, " id ", player.id)
	return this.sTree.Tick(player, tick.Blackboard)
}

//CheckDis
type CheckDisNode struct {
	b3core.Action
	index string
	dis   float64
}

func (this *CheckDisNode) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
	this.dis = setting.GetProperty("dis")
	// var err error
	// this.dis, err = strconv.ParseFloat(disStr, 64)
	// if nil != err {
	// 	glog.Error("CheckDisNode::Initialize failed, ", err)
	// }
}

func (this *CheckDisNode) OnTick(tick *b3core.Tick) b3.Status {
	// glog.Info("CheckDisNode::OnTick ", this.index, this.dis)
	id := tick.Blackboard.GetUInt32(this.index, "", "")
	if id < 1 {
		// glog.Info(1)
		return b3.FAILURE
	}
	player := tick.GetTarget().(*plr.ScenePlayer)
	tball := player.FindVeiwAnimal(uint32(id))
	if tball == nil {
		// glog.Info(2)
		tick.Blackboard.SetTree(this.index, uint32(0), "")
		return b3.FAILURE
	}

	//	glog.Info("dis ", player.SelfAnimal.pos.SqrMagnitudeTo(&tball.pos), " ", this.dis*this.dis)
	if player.SelfAnimal.GetPosV().SqrMagnitudeTo(tball.GetPosV()) > this.dis*this.dis {
		// glog.Info(3)
		return b3.FAILURE
	}
	return b3.SUCCESS
}

//放技能
type BBCastSkillNode struct {
	b3core.Action
	skillId uint32
	index   string
}

func (this *BBCastSkillNode) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.skillId = uint32(setting.GetPropertyAsInt("skillid"))
	this.index = setting.GetPropertyAsString("index")
}

func (this *BBCastSkillNode) OnTick(tick *b3core.Tick) b3.Status {
	player := tick.GetTarget().(*plr.ScenePlayer)
	//	targetId := tick.Blackboard.GetUInt32(this.index, "", "")
	//	glog.Info("bbcast targetid ", targetId)
	if player.Skill.CastSkill(this.skillId, player.Face) {
		//		glog.Info(player.id, " bbcastskill ", this.skillId, " success")
		tick.Blackboard.Set("lastAttackTime", time.Now().UnixNano()/1e6, "", "")
		return b3.SUCCESS
	} else {
		//		glog.Info(player.id, " bbcastskill ", this.skillId, " fail")
	}
	return b3.FAILURE
}

//随机
type RandomComposite struct {
	b3core.Composite
}

func (this *RandomComposite) OnOpen(tick *b3core.Tick) {
	tick.Blackboard.Set("runningChild", -1, tick.GetTree().GetID(), this.GetID())
}

func (this *RandomComposite) OnTick(tick *b3core.Tick) b3.Status {
	//	glog.Info("RandomComposite")
	var child = tick.Blackboard.GetInt("runningChild", tick.GetTree().GetID(), this.GetID())
	if -1 == child {
		child = int(rand.Uint32()) % this.GetChildCount()
	}

	//	glog.Info("random child ", child, " ", this.GetChildCount())
	var status = this.GetChild(child).Execute(tick)
	if status == b3.RUNNING {
		tick.Blackboard.Set("runningChild", child, tick.GetTree().GetID(), this.GetID())
	} else {
		tick.Blackboard.Set("runningChild", -1, tick.GetTree().GetID(), this.GetID())
	}
	return status
}

//HpMoreThan
type HpMoreThan struct {
	b3core.Condition
	rate float32
}

func (this *HpMoreThan) Initialize(setting *b3config.BTNodeCfg) {
	this.Condition.Initialize(setting)
	this.rate = float32(setting.GetProperty("rate"))
}

func (this *HpMoreThan) OnTick(tick *b3core.Tick) b3.Status {
	player := tick.GetTarget().(*plr.ScenePlayer)
	rate := float32(player.SelfAnimal.GetHP()) / float32(player.SelfAnimal.GetHpMax())
	//	glog.Info("rate1 ", rate, " rate2 ", this.rate, " curhp ", player.SelfAnimal.GetHP(), " maxhp ", player.SelfAnimal.GetAttr(AttrHpMax))
	if rate >= this.rate {
		return b3.SUCCESS
	}
	return b3.FAILURE
}

//EnemyToAttackTarget
type EnemyToAttackTarget struct {
	b3core.Condition
	index1 string
	index2 string
}

func (this *EnemyToAttackTarget) Initialize(setting *b3config.BTNodeCfg) {
	this.Condition.Initialize(setting)
	this.index1 = setting.GetPropertyAsString("index1")
	this.index2 = setting.GetPropertyAsString("index2")
}

func (this *EnemyToAttackTarget) OnTick(tick *b3core.Tick) b3.Status {
	id1 := tick.Blackboard.GetUInt32(this.index1, "", "")
	//	id2 := tick.Blackboard.GetUInt32(this.index2, "", "")
	if id1 < 1 {
		return b3.FAILURE
	}

	//	tick.Blackboard.Set(this.index1, id2, "", "")
	tick.Blackboard.Set(this.index2, id1, "", "")
	var currTime int64 = time.Now().UnixNano() / 1000000
	tick.Blackboard.Set("attackTime", currTime, "", "")
	player := tick.GetTarget().(*plr.ScenePlayer)
	player.Face = id1

	return b3.SUCCESS
}

//Parallel
type ParallelComposite struct {
	b3core.Composite
	failCond int //1有一个失败就失败 0全失败才失败
	succCond int //1有一个成功就成功 0全成功才成功
	//如果不能确定状态 那就有running返回running，不然失败
}

func (this *ParallelComposite) Initialize(setting *b3config.BTNodeCfg) {
	this.Composite.Initialize(setting)
	this.failCond = setting.GetPropertyAsInt("fail_cond")
	this.succCond = setting.GetPropertyAsInt("succ_cond")
}

func (this *ParallelComposite) OnTick(tick *b3core.Tick) b3.Status {
	var failCount int
	var succCount int
	var hasRunning bool
	for i := 0; i < this.GetChildCount(); i++ {
		var status = this.GetChild(i).Execute(tick)
		if status == b3.FAILURE {
			failCount++
		} else if status == b3.SUCCESS {
			succCount++
		} else {
			hasRunning = true
		}
	}
	if (this.failCond == 0 && failCount == this.GetChildCount()) || (this.failCond == 1 && failCount > 0) {
		return b3.FAILURE
	}
	if (this.succCond == 0 && succCount == this.GetChildCount()) || (this.succCond == 1 && succCount > 0) {
		return b3.FAILURE
	}
	if hasRunning {
		return b3.RUNNING
	}
	return b3.FAILURE
}

//MoveCtrl
type MoveCtrl struct {
	b3core.Action
	IsOn int
}

func (this *MoveCtrl) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.IsOn = setting.GetPropertyAsInt("isOn")
}

func (this *MoveCtrl) OnTick(tick *b3core.Tick) b3.Status {
	player := tick.GetTarget().(*plr.ScenePlayer)

	if this.IsOn > 0 {
		player.Power = 1
	} else {
		player.Power = 0
	}
	//	glog.Info("player.Power ", player.Power)
	return b3.SUCCESS
}

//CheckDis2 边缘距离
type CheckDis2Node struct {
	b3core.Action
	index string
	dis   float64
}

func (this *CheckDis2Node) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
	this.dis = setting.GetProperty("dis")
	// var err error
	// this.dis, err = strconv.ParseFloat(disStr, 64)
	// if nil != err {
	// 	glog.Error("CheckDisNode::Initialize failed, ", err)
	// }
}

func (this *CheckDis2Node) OnTick(tick *b3core.Tick) b3.Status {
	// glog.Info("CheckDisNode::OnTick ", this.index, this.dis)
	id := tick.Blackboard.GetUInt32(this.index, "", "")
	if id < 1 {
		// glog.Info(1)
		return b3.FAILURE
	}
	player := tick.GetTarget().(*plr.ScenePlayer)
	tball := player.FindVeiwAnimal(uint32(id))
	if tball == nil {
		// glog.Info(2)
		tick.Blackboard.SetTree(this.index, uint32(0), "")
		return b3.FAILURE
	}

	dis := this.dis + player.SelfAnimal.GetRadius() + tball.GetRadius()
	//	glog.Info("dis ", player.SelfAnimal.pos.SqrMagnitudeTo(&tball.pos), " ", dis*dis)
	if player.SelfAnimal.GetPosV().SqrMagnitudeTo(tball.GetPosV()) > dis*dis {
		// glog.Info(3)
		return b3.FAILURE
	}
	return b3.SUCCESS
}

//MoveBack
type MoveBack struct {
	b3core.Action
	index string
}

func (this *MoveBack) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.index = setting.GetPropertyAsString("index")
}

func (this *MoveBack) OnTick(tick *b3core.Tick) b3.Status {
	player := tick.GetTarget().(*plr.ScenePlayer)
	if time.Now().UnixNano()/1e6-tick.Blackboard.GetInt64("lastAttackTime", "", "") > 1000 {
		return b3.FAILURE
	}
	id := tick.Blackboard.GetUInt32(this.index, "", "")
	if id < 1 {
		return b3.FAILURE
	}

	tball := player.FindVeiwAnimal(uint32(id))
	if tball == nil {
		tick.Blackboard.Set(this.index, uint32(0), "", "")
		return b3.FAILURE
	}

	return b3.SUCCESS
}

//WaitSkillIdle
type WaitSkillIdle struct {
	b3core.Action
}

func (this *WaitSkillIdle) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
}

func (this *WaitSkillIdle) OnTick(tick *b3core.Tick) b3.Status {
	player := tick.GetTarget().(*plr.ScenePlayer)
	if player.Skill.GetCurSkillId() == 0 {
		return b3.SUCCESS
	}
	return b3.RUNNING
}
