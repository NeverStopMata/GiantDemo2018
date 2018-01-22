package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

var gConn *net.TCPConn

func TcpClient() {
	addr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:8000")
	conn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		panic(err)
	}

	fmt.Println("connect to ", conn.RemoteAddr().String())

	gConn = conn
	go handler()
}

func handler() {
	for {
		tempbuf := make([]byte, 1024)
		readnum, err := io.ReadAtLeast(gConn, tempbuf, 1)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(tempbuf[:readnum]))

		go readline()
	}
}

func main() {
	TcpClient()

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
		gConn.Write(b)
	}
}
