package redis

/////////////////////////////////
//使用的必要步骤
//1.获取Slots实例
//  NewRedisCluster
//2.增加RedisGroup
//  AddGroup 或者 SetDefaultGroup
//3.执行Redis命令
//  AsyncDo 或者 SyncDo
////////////////////////////////

import (
	"base/glog"
	"errors"
	consul "github.com/hashicorp/consul/api"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const SlotsSepFlag = "|" //redisGroup对应的slots在consul中的分隔符

const (
	redisTypeMain  = 0
	redisTypeSlave = 1
)

const (
	ClusterStatusExit   = 0 //退出状态
	ClusterStatusInit   = 1 //cluster刚启动的时候的初始状态
	ClusterStatusNor    = 2 //正常状态
	ClusterStatusExpIng = 3 //正在扩容
	ClusterStatusExpEd  = 4 //扩容完了正在做后续处理
	ClusterStatusFatal  = 10
)

const maxTmpChanLen = 100000

const MAXSLOTNUM uint32 = 10240

const InvalidSlotValue = -1

var (
	ErrSlotLocked      = errors.New("Slot Locked")
	ErrSlotRedisNotSet = errors.New("Empty redis group mapped to slot")
	ErrInvalidArgs     = errors.New("Invalid args")
)

type Cluster struct {
	addrRedis        map[string]*RedisGroup //用RedisGroup中的MainRedis的地址表示这个RedisGroup
	slotsRedis       [MAXSLOTNUM]*RedisGroup
	newSlotsRedis    [MAXSLOTNUM]*RedisGroup
	defaultPasswd    string
	defaultWaitChLen int
	defaultCmdChLen  int

	tmpResultChan chan chan *Result

	clusterStatus int32
	slotMoving    int32

	consulSvr string
	name      string
	dbsType   string
}

//通过ch获取结果,如果不关心结果则ch为nil
func (this *Cluster) AsyncDo(ch chan *Result, cmd string, args ...interface{}) error {
	if CmdType(cmd) == ReadCmd {
		return this.SendToSlave(ch, cmd, &args)
	} else {
		return this.SendToMain(ch, cmd, &args)
	}
	return nil
}

//直接获取结果
func (this *Cluster) SyncDo(cmd string, args ...interface{}) (interface{}, error) {
	if CmdType(cmd) == ReadCmd {
		return this.DoOnSlave(cmd, args...)
	} else {
		return this.DoOnMain(cmd, args...)
	}
	return nil, nil
}

func (this *Cluster) getTmpRetChan() chan *Result {
	return <-this.tmpResultChan
}

func (this *Cluster) releaseTmpRetChan(ch chan *Result) {
	this.tmpResultChan <- ch
}

func (this *Cluster) DoOnMain(cmd string, args ...interface{}) (data interface{}, err error) {
	ch := make(chan *Result)
	err = this.SendToMain(ch, cmd, &args)
	if err != nil {
		close(ch)
		return
	}
	tmpResult := <-ch
	data = tmpResult.Data
	err = tmpResult.Err
	close(ch)
	return data, err
}

func (this *Cluster) DoOnMainArg1(cmd string, arg interface{}) (data interface{}, err error) {
	ch := make(chan *Result)
	err = this.SendToMainArg1(ch, cmd, arg)
	if err != nil {
		close(ch)
		return
	}
	tmpResult := <-ch
	data = tmpResult.Data
	err = tmpResult.Err
	close(ch)
	return data, err
}

func (this *Cluster) DoOnMainArg2(cmd string, arg1 interface{}, arg2 interface{}) (data interface{}, err error) {
	ch := make(chan *Result)
	err = this.SendToMainArg2(ch, cmd, arg1, arg2)
	if err != nil {
		close(ch)
		return
	}
	tmpResult := <-ch
	data = tmpResult.Data
	err = tmpResult.Err
	close(ch)
	return data, err
}

func (this *Cluster) DoOnMainArg3(cmd string, arg1 interface{}, arg2 interface{}, arg3 interface{}) (data interface{}, err error) {
	ch := make(chan *Result)
	err = this.SendToMainArg3(ch, cmd, arg1, arg2, arg3)
	if err != nil {
		close(ch)
		return
	}
	tmpResult := <-ch
	data = tmpResult.Data
	err = tmpResult.Err
	close(ch)
	return data, err
}

func (this *Cluster) sendArg0(redisCon *RedisCon, ch chan *Result, cmd *string) error {
	if redisCon == nil {
		ch <- &Result{
			Err: errors.New("RedisCon is nil"),
		}
		return errors.New("Can't get RedisCon for cmd")
	}
	if redisCon.IsNetClose() {
		return errors.New("The connection to redis was closed")
	}
	return redisCon.sendToCmdChanArgsNum0(ch, cmd)
}

func (this *Cluster) sendArg1(redisCon *RedisCon, ch chan *Result, cmd *string, arg interface{}) error {
	if redisCon == nil {
		ch <- &Result{
			Err: errors.New("RedisCon is nil"),
		}
		return errors.New("Can't get RedisCon for cmd")
	}
	if redisCon.IsNetClose() {
		return errors.New("The connection to redis was closed")
	}
	return redisCon.sendToCmdChanArgsNum1(ch, cmd, arg)
}

func (this *Cluster) sendArg2(redisCon *RedisCon, ch chan *Result, cmd *string, arg1 interface{}, arg2 interface{}) error {
	if redisCon == nil {
		ch <- &Result{
			Err: errors.New("RedisCon is nil"),
		}
		return errors.New("Can't get RedisCon for cmd")
	}
	if redisCon.IsNetClose() {
		return errors.New("The connection to redis was closed")
	}
	return redisCon.sendToCmdChanArgsNum2(ch, cmd, arg1, arg2)
}

func (this *Cluster) sendArg3(redisCon *RedisCon, ch chan *Result, cmd *string, arg1 interface{}, arg2 interface{}, arg3 interface{}) error {
	if redisCon == nil {
		ch <- &Result{
			Err: errors.New("RedisCon is nil"),
		}
		return errors.New("Can't get RedisCon for cmd")
	}
	if redisCon.IsNetClose() {
		return errors.New("The connection to redis was closed")
	}
	return redisCon.sendToCmdChanArgsNum3(ch, cmd, arg1, arg2, arg3)
}

func (this *Cluster) send(redisCon *RedisCon, ch chan *Result, cmd *string, args *[]interface{}) error {
	if redisCon == nil {
		ch <- &Result{
			Err: errors.New("RedisCon is nil"),
		}
		return errors.New("Can't get RedisCon for cmd")
	}
	if redisCon.IsNetClose() {
		return errors.New("The connection to redis was closed")
	}
	return redisCon.SendToCmdChan(ch, cmd, args)
}

func (this *Cluster) SendToMainForMoveData(ch chan *Result, cmd string, args *[]interface{}) error {
	if args == nil || len(*args) == 0 {
		return ErrInvalidArgs
	}
	redisCon, err := this.getRedis(redisTypeMain, &cmd, (*args)[0])
	if err != nil {
		//去掉通过chan返回错误
		//if ch != nil{
		//    ch<-&Result{
		//        Err:err,
		//        Data:nil,
		//    }
		//}
		return err
	}
	return this.send(redisCon, ch, &cmd, args)
}

//需要判断函数返回值，是否为nil为nil再等待结果
func (this *Cluster) SendToMain(ch chan *Result, cmd string, args ...interface{}) error {
	if len(args) == 0 {
		return ErrInvalidArgs
	}
	redisCon, err := this.getRedis(redisTypeMain, &cmd, args[0])
	if err != nil {
		//去掉通过chan返回错误
		//if ch != nil{
		//    ch<-&Result{
		//        Err:err,
		//        Data:nil,
		//    }
		//}
		return err
	}
	return this.send(redisCon, ch, &cmd, &args)
}

func (this *Cluster) SendToMainArg0(ch chan *Result, cmd string) error {
	return nil
	redisCon, err := this.getRedis(redisTypeMain, &cmd, nil)
	if err != nil {
		return err
	}
	return this.sendArg0(redisCon, ch, &cmd)
}

func (this *Cluster) SendToMainArg1(ch chan *Result, cmd string, arg interface{}) error {
	redisCon, err := this.getRedis(redisTypeMain, &cmd, arg)
	if err != nil {
		return err
	}
	return this.sendArg1(redisCon, ch, &cmd, arg)
}

func (this *Cluster) SendToMainArg2(ch chan *Result, cmd string, arg1 interface{}, arg2 interface{}) error {
	redisCon, err := this.getRedis(redisTypeMain, &cmd, arg1)
	if err != nil {
		return err
	}
	return this.sendArg2(redisCon, ch, &cmd, arg1, arg2)
}

func (this *Cluster) SendToMainArg3(ch chan *Result, cmd string, arg1 interface{}, arg2 interface{}, arg3 interface{}) error {
	redisCon, err := this.getRedis(redisTypeMain, &cmd, arg1)
	if err != nil {
		return err
	}
	return this.sendArg3(redisCon, ch, &cmd, arg1, arg2, arg3)
}

func (this *Cluster) GetResult(ch chan *Result, cmd string, args ...interface{}) (interface{}, error) {
	//result := <-ch
	//if result.Err == nil {
	//	return result.Data, result.Err
	//}
	//err := this.MoveErr(ch, result, cmd, &args)
	//if err != nil {
	//	return nil, err
	//}
	//result = <-ch
	//return result.Data, result.Err
	return nil, nil
}

func (this *Cluster) DoOnSlave(cmd string, args ...interface{}) (data interface{}, err error) {
	ch := make(chan *Result, 1)
	err = this.SendToSlave(ch, cmd, args...)
	if err != nil {
		close(ch)
		return
	}
	tmpResult := <-ch
	close(ch)
	data = tmpResult.Data
	err = tmpResult.Err
	return
}
func (this *Cluster) DoOnSlaveArg1(cmd string, arg interface{}) (data interface{}, err error) {
	ch := make(chan *Result, 1)
	err = this.SendToSlaveArg1(ch, cmd, arg)
	if err != nil {
		close(ch)
		return
	}
	tmpResult := <-ch
	close(ch)
	data = tmpResult.Data
	err = tmpResult.Err
	return
}
func (this *Cluster) DoOnSlaveArg2(cmd string, arg1 interface{}, arg2 interface{}) (data interface{}, err error) {
	ch := make(chan *Result, 1)
	err = this.SendToSlaveArg2(ch, cmd, arg1, arg2)
	if err != nil {
		close(ch)
		return
	}
	tmpResult := <-ch
	close(ch)
	data = tmpResult.Data
	err = tmpResult.Err
	return
}
func (this *Cluster) DoOnSlaveArg3(cmd string, arg1 interface{}, arg2 interface{}, arg3 interface{}) (data interface{}, err error) {
	ch := make(chan *Result, 1)
	err = this.SendToSlaveArg3(ch, cmd, arg1, arg2, arg3)
	if err != nil {
		close(ch)
		return
	}
	tmpResult := <-ch
	close(ch)
	data = tmpResult.Data
	err = tmpResult.Err
	return
}

func (this *Cluster) SendToSlave(ch chan *Result, cmd string, args ...interface{}) error {
	if len(args) == 0 {
		return ErrInvalidArgs
	}
	redisCon, err := this.getRedis(redisTypeSlave, &cmd, args[0])
	if err != nil {
		if err.Error() == ErrSlotLocked.Error() {
			this.movingSlotErr(ch, &cmd, &args)
		}
		return err
	}
	return this.send(redisCon, ch, &cmd, &args)
}

func (this *Cluster) SendToSlaveArg1(ch chan *Result, cmd string, arg interface{}) error {
	redisCon, err := this.getRedis(redisTypeSlave, &cmd, arg)
	if err != nil {
		return err
	}
	return this.sendArg1(redisCon, ch, &cmd, arg)
}

func (this *Cluster) SendToSlaveArg2(ch chan *Result, cmd string, arg1 interface{}, arg2 interface{}) error {
	redisCon, err := this.getRedis(redisTypeSlave, &cmd, arg1)
	if err != nil {
		return err
	}
	return this.sendArg2(redisCon, ch, &cmd, arg1, arg2)
}

func (this *Cluster) SendToSlaveArg3(ch chan *Result, cmd string, arg1 interface{}, arg2 interface{}, arg3 interface{}) error {
	redisCon, err := this.getRedis(redisTypeSlave, &cmd, arg1)
	if err != nil {
		return err
	}
	return this.sendArg3(redisCon, ch, &cmd, arg1, arg2, arg3)
}

func (this *Cluster) movingSlotErr(ch chan *Result, cmd *string, args *[]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			glog.Error(r)
		}
	}()
	if ch != nil {
		ch <- &Result{
			Data: nil,
			Err:  ErrSlotLocked,
		}
	}
	//TODO 记录redis命令
}

