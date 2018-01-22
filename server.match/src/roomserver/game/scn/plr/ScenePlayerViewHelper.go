package plr

// 玩家视野相关 辅助类

import (
	"math"
	"roomserver/conf"
	"roomserver/game/bll"
	"roomserver/game/cll"
	"roomserver/game/interfaces"
	"roomserver/util"
	"usercmd"
)

type ScenePlayerViewHelper struct {
	ViewRect       *util.Square              // 玩家视野
	RealViewRect   util.Square               // 玩家视野（根据玩家原始视野得到所有cell的外边框）
	LookCells      map[int]*cll.Cell         // 根据玩家原始视野得到所有cell集合
	LookFeeds      map[uint32]*bll.BallFeed  // 视野中的feed（用于sendSceneMsg）
	LookBallSkill  map[uint32]*bll.BallSkill // 视野中的技能球（用于sendSceneMsg）
	LookBallPlayer map[uint64]*ScenePlayer   // 视野中的玩家 （用于sendSceneMsg）
	LookBallFoods  map[uint32]*bll.BallFood  // 视野中的food（用于sendSceneMsg）
	Others         map[uint64]*ScenePlayer   // 视野中其它玩家
	RoundPlayers   []*ScenePlayer            // 周围玩家，包含死亡玩家
}

func (this *ScenePlayerViewHelper) Init() {
	this.ViewRect = &util.Square{}
	this.LookCells = make(map[int]*cll.Cell)
	this.LookFeeds = make(map[uint32]*bll.BallFeed)
	this.LookBallSkill = make(map[uint32]*bll.BallSkill)
	this.LookBallPlayer = make(map[uint64]*ScenePlayer)
	this.LookBallFoods = make(map[uint32]*bll.BallFood)
	this.Others = make(map[uint64]*ScenePlayer)
}

// 获取视野大小, 更新视野
func (this *ScenePlayerViewHelper) UpdateView(scene IScene, selfBall *bll.BallPlayer, roomSize float64, cellNumX, cellNumY int) {
	// 没有动过，直接返回
	if math.Abs(selfBall.GetRect().X-this.ViewRect.X) < util.EPSILON &&
		math.Abs(selfBall.GetRect().Y-this.ViewRect.Y) < util.EPSILON {
		return
	}

	// TODO : 视野调整，具体需要等客户端确定好摄像机后，再根据手机各分辨率下，找一个最大包含屏幕的区域大小。
	//        目前为了可以适配py_guiclient，暂时修改如下
	this.ViewRect.CopyFrom(selfBall.GetRect())
	this.ViewRect.SetRadius(9)

	if conf.ConfigMgr_GetMe().Global.AllMap > 0 {
		this.ViewRect.SetRadius(roomSize * 10)
	}

	this.RealViewRect.CopyFrom(this.ViewRect)
	minX := int(math.Max(math.Floor(this.RealViewRect.Left/cll.CellWidth)*cll.CellWidth, 0))
	maxX := int(math.Min(math.Floor(this.RealViewRect.Right/cll.CellWidth)*cll.CellWidth, float64(cellNumX-1)*cll.CellWidth))
	minY := int(math.Max(math.Floor(this.RealViewRect.Bottom/cll.CellHeight)*cll.CellHeight, 0))
	maxY := int(math.Min(math.Floor(this.RealViewRect.Top/cll.CellHeight)*cll.CellHeight, float64(cellNumY-1)*cll.CellHeight))
	this.RealViewRect.Left = float64(minX)
	this.RealViewRect.Right = float64(maxX) + cll.CellWidth
	this.RealViewRect.Bottom = float64(minY)
	this.RealViewRect.Top = float64(maxY) + cll.CellHeight

	newCells := scene.GetAreaCells(this.ViewRect)
	this.LookCells = make(map[int]*cll.Cell)
	for _, newCell := range newCells {
		this.LookCells[newCell.ID()] = newCell
	}
}
func (this *ScenePlayerViewHelper) ResetMsg() {
	this.LookBallPlayer = make(map[uint64]*ScenePlayer)
	for k, v := range this.Others {
		this.LookBallPlayer[k] = v
	}
}

func (this *ScenePlayerViewHelper) UpdateVeiwFeeds() (addFeeds []*usercmd.MsgBall, delFeeds []uint32) {
	newFeeds := make(map[uint32]*bll.BallFeed)
	for _, cell := range this.LookCells {
		for _, feed := range cell.Feeds {
			newFeeds[feed.GetID()] = feed
		}
	}
	//add
	for _, feed := range newFeeds {
		if _, ok := this.LookFeeds[feed.GetID()]; !ok {
			addFeeds = append(addFeeds, bll.FeedToMsgBall(feed))
		}
	}

	//del
	for _, feed := range this.LookFeeds {
		if _, ok := newFeeds[feed.GetID()]; !ok {
			delFeeds = append(delFeeds, feed.GetID())
		}
	}

	this.LookFeeds = newFeeds
	return
}

