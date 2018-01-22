package chart

// 用于统计在线人数、机器人数、房间数

import (
	"base/env"
	"github.com/fananchong/gochart"
)

type ChartOnline struct {
	gochart.ChartTime
	funcPlayerNum      func() int
	funcRoomNum        func() int
	funcScenePlayerNum func() int
}

func NewChartOnline(funcPlayerNum, funcRoomNum, funcScenePlayerNum func() int) *ChartOnline {
	this := &ChartOnline{funcPlayerNum: funcPlayerNum, funcRoomNum: funcRoomNum, funcScenePlayerNum: funcScenePlayerNum}
	this.RefreshTime = DEFAULT_REFRESH_TIME
	this.SampleNum = DEFAULT_SAMPLE_NUM
	this.ChartType = "line"
	this.Title = "在线统计"
	this.SubTitle = env.Get("room", "local")
	this.YAxisText = "Num"
	this.YMax = "5000"
	return this
}

func (this *ChartOnline) Update(now int64) map[string][]interface{} {
	datas := make(map[string][]interface{})
	playernum := this.funcPlayerNum()
	roomnum := this.funcRoomNum()
	robotnum := this.funcScenePlayerNum() - playernum
	if robotnum < 0 {
		robotnum = 0
	}
	datas["player"] = []interface{}{playernum}
	datas["robot"] = []interface{}{robotnum}
	datas["room"] = []interface{}{roomnum}
	return datas
}
