package chart

// 用于统计数据，并使用图表的方式呈现

import (
	"base/env"
	"base/glog"
	"runtime/debug"

	"github.com/fananchong/gochart"
)

type ChartMgr struct {
	chartSvr *gochart.ChartServer
}

var chartMgr *ChartMgr

func ChartMgr_GetMe() *ChartMgr {
	if chartMgr == nil {
		chartMgr = &ChartMgr{}
	}
	return chartMgr
}

func (this *ChartMgr) Init(funcPlayerNum, funcRoomNum, funcScenePlayerNum func() int) {
	addr := env.Get("room", "chartAddr")
	if addr != "" {
		s := &gochart.ChartServer{}
		this.chartSvr = s

		s.AddChart("online", NewChartOnline(funcPlayerNum, funcRoomNum, funcScenePlayerNum), true)
		s.AddChart("cpu", NewChartCPU(), true)
		s.AddChart("mem", NewChartMemory(), true)
		s.AddChart("net", NewChartNetwork(), true)

		go func() {
			defer func() {
				if err := recover(); err != nil {
					glog.Error("[异常] 报错 ", err, "\n", string(debug.Stack()))
				}
			}()
			glog.Infoln(s.ListenAndServe(addr).Error())
		}()
	}
}

func (this *ChartMgr) IsEnabled() bool {
	addr := env.Get("room", "chartAddr")
	enabled := false
	if addr != "" {
		enabled = true
	}
	return enabled
}

func (this *ChartMgr) AddChart(chartName string, chartObj gochart.IChartInner) {
	this.chartSvr.AddChart(chartName, chartObj, true)
}

const (
	DEFAULT_REFRESH_TIME = 5
	DEFAULT_SAMPLE_NUM   = 3600 / DEFAULT_REFRESH_TIME
)
