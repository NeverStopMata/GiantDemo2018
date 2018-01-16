package chart

// 用于统计zlib压缩率情况

import (
	"base/env"
	"sync"

	"github.com/fananchong/gochart"
)

type ChartCompressRatio struct {
	gochart.ChartTime
	beforeSize []interface{}
	afterSize  []interface{}
	m          sync.Mutex
}

func NewChartCompressRatio() *ChartCompressRatio {
	this := &ChartCompressRatio{
		beforeSize: make([]interface{}, 0),
		afterSize:  make([]interface{}, 0),
	}
	this.RefreshTime = DEFAULT_REFRESH_TIME
	this.SampleNum = DEFAULT_SAMPLE_NUM
	this.ChartType = "line"
	this.Title = "消息协议压缩统计"
	this.SubTitle = env.Get("room", "local")
	this.YAxisText = "Byte"
	this.YMax = "1000"
	return this
}

func (this *ChartCompressRatio) Update(now int64) map[string][]interface{} {
	this.m.Lock()
	defer this.m.Unlock()

	datas := make(map[string][]interface{})
	datas["before"] = make([]interface{}, 0)
	datas["after"] = make([]interface{}, 0)
	datas["before"] = append(datas["before"], this.beforeSize...)
	datas["after"] = append(datas["after"], this.afterSize...)

	this.beforeSize = make([]interface{}, 0)
	this.afterSize = make([]interface{}, 0)

	return datas
}

func (this *ChartCompressRatio) AddCompressInfo(before, after int) {
	this.m.Lock()
	defer this.m.Unlock()
	this.beforeSize = append(this.beforeSize, before)
	this.afterSize = append(this.afterSize, after)
}
