package redis

import (
	"encoding/json"
	"errors"
	"strings"
)

type RedisOp struct {
	conn *RedisCon
}

func NewRedisOp(key int) *RedisOp {
	pool := Get(key)
	if pool == nil {
		return nil
	}
	return &RedisOp{conn: pool}
}

//func (this *RedisOp) Flush(isclose bool) error {
//	if isclose {
//		defer this.conn.D()
//	}
//	return this.conn.Flush()
//}

//func (this *RedisOp) Close() error {
//	return this.conn.Close()
//}

func (this *RedisOp) GetConn() *RedisCon {
	return this.conn
}

func (this *RedisOp) Set(key string, val interface{}) error {
	//return this.conn.Send("SET", key, val)
	return this.conn.SendArg2("SET", key, val)
}

func (this *RedisOp) SetEX(key string, seconds int64, val interface{}) error {
	return this.conn.SendArg3("SETEX", key, seconds, val)
}

func (this *RedisOp) SetObject(key string, obj interface{}) error {
	return this.conn.Send("HMSET", Args{}.Add(key).AddFlat(obj)...)
}

func (this *RedisOp) SetField(key string, fieldKey string, obj interface{}) error {
	//return this.conn.Send("HSET", Args{}.Add(key).Add(fieldKey).AddFlat(obj)...)
	return this.conn.SendArg3("HSET", key, fieldKey, obj)
}

func (this *RedisOp) Get(key string) (interface{}, error) {
	//return this.conn.Do("GET", key)
	return this.conn.DoArg1("GET", key)
}

func (this *RedisOp) GetObject(key string, obj interface{}) error {
	v, err := Values(this.conn.DoArg1("HGETALL", key))
	if err != nil {
		return err
	}
	return ScanStruct(v, obj)
}

func (this *RedisOp) IsSortedSetExist(key string, obj interface{}) bool {
	v, err := this.conn.Do("ZRANK", key, obj)
	if err != nil {
		return false
	}
	return v != nil
}
func (this *RedisOp) AddSortedSet(key string, score int64, obj interface{}) error {
	return this.conn.Send("ZADD", key, score, obj)
}
func (this *RedisOp) GetSortedGetMin(key string) (interface{}, int64, error) {
	v, err := this.conn.Do("ZRANGE", key, 0, 0, "WITHSCORES")
	if err != nil {
		return nil, 0, err
	}
	bulks, err := MultiBulk(v, err)
	if err != nil {
		return nil, 0, err
	}
	if bulks == nil || len(bulks) <= 1 {
		return nil, 0, errors.New("empty set")
	}
	score, err := Int64(bulks[1], err)
	if err != nil {
		return nil, 0, err
	}
	return bulks[0], score, nil
}
func (this *RedisOp) GetSortedReverseRange(key string, indexStart int, indexStop int) ([]interface{}, error) {
	v, err := this.conn.Do("ZREVRANGE", key, indexStart, indexStop)
	if err != nil {
		return nil, err
	}
	bulks, err := MultiBulk(v, err)
	if err != nil {
		return nil, err
	}
	return bulks, err
}
func (this *RedisOp) GetSortedGetMax(key string) (interface{}, int64, error) {
	v, err := this.conn.Do("ZREVRANGE", key, 0, 0, "WITHSCORES")
	if err != nil {
		return nil, 0, err
	}
	bulks, err := MultiBulk(v, err)
	if err != nil {
		return nil, 0, err
	}
	if bulks == nil || len(bulks) <= 1 {
		return nil, 0, errors.New("empty set")
	}
	score, err := Int64(bulks[1], err)
	if err != nil {
		return nil, 0, err
	}
	return bulks[0], score, nil
}

func (this *RedisOp) GetSortedBetween(key string, min, max int64) ([]interface{}, error) {
	v, err := this.conn.Do("ZRANGEBYSCORE", key, min, max)
	if err != nil {
		return nil, err
	}
	bulks, err := MultiBulk(v, err)
	if err != nil {
		return nil, err
	}
	return bulks, err
}

