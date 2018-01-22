package redis

import (
	"base/glog"
)

type RedisGroup struct {
	mainRedis    *RedisCon
	slaveAddrMap map[string]*RedisCon
	slaveRedis   []*RedisCon
	index        int //遍历SlaveRedis的索引，如果有多个redis，则循环使用slaveRedis
	slaveNum     int //slaveRedis的个数，保存起来不用每次计算

	waitChanLen int
	cmdChanLen  int
}

func NewRedisGroup(server string, passwd string, cmdChanlen, waitChanLen int) (*RedisGroup, error) {
	redisCon, err := NewRedisCon(server, passwd, cmdChanlen, waitChanLen)
	if err != nil {
		glog.Error("Make connection with redis failed")
		return nil, err
	}

	return &RedisGroup{
		mainRedis:   redisCon,
		index:       0,
		slaveNum:    0,
		slaveRedis:  make([]*RedisCon, 0, 1),
		cmdChanLen:  waitChanLen,
		waitChanLen: waitChanLen,
	}, nil
}

func (this *RedisGroup) reset() {
	this.slaveRedis = make([]*RedisCon, 0, 1)
	this.slaveNum = 0
	this.index = 0
}

func (this *RedisGroup) AddSlave(server string, passwd string) (err error) {
	var redisCon *RedisCon
	err = nil
	if redisCon = this.slaveAddrMap[server]; redisCon != nil {
		goto SETSLAVE
	}
	redisCon, err = NewRedisCon(server, passwd, this.cmdChanLen, this.waitChanLen)
	if err != nil {
		glog.Error("Make connection with redis failed")
		return err
	}

SETSLAVE:
	this.slaveRedis = append(this.slaveRedis, redisCon)
	this.slaveNum += 1
	return nil
}

func (this *RedisGroup) GetMain() *RedisCon {
	if this.mainRedis != nil && !this.mainRedis.IsNetClose() {
		return this.mainRedis
	}
	return nil
}

func (this *RedisGroup) GetSlave() *RedisCon {
	if this.slaveNum == 0 { //如果没有slave则返回mainRedis
		return this.GetMain()
	}

	nowIndex := this.index
	for {
		redis := this.slaveRedis[this.index]
		if redis != nil && !redis.IsNetClose() {
			this.index += 1
			if this.index == this.slaveNum {
				this.index = 0
			}
			return redis
		}
		this.index += 1
		if this.index == this.slaveNum {
			this.index = 0
		}

		if this.index == nowIndex { //找了一遍，所有slaveRedis都不在线
			return this.GetMain()
		}
	}
}