func (this *ScenePlayerViewHelper) UpdateVeiwBallSkill() (adds []*usercmd.MsgBall, dels []uint32) {
	news := make(map[uint32]*bll.BallSkill)
	for _, cell := range this.LookCells {
		for _, ball := range cell.Skills {
			news[ball.GetID()] = ball
		}
	}
	//add
	for _, ball := range news {
		if _, ok := this.LookBallSkill[ball.GetID()]; !ok {
			adds = append(adds, bll.SkillToMsgBall(ball))
		}
	}

	//del
	for _, ball := range this.LookBallSkill {
		if _, ok := news[ball.GetID()]; !ok {
			dels = append(dels, ball.GetID())
		}
	}

	this.LookBallSkill = news
	return
}

func (this *ScenePlayerViewHelper) updateViewBallPlayer() (adds []*usercmd.MsgPlayerBall, dels []uint32) {
	//add
	for _, ball := range this.Others {
		if _, ok := this.LookBallPlayer[ball.ID]; !ok {
			adds = append(adds, bll.PlayerBallToMsgBall(ball.SelfAnimal))
		}
	}
	for _, ball := range this.LookBallPlayer {
		if _, ok := this.Others[ball.ID]; !ok {
			dels = append(dels, ball.SelfAnimal.GetID())
		}
	}
	return
}

func (this *ScenePlayerViewHelper) UpdateVeiwFoods() (addFoods []*usercmd.MsgBall, delFoods []uint32) {
	newFoods := make(map[uint32]*bll.BallFood)
	for _, cell := range this.LookCells {
		for _, food := range cell.Foods {
			newFoods[food.GetID()] = food
		}
	}
	//add
	for _, food := range newFoods {
		if _, ok := this.LookBallFoods[food.GetID()]; !ok {
			addFoods = append(addFoods, bll.FoodToMsgBall(food))
		}
	}

	//del
	for _, food := range this.LookBallFoods {
		if _, ok := newFoods[food.GetID()]; !ok {
			delFoods = append(delFoods, food.GetID())
		}
	}

	this.LookBallFoods = newFoods
	return
}

//更新玩家当前帧视野
func (this *ScenePlayerViewHelper) UpdateViewPlayers(scene IScene, selfBall *bll.BallPlayer) {
	this.Others = make(map[uint64]*ScenePlayer)
	this.RoundPlayers = this.RoundPlayers[:0]
	for _, player := range scene.GetPlayers() {
		if selfBall.GetPlayerId() != player.ID {
			_, _, ok1 := this.RealViewRect.ContainsCircle(player.SelfAnimal.Pos.X, player.SelfAnimal.Pos.Y, 0)
			if ok1 {
				if player.IsLive {
					this.Others[player.ID] = player
				}
			}

			_, _, ok2 := player.RealViewRect.ContainsCircle(selfBall.Pos.X, selfBall.Pos.Y, 0)
			if ok2 {
				this.RoundPlayers = append(this.RoundPlayers, player)
			}
		}
	}
}

//寻找最近的类型目标
func (this *ScenePlayerViewHelper) FindNearBallByKind(selfBall *bll.BallPlayer, kind interfaces.BallKind, dir *util.Vector2, cells []*cll.Cell, ballType uint32) (interfaces.IBall, float64) {
	var minball interfaces.IBall
	var min float64 = 10000

	if kind == interfaces.BallKind_None && ballType == 0 {
		return nil, min
	}

	if kind == interfaces.BallKind_None {
		kind = interfaces.BallTypeToKind(usercmd.BallType(ballType))
	}

	//寻找最近目标
	pos := selfBall.GetPosV()
	if kind == interfaces.BallKind_Player {
		for _, o := range this.Others {
			if o.IsLive == false {
				continue
			}
			ball := o.SelfAnimal
			if dir != nil && util.IsSameDirTotally(dir, ball.GetPosV(), selfBall.GetPosV()) == false {
				continue
			}
			dis := ball.Pos.SqrMagnitudeTo(pos)
			if minball == nil || dis < min {
				min = dis
				minball = ball
			}
		}
	} else {
		if cells != nil {
			for _, cell := range cells {
				ball, dis := cell.FindNearBallByKind(selfBall, pos, kind, dir, ballType)
				if ball != nil && (minball == nil || dis < min) {
					minball = ball
					min = dis
				}
			}
		} else {
			for _, cell := range this.LookCells {
				ball, dis := cell.FindNearBallByKind(selfBall, pos, kind, dir, ballType)
				if ball != nil && (minball == nil || dis < min) {
					minball = ball
					min = dis
				}
			}
		}
	}
	return minball, min
}

func (this *ScenePlayerViewHelper) FindNearBall(id uint32) interfaces.IBall {
	ani := this.FindVeiwAnimal(id)
	if ani != nil {
		return ani
	}

	for _, cell := range this.LookCells {
		if tball, ok := cell.NoTypeFind(id); ok {
			return tball
		}
	}
	return nil
}

func (this *ScenePlayerViewHelper) FindVeiwAnimal(ballId uint32) *bll.BallPlayer {
	for _, viewPlayer := range this.Others {
		if viewPlayer.SelfAnimal.GetID() == ballId {
			return viewPlayer.SelfAnimal
		}
	}
	return nil
}

func (this *ScenePlayerViewHelper) GetViewRect() *util.Square {
	return this.ViewRect
}