func (this *Cluster) getRedis(redisType int, cmd *string, arg interface{}) (redisCon *RedisCon, err error) {
	var key string
	key, err = String(arg, nil)
	if err != nil {
		glog.Error(err)
		glog.Error(arg)
		return
	}

	slotValue := int32(Crc32Slot(&key))
	//如果此时因为处理扩容的后续处理，不能write
	if atomic.LoadInt32(&(this.clusterStatus)) != ClusterStatusNor && slotValue == atomic.LoadInt32(&(this.slotMoving)) && CmdType(*cmd) == WriteCmd {
		//glog.Error("The slot is lock for write:", slotValue)
		//TODO 将命令保存起来
		return nil, ErrSlotLocked
	}
	redisGroup := this.slotsRedis[slotValue]
	if redisGroup != nil {
		if redisType == redisTypeMain {
			redisCon = redisGroup.GetMain()
		} else {
			redisCon = redisGroup.GetSlave()
		}
		return
	} else {
		redisCon = nil
		err = ErrSlotRedisNotSet
		glog.Error("Redis group was not set for slot:", slotValue)
		return
	}
}

func (this *Cluster) GetGroupOfAddr(server string) *RedisGroup {
	return this.addrRedis[server]
}

func (this *Cluster) GetGroupOfSlot(slotValue uint32) *RedisGroup {
	return this.slotsRedis[slotValue]
}

