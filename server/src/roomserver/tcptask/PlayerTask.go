package tcptask

// 玩家网络会话类

import (
	"net"
	"runtime/debug"
	"sync/atomic"
	"time"

	"base/glog"
	"base/gonet"
	"base/locate"
	"common"
	"roomserver/client/dbclient"
	"roomserver/client/rcenterclient"
	"roomserver/interfaces"
	"roomserver/redismgr"
	"roomserver/udp"
	"usercmd"
)

const (
	MAX_NAME_SIZE = 14 // 最大名字长度
)

type IRoomMgr interface {
	AddPlayer(player interfaces.IPlayerTask) bool
}

type IScenePlayerMgr interface {
	GetUDataFromKey(key string) *common.UserData
}

type PlayerTask struct {
	gonet.TcpTask                  // tcp会话
	udpConn       *udp.UdpSess     // udp会话
	key           string           // 登录RoomServer令牌
	id            uint64           // ID
	name          string           // 名字
	room          interfaces.IRoom // 房间
	udata         *common.UserData // 玩家数据
	activeTime    int64            // 用于保持网络会话有效
}

var ScenePlayerMgr IScenePlayerMgr
var RoomMgr IRoomMgr

func NewPlayerTask(conn net.Conn) *PlayerTask {
	s := &PlayerTask{
		TcpTask:    *gonet.NewTcpTask(conn),
		activeTime: time.Now().UnixNano(),
	}
	s.Derived = s

	// 设置发送缓冲区限制
	s.SetSendBuffSizeLimt(256 * 1024)
	return s
}

func (this *PlayerTask) ParseMsg(data []byte, flag byte) bool {
	cmd := usercmd.MsgTypeCmd(common.GetCmd(data))
	if !this.IsVerified() {
		return this.LoginVerify(cmd, data, flag)
	}

	atomic.StoreInt64(&this.activeTime, time.Now().UnixNano())

	if cmd == usercmd.MsgTypeCmd_HeartBeat {
		this.AsyncSend(data, flag)
		return true
	}

	if this.udata.Model == uint32(common.UserModelWatch) {
		glog.Error("[观战] 观战模式,仅能接收消息")
		return false
	}

	if this.room == nil || this.room.IsClosed() {
		return false
	}

	switch cmd {
	default:
		this.room.PostPlayerCmd(this.id, cmd, data, flag)
	}

	return true
}

func (this *PlayerTask) LoginVerify(cmd usercmd.MsgTypeCmd, data []byte, flag byte) bool {
	if cmd != usercmd.MsgTypeCmd_Login {
		glog.Error("[登录] 不是登录验证指令 ", cmd)
		return false
	}

	revCmd, ok := common.DecodeCmd(data, flag, &usercmd.MsgLogin{}).(*usercmd.MsgLogin)
	if !ok {
		this.RetErrorMsg(int(common.ErrorCodeDecode))
		return false
	}

	nickname := []rune(revCmd.Name)
	if len(nickname) > MAX_NAME_SIZE {
		nickname = nickname[:MAX_NAME_SIZE]
	}

	location, cnet := locate.GetLoc(common.GetIP(this.Conn.RemoteAddr().String()))

	glog.Info("[登录] 收到登录请求 ", this.Conn.RemoteAddr(), ",", revCmd.Key, ",", string(nickname), ",", location, ",", cnet)

	// 判断内存中是否有key
	prevUData := ScenePlayerMgr.GetUDataFromKey(revCmd.Key)
	if prevUData != nil {
		//this.OnClose() //老玩家下线
		this.udata = prevUData
	}

	if this.udata == nil {
		udata := &common.UserData{}
		if !redismgr.GetMe().LoadFromRedis(revCmd.Key, udata) {
			this.RetErrorMsg(int(common.ErrorCodeVerify))
			glog.Error("[登录] 验证失败 ", this.Conn.RemoteAddr(), ",", string(nickname), ",", revCmd.Key)
			return false
		}
		this.udata = udata
	}
	this.id = this.udata.Id
	adata, ok := dbclient.GetMe().GetUserById(this.udata.Id)
	if !ok {
		glog.Error("[登录] 操作失败 GetAccInfo ")
		this.RetErrorMsg(int(common.ErrorCodeVerify))
		return false
	}
	//检查是否重复连接
	otask := PlayerTaskMgr_GetMe().GetTask(this.id)
	if otask != nil {
		otask.RetErrorMsg(common.ErrorCodeReLogin)
		otask.Stop()
		otask.Close()
		if nil != otask.room {
			otask.room.ResetPlayerTask(this.id)
		}
		PlayerTaskMgr_GetMe().remove(otask)
		glog.Info("[登录] 发现重复登录 ", otask.id, ",", otask.udata.Account, ",", otask.name, ",", otask.key, ",old:", otask.Conn.RemoteAddr(), " ,new:", this.Conn.RemoteAddr())
		otask = nil
	}

	this.udata.Icon = adata.Icon
	this.udata.PassIcon = adata.PassIcon
	this.udata.PlayNum = adata.PlayNum
	this.udata.Level = adata.Level
	this.udata.HideScore = adata.HideScores

	this.key = revCmd.Key
	this.name = string(nickname)

	this.Verify()
	PlayerTaskMgr_GetMe().Add(this)
	glog.Info("[登录] 验证账号完成 ", this.Conn.RemoteAddr(), ",", this.udata.Id, ",", this.udata.Account, ",", this.name, ",", this.key)
	if RoomMgr.AddPlayer(this) {
		this.online()
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				glog.Error("[异常] 报错 ", err, "\n", string(debug.Stack()))
			}
		}()
		//更新最后玩游戏时间
		dbclient.GetMe().SetLastPlayTime(this.id)
	}()
	glog.Info("[登录] 登录验证成功 ", this.Conn.RemoteAddr(), ",", this.udata.Id, ",", this.udata.Account, ",", this.name, ",", this.key)
	return true
}

