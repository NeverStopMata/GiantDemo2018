package bll

// 食物球

import (
	"roomserver/conf"
	"roomserver/game/consts"
	"roomserver/game/interfaces"
	"roomserver/util"
	"usercmd"
)

type BallFood struct {
	id         uint32           //动态id
	typeID     uint16           //xml表里id
	BallType   usercmd.BallType //大类型
	Pos        util.Vector2
	radius     float64
	rect       util.Square
	birthPoint interfaces.IBirthPoint
	exp        int32
	HLState    int32 //2：表示在上；1：表示在下；-2：表示玩家在上升的过程中，处于虚空！；-1表示玩家在下降过程中
}

func NewBallFood(id uint32, typeId uint16, x, y float64, scene IScene) *BallFood {
	var radius float32 = conf.ConfigMgr_GetMe().GetFoodSize(scene.SceneID(), typeId)
	ballType := conf.ConfigMgr_GetMe().GetFoodBallType(scene.SceneID(), typeId)
	ball := &BallFood{
		id:       id,
		typeID:   typeId,
		Pos:      util.Vector2{x, y},
		BallType: ballType,
		radius:   float64(radius),
	}
	ball.ResetRect()
	ball.SetExp(consts.DefaultBallFoodExp)
	scene.AddBall(ball)
	return ball
}

func (ball *BallFood) GetRect() *util.Square {
	return &ball.rect
}

func (ball *BallFood) OnReset() {

}

func (ball *BallFood) GetID() uint32 {
	return ball.id
}

func (ball *BallFood) SetID(id uint32) {
	ball.id = id
}

func (ball *BallFood) GetTypeId() uint16 {
	return ball.typeID
}

func (ball *BallFood) GetType() usercmd.BallType {
	return ball.BallType
}

func (ball *BallFood) GetPos() (float64, float64) {
	return ball.Pos.X, ball.Pos.Y
}

func (ball *BallFood) SetPos(x, y float64) {
	ball.Pos.X = x
	ball.Pos.Y = y
}

func (ball *BallFood) GetPosV() *util.Vector2 {
	return &ball.Pos
}

func (this *BallFood) SetPosV(pos util.Vector2) {
	this.Pos = pos
}

func (ball *BallFood) SetExp(exp int32) {
	ball.exp = exp
}

func (ball *BallFood) ResetRect() {
	ball.rect.X = ball.Pos.X
	ball.rect.Y = ball.Pos.Y
	ball.rect.SetRadius(ball.radius)
}

func (ball *BallFood) SetBirthPoint(birthPoint interfaces.IBirthPoint) {
	ball.birthPoint = birthPoint
}

func (ball *BallFood) GetBirthPoint() interfaces.IBirthPoint {
	return ball.birthPoint
}

func (ball *BallFood) GetRadius() float64 {
	return ball.radius
}
