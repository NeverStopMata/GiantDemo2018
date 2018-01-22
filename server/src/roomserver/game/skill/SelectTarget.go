package skill

// 几种选取目标的方式。如用于释放技能时

import (
	b3core "base/behavior3go/core"
	_ "base/glog"
	"roomserver/game/bll"
	"roomserver/game/interfaces"
	"roomserver/game/scn/plr"
	"roomserver/util"
)

// 获取朝向上最近的目标
func FindNearTarget(tick *b3core.Tick, player *plr.ScenePlayer) (interfaces.IBall, interfaces.BallKind) {
	//	angleVel := GetPlayerDir(tick, player)

	//	var rect util.Square
	//	rect.CopyFrom(player.GetViewRect())
	//	rect.SetRadius(GetAttackRange(tick, player))
	//	cells := player.GetScene().GetAreaCells(&rect)

	//	minball_feed, min_feed := player.FindNearBallByKind(interfaces.BallKind_Feed, angleVel, cells, 0)
	//	minball_player, min_player := player.FindNearBallByKind(interfaces.BallKind_Player, angleVel, cells, 0)
	//	minball_ballskill, min_ballskill := player.FindNearBallByKind(interfaces.BallKind_Skill, angleVel, cells, 0)
	//	if minball_player == nil && minball_feed == nil && minball_ballskill == nil {
	//		return nil, interfaces.BallKind_None
	//	}

	//	if min_feed <= min_player && min_feed <= min_ballskill {
	//		return minball_feed, interfaces.BallKind_Feed
	//	} else if min_ballskill <= min_player && min_ballskill <= min_feed {
	//		return minball_ballskill, interfaces.BallKind_Skill
	//	} else {
	//		return minball_player, interfaces.BallKind_Player
	//	}//mata
	angleVel := GetPlayerDir(tick, player)

	var rect util.Square
	rect.CopyFrom(player.GetViewRect())
	rect.SetRadius(GetAttackRange(tick, player))
	cells := player.GetScene().GetAreaCells(&rect)

	minball_player, _ := player.FindNearBallByKind(interfaces.BallKind_Player, angleVel, cells, 0)
	if minball_player == nil {
		return nil, interfaces.BallKind_None
	}
	return minball_player, interfaces.BallKind_Player

}

// 获取朝向上所有目标
func FindTarget_SemiCircle(tick *b3core.Tick, player *plr.ScenePlayer) ([]interfaces.IBall, []interfaces.BallKind) {
	var balllist []interfaces.IBall
	var balltype []interfaces.BallKind

	dir := GetPlayerDir(tick, player)

	// player
	for _, o := range player.Others {
		if o.IsLive == false {
			continue
		}
		ball := o.SelfAnimal
		if util.IsSameDir(dir, ball.GetPosV(), player.SelfAnimal.GetPosV()) == false {
			continue
		}
		balllist = append(balllist, ball)
		balltype = append(balltype, interfaces.BallKind_Player)
	}

	var rect util.Square
	rect.CopyFrom(player.GetViewRect())
	rect.SetRadius(GetAttackRange(tick, player))
	cells := player.GetScene().GetAreaCells(&rect)

	// ballskill
	for _, cell := range cells {
		for _, ball := range cell.Skills {
			if util.IsSameDir(dir, ball.GetPosV(), player.SelfAnimal.GetPosV()) == false {
				continue
			}
			balllist = append(balllist, ball)
			balltype = append(balltype, interfaces.BallKind_Skill)
		}
	}

	// feed
	for _, cell := range cells {
		for _, ball := range cell.Feeds {
			if util.IsSameDir(dir, ball.GetPosV(), player.SelfAnimal.GetPosV()) == false {
				continue
			}
			balllist = append(balllist, ball)
			balltype = append(balltype, interfaces.BallKind_Feed)
		}
	}

	return balllist, balltype
}

// 获取所有目标
func FindTarget_Circle(tick *b3core.Tick, player *plr.ScenePlayer) ([]interfaces.IBall, []interfaces.BallKind) {
	var balllist []interfaces.IBall
	var balltype []interfaces.BallKind

	// player
	for _, o := range player.Others {
		if o.IsLive == false {
			continue
		}
		ball := o.SelfAnimal
		balllist = append(balllist, ball)
		balltype = append(balltype, interfaces.BallKind_Player)
	}

	var rect util.Square
	rect.CopyFrom(player.GetViewRect())
	rect.SetRadius(GetAttackRange(tick, player))
	cells := player.GetScene().GetAreaCells(&rect)

	// ballskill
	for _, cell := range cells {
		for _, ball := range cell.Skills {
			balllist = append(balllist, ball)
			balltype = append(balltype, interfaces.BallKind_Skill)
		}
	}

	// feed
	for _, cell := range cells {
		for _, ball := range cell.Feeds {
			balllist = append(balllist, ball)
			balltype = append(balltype, interfaces.BallKind_Feed)
		}
	}
	return balllist, balltype
}

// 获取玩家朝向
func GetPlayerDir(tick *b3core.Tick, player *plr.ScenePlayer) *util.Vector2 {
	angleVel := &util.Vector2{}
	usedefault := true
	targetId := tick.Blackboard.GetUInt32("skillTargetId", "", "")
	if 0 != targetId {
		tball := player.FindVeiwAnimal(targetId)
		if tball != nil {
			x, y := tball.GetPos()
			angleVel.X = x - player.SelfAnimal.GetPosV().X
			angleVel.Y = y - player.SelfAnimal.GetPosV().Y
			usedefault = false
		}
	}
	if usedefault {
		angleVel.X = player.SelfAnimal.GetAngleVel().X
		angleVel.Y = player.SelfAnimal.GetAngleVel().Y
	}

	return angleVel
}

// 攻击范围
func GetAttackRange(tick *b3core.Tick, player *plr.ScenePlayer) float64 {
	attackRange := tick.Blackboard.Get("attackRange", "", "")
	if attackRange != nil {
		r := attackRange.(float64) * float64(10)
		if r >= 0 {
			return r * player.SelfAnimal.GetSizeScale()
		}
	}
	return player.SelfAnimal.GetEatRange()
}

// 是否可以攻击
func IsCanAttack(tick *b3core.Tick, player *plr.ScenePlayer, target interfaces.IBall) bool {
	distance := player.SelfAnimal.SqrMagnitudeTo(target)
	eatRange := GetAttackRange(tick, player)
	return distance <= (eatRange+target.GetRect().Radius)*(eatRange+target.GetRect().Radius)
}

func IsCanAttackPlayer(tick *b3core.Tick, player *plr.ScenePlayer, target *bll.BallPlayer) bool {
	if player.SelfAnimal.PreTryHit(target) == false {
		return false
	}
	distance := player.SelfAnimal.SqrMagnitudeTo(target)
	eatRange := GetAttackRange(tick, player)
	return distance <= (eatRange+target.GetRect().Radius)*(eatRange+target.GetRect().Radius)
}