// 上线
func (this *PlayerTask) online() {
	// 玩家上线/刷新房间/刷新服务器
	room := this.room
	if room != nil {
		if room.IsClosed() {
			return
		}
		redismgr.GetMe().UserOnline(this.key, this.id, rcenterclient.GetMe().Id, room.RoomType(), room.ID(), this.udata.Model, this.udata.SceneId, this.udata.IsCustom)
		room.AddLoginUser(this.id)
	}
	rcenterclient.GetMe().UpdateServer()
}

// 下线
func (this *PlayerTask) offline() {
	// 玩家下线/刷新房间/刷新服务器
	redismgr.GetMe().UserOffline(this.id)
	rcenterclient.GetMe().UpdateServer()
}

func (this *PlayerTask) OnClose() {
	if !this.IsVerified() {
		return
	}
	// 下线从房间删除
	if this.room != nil {
		this.room.DecPlayerNum()
		if !this.room.IsClosed() {
			this.room.PostToRemovePlayerById(this.id)
		}
	}
	if !PlayerTaskMgr_GetMe().remove(this) {
		glog.Info("[注销] 玩家重复登录 ", this.id, ",", this.udata.Account, ",", this.name, ",", this.key)
		return
	}
	this.offline()

	if this.udpConn != nil {
		this.udpConn.Close()
		if !this.room.IsClosed() {
			this.room.PostBindUdpSession(nil)
		}
	}
	glog.Info("[注销] 下线完成 ", this.udata.Id, ",", this.udata.Account, ",", this.name, ",", this.key)
}

func (this *PlayerTask) SendCmd(cmd usercmd.MsgTypeCmd, msg common.Message) error {
	data, flag, err := common.EncodeGoCmd(uint16(cmd), msg)
	if err != nil {
		glog.Info("[玩家] 发送失败 cmd:", cmd, ",len:", len(data), ",err:", err)
		return err
	}
	this.AsyncSend(data, flag)
	return nil
}

// 优先采用可丢包的原生UDP信道发送, 原生UDP未发送成功则尝试采用KCP或TCP信道发送
func (this *PlayerTask) SendUDPCmd(cmd usercmd.MsgTypeCmd, msg common.Message) error {
	if this.udpConn != nil {
		bSent := this.udpConn.SendUDPCmd(cmd, msg)
		if bSent {
			return nil // UDP 发送成功
		}
	}
	// 原生UDP未发送成功则尝试采用KCP或TCP信道发送
	return this.SendCmd(cmd, msg)
}

func (this *PlayerTask) BindUdpSession(sess *udp.UdpSess) {
	this.udpConn = sess
}

func (this *PlayerTask) RetErrorMsg(ecode int) {
	retCmd := &usercmd.RetErrorMsgCmd{
		RetCode: uint32(ecode),
	}
	this.SendCmd(usercmd.MsgTypeCmd_ErrorMsg, retCmd)
}

func (this *PlayerTask) Name() string {
	return this.name
}

func (this *PlayerTask) ID() uint64 {
	return this.id
}

func (this *PlayerTask) RemoteAddrStr() string {
	return this.Conn.RemoteAddr().String()
}

func (this *PlayerTask) UData() *common.UserData {
	return this.udata
}

func (this *PlayerTask) IsTimeout() bool {
	preTime := atomic.LoadInt64(&this.activeTime)
	timeout := float64(time.Now().UnixNano()-preTime) / float64(time.Minute)
	return timeout > float64(_TASK_MAX_TIMEOUT)
}

func (this *PlayerTask) SetRoom(room interfaces.IRoom) {
	this.room = room
}

// 登录RoomServer令牌
func (this *PlayerTask) Key() string {
	return this.key
}
