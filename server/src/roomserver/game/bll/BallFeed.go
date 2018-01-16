package bll

// 蘑菇球

import (
	"base/ape"
	"roomserver/conf"
	"roomserver/util"
)

type BallFeed struct {
	BallMove
}

func NewBallFeed(scene IScene, typeId uint16, id uint32, x, y float64) *BallFeed {
	radius := float64(conf.ConfigMgr_GetMe().GetFoodSize(scene.SceneID(), typeId))
	ballType := conf.ConfigMgr_GetMe().GetFoodBallType(scene.SceneID(), typeId)
	ball := &BallFeed{
		BallMove: BallMove{
			BallFood: BallFood{
				id:       id,
				typeID:   typeId,
				BallType: ballType,
				Pos:      util.Vector2{float64(x), float64(y)},
				radius:   float64(radius),
			},
			PhysicObj: ape.NewCircleParticle(float32(x), float32(y), float32(radius)),
		},
	}
	ball.ResetRect()
	ball.PhysicObj.SetFixed(true)
	scene.AddBall(ball)
	scene.AddFeedPhysic(ball.PhysicObj)
	return ball
}
