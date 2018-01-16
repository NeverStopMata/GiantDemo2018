package main

import (
	"fmt"
	"io"
	"net"
)

func TcpServer() {
	addr, _ := net.ResolveTCPAddr("tcp", ":8000")
	lis, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := lis.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("on connect. addr =", conn.RemoteAddr())

		go handler(conn)
	}
}

func handler(conn *net.TCPConn) {
	for {
		tempbuf := make([]byte, 1024)
		readnum, err := io.ReadAtLeast(conn, tempbuf, 1)
		if err != nil {
			return
		}
		fmt.Println("recv data:", string(tempbuf[:readnum]))
		_, err = conn.Write(tempbuf[:readnum])
		if err != nil {
			return
		}
	}
}

func main() {
	TcpServer()
}
