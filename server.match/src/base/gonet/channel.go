package gonet

import (
	"base/util"
	"net"
	"strings"
	"sync"
)

type IChannel interface {
	OnAdd()
	OnLeave(node *ChannelNode)
}

type ChannelNode struct {
	TcpTask
	Id   uint64
	Addr string
}

type Channel struct {
	Id        uint64
	sendMutex sync.RWMutex
	litener   *TcpServer
	onlines   []*ChannelNode
	handle    IChannel
}

func NewChannelNode(conn *net.TCPConn, addr string) (node *ChannelNode) {
	node = new(ChannelNode)
	node.Addr = addr
	node.TcpTask.Conn = conn
	node.Id = addrToId(addr)
	return
}

func NewChannel() (c *Channel) {
	return &Channel{
		litener: &TcpServer{},
	}
}

func (this *Channel) Bind(addr string) error {
	this.Id = addrToId(addr)
	err := this.litener.Bind(addr)
	if err != nil {
		return err
	}
	go this.Loop()
	return nil
}

func (this *Channel) Loop() {
	for {
		conn, err := this.litener.Accept()
		if err != nil {
			this.Add(NewChannelNode(conn, ""))
		}
	}
}

func (this *Channel) Join(addr string) error {
	conn, err := TcpDial(addr)
	if err != nil {
		return err
	}
	this.Add(NewChannelNode(conn, addr))
	return nil
}

func (this *Channel) Add(node *ChannelNode) {
	this.sendMutex.Lock()
	this.onlines = append(this.onlines, node)
	this.sendMutex.Unlock()
}

func (this *Channel) Remove(node *ChannelNode) {
	this.sendMutex.Lock()
	defer this.sendMutex.Unlock()
	for idx, n := range this.onlines {
		if n.Id == node.Id {
			this.onlines = append(this.onlines[:idx], this.onlines[idx+1:]...)
			return
		}
	}
}

func (this *Channel) SendMsg(msg []byte, flag byte) {
	this.sendMutex.RLock()
	defer this.sendMutex.RUnlock()
	for _, node := range this.onlines {
		node.AsyncSend(msg, flag)
	}
}

func addrToId(addr string) uint64 {
	var (
		ip   uint32
		port uint32
	)
	idx := strings.Index(addr, ":")
	if idx == 0 {
		ip = util.IPStrToUInt(util.GetTopInnerIP().String())
	} else {
		ip = util.IPStrToUInt(addr[:idx])
	}
	port = util.Dtoi(addr[idx+1:])

	return uint64(ip)<<16 + uint64(port)
}
