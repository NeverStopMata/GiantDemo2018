package redis

import (
	"base/glog"
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

const (
	defaultLinkTimeout = time.Second * 5
)

type NetStatus int

const (
	statusReady NetStatus = iota
	statusLinking
	statusOk
	statusClose
)
const (
	recvPacket = 0
	sendPacket = 0
	closeConn  = 1
)

var (
	ErrInvalidAddress = errors.New("tcp address is invalid")
)

type Conn struct {
	conn      *net.TCPConn
	address   string
	isclosed  int32
	signal    chan int
	queue     *CmdQueue
	waitQueue chan ICommond
}

func NewConn(address string, autoConnect bool) (conn *Conn, err error) {
	conn = &Conn{address: address}
	conn.signal = make(chan int, 1)
	conn.queue = NewCmdQueue(10000)
	conn.waitQueue = make(chan ICommond, 100)
	if autoConnect {
		err = conn.Connect()
		for err != nil {
			glog.Error("connect fail, try reconnect...")
			time.Sleep(time.Second)
			err = conn.Connect()
		}
	}
	return
}

func (c *Conn) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, defaultLinkTimeout)
	if err == nil {
		c.conn = conn.(*net.TCPConn)
		go c.loop()
	}
	return err
}

func (c *Conn) ReConnect() (err error) {
	err = c.Connect()
	for err != nil {
		glog.Error("connect fail, try reconnect...")
		time.Sleep(time.Second)
		err = c.Connect()
	}
	return
}

func (c *Conn) Send(cmd ICommond) (err error) {
	if c.IsClose() {
		cmd.SetData(nil, ErrRedisConClose)
		return ErrRedisConClose
	}
	if err = c.queue.Push(cmd); err != nil {
		// TODO when cmd handle is full, you can drop or wait
		//		time.Sleep(time.Second)
		//		c.Send(cmd)
		//fmt.Println("send cmd.GetBytes:", len(cmd.GetBytes()))
		//c.waitQueue <- cmd
		cmd.SetData(nil, err)
	} else {
		c.signalSend()
	}
	return
}

func (c *Conn) signalSend() {
	select {
	case c.signal <- sendPacket:
	default:
	}
}

func (c *Conn) IsClose() bool {
	return atomic.LoadInt32(&c.isclosed) == 1
}

func (c *Conn) Close() {

}

func (c *Conn) beClosed() {
	if atomic.CompareAndSwapInt32(&c.isclosed, 0, 1) {
		c.conn.Close()
		fmt.Println("[Conn] conn be close start")
		glog.Info("[Conn] conn be close start")
		select {
		case c.signal <- closeConn:
		default:
		}
		cmd := c.queue.POP()
		for cmd != nil {
			cmd.Done()
			c.queue.Adv()
			cmd = c.queue.POP()
		}
		fmt.Println("[Conn] conn be close")
		glog.Info("[Conn] conn be close")
	}
}

func (c *Conn) loop() {
	go c.onRecv()
	var buf = NewPacketSize(8192)
	var err error
	var cmd ICommond
	for {
		if c.IsClose() {
			return
		}
		if cmd == nil {
			cmd = c.queue.Get()
		}
		if cmd == nil {
			r := <-c.signal
			if r != sendPacket {
				c.beClosed()
				return
			}
		}
		if cmd == nil {
			cmd = c.queue.Get()
		}
		//		fmt.Println("send queue:", c.queue.String())
		for buf.w < 8192 && cmd != nil {
			err = buf.Bytes(cmd.GetBytes())
			if err != nil {
				fmt.Println("get bytes err:", err)
				fmt.Println("send queue:", c.queue.String())
			}
			cmd = c.queue.Get()
		}
		//		fmt.Println("send queue ok:", c.queue.String())
		if err = buf.Flush(c.conn); err != nil {
			c.beClosed()
			return
		}
	}
}

