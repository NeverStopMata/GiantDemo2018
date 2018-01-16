package bll

// 球移动 基类

import (
	"base/ape"
	"roomserver/game/bll/internal"
	"roomserver/game/interfaces"
	"roomserver/util"
)

type BallMove struct {
	BallFood
	internal.Force
	speed     util.Vector2        //速度
	angleVel  util.Vector2        //单位速度向量
	PhysicObj *ape.CircleParticle //物理体
}

func (ball *BallMove) GetSpeed() *util.Vector2 {
	return &ball.speed
}

func (ball *BallMove) SetSpeed(v *util.Vector2) {
	ball.speed.X = v.X
	ball.speed.Y = v.Y
}

func (this *BallPlayer) GetAngleVel() *util.Vector2 {
	return &this.angleVel
}

func (ball *BallMove) SqrMagnitudeTo(target interfaces.IBall) float64 {
	x, y := target.GetPos()
	return (ball.Pos.X-x)*(ball.Pos.X-x) + (ball.Pos.Y-y)*(ball.Pos.Y-y)
}
