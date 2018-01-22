// Package dbclient has the client to dbserver.
package dbclient

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"base/glog"
	"base/rpc"
	"common"
)

var ErrDbClient = errors.New("db error")

type DbClient struct {
	client   *rpc.Client
	isclose  int32
	curUrl   string
	mainUrl  string
	slaveUrl string
	urlCount int
}

var dclientm *DbClient

func GetMe() *DbClient {
	if dclientm == nil {
		dclientm = &DbClient{
			isclose: 1,
		}
	}
	return dclientm
}

func (this *DbClient) SetUrl(str string) bool {
	dbaddrs := strings.Split(str, "/")
	if len(dbaddrs) == 0 {
		return false
	}
	this.mainUrl = dbaddrs[0]
	this.curUrl = this.mainUrl
	this.urlCount = 10

	if len(dbaddrs) > 1 {
		this.slaveUrl = dbaddrs[1]
	}

	return true
}

func (this *DbClient) Connect() bool {
	if atomic.LoadInt32(&this.isclose) == 0 {
		return false
	}
	client, err := rpc.Dial("tcp", this.curUrl, this.ReConnect)
	if err != nil {
		glog.Info("[RPC] 连接失败 ", err)
		return false
	}
	this.client = client
	atomic.StoreInt32(&this.isclose, 0)
	glog.Info("[RPC] 连接服务器成功 ", this.curUrl)
	return true
}

func (this *DbClient) ReConnect() bool {
	if atomic.CompareAndSwapInt32(&this.isclose, 0, 1) {
		for {
			glog.Info("[RPC] 重连中...")
			if this.Connect() {
				this.urlCount = 10
				break
			}
			this.urlCount--
			if this.urlCount < 0 && this.slaveUrl != "" {
				if this.curUrl == this.mainUrl {
					this.curUrl = this.slaveUrl
				} else {
					this.curUrl = this.mainUrl
				}
			}
			time.Sleep(time.Second)
		}
		return true
	}
	return false
}

func (this *DbClient) RemoteCall(serviceMethod string, args interface{}, reply interface{}) error {
	if this.client == nil {
		glog.Error("[RPC] 未初始化 ", serviceMethod)
		return ErrDbClient
	}
	err := this.client.Call(serviceMethod, args, reply)
	if err != nil {
		if err == rpc.ErrShutdown {
			if this.ReConnect() {
				return this.client.Call(serviceMethod, args, reply)
			} else {
				for i := 0; i < 30; i++ {
					glog.Error("[RPC] 调用等待 ", serviceMethod)
					time.Sleep(time.Second)
					if this.client == nil {
						continue
					}
					err = this.client.Call(serviceMethod, args, reply)
					if err != rpc.ErrShutdown {
						break
					}
				}
			}
		}
		//glog.Error("[RPC] 调用失败 ", serviceMethod, ",", err)
		return err
	}
	return nil
}

// 刷新房间数据
func (this *DbClient) SyncDatas(userids []uint64) ([]common.SyncData, bool) {
	reply := common.RetSyncDatas{}
	err := this.RemoteCall("RPCTask.SyncDatas", &common.ReqSyncDatas{userids}, &reply)
	if err != nil {
		return nil, false
	}
	return reply.SData, true
}

// 主动退出无尽模式
func (this *DbClient) OnExitFreeRoom(datas *common.ReqRoomInc) bool {
	reply := common.RetRoomInc{}
	err := this.RemoteCall("RPCTask.OnExitFreeRoom", datas, &reply)
	if err != nil {
		return false
	}
	return true
}

// 自由/组队结算
func (this *DbClient) RefreshRoomInc(datas *common.ReqRoomInc) bool {
	reply := common.RetRoomInc{}
	err := this.RemoteCall("RPCTask.RefreshRoomInc", datas, &reply)
	if err != nil {
		return false
	}
	return true
}

// 闪电战结算
func (this *DbClient) EndQRoom(datas *common.ReqEndQRoom) bool {
	reply := common.RetEndQRoom{}
	err := this.RemoteCall("RPCTask.EndQRoom", datas, &reply)
	if err != nil {
		return false
	}
	return true
}