func (c *Conn) onRecv() {
	var (
		buf  = NewPacketSize(8 * 1024)
		cmd  ICommond
		err  error
		data interface{}
	)
	for {
		data, err = c.readReply(buf)
		if err != nil {
			glog.Error(err)
			c.beClosed()
			return
		}
		cmd = c.queue.POP()
		if cmd == nil {
			//fmt.Println("recv nil cmd from redis:")
			//fmt.Println("recv cmd from redis:", c.queue.String())
			continue
		}
		if cmd.SetData(data, err) {
			//			fmt.Println("recv cmd from redis ok:", c.queue.String())
			//			fmt.Println("recv cmd from redis ok 2:", c.queue.String())
			//if len(c.waitQueue) > 0 {
			//	//fmt.Println("recv from wait:", len(c.waitQueue))
			//	//cmd2 := <-c.waitQueue
			//	//fmt.Println("recv from wait ok:", len(c.waitQueue))
			//	//fmt.Println("cmd2.GetBytes:", len(cmd2.GetBytes()))
			//	//c.queue.ReplaceLast(cmd2)
			//	//fmt.Println("Replace :", c.queue.String())
			//	c.signalSend()
			//} else {
			//	c.queue.Adv()
			//}
			c.queue.Adv()
			cmd.Done()
		}
	}
	fmt.Println("over")
}

func (c *Conn) readReply(buf *Packet) (data interface{}, err error) {
	var line []byte
	var n int
	if line, err = buf.ReadLine(c.conn); err != nil {
		glog.Error("read error:", err)
		return
	}
	if line == nil || len(line) == 0 {
		return nil, errors.New("commond replay is nil")
	}
	switch line[0] {
	case '+':
		switch {
		case len(line) == 3 && line[1] == 'O' && line[2] == 'K':
			data = "OK"
		case len(line) == 5 && line[1] == 'P' && line[2] == 'O' && line[3] == 'N' && line[4] == 'G':
			data = "PONG"
		default:
			data = line[1:]
		}
	case '-':
		err = errors.New(string(line[1:]))
	case ':':
		return parseToInt(line[1:])
	case '$':
		n, err = parseToLen(line[1:])
		if n < 0 || err != nil {
			return
		}
		p := make([]byte, n+2)
		if err = buf.ReadBytes(c.conn, p, true); err != nil {
			return
		}
		if p[n] != '\r' || p[n+1] != '\n' {
			return nil, errors.New("replay format is error")
		}
		data = p[:n]
	case '*':
		n, err = parseToLen(line[1:])
		if n < 0 || err != nil {
			return
		}
		r := make([]interface{}, n)
		for i := range r {
			r[i], err = c.readReply(buf)
			if err != nil {
				return nil, err
			}
		}
		data = r
	}
	return
}

// parseInt parses an integer reply.
func parseToInt(p []byte) (int64, error) {
	if len(p) == 0 {
		return 0, protocolError("malformed integer")
	}

	var negate bool
	if p[0] == '-' {
		negate = true
		p = p[1:]
		if len(p) == 0 {
			return 0, protocolError("malformed integer")
		}
	}

	var n int64
	for _, b := range p {
		n *= 10
		if b < '0' || b > '9' {
			return 0, protocolError("illegal bytes in length")
		}
		n += int64(b - '0')
	}

	if negate {
		n = -n
	}
	return n, nil
}

// parseLen parses bulk string and array lengths.
func parseToLen(p []byte) (int, error) {
	if len(p) == 0 {
		return -1, protocolError("malformed length")
	}

	if p[0] == '-' && len(p) == 2 && p[1] == '1' {
		// handle $-1 and $-1 null replies.
		return -1, nil
	}

	var n int
	for _, b := range p {
		n *= 10
		if b < '0' || b > '9' {
			return -1, protocolError("illegal bytes in length")
		}
		n += int(b - '0')
	}

	return n, nil
}
