package main

import (
	"base/gonet"
	"fmt"
	"net"
)

type EchoTask struct {
	gonet.TcpTask
	roomId int // 房间ID
}

func NewEchoTask(conn net.Conn) *EchoTask {
	s := &EchoTask{
		TcpTask: *gonet.NewTcpTask(conn),
	}
	s.Derived = s
	return s
}

func (this *EchoTask) ParseMsg(data []byte, flag byte) bool {

	if this.IsVerified() == false {
		this.Verify()
	}

	fmt.Println("recv data:", string(data))

	this.AsyncSend(data, flag)

	// 某人聊天
	chanChat <- &ChatInfo{this, data}

	return true
}

func (this *EchoTask) OnClose() {
	// 某人离开聊天
	chanRemove <- this
}

type EchoServer struct {
	gonet.Service
	tcpser *gonet.TcpServer
}

var serverm *EchoServer

func EchoServer_GetMe() *EchoServer {
	if serverm == nil {
		serverm = &EchoServer{
			tcpser: &gonet.TcpServer{},
		}
		serverm.Derived = serverm
	}
	return serverm
}

func (this *EchoServer) Init() bool {
	err := this.tcpser.Bind(":8000")
	if err != nil {
		fmt.Println("绑定端口失败")
		return false
	}

	// 开始聊天协程
	go doChat()

	return true
}

func (this *EchoServer) MainLoop() {
	conn, err := this.tcpser.Accept()
	if err != nil {
		return
	}
	task := NewEchoTask(conn)
	task.Start()

	// 某人加入聊天
	chanAdd <- task
}

func (this *EchoServer) Reload() {

}

func (this *EchoServer) Final() bool {
	this.tcpser.Close()
	return true
}

func main() {

	EchoServer_GetMe().Main()

}