func (this *RedisOp) RemoveSortedSet(key string, obj interface{}) error {
	_, err := this.conn.Do("ZREM", key, obj)
	return err
}
func (this *RedisOp) GetSortedSetCout(key string) (int, error) {
	count, errCount := Int(this.conn.Do("ZCARD", key))
	if errCount != nil {
		return 0, errCount
	}
	return count, nil
}

func (this *RedisOp) UpdateSortedSetScore(key string, score int64, obj interface{}) (int64, error) {
	return Int64(this.conn.Do("ZINCRBY", key, score, obj))
}

func (this *RedisOp) Incrby(key string, val interface{}) error {
	return this.conn.SendArg2("INCRBY", key, val)
}

func (this *RedisOp) FieldIncrby(key string, fieldKey interface{}, val interface{}) error {
	return this.conn.SendArg3("HINCRBY", key, fieldKey, val)
}

func (this *RedisOp) GetField(objKey string, fieldKey interface{}) (interface{}, error) {
	return this.conn.DoArg2("HGET", objKey, fieldKey)
}

func (this *RedisOp) Exist(key string) bool {
	v, err := Bool(this.conn.DoArg1("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

func (this *RedisOp) ExistField(objKey string, fieldKey interface{}) bool {
	v, err := Bool(this.conn.DoArg2("HEXISTS", objKey, fieldKey))
	if err != nil {
		return false
	}
	return v
}

func (this *RedisOp) DelField(objKey, fieldKey interface{}) bool {
	v, err := Bool(this.conn.DoArg2("HDEL", objKey, fieldKey))
	if err != nil {
		return false
	}
	return v
}

//pitt 20170522 添加
func (this *RedisOp) DelFields(objKey string, args ...interface{}) (interface{}, error) {
	return this.Do("HDEL", Args{}.Add(objKey).AddFlat(args)...)
}

func (this *RedisOp) Del(key string) bool {
	v, err := Bool(this.conn.DoArg1("DEL", key))
	if err != nil {
		return false
	}
	return v
}

func (this *RedisOp) SetExpire(key string, second int) error {
	return this.conn.SendArg2("EXPIRE", key, second)
}

func (this *RedisOp) SetExpireAt(key string, timestamp int64) error {
	return this.conn.SendArg2("EXPIREAT", key, timestamp)
}

func (this *RedisOp) Do(commandName string, args ...interface{}) (interface{}, error) {
	return this.conn.Do(commandName, args...)
}

func (this *RedisOp) Send(commandName string, args ...interface{}) error {
	return this.conn.Send(commandName, args...)
}

func (this *RedisOp) GetFields(objKey string, args ...interface{}) (interface{}, error) {
	return this.Do("HMGET", Args{}.Add(objKey).AddFlat(args)...)
}

func (this *RedisOp) SetFields(objKey string, args ...interface{}) (interface{}, error) {
	return this.Do("HMSET", Args{}.Add(objKey).AddFlat(args)...)
}

type RedisObj struct {
	Name string
	conn *RedisPool
}

func NewRedisObj(name string, conn *RedisPool) *RedisObj {
	if conn == nil {
		return nil
	}
	return &RedisObj{name, conn}
}

func (this *RedisObj) GetConn() *RedisPool {
	return this.conn
}

//func (this *RedisObj) Flush(isclose bool) error {
//	if isclose {
//		defer this.conn.Close()
//	}
//	return this.conn.Flush()
//}

//func (this *RedisObj) Close() error {
//	return this.conn.Close()
//}

func (this *RedisObj) SetExpire(second int) error {
	return this.conn.Send("EXPIRE", this.Name, second)
}

func (this *RedisObj) Set(val interface{}) error {
	return this.conn.Send("HMSET", Args{}.Add(this.Name).AddFlat(val)...)
}

func (this *RedisObj) SetField(fieldKey, val interface{}) error {
	return this.conn.SendArg3("HSET", this.Name, fieldKey, val)
}

func (this *RedisObj) Get(val interface{}) error {
	v, err := Values(this.conn.DoArg1("HGETALL", this.Name))
	if err != nil {
		return err
	}
	return ScanStruct(v, val)
}

func (this *RedisObj) GetField(fieldKey string) (interface{}, error) {
	b, err := this.conn.DoArg2("HGET", this.Name, fieldKey)
	return b, err
}

func (this *RedisObj) Incrby(fieldKey string) error {
	return this.conn.SendArg3("HINCRBY", this.Name, fieldKey, 1)
}

func (this *RedisObj) Size() (int, error) {
	return Int(this.conn.DoArg1("HLEN", this.Name))
}

func (this *RedisObj) Remove(fieldKey string) error {
	return this.conn.SendArg2("HDEL", this.Name, fieldKey)
}

func (this *RedisObj) HasField(fieldKey string) bool {
	v, err := Bool(this.conn.DoArg2("HEXISTS", this.Name, fieldKey))
	if err != nil {
		return false
	}
	return v
}

func (this *RedisObj) Exist() bool {
	v, err := Bool(this.conn.DoArg1("EXISTS", this.Name))
	if err != nil {
		return false
	}
	return v
}

func (this *RedisObj) Clear() error {
	return this.conn.SendArg1("DEL", this.Name)
}

type SortedSet struct {
	Name string
	conn *RedisPool
}

func NewSortedSet(name string, conn *RedisPool) *SortedSet {
	if conn == nil {
		return nil
	}
	return &SortedSet{name, conn}
}

func (this *SortedSet) SetExpire(second int) error {
	return this.conn.SendArg2("EXPIRE", this.Name, second)
}

func (this *SortedSet) AddObject(score float64, v interface{}) error {
	//this.conn.Send("HMSET", redis.Args{}.Add(key).AddFlat(val))
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return this.conn.Send("ZADD", this.Name, score, b)
}

func (this *SortedSet) Set(score float64, v string) error {
	return this.conn.SendArg3("ZADD", this.Name, score, []byte(v))
}

func (this *SortedSet) AddString(score float64, v string) error {
	return this.conn.SendArg3("ZADD", this.Name, score, v)
}

func (this *SortedSet) Size() int {
	b, err := Int(this.conn.DoArg1("ZCARD", this.Name))
	if err != nil {
		return -1
	}
	return b
}

func (this *SortedSet) SizeByScore(min, max float64) int {
	b, err := Int(this.conn.DoArg3("ZCOUNT", this.Name, min, max))
	if err != nil {
		return -1
	}
	return b
}

func (this *SortedSet) GetObject(index int, clazz interface{}) error {
	b, err := Bytes(this.conn.DoArg3("ZRANGE", this.Name, index, index+1))
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, clazz)
	return err
}

func (this *SortedSet) Get(index int) (string, error) {
	b, err := Bytes(this.conn.DoArg3("ZRANGE", this.Name, index, index+1))
	if err != nil {
		return "", err
	}
	return string(b), err
}

//func (this *SortedSet) GetObjects(clazz []interface{}, start, limit int) error {
//  b, err := redis.MultiBulk(Do("ZRANGE", this.Name, start, start+limit))
//  if err != nil {
//      return err
//  }
//  for i, v := range b {
//      bb, err := redis.Bytes(v, nil)
//      if err != nil {
//          break
//      }
//      err = json.Unmarshal(bb, &clazz[i])
//      if err != nil {
//          break
//      }
//  }
//  return err
//}

func (this *SortedSet) RemoveObject(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return this.conn.Send("ZREM", this.Name, b)
}

func (this *SortedSet) Remove(v string) error {
	return this.conn.SendArg2("ZREM", this.Name, []byte(v))
}

func (this *SortedSet) GetString(index int) (string, error) {
	str, err := String(this.conn.DoArg3("ZRANGE", this.Name, index, index+1))
	if err == nil {
		str = strings.Trim(str, "\"")
	}
	return str, err
}

func (this *SortedSet) GetAllStrings() ([]string, error) {
	return this.GetStrings(0, -1)
}

func (this *SortedSet) FindAll() ([]string, error) {
	return this.Find(0, -1)
}

func (this *SortedSet) FindAllRev() ([]string, error) {
	return this.FindRev(0, -1)
}

func (this *SortedSet) GetStrings(start, limit int) ([]string, error) {
	a, err := this.conn.DoArg3("ZRANGE", this.Name, start, start+limit-1)
	if err != nil {
		return nil, err
	}
	b, err := MultiBulk(a, err)
	if err != nil {
		return nil, err
	}

	var list = make([]string, 0)
	for _, v := range b {
		s, err := String(v, nil)
		if err != nil {
			break
		}
		s = strings.Trim(s, "\"")
		list = append(list, s)
	}
	return list, err
}

func (this *SortedSet) Find(start, limit int) ([]string, error) {
	a, err := this.conn.DoArg3("ZRANGE", this.Name, start, start+limit-1)
	if err != nil {
		return nil, err
	}
	b, err := MultiBulk(a, err)
	if err != nil {
		return nil, err
	}

	var list = make([]string, 0)
	for _, v := range b {
		b, err := Bytes(v, nil)
		if err != nil {
			break
		}
		list = append(list, string(b))
	}
	return list, err
}

func (this *SortedSet) GetStringsRev(start, limit int) ([]string, error) {
	a, err := this.conn.DoArg3("ZREVRANGE", this.Name, start, start+limit-1)
	if err != nil {
		return nil, err
	}
	b, err := MultiBulk(a, err)
	if err != nil {
		return nil, err
	}

	var list = make([]string, 0)
	for _, v := range b {
		s, err := String(v, nil)
		if err != nil {
			break
		}
		s = strings.Trim(s, "\"")
		list = append(list, s)
	}
	return list, err
}

func (this *SortedSet) FindRev(start, limit int) ([]string, error) {
	a, err := this.conn.DoArg3("ZREVRANGE", this.Name, start, start+limit-1)
	if err != nil {
		return nil, err
	}
	b, err := MultiBulk(a, err)
	if err != nil {
		return nil, err
	}

	var list = make([]string, 0)
	for _, v := range b {
		b, err := Bytes(v, nil)
		if err != nil {
			break
		}
		list = append(list, string(b))
	}
	return list, err
}

func (this *SortedSet) RemoveString(v string) error {
	return this.conn.SendArg2("ZREM", this.Name, v)
}

func (this *SortedSet) RemoveRange(start, limit int) error {
	return this.conn.SendArg3("ZREMRANGEBYRANK", this.Name, start, start+limit-1)
}

func (this *SortedSet) RemoveIndex(index int) error {
	return this.conn.SendArg3("ZREMRANGEBYRANK", this.Name, index, index+1)
}

func (this *SortedSet) ObjectScore(v interface{}) int {
	b, err := json.Marshal(v)
	if err != nil {
		return -1
	}
	r, err := Int(this.conn.Do("ZINCRBY", this.Name, b))
	if err != nil {
		return -1
	}
	return r
}

func (this *SortedSet) StringScore(v string) int {
	r, err := Int(this.conn.DoArg2("ZINCRBY", this.Name, v))
	if err != nil {
		return -1
	}
	return r
}

func (this *SortedSet) Score(v string) int {
	r, err := Int(this.conn.DoArg2("ZINCRBY", this.Name, []byte(v)))
	if err != nil {
		return -1
	}
	return r
}

func (this *SortedSet) Clear() error {
	return this.conn.SendArg1("DEL", this.Name)
}
