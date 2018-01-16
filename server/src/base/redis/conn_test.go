package redis

import (
	"fmt"
	"sync"
	//	"sync"
	"testing"
	"time"
)

var (
	conn *Conn
)

func init() {
	var err error
	conn, err = NewConn("192.168.124.130:6379", true)
	if err != nil {
		fmt.Println("link error:", err)
	}
	fmt.Println("link ok")
	//	cmd = NewBenchCommond(conn)
}

func TestChan(t *testing.T) {
	b := make(chan ICommond, 10)
	go func() {
		for {
			cmd := NewBenchCommond(nil)
			b <- cmd
			cmd.Set("pse", "srere")
			time.Sleep(time.Second)
		}
	}()
	time.Sleep(time.Second)

	for {
		cmd := <-b
		fmt.Println(len(cmd.GetBytes()))
		time.Sleep(time.Second)
	}

}

func TestConnLink(t *testing.T) {
	t.Log("Start")
	return
	w := new(sync.WaitGroup)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.999999999"))
	for j := 0; j < 10000; j++ {
		w.Add(1)
		//		k := j
		go func() {
			count := 5
			for k := 0; k < 100; k++ {
				cmd := NewBenchCommond(conn)
				for i := 0; i < count; i++ {
					cmd.Set("aa", "fdsaf")
				}
				cmd.Flush()
			}
			//			fmt.Println("finish!", k)
			w.Done()
		}()
	}
	w.Wait()
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.999999999"))
	//	b, err := String(cmd.Get("aa"))
	//	fmt.Println(b, err)
	//	time.Sleep(time.Second * 10)
	//conn.Send()
	//conn.beClosed()
}

func TestOldConn(t *testing.T) {
	return
	conn, err := NewRedisCon("192.168.124.130:6379", "", 1024, 1024)
	if err != nil {
		t.Log(err)
	}
	count := 5000000
	w := sync.WaitGroup{}
	w.Add(1)
	c := make(chan *Result, count)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.999999999"))
	go func() {
		for i := 0; i < count; i++ {
			<-c
		}
		w.Done()
	}()
	for i := 0; i < count; i++ {
		conn.Send(c, "set", "aa", "fdsaf")
	}
	w.Wait()
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.999999999"))
}

func BenchmarkConn(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	//	for i := 0; i < b.N; i++ {
	//		cmd.Set("aa", "12614621")
	//	}
}
