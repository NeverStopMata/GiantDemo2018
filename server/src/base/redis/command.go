package redis

import (
	//	"fmt"
	"sync"
)

const (
	WriteCmd       = 0
	ReadCmd        = 1
	defaultCmdSize = 1024
)

var (
	flagAuth   = "AUTH"
	flagPing   = "PING"
	flagSelect = "SELECT"
	// keys
	flagDel    = "DEL"
	flagGet    = "get"
	flagSet    = "set"
	flagExist  = "EXISTS"
	flagDump   = "DUMP"
	flagExpire = "EXPIRE"
	flagKeys   = "KEYS"
	flagRename = "RENAME"
	// hashes
	flagHdel         = "HDEL"
	flagHexist       = "HEXISTS"
	flagHget         = "HGET"
	flagHset         = "HSET"
	flagHgetall      = "HGETALL"
	flagHincrby      = "HINCRBY"
	flagHincrbyFloat = "HINCRBYFLOAT"
	flagHkeys        = "HKEYS"
	flagHlen         = "HLEN"
	flagHmget        = "HMGET"
	flagHmset        = "HMSET"
	flagHsetnx       = "HSETNX"
	flagHstrlen      = "HSTRLEN"
	flagHvals        = "HVALS"
	flagHscan        = "HSCAN"
	// lists
	flagBlpop      = "BLPOP"
	flagBrpop      = "BRPOP"
	flagBrpoplpush = "BRPOPLPUSH"
	flagLindex     = "LINDEX"
	flagLinsert    = "LINSERT"
	flagLlen       = "LLEN"
	flagLpop       = "LPOP"
	flagLpush      = "LPUSH"
	flagLpushx     = "LPUSHX"
	flagLrange     = "LRANGE"
	flagLrem       = "LREM"
	flagLset       = "LSET"
	flagLtrim      = "LTRIM"
	flagRpop       = "RPOP"
	flagRpoplpush  = "RPOPLPUSH"
	flagRpush      = "RPUSH"
	flagRpushx     = "RPUSHX"
	// sets
	// sorted sets
)

var cmdType map[string]int = make(map[string]int)

func InitCmdType() {
	cmdType["GET"] = ReadCmd
	cmdType["MGET"] = ReadCmd
	cmdType["HMGET"] = ReadCmd
	cmdType["SMEMBERS"] = ReadCmd
	cmdType["SISMEMBER"] = ReadCmd
	cmdType["TTL"] = ReadCmd
	cmdType["ZSCORE"] = ReadCmd
	cmdType["ZREVRANGE"] = ReadCmd
	cmdType["ZRANGE"] = ReadCmd
	cmdType["ZCARD"] = ReadCmd
}

func CmdType(cmd string) int {
	return cmdType[cmd]
}

type RResult struct {
	Data interface{}
	Err  error
}

type ICommond interface {
	GetBytes() []byte
	SetData(t interface{}, err error) bool
	Done()
}

type RCommond struct {
	buf        *Packet
	conn       *Conn
	done       sync.WaitGroup
	c          chan bool
	result     RResult
	isAutoSend bool
	max        int
}

func NewCommond(conn *Conn) (cmd *RCommond) {
	cmd = &RCommond{}
	cmd.buf = NewPacketSize(defaultCmdSize)
	cmd.c = make(chan bool, 1)
	cmd.conn = conn
	cmd.isAutoSend = true
	return
}

func (cmd *RCommond) DoForCluster(c string, args []interface{}) {
	cmd.buf.WriteCommond(&c, &args)
}

func (cmd *RCommond) Do(c string, args ...interface{}) {
	//cmd.buf.WriteBytes()
}

func (cmd *RCommond) Scan(c string, arg1 int, arg2 string, arg3 int) {
	cmd.buf.WriteCmd(&c, 4)
	cmd.buf.WriteInt64(int64(arg1))
	cmd.buf.WriteString(&arg2)
	cmd.buf.WriteInt64(int64(arg3))
}

func (cmd *RCommond) Reset() {
	cmd.buf.Reset()
	cmd.result.Err = ErrRedisConClose
}

func (cmd *RCommond) GetBytes() (b []byte) {
	b = cmd.buf.buf[:cmd.buf.w]
	//	cmd.Reset()
	return
}

func (cmd *RCommond) SetData(data interface{}, err error) bool {
	cmd.result.Data = data
	cmd.result.Err = err
	return true
}

func (cmd *RCommond) Done() {
	//cmd.done.Done()
	select {
	case cmd.c <- true:
		// ok
	default:
	}
}

func (cmd *RCommond) waitConn() {
	if cmd.isAutoSend {
		cmd.conn.Send(cmd)
	} else {
		cmd.max++
	}
}

func (cmd *RCommond) doint1(c *string, arg1 *string, arg2 int64) {
	cmd.buf.WriteCmd(c, 3)
	cmd.buf.WriteString(arg1)
	cmd.buf.WriteInt64(arg2)
	cmd.waitConn()
}

