package internal

import (
	"usercmd"
)

//玩家逻辑消息
type PlayerCmd struct {
	PlayerID uint64
	Cmd      usercmd.MsgTypeCmd
	Data     []byte
	Flag     byte
}