//func (this *Cluster) AddGroup(server string, passwd string, slotValue, cmdChanLen, waitChanLen int) (redisGroup *RedisGroup, err error) {
func (this *Cluster) AddGroup(server string, passwd string, cmdChanLen, waitChanLen int) (redisGroup *RedisGroup, err error) {
	//if slotValue < 0 || slotValue >= int(MAXSLOTNUM) {
	//	return nil, errors.New("Invalid slotValue" + strconv.Itoa(slotValue))
	//}

	err = nil
	if redisGroup = this.addrRedis[server]; redisGroup != nil {
		goto SETSLOT
	}
	redisGroup, err = NewRedisGroup(server, passwd, cmdChanLen, waitChanLen)
	if err != nil {
		glog.Error("[Redis]NewRedisGroup failed,server:", server, "passwd:", passwd)
		return nil, err
	}
	this.addrRedis[server] = redisGroup

SETSLOT:
	//this.slotsRedis[slotValue] = redisGroup
	this.defaultPasswd = passwd
	this.defaultWaitChLen = waitChanLen
	this.defaultCmdChLen = cmdChanLen
	return redisGroup, nil
}

var cluster *Cluster

func GetRedisCluster() *Cluster {
	return cluster
}

//dbserverName:此dbserver在consul中的名字
//consulSvr:此cluster连接的consul的server
func NewRedisCluster(dbserverType string, dbserverName string, consulSvr string) *Cluster {
	if cluster == nil {
		InitCmdType()
		cluster = new(Cluster)
		cluster.addrRedis = make(map[string]*RedisGroup)
		cluster.consulSvr = consulSvr
		cluster.name = NamePre + dbserverName //"dbs_r_c_"+dbserverName
		cluster.clusterStatus = ClusterStatusInit
		cluster.dbsType = ClusterConsulServerName + dbserverType
		cluster.setSlotMoving(InvalidSlotValue)
		var wait sync.WaitGroup
		wait.Add(1)
		go clusterMonitor(&wait, cluster)
		wait.Wait()
	}
	return cluster
}

