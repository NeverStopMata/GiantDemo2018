package bll

import (
	"base/ape"
	"roomserver/game/interfaces"
)

type IScene interface {
	SceneID() uint32
	AddBall(ball interfaces.IBall)
	AddFeedPhysic(feed ape.IAbstractParticle)
	AddAnimalPhysic(animal ape.IAbstractParticle)
	GetRandomPos() (x, y float64, op bool)
	RoomSize() float64
	ReturnId(id uint32)
	UpdateSkillBallCell(ball *BallSkill, oldCellID int)
}
