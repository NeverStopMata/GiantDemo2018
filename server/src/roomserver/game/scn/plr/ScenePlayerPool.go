package plr

// 玩家用的对象池

import (
	"usercmd"
)

type ScenePlayerPool struct {
	MsgBallMove usercmd.BallMove
	MsgEats     []*usercmd.BallEat
	MsgHits     []*usercmd.HitMsg
	eatmsg      []*usercmd.BallEat
	hitmsg      []*usercmd.HitMsg
}

func (this *ScenePlayerPool) Init() {
	for index := 0; index < 20; index++ {
		this.hitmsg = append(this.hitmsg, &usercmd.HitMsg{})
		this.eatmsg = append(this.eatmsg, &usercmd.BallEat{})
	}
}

func (this *ScenePlayerPool) ResetMsg() {
	this.eatmsg = append(this.eatmsg, this.MsgEats...)
	this.MsgEats = make([]*usercmd.BallEat, 0)
	this.hitmsg = append(this.hitmsg, this.MsgHits...)
	this.MsgHits = make([]*usercmd.HitMsg, 0)
}

func (this *ScenePlayerPool) AddEatMsg(ballid, beEat uint32) {
	if len(this.eatmsg) < 0 {
		this.eatmsg = append(this.eatmsg, &usercmd.BallEat{})
	}
	msg := this.eatmsg[len(this.eatmsg)-1]
	this.eatmsg = append(this.eatmsg[0 : len(this.eatmsg)-1])
	msg.Source = ballid
	msg.Target = beEat
	this.MsgEats = append(this.MsgEats, msg)
}

func (this *ScenePlayerPool) AddHitMsg(ballid, beEat uint32, addhp int32, curhp uint32, scene IScene) {
	if len(this.hitmsg) < 0 {
		this.hitmsg = append(this.hitmsg, &usercmd.HitMsg{})
	}
	msg := this.hitmsg[len(this.hitmsg)-1]
	this.hitmsg = append(this.hitmsg[0 : len(this.hitmsg)-1])
	msg.Source = ballid
	msg.Target = beEat
	msg.AddHp = addhp
	msg.CurHp = curhp
	this.MsgHits = append(this.MsgHits, msg)
}
