package scn

import (
	"common"
	"roomserver/conf"
	"roomserver/game/scn/plr"
	"roomserver/game/team"
	"usercmd"
)

type IRoom interface {
	NewScenePlayer(udata *common.UserData, name string, isRobot bool) *plr.ScenePlayer
	RoomType() uint32
	ID() uint32
	BroadcastMsg(msgNo usercmd.MsgTypeCmd, msg common.Message)
	StartTime() int64
	EndTime() int64
	Name() string
	GetTeam(teamid uint32) *team.Team
	GetRoomTypeConfig() *conf.XmlRoomModel
	SaveRoomData(p *plr.ScenePlayer)
	RemovePlayerById(playerId uint64)
	SceneID() uint32
	RoomSize() float64
	NewRobotUID() uint64
}
