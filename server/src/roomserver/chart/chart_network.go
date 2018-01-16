package chart

import (
	"base/env"
	"fmt"
	"github.com/fananchong/gochart"
	"github.com/shirou/gopsutil/net"
	"strconv"
)

type ChartNetwork struct {
	gochart.ChartTime
	presend uint64
	prerecv uint64
}

func NewChartNetwork() *ChartNetwork {
	this := &ChartNetwork{}
	this.RefreshTime = DEFAULT_REFRESH_TIME
	this.SampleNum = DEFAULT_SAMPLE_NUM
	this.ChartType = "line"
	this.Title = "网络带宽"
	this.SubTitle = env.Get("room", "local")
	this.YAxisText = "net"
	this.YMax = "1000"
	this.ValueSuffix = "Mbps"
	return this
}

func (this *ChartNetwork) Update(now int64) map[string][]interface{} {
	datas := make(map[string][]interface{})
	nv, _ := net.IOCounters(false)
	if this.presend == 0 {
		this.presend = nv[0].BytesSent
	}
	if this.prerecv == 0 {
		this.prerecv = nv[0].BytesRecv
	}
	v1, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(nv[0].BytesSent-this.presend)*8/float64(1024*1024)), 64)
	datas["Sent"] = []interface{}{v1}
	v2, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(nv[0].BytesRecv-this.prerecv)*8/float64(1024*1024)), 64)
	datas["Recv"] = []interface{}{v2}
	this.presend = nv[0].BytesSent
	this.prerecv = nv[0].BytesRecv
	return datas
}
