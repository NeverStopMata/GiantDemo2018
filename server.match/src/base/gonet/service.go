package gonet

import (
	"base/glog"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
)

type IService interface {
	Init() bool
	MainLoop()
	Reload()
	Final() bool
}

type Service struct {
	terminate bool
	Derived   IService
}

func (this *Service) Terminate() {
	this.terminate = true
}

func (this *Service) isTerminate() bool {
	return this.terminate
}

func (this *Service) SetCpuNum(num int) {
	if num > 0 {
		runtime.GOMAXPROCS(num)
	} else if num == -1 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
}

func (this *Service) Main() bool {

	defer func() {
		if err := recover(); err != nil {
			glog.Error("[异常] ", err, "\n", string(debug.Stack()))
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGPIPE, syscall.SIGHUP)
	go func() {
		for sig := range ch {
			switch sig {
			case syscall.SIGHUP:
				this.Derived.Reload()
			case syscall.SIGPIPE:
			default:
				this.Terminate()
			}
			glog.Info("[服务] 收到信号 ", sig)
		}
	}()

	runtime.GOMAXPROCS(runtime.NumCPU())

	if !this.Derived.Init() {
		return false
	}

	for !this.isTerminate() {
		this.Derived.MainLoop()
	}

	this.Derived.Final()
	return true
}
