// Package rcenterclient has the client to room center server.
package rcenterclient

import (
	"time"

	"github.com/gogo/protobuf/proto"

	"base/env"
	"base/glog"
	"base/gonet"
	"common"
	"roomserver/conf"
	"usercmd"
)

type RCenterClient struct {
	gonet.TcpTask
	mclient *gonet.TcpClient
	Id      uint16
}

var RoomNumGetter interface { // of RoomMgr.GetNum()
	GetNum() int
}
var UserNumGetter interface { // of PlayerTaskMgr.GetNum()
	GetNum() int
}
var ReloginChecker interface {
	CheckRelogin(revCmd *usercmd.ReqCheckRelogin) bool
}
var Terminator interface {
	Terminate()
}

var lclientm *RCenterClient

func GetMe() *RCenterClient {
	if lclientm == nil {
		lclientm = &RCenterClient{
			TcpTask: *gonet.NewTcpTask(nil),
			mclient: &gonet.TcpClient{},
		}
		lclientm.Derived = lclientm
	}
	return lclientm
}

func (this *RCenterClient) Connect() bool {
	loginaddr := env.Get("room", "rcenter")
	if loginaddr == "" {
		loginaddr = env.Get("global", "rcenter")
	}

	conn, err := this.mclient.Connect(loginaddr)
	if err != nil {
		glog.Error("[启动] 连接失败 ", loginaddr)
		return false
	}

	this.Conn = conn
	glog.Info("准备打开room")
	this.Start()

	if env.Get("global", "useoldwatcher") != "true" && env.Get("room", "wlocal") == "" {
		glog.Error("必须配置观战监听端口")
		return false
	}

	this.SendCmd(usercmd.CmdType_Login, &usercmd.ReqServerLogin{
		Address:  env.Get("room", "local"),
		WAddress: env.Get("room", "wlocal"),
		Key:      env.Get("global", "key"),
		SerType:  common.ServerTypeRoom,
	})

	glog.Info("[启动] 连接服务器成功 ", loginaddr)
	return true
}

func (this *RCenterClient) ParseMsg(data []byte, flag byte) bool {

	cmd := usercmd.CmdType(common.GetCmd(data))

	switch cmd {
	case usercmd.CmdType_Login:
		{
			revCmd, ok := common.DecodeCmd(data, flag, &usercmd.RetServerLogin{}).(*usercmd.RetServerLogin)
			if !ok {
				return false
			}

			this.Id = uint16(revCmd.Id)
			this.Verify()
			this.UpdateServer()

			glog.Info("[启动] 连接验证成功 ", env.Get("room", "local"), ",", this.Id)
		}
	case usercmd.CmdType_ChkReLogin:
		{
			revCmd, ok := common.DecodeCmd(data, flag, &usercmd.ReqCheckRelogin{}).(*usercmd.ReqCheckRelogin)
			if !ok {
				return false
			}
			if !ReloginChecker.CheckRelogin(revCmd) { // PlayerTaskMgr_GetMe().CheckRelogin()
				return false
			}
		}
	case usercmd.CmdType_LoadConfig:
		{
			if !conf.ReloadConfig() {
				glog.Info("[GM] 加载配置失败")
				return false
			}
			glog.Info("[GM] 加载配置成功")
		}
	default:
		glog.Error("[服务] 收到未知指令 ", cmd, ",", data, ",", flag)
	}

	return true
}

func (this *RCenterClient) OnClose() {

	if !this.Reset() {
		Terminator.Terminate()
		glog.Error("[服务] 重连失败,服务器启动失败", this.Conn.RemoteAddr())
		return
	}

	glog.Info("[服务] 与管理服务器断开连接,开始重连..")

	for {
		glog.Info("[服务] 重连中..")
		if this.Connect() {
			break
		}
		time.Sleep(time.Second * 3)
	}

	glog.Info("[服务] 重连服务器成功 ", this.Conn.RemoteAddr())
}

func (this *RCenterClient) SendCmd(cmd usercmd.CmdType, msg proto.Message) bool {
	data, flag, err := common.EncodeCmd(uint16(cmd), msg)
	if err != nil {
		glog.Info("[服务] 发送失败 cmd:", cmd, ",len:", len(data), ",err:", err)
		return false
	}
	this.AsyncSend(data, flag)
	return true
}

func (this *RCenterClient) GetId() uint16 {
	return this.Id
}

func (this *RCenterClient) SendCmdToServer(serverid uint16, cmd usercmd.CmdType, msg proto.Message) bool {
	data, flag, err := common.EncodeCmd(uint16(cmd), msg)
	if err != nil {
		glog.Info("[服务] 发送失败 cmd:", cmd, ",len:", len(data), ",err:", err)
		return false
	}
	reqCmd := &usercmd.S2SCmd{
		ServerId: uint32(serverid),
		Flag:     uint32(flag),
		Data:     data,
	}
	return this.SendCmd(usercmd.CmdType_S2S, reqCmd)
}

func (this *RCenterClient) AddRoom(roomtype, roomid, endtime, sceneid uint32, iscoop bool, hscores uint32, uncoop uint32, robot int, isNew bool, level uint32) bool {
	reqCmd := &usercmd.ReqAddRoom{
		RoomType: roomtype,
		RoomId:   roomid,
		EndTime:  endtime,
		IsCoop:   iscoop,
		HScores:  hscores,
		UnCoop:   uncoop,
		Robot:    uint32(robot),
		IsNew:    isNew,
		SceneId:  sceneid,
		Level:    level,
	}
	return this.SendCmd(usercmd.CmdType_AddRoom, reqCmd)
}

func (this *RCenterClient) RemoveRoom(roomtype, roomid uint32) bool {
	reqCmd := &usercmd.ReqRemoveRoom{
		RoomType: roomtype,
		RoomId:   roomid,
	}
	return this.SendCmd(usercmd.CmdType_RemoveRoom, reqCmd)
}

func (this *RCenterClient) UpdateRoom(roomtype, roomid uint32, usernum int32, ustate int, iscoop bool, hscores uint32, robotlevel uint32) bool {
	reqCmd := &usercmd.ReqUpdateRoom{
		RoomType: roomtype,
		RoomId:   roomid,
		UserNum:  usernum,
		UState:   int32(ustate),
		IsCoop:   iscoop,
		HScores:  hscores,
		Robot:    robotlevel,
	}
	return this.SendCmd(usercmd.CmdType_UpdateRoom, reqCmd)
}

func (this *RCenterClient) EndGame(roomid uint32, userid uint64) bool {
	reqCmd := &usercmd.ReqEndGame{
		RoomId: roomid,
		UserId: userid,
	}
	return this.SendCmd(usercmd.CmdType_EndGame, reqCmd)
}

func (this *RCenterClient) UpdateServer() bool {
	reqCmd := &usercmd.ReqUpdateServer{
		RoomNum: uint32(RoomNumGetter.GetNum()),
		UserNum: uint32(UserNumGetter.GetNum()),
	}
	return this.SendCmd(usercmd.CmdType_UpdateServer, reqCmd)
}
