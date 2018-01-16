package tcptask

// 玩家网络会话管理类

import (
	"runtime/debug"
	"sync"
	"time"

	"base/glog"
	"common"
	"roomserver/conf"
	"roomserver/udp"
	"usercmd"
)

const (
	_TASK_MAX_TIMEOUT = 1 // 连接超时时间（单位：分）
)

type PlayerTaskMgr struct {
	mutex sync.RWMutex
	tasks map[uint64]*PlayerTask
}

var (
	ptaskm             *PlayerTaskMgr
	playertaskmgr_once sync.Once
)

func PlayerTaskMgr_GetMe() *PlayerTaskMgr {
	if ptaskm == nil {
		playertaskmgr_once.Do(func() {
			ptaskm = &PlayerTaskMgr{
				tasks: make(map[uint64]*PlayerTask),
			}
			ptaskm.Init()
		})
	}

	return ptaskm
}

func (this *PlayerTaskMgr) Init() {
	go func() {
		tick := time.NewTicker(time.Second * 3)
		defer func() {
			tick.Stop()
			if err := recover(); err != nil {
				glog.Error("[异常] 心跳线程出错 ", err, "\n", string(debug.Stack()))
			}
		}()
		for {
			select {
			case <-tick.C:
				this.timeAction()
			}
		}
	}()
}

func (this *PlayerTaskMgr) Add(task *PlayerTask) bool {
	if task == nil {
		return false
	}
	this.mutex.Lock()
	this.tasks[task.id] = task
	this.mutex.Unlock()
	return true
}

func (this *PlayerTaskMgr) remove(task *PlayerTask) bool {
	if task == nil {
		glog.Info("[PlayerTaskMgr_Remove_11] ")
		return false
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	t, ok := this.tasks[task.id]
	if !ok {
		glog.Info("[PlayerTaskMgr_Remove_12] ")
		return false
	}
	if t.key != task.key {
		glog.Info("[PlayerTaskMgr_Remove_13] ", task.id, ",", task.udata.Account, ",", task.name, ",", task.key, " ,:", t.key)
		return false
	}
	delete(this.tasks, task.id)
	return true
}

func (this *PlayerTaskMgr) GetTask(uid uint64) *PlayerTask {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	user, ok := this.tasks[uid]
	if !ok {
		return nil
	}
	return user
}

func (this *PlayerTaskMgr) GetNum() int {
	this.mutex.RLock()
	tasknum := len(this.tasks)
	this.mutex.RUnlock()
	return tasknum
}

func (this *PlayerTaskMgr) timeAction() {

	if conf.ConfigMgr_GetMe().Global.Pystress != 0 {
		// 压力测试，方便测试客户端不发心跳包
		return
	}

	var ptasks []*PlayerTask
	this.mutex.RLock()
	for _, t := range this.tasks {
		if t.IsTimeout() {
			ptasks = append(ptasks, t)
		}
	}
	this.mutex.RUnlock()
	for _, t := range ptasks {
		if !t.Stop() {
			this.remove(t)
		}
		glog.Info("[玩家] 连接超时 ", t.id, ",", t.udata.Account, ",", t.key)
	}
}

// CheckRelogin 检测重登, 如果检测到重登则踢人并返回true.
func (this *PlayerTaskMgr) CheckRelogin(revCmd *usercmd.ReqCheckRelogin) bool {
	task := this.GetTask(revCmd.Id)
	if task == nil {
		glog.Error("[登录] 重复登录检查,未找到玩家 ", revCmd.Id, ",", revCmd.Key)
		return false
	}
	if task.key != revCmd.Key {
		glog.Error("[登录] 重复登录检查,key不一致 ", revCmd.Id, ",", revCmd.Key, ",", task.key)
		return false
	}
	task.RetErrorMsg(common.ErrorCodeReLogin)
	task.Stop()
	glog.Info("[登录] 重复登录检查,被踢下线 ", revCmd.Id, ",", revCmd.Key, ",", task.name)
	return true
} // CheckRelogin()

// AddUdpSess 添加Udp会话到相应的PlayerTask
func (this *PlayerTaskMgr) AddUdpSess(
	revCmd *usercmd.MsgBindTCPSession) *udp.UdpSess {

	uid := revCmd.GetId()
	user := this.GetTask(uid)
	// udpConn 绑定时验证room key，防止劫持
	if user == nil || user.key != revCmd.GetKey() {
		glog.Error("[UDPSession] 要绑定的用户不存在或room key不正确")
		return nil
	}

	if user.room == nil || user.room.IsClosed() == true {
		return nil
	}

	udpSess := udp.NewUdpSess(uid, user.key)
	data := make(map[uint64]*udp.UdpSess)
	data[uid] = udpSess
	user.room.PostBindUdpSession(data)
	glog.Info("[UDPSession] 绑定TCPSession成功，uid: ", uid)
	return udpSess
}
