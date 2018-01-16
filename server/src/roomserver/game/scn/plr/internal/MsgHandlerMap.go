package internal

import (
	"base/glog"
	"usercmd"
)

// 网络消息处理器
type MsgHandler func(data []byte, flag byte)

type MsgHandlerMap struct {
	handlerMap map[usercmd.MsgTypeCmd]MsgHandler
}

func (this *MsgHandlerMap) Init() {
	this.handlerMap = make(map[usercmd.MsgTypeCmd]MsgHandler)
}

func (this *MsgHandlerMap) RegisterHandler(cmd usercmd.MsgTypeCmd, cb MsgHandler) {
	this.handlerMap[cmd] = cb
}

func (this *MsgHandlerMap) Call(cmd usercmd.MsgTypeCmd, data []byte, flag byte) {
	cb, ok := this.handlerMap[cmd]
	if ok {
		cb(data, flag)
	} else {
		glog.Error("MsgHandlerMap.Call: unknow cmd,", cmd)
	}
}
