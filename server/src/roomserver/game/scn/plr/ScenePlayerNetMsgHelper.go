package plr

// 房间玩家协议处理 辅助类

import (
	"time"

	"base/glog"
	"common"
	"roomserver/game/scn/plr/internal"
	"roomserver/redismgr"
	"usercmd"
)

type ScenePlayerNetMsgHelper struct {
	msgHandlerMap internal.MsgHandlerMap // 玩家协议处理器
	selfPlayer    *ScenePlayer           // 玩家自身的引用
}

func (this *ScenePlayerNetMsgHelper) Init(selfPlayer *ScenePlayer) {
	this.selfPlayer = selfPlayer
	this.msgHandlerMap.Init()
	this.RegCmds()
}

//注册网络消息
func (this *ScenePlayerNetMsgHelper) RegCmds() {
	this.msgHandlerMap.RegisterHandler(usercmd.MsgTypeCmd_SceneChat, this.OnSceneChat)
	this.msgHandlerMap.RegisterHandler(usercmd.MsgTypeCmd_Run, this.OnRun)
	this.msgHandlerMap.RegisterHandler(usercmd.MsgTypeCmd_CastSkill, this.OnCastSkill)
	this.msgHandlerMap.RegisterHandler(usercmd.MsgTypeCmd_ActCloseSocket, this.OnCloseSocket)
	this.msgHandlerMap.RegisterHandler(usercmd.MsgTypeCmd_TeamNotice, this.OnTeamNotice)
	this.msgHandlerMap.RegisterHandler(usercmd.MsgTypeCmd_Move, this.OnNetMove)
	this.msgHandlerMap.RegisterHandler(usercmd.MsgTypeCmd_ReLife, this.OnNetReLife)
	this.msgHandlerMap.RegisterHandler(usercmd.MsgTypeCmd_ToSpeak, this.OnNetToSpeak)
	this.msgHandlerMap.RegisterHandler(usercmd.MsgTypeCmd_DoNothing, this.OnNetDoNothing)
}

func (this *ScenePlayerNetMsgHelper) OnNetDoNothing(data []byte, flag byte) {
	op, ok := common.DecodeCmd(data, flag, &usercmd.MsgDoNothing{}).(*usercmd.MsgDoNothing)
	if !ok {
		glog.Error("DecodeCmd error: OnNetDoNothing")
		return
	}
	glog.Info("player ", op.Id, " Says ", op.Hello)
}

//收到玩家消息
func (this *ScenePlayerNetMsgHelper) OnRecvPlayerCmd(
	cmd usercmd.MsgTypeCmd, data []byte, flag byte) {
	this.msgHandlerMap.Call(cmd, data, flag)
}

//释放技能
func (this *ScenePlayerNetMsgHelper) OnCastSkill(data []byte, flag byte) {
	op, ok := common.DecodeCmd(data, flag, &usercmd.MsgCastSkill{}).(*usercmd.MsgCastSkill)
	if !ok {
		glog.Error("DecodeCmd error: OnCastSkill")
		return
	}
	this.selfPlayer.CastSkill(op)
}

func (this *ScenePlayerNetMsgHelper) OnNetMove(data []byte, flag byte) {
	op, ok := common.DecodeCmd(data, flag, &usercmd.MsgMove{}).(*usercmd.MsgMove)
	if !ok {
		glog.Error("DecodeCmd error: OnNetMove")
		return
	}

	if power, angle, face, ok := this.selfPlayer.CheckMoveMsg(float64(op.Power), float64(op.Angle), op.Face); ok {
		this.selfPlayer.Move(power, angle, face)
	}
}

func (this *ScenePlayerNetMsgHelper) OnNetReLife(data []byte, flag byte) {
	this.selfPlayer.Relife()
}

func (this *ScenePlayerNetMsgHelper) OnNetToSpeak(data []byte, flag byte) {
	dban, tban := redismgr.GetMe().IsBanVoice(this.selfPlayer.ID)
	var banVoice bool = false
	if dban || tban {
		banVoice = true
	}
	this.selfPlayer.room.ReqToSpeak(this.selfPlayer.ID, banVoice, this.selfPlayer.room.GetNewLoginUsers())
}

func (this *ScenePlayerNetMsgHelper) OnCloseSocket(data []byte, flag byte) {
	this.selfPlayer.ClientCloseSocket()
}

// 奔跑
func (this *ScenePlayerNetMsgHelper) OnRun(data []byte, flag byte) {
	op := &usercmd.MsgRun{}
	if common.DecodeGoCmd(data, flag, op) != nil {
		glog.Error("DecodeCmd error:OnSpeedUp ", this.selfPlayer.Name)
		return
	}
	this.selfPlayer.Run(op)
}

// 聊天
func (this *ScenePlayerNetMsgHelper) OnSceneChat(data []byte, flag byte) {
	op, ok := common.DecodeCmd(data, flag, &usercmd.MsgSceneChat{}).(*usercmd.MsgSceneChat)
	if ok {
		glog.Info("chat ", this.selfPlayer.Name, " size:", len(data), " flag:", flag, " mes:", op)
	}
	this.selfPlayer.AsyncRoundMsg(data, flag) //直接转发
}

//队伍通知
func (this *ScenePlayerNetMsgHelper) OnTeamNotice(data []byte, flag byte) {
	op, _ := common.DecodeCmd(data, flag, &usercmd.TeamNoticeMsg{}).(*usercmd.TeamNoticeMsg)
	team := this.selfPlayer.room.GetTeam(this.selfPlayer.udata.TeamId)
	if team != nil {
		team.NoticeTime = time.Now().Unix()
		this.selfPlayer.room.BroadcastTeamMsg(this.selfPlayer.udata.TeamId, usercmd.MsgTypeCmd_TeamNotice, op)
	}
}
