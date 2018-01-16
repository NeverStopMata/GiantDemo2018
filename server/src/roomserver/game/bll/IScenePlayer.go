package bll

import (
	"common"
	"roomserver/game/interfaces"
)

type IScenePlayer interface {
	GetBallScene() IScene
	SceneID() uint32
	GetID() uint64
	GetPower() float64
	IsRunning() bool
	GetIsLive() bool
	RoomType() uint32
	UData() *common.UserData
	KilledByPlayer(killer IScenePlayer)
	RefreshPlayer()
	UpdateExp(addexp int32) bool
	NewSkillBall(sb *BallSkill) interfaces.ISkillBall
	GetFrame() uint32
	GetAngle() float64
	GetFace() uint32
	GetIsRobot() bool
}