func (cmd *RCommond) doint2(c *string, arg1 *string, arg2 int64, arg3 int64) {
	cmd.buf.WriteCmd(c, 3)
	cmd.buf.WriteString(arg1)
	cmd.buf.WriteInt64(arg2)
	cmd.buf.WriteInt64(arg3)
	cmd.waitConn()
}

func (cmd *RCommond) dostr1(c *string, arg1 *string) {
	cmd.buf.WriteCmd(c, 2)
	cmd.buf.WriteString(arg1)
	cmd.waitConn()
}

func (cmd *RCommond) dostr2(c *string, arg1 *string, arg2 *string) {
	cmd.buf.WriteCmd(c, 3)
	cmd.buf.WriteString(arg1)
	cmd.buf.WriteString(arg2)
	cmd.waitConn()
}

func (cmd *RCommond) dostr3(c *string, arg1 *string, arg2 *string, arg3 *string) {
	cmd.buf.WriteCmd(c, 4)
	cmd.buf.WriteString(arg1)
	cmd.buf.WriteString(arg2)
	cmd.buf.WriteString(arg3)
	cmd.waitConn()
}

func (cmd *RCommond) dostrmore(c *string, args *[]string) {
	cmd.buf.WriteCmd(c, 1+len(*args))
	for _, arg := range *args {
		cmd.buf.WriteString(&arg)
	}
	cmd.waitConn()
}

//////////////////////////////////////////////////////////
//////////////////        keys       /////////////////////
//////////////////////////////////////////////////////////
func (cmd *RCommond) Del(key string) error {
	cmd.dostr1(&flagDel, &key)
	return cmd.result.Err
}

func (cmd *RCommond) Get(key string) (interface{}, error) {
	cmd.dostr1(&flagGet, &key)
	return cmd.result.Data, cmd.result.Err
}

func (cmd *RCommond) Set(key string, val string) error {
	cmd.dostr2(&flagSet, &key, &val)
	return cmd.result.Err
}

func (cmd *RCommond) SetInt(key string, val int64) error {
	cmd.doint1(&flagSet, &key, val)
	return cmd.result.Err
}

func (cmd *RCommond) Exist(key string) error {
	cmd.dostr1(&flagExist, &key)
	return cmd.result.Err
}

func (cmd *RCommond) Expire(key string, val int) error {
	cmd.doint1(&flagExpire, &key, int64(val))
	return cmd.result.Err
}

func (cmd *RCommond) Keys(key string) (interface{}, error) {
	cmd.dostr1(&flagKeys, &key)
	return cmd.result.Data, cmd.result.Err
}

func (cmd *RCommond) Rename(key string, name string) {
	cmd.dostr2(&flagRename, &key, &name)
}

func (cmd *RCommond) Dump(key string) {
	cmd.dostr1(&flagDump, &key)
}

//////////////////////////////////////////////////////////
//////////////////        hashs       ////////////////////
//////////////////////////////////////////////////////////

func (cmd *RCommond) HSet(key string, field string, val string) {
	cmd.buf.WriteCmd(&flagHset, 4)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteString(&field)
	cmd.buf.WriteString(&val)
}

func (cmd *RCommond) HSetInt(key string, field string, val int64) {
	cmd.buf.WriteCmd(&flagHset, 4)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteString(&field)
	cmd.buf.WriteInt64(val)
}

func (cmd *RCommond) HSetVal(key string, field string, val interface{}) {
	cmd.buf.WriteCmd(&flagHset, 4)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteString(&field)
	cmd.buf.WriteArg(val)
}

type BenchCommond struct {
	*RCommond
	results []RResult
	cur     int
}

func NewBenchCommond(conn *Conn) (cmd *BenchCommond) {
	cmd = &BenchCommond{RCommond: NewCommond(conn)}
	cmd.isAutoSend = false
	return
}

func (cmd *BenchCommond) SetData(data interface{}, err error) bool {
	cmd.results = append(cmd.results, RResult{data, err})
	cmd.cur++
	//	fmt.Println("bench commond setdata:", cmd.cur, cmd.max)
	return cmd.cur == cmd.max
}

func (cmd *BenchCommond) Flush() (res []RResult) {
	//	cmd.done.Add(1)
	if cmd.conn.Send(cmd) == nil {
		<-cmd.c
		res = cmd.results
	}
	//	cmd.done.Wait()
	cmd.Reset()
	return
}

func (cmd *BenchCommond) Reset() {
	cmd.RCommond.Reset()
	cmd.results = make([]RResult, 0)
	cmd.cur, cmd.max = 0, 0
}

type SyncCommond struct {
	*RCommond
	results []RResult
	cur     int
	c       chan *RResult
}

func NewSyncCommond(conn *Conn, c chan *RResult) (cmd *SyncCommond) {
	cmd = &SyncCommond{RCommond: NewCommond(conn)}
	cmd.isAutoSend = false
	cmd.c = c
	return
}

func (cmd *SyncCommond) SetData(data interface{}, err error) bool {
	cmd.c <- &RResult{data, err}
	return true
}

func (cmd *SyncCommond) Done() {
}

func (cmd *SyncCommond) Flush() error {
	return cmd.conn.Send(cmd)
}
