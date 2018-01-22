package gonet

import (
	"base/glog"
	"net"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type TcpClient struct {
}

func (this *TcpClient) Connect(address string) (*net.TCPConn, error) {
	return TcpDial(address)
}

func TcpDial(address string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		glog.Error("[连接] 解析失败 ", address)
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		glog.Error("[连接] 连接失败 ", address)
		return nil, err
	}

	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(1 * time.Minute)
	conn.SetNoDelay(true)
	conn.SetWriteBuffer(128 * 1024)
	conn.SetReadBuffer(128 * 1024)

	glog.Info("[连接] 连接成功 ", address)
	return conn, nil
}

type WebClient struct {
}

func (this *WebClient) Connect(address string) (*websocket.Conn, error) {

	u := url.URL{Scheme: "ws", Host: address, Path: "/"}
	glog.Info("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		glog.Error("connecting to %s", u.String())
		return nil, err
	}
	glog.Info("[连接] 连接成功 ", address)
	return conn, nil
}
