package bll

// 可以动态生成的障碍物，也可以在runtime移除

import (
	"base/ape"
	"roomserver/game/bll/internal"
	"roomserver/game/interfaces"
	"roomserver/util"
)

type DynamicObstacle struct {
	BallFood
	internal.Force
	speed     util.Vector2           //速度
	angleVel  util.Vector2           //单位速度向量
	PhysicObj *ape.RectangleParticle //物理体
}

func (ball *DynamicObstacle) GetSpeed() *util.Vector2 {
	return &ball.speed
}

func (ball *DynamicObstacle) SetSpeed(v *util.Vector2) {
	ball.speed.X = v.X
	ball.speed.Y = v.Y
}

//func (this *BallPlayer) GetAngleVel() *util.Vector2 {
//	return &this.angleVel
//}

func (ball *DynamicObstacle) SqrMagnitudeTo(target interfaces.IBall) float64 {
	x, y := target.GetPos()
	return (ball.Pos.X-x)*(ball.Pos.X-x) + (ball.Pos.Y-y)*(ball.Pos.Y-y)
}