func (this *DbClient) AddMoney(uid uint64, mtype, mnum uint32, text string, event string) bool {
	if mnum == 0 {
		return true
	}
	reply := common.RetAddMoney{}
	err := this.RemoteCall("RPCTask.AddMoney", &common.ReqAddMoney{uid, mtype, mnum, text, event}, &reply)
	if err != nil {
		return false
	}
	return true
}

func (this *DbClient) SubMoney(uid uint64, mtype, mnum uint32, text string, event string) bool {
	if mnum == 0 {
		return true
	}
	reply := common.RetSubMoney{}
	err := this.RemoteCall("RPCTask.SubMoney", &common.ReqSubMoney{uid, mtype, mnum, text, event}, &reply)
	if err != nil {
		return false
	}
	return true
}

func (this *DbClient) UpdateSeasonCourse(userid uint64, killNum, TopExp, sceneId, animalid uint32) bool {
	reply := common.RetUpdateSeasonCourse{}
	err := this.RemoteCall("RPCTask.UpdateSeasonCourse", &common.ReqUpdateSeasonCourse{userid, common.SEASON_ROOM_DATA, killNum, TopExp, sceneId, animalid, 0, 0}, &reply)
	if err != nil {
		return false
	}
	return true
}

func (this *DbClient) AddHideExp(userid uint64, hideexp uint32) bool {
	if hideexp == 0 {
		return true
	}
	reply := common.RetAddHideExp{}
	err := this.RemoteCall("RPCTask.AddHideExp", &common.ReqAddHideExp{userid, hideexp}, &reply)
	if err != nil {
		return false
	}
	return true
}

//获取社区用户信息
func (this *DbClient) GetUserInfo(userid uint64, toid uint64) (*common.SnsUserData, error) {
	reply := &common.SnsUserData{}
	err := this.RemoteCall("RPCTask.GetUserInfo", &common.ReqGetSnsInfo{FromId: userid, ToId: toid}, reply)
	if err != nil {
		glog.Error("[用户] 出错:", userid, err)
		return reply, err
	}
	return reply, nil
}

func (this *DbClient) SetLastPlayTime(UID uint64) bool {
	reply := common.RspSetLastPlayTime{}
	err := this.RemoteCall("RPCTask.SetLastPlayTime", &common.ReqSetLastPlayTime{UID: UID}, &reply)
	if err != nil {
		return false
	}
	return true
}

//写ulog到单点，方便汇总
func (this *DbClient) WriteULog(args ...interface{}) {
	info := fmt.Sprint(args...)
	reply := common.RspWriteULog{}
	this.RemoteCall("RPCTask.WriteULog", &common.ReqWriteULog{Info: info}, &reply)
}

func (this *DbClient) SetOnlineNotifyMessage(userid uint64, itemid, itemnum uint32) {
	reply := common.RspMessage{}
	this.RemoteCall("RPCTask.AddMessage", &common.ReqAddMessage{UID: userid, Itemid: itemid, ItemNum: itemnum}, &reply)
	return
}

// 通过id获取数据
func (this *DbClient) GetUserById(userid uint64) (*common.RetUserInfo, bool) {
	reply := common.RetUserInfo{}
	err := this.RemoteCall("RPCTask.GetUserById", &common.ReqUserById{Id: userid}, &reply)
	if err != nil {
		return nil, false
	}
	return &reply, true
}

// 返回玩家资源信息你
func (this *DbClient) GetResInfoById(userid uint64) (*common.RetResInfo, bool) {
	reply := common.RetResInfo{}
	err := this.RemoteCall("RPCTask.GetResInfoById", &common.ReqResInfo{UserId: userid}, &reply)
	if err != nil {
		return nil, false
	}
	return &reply, true
}

// 设置玩家资源信息
func (this *DbClient) SetResInfoById(userid uint64, req *common.ReqSetResInfo) (uint32, bool) {
	reply := common.RetSetResInfo{}
	err := this.RemoteCall("RPCTask.SetResInfoById", &common.ReqSetResInfo{UserId: userid, ResList: req.ResList}, &reply)
	if err != nil {
		return 0, false
	}
	return reply.Ret, true
}

//添加资源
func (this *DbClient) AddRes(userid uint64, resid, num uint32) (uint32, bool) {
	reply := &common.RetAddRes{}
	err := this.RemoteCall("RPCTask.AddRes", &common.ReqAddRes{UserId: userid, Resid: resid, Num: num}, reply)
	if err != nil {
		return 0, false
	}
	return reply.Ret, true
}
