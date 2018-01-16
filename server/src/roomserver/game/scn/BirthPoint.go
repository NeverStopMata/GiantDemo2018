package scn

import (
	"math"

	"base/glog"
	bmath "base/math"
	"roomserver/game/bll"
	"roomserver/game/interfaces"
	"usercmd"
)

type BirthPoint struct {
	id             uint32
	scene          *Scene
	pos            bmath.Vector2
	ballTypeId     uint16
	ballType       uint16
	birthTime      float64
	birthMax       uint32
	birthRadiusMin float32
	birthRadiusMax float32
	childrenCount  uint32
	birthTimer     float64
}

//创建动态出生点 食物、 动态障碍物 (BallFood、 BallFeed)
func NewBirthPoint(id uint32, x, y, rMin, rMax float32, ballTypeId uint16, ballType uint16, birthTime float64, birthMax uint32, scene *Scene) *BirthPoint {

	point := &BirthPoint{
		id:         id,
		pos:        bmath.Vector2{x, y},
		ballTypeId: ballTypeId,
		ballType:   ballType,
		birthTime:  birthTime,
		birthMax:   birthMax,
		scene:      scene,
	}
	point.birthRadiusMin = rMin
	point.birthRadiusMax = rMax
	point.Init()
	return point
}

func (this *BirthPoint) Init() {
	var i uint32 = 0
	for ; i < this.birthMax; i++ {
		this.CreateUnit()
	}
}

func (this *BirthPoint) CreateUnit() interfaces.IBall {
	this.childrenCount++
	scene := this.scene
	var ball interfaces.IBall
	ballType := interfaces.BallTypeToKind(usercmd.BallType(this.ballType))
	switch ballType {
	case interfaces.BallKind_Food:
		posNew := BallFood_InitPos(&this.pos, usercmd.BallType(this.ballType), this.birthRadiusMin, this.birthRadiusMax)
		ball = bll.NewBallFood(this.id, this.ballTypeId, float64(posNew.X), float64(posNew.Y), scene)
	case interfaces.BallKind_Feed:
		x := math.Floor(float64(this.pos.X)) + 0.25
		y := math.Floor(float64(this.pos.Y)) + 0.25
		ball = bll.NewBallFeed(scene, this.ballTypeId, this.id, x, y)
	default:
		glog.Error("CreateUnit unknow ballType:", ballType, "  typeid:", this.ballTypeId)
	}

	ball.SetBirthPoint(this)
	return ball
}

func (this *BirthPoint) Refresh(perTime float64, scene *Scene) {
	if this.childrenCount >= this.birthMax {
		return
	}
	if this.birthTimer >= this.birthTime {
		this.birthTimer = 0
		this.CreateUnit()
	} else {
		this.birthTimer += perTime
	}
}

func (this *BirthPoint) OnChildRemove(ball interfaces.IBall) {
	this.childrenCount--
}
