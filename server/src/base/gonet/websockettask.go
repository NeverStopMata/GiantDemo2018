package gonet

import (
	"base/glog"
	"container/list"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type IWebSocketTask interface {
	ParseMsg(data []byte, flag byte) bool
	OnClose()
}

type WebSocketTask struct {
	closed      int32
	verified    bool
	stopedChan  chan bool
	sendmsglist *list.List
	sendMutex   sync.Mutex
	Conn        *websocket.Conn
	Derived     IWebSocketTask
	msgchan     chan []byte
	signal      chan int
}

func NewWebSocketTask(conn *websocket.Conn) *WebSocketTask {
	return &WebSocketTask{
		closed:      -1,
		verified:    false,
		Conn:        conn,
		stopedChan:  make(chan bool, 1),
		sendmsglist: list.New(),
		msgchan:     make(chan []byte, 1024),
		signal:      make(chan int, 1),
	}
}

func (this *WebSocketTask) Signal() {
	select {
	case this.signal <- 1:
	default:
	}
}

func (this *WebSocketTask) Stop() {
	if !this.IsClosed() && len(this.stopedChan) == 0 {
		this.stopedChan <- true
	} else {
		glog.Info("[连接] 连接关闭失败")
	}
}

func (this *WebSocketTask) Start() {
	if atomic.CompareAndSwapInt32(&this.closed, -1, 0) {
		glog.Info("[连接] 收到连接 ", this.Conn.RemoteAddr())
		go this.sendloop()
		go this.recvloop()
	}
}

func (this *WebSocketTask) Close() {
	if atomic.CompareAndSwapInt32(&this.closed, 0, 1) {
		glog.Info("[连接] 断开连接 ", this.Conn.RemoteAddr())
		this.Conn.Close()
		close(this.stopedChan)
		this.Derived.OnClose()
	}
}

func (this *WebSocketTask) Reset() {
	if atomic.LoadInt32(&this.closed) == 1 {
		glog.Info("[连接] 重置连接 ", this.Conn.RemoteAddr())
		this.closed = -1
		this.verified = false
		this.stopedChan = make(chan bool)
	}
}

func (this *WebSocketTask) IsClosed() bool {
	return atomic.LoadInt32(&this.closed) != 0
}

func (this *WebSocketTask) Verify() {
	this.verified = true
}

func (this *WebSocketTask) IsVerified() bool {
	return this.verified
}

func (this *WebSocketTask) Terminate() {
	this.Close()
}

func (this *WebSocketTask) AsyncSend(buffer []byte, flag byte) bool {
	if this.IsClosed() {
		return false
	}

	bsize := len(buffer)

	//glog.Info("[AsynSend] raw", buffer, bsize)

	totalsize := bsize + 4
	sendbuffer := make([]byte, 0, totalsize)
	sendbuffer = append(sendbuffer, byte(bsize), byte(bsize>>8), byte(bsize>>16), flag)
	sendbuffer = append(sendbuffer, buffer...)
	this.msgchan <- sendbuffer

	//glog.Info("[AsynSend] final", sendbuffer, bsize)

	return true
}

func (this *WebSocketTask) recvloop() {
	defer func() {
		if err := recover(); err != nil {
			glog.Error("[异常] ", err, "\n", string(debug.Stack()))
		}
	}()
	defer this.Close()

	var (
		datasize int
	)

	for {

		_, bytemsg, err := this.Conn.ReadMessage()
		if nil != err {
			glog.Error("[WebSocket] 接收失败 ", this.Conn.RemoteAddr(), ",", err)
			return
		}

		datasize = int(bytemsg[0]) | int(bytemsg[1])<<8 | int(bytemsg[2])<<16
		if datasize > cmd_max_size {
			glog.Error("[WebSocket] 数据超过最大值 ", this.Conn.RemoteAddr(), ",", datasize)
			return
		}

		//glog.Error("[WebSocketServer] the whole msg recieved", bytemsg)
		this.Derived.ParseMsg(bytemsg[cmd_header_size:], bytemsg[3])
	}
}

func (this *WebSocketTask) sendloop() {
	defer func() {
		if err := recover(); err != nil {
			glog.Error("[异常] ", err, "\n", string(debug.Stack()))
		}
	}()
	defer this.Close()

	var (
		timeout = time.NewTimer(time.Second * cmd_verify_time)
	)

	defer timeout.Stop()

	for {
		select {
		case bytemsg := <-this.msgchan:
			if nil != bytemsg && len(bytemsg) > 0 {
				err := this.Conn.WriteMessage(websocket.BinaryMessage, bytemsg)
				if nil != err {
					glog.Error("[WebSocket] 发送失败 ", this.Conn.RemoteAddr(), ",", err)
					return
				}
				//glog.Info("[连接] the bytemsg send to client final", bytemsg)
			} else {
				glog.Error("[WebSocket] byte msg in the send chan is wrong !", bytemsg)
				return
			}
		case <-this.stopedChan:
			return
		case <-timeout.C:
			if !this.IsVerified() {
				glog.Error("[WebSocket] 验证超时 ", this.Conn.RemoteAddr())
				return
			}
		}
	}
}
