package main

import (
	"base/gonet"
	"fmt"
	"time"
)

type Client struct {
	gonet.TcpTask
	mclient *gonet.TcpClient
}

func NewClient() *Client {
	s := &Client{
		TcpTask: *gonet.NewTcpTask(nil),
	}
	s.Derived = s
	return s
}

func (this *Client) Connect(addr string) bool {

	conn, err := this.mclient.Connect(addr)
	if err != nil {
		fmt.Println("连接失败 ", addr)
		return false
	}

	this.Conn = conn

	this.Start()

	fmt.Println("连接成功 ", addr)
	return true
}

func (this *Client) ParseMsg(data []byte, flag byte) bool {

	if this.IsVerified() == false {
		this.Verify()
	}

	fmt.Println(string(data))

	go readline()

	return true
}

func (this *Client) OnClose() {

}

var gClient *Client

func main() {

	gClient = NewClient()

	if !gClient.Connect("127.0.0.1:8000") {
		return
	}

	go readline()

	for {
		time.Sleep(100 * time.Second)
	}
}

func readline() {
	fmt.Print("> ")
	var c byte
	var err error
	var b []byte
	for {
		_, err = fmt.Scanf("%c", &c)
		if err == nil && c != '\n' {
			b = append(b, c)
		} else {
			break
		}
	}
	if len(b) > 0 {
		gClient.AsyncSend(b, 0)
	}
}