func ClusterExit() {
	csCfg := consul.DefaultConfig()
	csCfg.Address = cluster.consulSvr
	csClient, err := consul.NewClient(csCfg)
	if err != nil {
		glog.Error("Get consul client failed:", err)
		return
	}

	kv := csClient.KV()
	p := &consul.KVPair{Key: RedisClusterStatusKey, Value: []byte(strconv.Itoa(ClusterStatusExit))}
	_, err = kv.Put(p, nil)
	if err != nil {
		glog.Error(err)
	}
	return
}

//用于和集群管理工具通信，获取节点信息,更新集群状态信息,获取redis是否在扩容等信息
func clusterMonitor(wait *sync.WaitGroup, c *Cluster) {
	csCfg := consul.DefaultConfig()
	csCfg.Address = c.consulSvr
	csClient, err := consul.NewClient(csCfg)
	if err != nil {
		glog.Error("Get consul client failed:", err)
		wait.Done()
		return
	}
	catalog := csClient.Catalog()
	reg := &consul.CatalogRegistration{
		Node:    c.name,
		Address: c.name,
		Service: &consul.AgentService{
			Service: c.dbsType,
		},
	}
	//将本节点的信息注册到consul中
	_, rErr := catalog.Register(reg, nil)
	if rErr != nil {
		glog.Error("Register server on consul failed:", rErr)
		wait.Done()
		return
	}

	kv := csClient.KV()

	var sValue int
	var kvPair *consul.KVPair
	err = c.getRedisClusterMap(kv, catalog)
	if err != nil {
		glog.Error("Set redis cluster map failed:", err)
		wait.Done()
		return
	}
	glog.Info("[Redis] RedisGroup is ready for work!")
	wait.Done()
	for {
		//获取集群的状态
		kvPair, _, err = kv.Get(RedisClusterStatusKey, nil)
		if err != nil {
			glog.Error("Get cluster status from consul failed:", err)
			return
		}
		sValue, err = strconv.Atoi(string(kvPair.Value))
		if err != nil {
			glog.Error("Convert redis cluster status failed:", err)
			return
		}
		//根据当前集群的状态做相应的处理
		switch sValue {
		case ClusterStatusNor: //正常状态
			if c.changeClusterStatus(ClusterStatusExpEd, ClusterStatusNor) || c.changeClusterStatus(ClusterStatusInit, ClusterStatusNor) { //如果redis扩容结束
				err = c.changeConsulDBSStatus(kv, ClusterStatusNor)
				if err != nil {
					glog.Error("Change status of the cluster in consul failed")
					return
				}
			}
			c.setSlotMoving(InvalidSlotValue)
			c.clusterStatus = int32(ClusterStatusNor)
			time.Sleep(time.Second * 1) //正常模式sleep时间长点
		case ClusterStatusExpIng: //正在进行扩容
			if c.changeClusterStatus(ClusterStatusNor, ClusterStatusExpIng) || c.changeClusterStatus(ClusterStatusInit, ClusterStatusExpIng) { //如果启动的时候或者运行的时候发现redis开始扩容
				if _, err = c.getMovingSlotFromConsul(kv); err != nil {
					glog.Error("RedisCluster status is moving slot, but get moved slot valued failed:", err)
					return
				}
				err = c.changeConsulDBSStatus(kv, ClusterStatusExpIng) //告知consul已经知道准备扩容
				if err != nil {
					glog.Error("Change status of the cluster in consul failed")
					return
				}
			}
			c.clusterStatus = int32(ClusterStatusExpIng)
			time.Sleep(time.Millisecond * 100) //异常模式(redis在扩容)sleep时间短一点，异常模式结束之后可以尽快回复状态
		case ClusterStatusExpEd: //扩容完成，做后续处理
			if c.changeClusterStatus(ClusterStatusExpIng, ClusterStatusExpEd) || c.changeClusterStatus(ClusterStatusInit, ClusterStatusExpEd) { //如果redis扩容之后在做后续处理
				if _, err = c.getMovingSlotFromConsul(kv); err != nil {
					glog.Error("RedisCluster status is moving slot, but get moved slot valued failed:", err)
					return
				}
				err = c.getRedisClusterMap(kv, catalog) //加载新的redis集群信息
				if err != nil {
					glog.Error("Moving slot finished, but reload new redis cluster info from consul failed:", err)
					return
				}
				err = c.changeConsulDBSStatus(kv, ClusterStatusExpEd)
				if err != nil {
					glog.Error("Change status of the cluster in consul failed")
					return
				}
			}
			c.clusterStatus = int32(ClusterStatusExpEd)
			time.Sleep(time.Millisecond * 100) //异常模式(redis在扩容)sleep时间短一点，异常模式结束之后可以尽快回复状态
		case ClusterStatusExit:
			glog.Error("Cluster is going to exit...")
			break
		default:
			c.clusterStatus = ClusterStatusFatal
			glog.Error("Invalid redis cluster status in consul")
		}
	}
}

