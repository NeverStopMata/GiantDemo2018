package chart

import (
	"base/env"
	"github.com/fananchong/gochart"
	"github.com/shirou/gopsutil/cpu"
	"strconv"
)

type ChartCPU struct {
	gochart.ChartTime
}

func NewChartCPU() *ChartCPU {
	this := &ChartCPU{}
	this.RefreshTime = DEFAULT_REFRESH_TIME
	this.SampleNum = DEFAULT_SAMPLE_NUM
	this.ChartType = "line"
	this.Title = "CPU占用"
	this.SubTitle = env.Get("room", "local")
	this.YAxisText = "cpu"
	this.YMax = "100"
	this.ValueSuffix = "%"
	return this
}

func (this *ChartCPU) Update(now int64) map[string][]interface{} {
	datas := make(map[string][]interface{})
	cc, _ := cpu.Percent(0, false)
	for i := 0; i < len(cc); i++ {
		datas["cpu"+strconv.Itoa(i)] = []interface{}{int(cc[i])}
	}
	return datas
}
