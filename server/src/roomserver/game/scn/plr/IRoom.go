package plr

import (
	"common"
	"roomserver/game/bll"
	"roomserver/game/team"
	"usercmd"
)

type IRoom interface {
	DeleteActOffline(acc string)
	SaveRoomData(p *ScenePlayer)
	RoomType() uint32
	IsTeamMemberLessThan(teamID uint32, count int) bool
	RemovePlayerById(playerId uint64)
	ID() uint32
	BroadcastMsg(msgNo usercmd.MsgTypeCmd, msg common.Message)
	Frame() uint32
	GetPlayerIScene() IScene
	GetBallIScene() bll.IScene
	SceneID() uint32
	RoomSize() float64
	CellNumX() int
	CellNumY() int
	ReqToSpeak(playerId uint64, banVoice bool, newLogin map[uint64]bool) bool
	GetNewLoginUsers() map[uint64]bool
	GetTeam(teamid uint32) *team.Team
	BroadcastTeamMsg(teamid uint32, msgNo usercmd.MsgTypeCmd, msg common.Message)
}
