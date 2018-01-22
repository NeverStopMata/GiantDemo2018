package udp

import (
	"base/glog"
	"base/gonet"
	"common"
	"net"
	"roomserver/chart"
	"runtime/debug"
	"sync"
	"usercmd"
)

const (
	MSG_HEADER_SIZE = 4
	MAX_PACKET_SIZE = 548
)

type UdpSess struct {
	gonet.UdpTask
	room_key    string
	uid         uint64
	udpSendBuff *gonet.ByteBuffer
}

type UdpSessAdder interface {
	AddUdpSess(*usercmd.MsgBindTCPSession) *UdpSess
}

var PlayerTaksMgr UdpSessAdder

func NewUdpSess(id uint64, key string) *UdpSess {
	s := &UdpSess{
		UdpTask:     *gonet.NewUdpTask(),
		uid:         id,
		room_key:    key,
		udpSendBuff: gonet.NewByteBuffer(),
	}
	s.Derived = s
	return s
}

func (this *UdpSess) ParseMsg(data []byte, flag int) {
	cmd := usercmd.MsgTypeCmd(common.GetCmd(data))
	glog.Info("[UdpSess] 收到消息 cmd:", cmd, ",len:", len(data))
	return
}

func (this *UdpSess) OnClose() {

}

func (this *UdpSess) SendUDPCmd(cmd usercmd.MsgTypeCmd, msg common.Message) bool {
	data, flag, err := common.EncodeGoCmdWithCompressSize(uint16(cmd), msg, MAX_PACKET_SIZE)
	if err != nil {
		glog.Info("[玩家] 发送失败 cmd:", cmd, ",len:", len(data), ",err:", err)
		return false
	}
	bsize := len(data)
	if compressChart != nil && flag == 1 {
		preSize := msg.Size()
		compressChart.AddCompressInfo(preSize, bsize)
	}

	// 放弃发送大包，由上层采用TCP发送
	if bsize+4 > MAX_PACKET_SIZE {
		return false
	}

	this.udpSendBuff.Reset()
	this.udpSendBuff.Append(byte(bsize), byte(bsize>>8), byte(bsize>>16), flag)
	this.udpSendBuff.Append(data...)
	this.AsyncSend(this.udpSendBuff.RdBuf()[:this.udpSendBuff.RdSize()], -1)
	return true
}

////////////////////////////////////UDP 服务管理器////////////////////////////////////
type UdpMgr struct {
	udpSvr *gonet.UdpServer
}

func handshake(data []byte, listener *net.UDPConn, addr *net.UDPAddr) *gonet.UdpTask {
	defer func() {
		if err := recover(); err != nil {
			glog.Error("[异常] ", err, "\n", string(debug.Stack()))
		}
	}()

	rsize := len(data)
	glog.Info("[UDPSession] 获取数据长度： ", rsize)
	if rsize <= MSG_HEADER_SIZE {
		return nil
	}

	buff := data[MSG_HEADER_SIZE:rsize]
	cmd := usercmd.MsgTypeCmd(common.GetCmd(buff))

	if cmd != usercmd.MsgTypeCmd_BindTCPSession {
		glog.Error("[UDPSession] 目前仅支持BindTCPSession协议 ", cmd)
		return nil
	}
	revCmd, ok := common.DecodeCmd(buff, 0, &usercmd.MsgBindTCPSession{}).(*usercmd.MsgBindTCPSession)
	if !ok {
		glog.Error("[UDPSession] BindTCPSession解码错误 ")
		return nil
	}

	udpsess := PlayerTaksMgr.AddUdpSess(revCmd)
	if udpsess == nil {
		return nil
	}
	return &udpsess.UdpTask
}

var (
	pudpm         *UdpMgr
	udpmgr_once   sync.Once
	compressChart *chart.ChartCompressRatio
)

func UdpMgr_GetMe() *UdpMgr {
	if pudpm == nil {
		udpmgr_once.Do(func() {
			pudpm = &UdpMgr{
				udpSvr: gonet.NewUdpServer(handshake, MAX_PACKET_SIZE),
			}

			if chart.ChartMgr_GetMe().IsEnabled() {
				compressChart = chart.NewChartCompressRatio()
				chart.ChartMgr_GetMe().AddChart("compress", compressChart)
			}
		})
	}
	return pudpm
}

func (this *UdpMgr) RunServer(bindAddr string) {
	err := this.udpSvr.BindAccept(bindAddr)
	if err != nil {
		glog.Error("[UDPMgr] BindAccept failed", err.Error())
	}
}
