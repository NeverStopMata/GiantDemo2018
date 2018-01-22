package interfaces

import (
	"roomserver/udp"
	"usercmd"
)

type IRoom interface {
	IsClosed() bool
	PostPlayerCmd(playerID uint64, cmd usercmd.MsgTypeCmd, data []byte, flag byte)
	ResetPlayerTask(PlayerTaskID uint64)
	RoomType() uint32
	ID() uint32
	AddLoginUser(UID uint64)
	DecPlayerNum()
	PostToRemovePlayerById(playerID uint64)
	PostBindUdpSession(data map[uint64]*udp.UdpSess)
}
