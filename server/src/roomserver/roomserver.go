package main

import (
	"base/env"
	"base/glog"
	"base/gonet"
	"base/httppprof"
	"math/rand"
	"net/http"
	"os"
	"runtime/pprof"
	"time"

	"roomserver/chart"
	"roomserver/client/dbclient"
	"roomserver/client/rcenterclient"
	"roomserver/client/voiceclient"
	"roomserver/conf"
	"roomserver/game/ai"
	"roomserver/game/roommgr"
	"roomserver/game/scn"
	"roomserver/game/scn/playermgr"
	"roomserver/game/scn/plr"
	"roomserver/game/skill"
	"roomserver/redismgr"
	"roomserver/tcptask"
	"roomserver/udp"
)

type RoomServer struct {
	gonet.Service
	roomser *gonet.TcpServer
}

func roomServerMain() {
	svr := RoomServer{
		roomser: &gonet.TcpServer{},
	}
	svr.Derived = &svr
	svr.Main() // -> gonet.Service.Main()
}

func (this *RoomServer) Init() bool {
	glog.Info("[启动] 开始初始化")

	rand.Seed(time.Now().Unix())

	// 性能检查
	pprofport := env.Get("room", "pprofport")
	if pprofport != "" {
		go func() {
			http.ListenAndServe(pprofport, nil)
		}()
	}

	// 全局配置
	if !conf.ConfigMgr_GetMe().Init() {
		glog.Error("[启动] 读取全局配置失败")
		return false
	}

	tcptask.ScenePlayerMgr = playermgr.GetMe()
	tcptask.RoomMgr = roommgr.GetMe()
	plr.NewIScenePlayerAI = ai.NewIScenePlayerAI
	plr.NewISkillPlayer = skill.NewISkillPlayer
	plr.NewISkillBall = skill.NewISkillBall
	scn.NewSceneAIPlayer = ai.NewSceneAIPlayer
	scn.NewCopyPlayerAI = ai.NewCopyPlayerAI

	if !redismgr.GetMe().Init() {
		return false
	}

	if !dbclient.GetMe().SetUrl(env.Get("room", "db")) {
		return false
	}

	if !dbclient.GetMe().Connect() {
		return false
	}

	voiceclient.RoomMgr = roommgr.GetMe()
	if !voiceclient.GetMe().Connect() {
		return false
	}

	initChart() // 在线人数画图

	// 绑定本地端口
	port := env.Get("room", "listen")
	err := this.roomser.Bind(port)
	if err != nil {
		glog.Error("[启动] 绑定端口失败")
		return false
	}

	// 如果配置了UDP端口，则运行UDP监听服务
	udp.PlayerTaksMgr = tcptask.PlayerTaskMgr_GetMe()
	udpPort := env.Get("room", "listenUDP")
	if udpPort != "" {
		udp.UdpMgr_GetMe().RunServer(udpPort)
	}

	if !this.initRCenterClient() {
		return false
	}
	if !ai.CreateBevAIMgr() {
		glog.Error("[启动]CreateBevAIMgr fail! ")
	}
	if !skill.LoadSkillBevTree() {
		glog.Error("[启动]LoadSkillBevTree fail! ")
	}
	if !conf.InitMapConfig() {
		glog.Error("[启动]InitMapConfig fail! ")
		return false
	}

	glog.Info("[启动] 完成初始化")
	startCPUProfile()
	return true
}

func (this *RoomServer) MainLoop() {
	conn, err := this.roomser.Accept()
	if err != nil {
		//glog.Error("roomser not Accept")
		return
	}
	tcptask.NewPlayerTask(conn).Start()
}

func (this *RoomServer) Final() bool {
	this.roomser.Close()
	roommgr.GetMe().Final() // AllUserOffline()
	redismgr.Final()
	stopCPUProfile()
	return true
}

func (this *RoomServer) Reload() {
	if !conf.ReloadConfig() {
		glog.Info("[配置] 加载失败")
	} else {
		glog.Info("[配置] 加载成功")
	}
}

func startCPUProfile() {
	if conf.ConfigMgr_GetMe().Global.Pystress != 0 {
		f, err := os.Create("cpu.prof")
		if err != nil {
			panic("open file failed")
		}
		pprof.StartCPUProfile(f)

		go httppprof.StartPProf()
	}
}

func stopCPUProfile() {
	if conf.ConfigMgr_GetMe().Global.Pystress != 0 {
		pprof.StopCPUProfile()
	}
}

func initChart() {
	func1 := tcptask.PlayerTaskMgr_GetMe().GetNum
	func2 := roommgr.GetMe().GetNum
	func3 := playermgr.GetMe().GetNum
	chart.ChartMgr_GetMe().Init(func1, func2, func3)
}

func (this *RoomServer) initRCenterClient() bool {
	rcenterclient.RoomNumGetter = roommgr.GetMe()
	rcenterclient.UserNumGetter = tcptask.PlayerTaskMgr_GetMe()
	rcenterclient.ReloginChecker = tcptask.PlayerTaskMgr_GetMe()
	rcenterclient.Terminator = this
	if !rcenterclient.GetMe().Connect() {
		return false
	}
	return true
}
