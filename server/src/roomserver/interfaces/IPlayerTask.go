package interfaces

import (
	"common"
	"roomserver/udp"
	"usercmd"
)

type IPlayerTask interface {
	SendCmd(cmd usercmd.MsgTypeCmd, msg common.Message) error
	Name() string
	RetErrorMsg(ecode int)
	ID() uint64
	Close()
	RemoteAddrStr() string
	UData() *common.UserData
	IsTimeout() bool
	SetRoom(room IRoom)
	Key() string
	SendUDPCmd(cmd usercmd.MsgTypeCmd, msg common.Message) error
	AsyncSend(buffer []byte, flag byte) bool
	BindUdpSession(sess *udp.UdpSess)
}
