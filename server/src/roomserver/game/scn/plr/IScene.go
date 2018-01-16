package plr

import (
	"base/ape"
	"roomserver/game/bll"
	"roomserver/game/cll"
	"roomserver/game/interfaces"
	"roomserver/util"
)

type IScene interface {
	NewBallPlayerId() uint32
	AddBall(ball interfaces.IBall)
	AddOffline(player *ScenePlayer)
	RemoveBall(ball interfaces.IBall)
	RemoveAnimalPhysic(animal ape.IAbstractParticle)
	Frame() uint32
	GetAreaCells(s *util.Square) (cells []*cll.Cell)
	GetPlayers() map[uint64]*ScenePlayer
	RoomSize() float64
	RoomType() uint32
	NewBallSkillId() uint32
	SceneID() uint32
	GetCell(px, py float64) (*cll.Cell, bool)
	RemoveFeed(feed *bll.BallFeed)
}