//更改consul中该dbser的状态
func (this *Cluster) changeConsulDBSStatus(kv *consul.KV, newStatus int) error {
	kvPair := &consul.KVPair{
		Key:   RedisClusterStatusKey + this.name, //此key的值需要实时保持和RedisClusterStatusKey的值同步
		Value: []byte(strconv.Itoa(newStatus)),
	}
	_, err := kv.Put(kvPair, nil)
	if err != nil {
		glog.Error("Set dbserver redis cluster status failed:", err)
	}
	return err
}

//更改Cluster中的clusterStatus状态
func (this *Cluster) changeClusterStatus(oldStatus int, newStatus int) (ret bool) {
	ret = atomic.CompareAndSwapInt32(&this.clusterStatus, int32(oldStatus), int32(newStatus))
	if ret {
		glog.Error("[Redis] Cluster status changed to:", newStatus, " from:", oldStatus)
	}
	return
}

func (this *Cluster) setRedisGroupSlots(kv *consul.KV, server string) (err error) {
	//RedisGroupSlotsKey
	var kvPair *consul.KVPair
	kvPair, _, err = kv.Get(RedisGroupSlotsKey+server, nil)
	if err != nil {
		glog.Error("Get redis slots info failed:", RedisGroupSlotsKey, server)
		return
	}
	if kvPair == nil {
		glog.Error("Can not find slots info for redis group:", server)
		return errors.New("Can not find slots info for redis group")
	}
	var szSlots string
	szSlots = string(kvPair.Value) //redisGroup的slots是以竖线作为分割的数字
	slotsArray := strings.Split(szSlots, SlotsSepFlag)
	var iSlot int
	for _, num := range slotsArray {
		if len(num) > 0 && num != " " {
			iSlot, err = strconv.Atoi(num)
			if err != nil {
				glog.Error("Convert slot into num failed:", num)
				return
			}
			_, err = this.SetSlots(iSlot, server)
			if err != nil {
				glog.Error(err)
				return
			}
		}
	}
	return
}

