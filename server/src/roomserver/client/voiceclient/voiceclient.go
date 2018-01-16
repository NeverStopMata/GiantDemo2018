// Package voiceclient has the client to voice server.
package voiceclient

import (
	"base/env"
	"base/glog"
	"base/rpc"
	"common"
	"sync/atomic"
	"time"
)

type VoiceClient struct {
	client   *rpc.Client
	isclosed int32
}

type IRoomMgr interface {
	SyncVoiceRoom()
}

var vclientm *VoiceClient

var RoomMgr IRoomMgr

func GetMe() *VoiceClient {
	if vclientm == nil {
		vclientm = &VoiceClient{
			isclosed: 1,
		}
	}
	return vclientm
}

func (this *VoiceClient) Connect() bool {
	if atomic.LoadInt32(&this.isclosed) != 1 {
		return false
	}
	var err error
	vaddr := env.Get("room", "voicetcp")
	this.client, err = rpc.Dial("tcp", vaddr, this.ReConnect)
	if err != nil {
		glog.Error("[RPC] 连接失败 ", err)
	}
	atomic.StoreInt32(&this.isclosed, 0)
	glog.Info("[RPC] 连接服务器成功 ", vaddr)
	return true
}

func (this *VoiceClient) ReConnect() bool {
	if atomic.LoadInt32(&this.isclosed) == 1 {
		return false
	}
	atomic.StoreInt32(&this.isclosed, 1)
	for {
		glog.Info("[RPC] 重连中...")
		if this.Connect() {
			break
		}
		time.Sleep(time.Second * 2)
	}
	RoomMgr.SyncVoiceRoom()
	return true
}

func (this *VoiceClient) RemoteCall(serviceMethod string, args interface{}, reply interface{}) bool {
	if this.client == nil {
		glog.Error("[RPC] 未初始化 ", serviceMethod)
		return false
	}
	err := this.client.Call(serviceMethod, args, reply)
	if err != nil {
		if err == rpc.ErrShutdown {
			this.ReConnect()
		}
		glog.Error("[RPC] 调用失败 ", serviceMethod, ",", err)
		return false
	}
	return true
}

// 创建房间
func (this *VoiceClient) NewRoom(roomid, lasttime uint32) bool {
	reply := common.RetNewRoom{}
	if !this.RemoteCall("RPCTask.NewRoom", &common.ReqNewRoom{RoomId: roomid, LastTime: lasttime}, &reply) {
		return false
	}
	return true
}

// 玩家发言
func (this *VoiceClient) ToSpeak(roomid uint32, userid uint64, tid, lasttime uint32, banVoice bool, newLogin map[uint64]bool) bool {
	reply := common.RetUserSpeak{}
	args := common.ReqUserSpeak{
		BanVoice: banVoice,
		RoomId:   roomid,
		UserId:   userid,
		TeamId:   tid,
		LastTm:   lasttime,
		NewLogin: newLogin,
	}
	if !this.RemoteCall("RPCTask.ToSpeak", &args, &reply) {
		return false
	}
	return true
}
