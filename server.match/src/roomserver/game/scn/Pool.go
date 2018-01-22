package scn

// 球、消息等 对象池
//    TODO: 增加 ballFeed、ballFood、ballSkill 对象池

import (
	"usercmd"
)

type MsgPool struct {
	MsgLoginResult usercmd.MsgLoginResult
	MsgAddPlayer   usercmd.MsgAddPlayer
	MsgSpeakUser   usercmd.RetSpeakUser
	MsgToSpeak     usercmd.RetToSpeak
	MsgTopRank     usercmd.MsgTop
}

func NewPool() (pool *MsgPool) {
	pool = &MsgPool{}
	pool.MsgAddPlayer.Player = &usercmd.MsgPlayer{}
	return
}