//从consul获取redis集群信息
func (this *Cluster) getRedisClusterMap(kv *consul.KV, cl *consul.Catalog) error {
	cataSvrs, _, e := cl.Service(RedisClusterMapService, "", nil)
	if e != nil {
		glog.Error("Get redis map from consul failed:", e)
		return e
	}
	var redisGroup *RedisGroup
	var err error
	for _, ctSvr := range cataSvrs {
		mainSvr := ctSvr.Node //main redis的server
		glog.Error("Main redis server:", mainSvr)
		if redisGroup = this.addrRedis[mainSvr]; redisGroup != nil {
			redisGroup.reset() //已经有这个redis的Group
		} else {
			redisGroup, err = this.AddGroup(mainSvr, "", 10000, 10000)
			if err != nil {
				glog.Error("AddGroup failed, redis server:", mainSvr)
				return err
			}
		}
		err = this.setRedisGroupSlots(kv, mainSvr)
		if err != nil {
			glog.Error("Set slots for redis group failed")
			return err
		}
		glog.Info("Main redis server:", mainSvr)
		var slaveSvr string
		for i := 0; i < len(ctSvr.ServiceTags); i++ {
			slaveSvr = ctSvr.ServiceTags[i]
			redisGroup.AddSlave(slaveSvr, "")
			//slaveSvr为slave redis的值
			glog.Error("Slave redis server:", slaveSvr)
		}
	}
	return nil
}

func (this *Cluster) SetSlotsRange(slotBegin int, slotEnd int, addr string) (*RedisGroup, error) {
	var redisGroup *RedisGroup
	var err error
	for slot := slotBegin; slot <= slotEnd; slot++ {
		if redisGroup, err = this.SetSlots(slot, addr); err != nil {
			break
		}
	}
	return redisGroup, err
}

//设置slot对应的redisGroup
func (this *Cluster) SetSlots(slotsNum int, addr string) (*RedisGroup, error) {
	if slotsNum < 0 || uint32(slotsNum) > MAXSLOTNUM {
		return nil, errors.New("Invalid slot")
	}
	redisGroup := this.addrRedis[addr]
	if redisGroup != nil {
		this.slotsRedis[slotsNum] = redisGroup
		return redisGroup, nil
	}
	return this.AddGroup(addr, this.defaultPasswd, this.defaultCmdChLen, this.defaultWaitChLen)
}

func (this *Cluster) UpdateClusterStatus(oldStatus int32, newStatus int32) bool {
	return atomic.CompareAndSwapInt32(&(this.clusterStatus), oldStatus, newStatus)
}

func (this *Cluster) setSlotMoving(slot int32) {
	this.slotMoving = slot
}

//如果正在对redis做扩容，则获取被被移动的slot
func (this *Cluster) getMovingSlotFromConsul(kv *consul.KV) (slot int, err error) {
	var kvPair *consul.KVPair
	slot = InvalidSlotValue
	kvPair, _, err = kv.Get(RedisClusterMovingSlot, nil)
	if err != nil {
		glog.Error("Get cluster status from consul failed:", err)
		return
	}
	if kvPair == nil {
		glog.Error("Can not get moving slot value from consul")
		return slot, errors.New("Can not get moving slot value from consul")
	}
	slot, err = strconv.Atoi(string(kvPair.Value))
	if err != nil {
		glog.Error("Convert redis cluster status failed:", err)
		return
	}
	this.slotMoving = int32(slot)
	glog.Error("RedisCluster is moving key of slot:", slot)
	return
}
