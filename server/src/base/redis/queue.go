package redis

import (
	"errors"
	"fmt"
	"sync"
)

var (
	noneResult = &RResult{Err: ErrRedisConClose}
)
var (
	// ring
	ErrRingEmpty = errors.New("ring buffer empty")
	ErrRingFull  = errors.New("ring buffer full")
)

type CmdQueue struct {
	sync.Mutex
	// read
	rp int
	rn int
	// write
	wp int
	wn int
	// control
	up   int
	un   int
	num  int
	data []ICommond
}

func NewCmdQueue(num int) (queue *CmdQueue) {
	queue = &CmdQueue{}
	queue.data = make([]ICommond, num)
	queue.num = num
	return
}

func (q *CmdQueue) Push(cmd ICommond) (err error) {
	q.Lock()
	if q.wn-q.rn >= q.num {
		q.Unlock()
		return ErrRingFull
	}
	q.data[q.wp] = cmd
	if q.wp++; q.wp >= q.num {
		q.wp = 0
	}
	q.wn++
	q.Unlock()
	return
}

func (q *CmdQueue) Get() (cmd ICommond) {
	if q.wn == q.un {
		return nil
	}
	cmd = q.data[q.up]
	if q.up++; q.up >= q.num {
		q.up = 0
	}
	q.un++
	return
}

func (q *CmdQueue) POP() (cmd ICommond) {
	if q.wn == q.rn {
		return nil
	}
	cmd = q.data[q.rp]
	return
}

func (q *CmdQueue) Adv() {
	q.data[q.rp] = nil
	if q.rp++; q.rp >= q.num {
		q.rp = 0
	}
	q.rn++
}

func (q *CmdQueue) ReplaceLast(cmd ICommond) {
	//	if q.wn-q.rn >= q.num {
	//		return ErrRingFull
	//	if q.wp != q.rp {
	//		fmt.Println("ReplaceLast error!")
	//	}
	//	q.Lock()
	q.data[q.rp] = nil
	q.data[q.rp] = cmd
	if q.wp++; q.wp >= q.num {
		q.wp = 0
	}
	if q.rp++; q.rp >= q.num {
		q.rp = 0
	}
	q.wn++
	q.rn++
	//	q.Unlock()
}

func (q *CmdQueue) Reset() {
	q.rn, q.rp, q.wn, q.wp, q.up, q.un = 0, 0, 0, 0, 0, 0
}

func (q *CmdQueue) String() string {
	return fmt.Sprintf("num:%d, rp:%d, rn:%d, wp:%d, wn:%d, up:%d, un:%d", q.num, q.rp, q.rn, q.wp, q.wn, q.up, q.un)
}
