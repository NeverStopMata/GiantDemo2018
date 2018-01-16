// Room server.
package main

import (
	"base/env"
	"base/glog"
	"flag"
	"runtime/debug"
)

var (
	logfile = flag.String("logfile", "", "Log file name")
	config  = flag.String("config", "config.json", "config path")
)

func main() {
	flag.Parse()
	defer func() {
		if err := recover(); err != nil {
			glog.Error("[异常] 报错 ", err, "\n", string(debug.Stack()))
		}
	}()

	if !env.Load(*config) {
		return
	}

	initLog()
	roomServerMain()
	glog.Info("[关闭] 房间服务器关闭完成")
	glog.Flush()
} // main()

func initLog() {
	loglevel := env.Get("global", "loglevel")
	if loglevel != "" {
		flag.Lookup("stderrthreshold").Value.Set(loglevel)
	}

	logtostderr := env.Get("global", "logtostderr")
	if logtostderr != "" {
		flag.Lookup("logtostderr").Value.Set(logtostderr)
	}

	if *logfile != "" {
		glog.SetLogFile(*logfile)
	} else {
		glog.SetLogFile(env.Get("room", "log"))
	}
	if env.Get("room", "logType") == "1" {
		glog.SetLogType(glog.LogTimeType_Day) //日志每天一换
	}
} // initLog()
