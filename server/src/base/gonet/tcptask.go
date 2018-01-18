package gonet

import (
	"base/glog"
	"io"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

type ITcpTask interface {
	ParseMsg(data []byte, flag byte) bool
	OnClose()
}

const (
	cmd_max_size     = 128 * 1024
	sendcmd_max_size = 64 * 1024
	cmd_header_size  = 4 // 3字节指令长度 1字节是否压缩
	cmd_verify_time  = 30
)

type TcpTask struct {
	closed            int32
	verified          bool
	stopedChan        chan bool
	sendBuff          *ByteBuffer
	sendMutex         sync.Mutex
	sendBuffSizeLimit int // -1 没有限制
	Conn              net.Conn
	Derived           ITcpTask
	signal            chan int
}

func NewTcpTask(conn net.Conn) *TcpTask {
	return &TcpTask{
		closed:            -1,
		verified:          false,
		Conn:              conn,
		stopedChan:        make(chan bool, 1),
		sendBuff:          NewByteBuffer(),
		signal:            make(chan int, 1),
		sendBuffSizeLimit: -1, //缺省不做限制（注意对外客户端链接的，需要做下限制，防止内存被给耗尽！）
	}
}

func (this *TcpTask) SetSendBuffSizeLimt(limit int) {
	this.sendBuffSizeLimit = limit
}

func (this *TcpTask) Signal() {
	select {
	case this.signal <- 1:
	default:
	}
}

func (this *TcpTask) Stop() bool {
	if this.IsClosed() {
		glog.Info("[连接] 关闭失败 ", this.Conn.RemoteAddr())
		return false
	}
	select {
	case this.stopedChan <- true:
	default:
		glog.Info("[连接] 关闭失败 ", this.Conn.RemoteAddr())
		return false
	}
	return true
}

func (this *TcpTask) Start() {
	if !atomic.CompareAndSwapInt32(&this.closed, -1, 0) {
		glog.Error(("一开始就错了"))
		return
	}
	glog.Info("[连接] 收到连接 ", this.Conn.RemoteAddr())
	go this.sendloop()
	go this.recvloop()
}

func (this *TcpTask) Close() {
	if !atomic.CompareAndSwapInt32(&this.closed, 0, 1) {
		return
	}
	glog.Info("[连接] 断开连接 ", this.Conn.RemoteAddr())
	this.Conn.Close()
	close(this.stopedChan)
	this.Derived.OnClose()
}

func (this *TcpTask) CloseIsNotCallback() {
	if !atomic.CompareAndSwapInt32(&this.closed, 0, 1) {
		return
	}
	glog.Info("[连接] 断开连接 ", this.Conn.RemoteAddr())
	this.Conn.Close()
	close(this.stopedChan)
}

func (this *TcpTask) Reset() bool {
	if atomic.LoadInt32(&this.closed) != 1 {
		return false
	}
	if !this.IsVerified() {
		return false
	}
	this.closed = -1
	this.verified = false
	this.stopedChan = make(chan bool)
	glog.Info("[连接] 重置连接 ", this.Conn.RemoteAddr())
	return true
}

func (this *TcpTask) IsClosed() bool {
	return atomic.LoadInt32(&this.closed) != 0
}

func (this *TcpTask) Verify() {
	this.verified = true
}

func (this *TcpTask) IsVerified() bool {
	return this.verified
}

func (this *TcpTask) Terminate() {
	this.Close()
}

func (this *TcpTask) CheckAndSend(buffer []byte, flag byte, threshold int) bool {
	if this.IsClosed() {
		return false
	}
	bsize := len(buffer)
	this.sendMutex.Lock()

	if this.sendBuffSizeLimit > 0 && this.sendBuff.RdSize()+bsize > this.sendBuffSizeLimit {
		// 缓冲区满。避免内存被耗尽
		glog.Errorln("send buff size limit. #1")
		this.Close()
		return false
	}

	this.sendBuff.Append(byte(bsize), byte(bsize>>8), byte(bsize>>16), flag)
	this.sendBuff.Append(buffer...)
	this.sendMutex.Unlock()
	this.Signal()
	return true
}

func (this *TcpTask) AsyncSend(buffer []byte, flag byte) bool {
	if this.IsClosed() {
		//glog.Error("AsyncSend ", this.Derived)
		return false
	}
	bsize := len(buffer)
	this.sendMutex.Lock()

	if this.sendBuffSizeLimit > 0 && this.sendBuff.RdSize()+bsize > this.sendBuffSizeLimit {
		// 缓冲区满。避免内存被耗尽
		glog.Errorln("send buff size limit. #2")
		this.Close()
		return false
	}

	this.sendBuff.Append(byte(bsize), byte(bsize>>8), byte(bsize>>16), flag)
	this.sendBuff.Append(buffer...)
	this.sendMutex.Unlock()
	this.Signal()
	return true
}

func (this *TcpTask) AsyncSendWithHead(head []byte, buffer []byte, flag byte) bool {
	if this.IsClosed() {
		return false
	}
	bsize := len(buffer) + len(head)
	this.sendMutex.Lock()

	// 该方法不会对外客户端，不用判断
	//	if this.sendBuffSizeLimit > 0 && this.sendBuff.RdSize()+bsize > this.sendBuffSizeLimit {
	//		// 缓冲区满。避免内存被耗尽
	//		this.Close()
	//		return false
	//	}

	this.sendBuff.Append(byte(bsize), byte(bsize>>8), byte(bsize>>16), flag)
	this.sendBuff.Append(head...)
	this.sendBuff.Append(buffer...)
	this.sendMutex.Unlock()
	this.Signal()
	return true
}

func (this *TcpTask) recvloop() {
	defer func() {
		if err := recover(); err != nil {
			glog.Error("[异常] ", err, "\n", string(debug.Stack()))
		}
	}()
	defer this.Close()

	var (
		recvBuff  *ByteBuffer = NewByteBuffer()
		neednum   int
		readnum   int
		err       error
		totalsize int
		datasize  int
		msgbuff   []byte
	)

	for {
		totalsize = recvBuff.RdSize()

		if totalsize < cmd_header_size {

			neednum = cmd_header_size - totalsize
			if recvBuff.WrSize() < neednum {
				recvBuff.WrGrow(neednum)
			}

			readnum, err = io.ReadAtLeast(this.Conn, recvBuff.WrBuf(), neednum)
			if err != nil {
				//glog.Error("[连接] 接收失败 ", this.Conn.RemoteAddr(), ",", err)
				return
			}

			recvBuff.WrFlip(readnum)
			totalsize = recvBuff.RdSize()
		}

		msgbuff = recvBuff.RdBuf()

		datasize = int(msgbuff[0]) | int(msgbuff[1])<<8 | int(msgbuff[2])<<16
		if datasize > cmd_max_size {
			glog.Error("[连接] 数据超过最大值 ", this.Conn.RemoteAddr(), ",", datasize)
			return
		}

		if totalsize < cmd_header_size+datasize {

			neednum = cmd_header_size + datasize - totalsize
			if recvBuff.WrSize() < neednum {
				recvBuff.WrGrow(neednum)
			}

			readnum, err = io.ReadAtLeast(this.Conn, recvBuff.WrBuf(), neednum)
			if err != nil {
				glog.Info("[连接] 接收失败 ", this.Conn.RemoteAddr(), ",", err)
				return
			}

			recvBuff.WrFlip(readnum)
			msgbuff = recvBuff.RdBuf()
		}

		this.Derived.ParseMsg(msgbuff[cmd_header_size:cmd_header_size+datasize], msgbuff[3])
		recvBuff.RdFlip(cmd_header_size + datasize)
	}
}

func (this *TcpTask) sendloop() {
	defer func() {
		if err := recover(); err != nil {
			glog.Error("[异常] ", err, "\n", string(debug.Stack()))
		}
	}()
	defer this.Close()

	var (
		tmpByte  = NewByteBuffer()
		timeout  = time.NewTimer(time.Second * cmd_verify_time)
		writenum int
		err      error
	)

	defer timeout.Stop()

	for {
		select {
		case <-this.signal:
			for {
				this.sendMutex.Lock()
				if this.sendBuff.RdReady() {
					tmpByte.Append(this.sendBuff.RdBuf()[:this.sendBuff.RdSize()]...)
					this.sendBuff.Reset()
				}
				this.sendMutex.Unlock()

				if !tmpByte.RdReady() {
					break
				}

				writenum, err = this.Conn.Write(tmpByte.RdBuf()[:tmpByte.RdSize()])
				if err != nil {
					glog.Info("[连接] 发送失败 ", this.Conn.RemoteAddr(), ",", err)
					return
				}
				tmpByte.RdFlip(writenum)
			}
		case <-this.stopedChan:
			return
		case <-timeout.C:
			if !this.IsVerified() {
				glog.Error("[连接] 验证超时 ", this.Conn.RemoteAddr())
				return
			}
		}
	}
}
