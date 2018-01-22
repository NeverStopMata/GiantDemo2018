package redis

func (this *Cluster) Del(key string) bool {
	//v, err := Bool(this.DoOnMain("DEL", key))
	v, err := Bool(this.DoOnMainArg1("DEL", key))
	if err != nil {
		return false
	}
	return v
}

func (this *Cluster) DelField(objKey, fieldKey string) bool {
	//v, err := Bool(this.DoOnMain("HDEL", objKey, fieldKey))
	v, err := Bool(this.DoOnMainArg2("HDEL", objKey, fieldKey))
	if err != nil {
		return false
	}
	return v
}

func (this *Cluster) ExistField(objKey string, fieldKey string) bool {
	//v, err := Bool(this.DoOnSlave("HEXISTS", objKey, fieldKey))
	v, err := Bool(this.DoOnSlaveArg2("HEXISTS", objKey, fieldKey))
	if err != nil {
		return false
	}
	return v
}

func (this *Cluster) Get(key string) (interface{}, error) {
	//return this.DoOnSlave("GET", key)
	return this.DoOnSlaveArg1("GET", key)
}

func (this *Cluster) GetOnMain(key string) (interface{}, error) {
	//return this.DoOnMain("GET", key)
	return this.DoOnMainArg1("GET", key)
}

func (this *Cluster) GetField(objKey string, fieldKey string) (interface{}, error) {
	//return this.DoOnSlave("HGET", objKey, fieldKey)
	return this.DoOnSlaveArg2("HGET", objKey, fieldKey)
}

func (this *Cluster) GetFieldOnMain(objKey string, fieldKey string) (interface{}, error) {
	//return this.DoOnMain("HGET", objKey, fieldKey)
	return this.DoOnMainArg2("HGET", objKey, fieldKey)
}

func (this *Cluster) GetObject(key string, obj interface{}) error {
	//v, err := Values(this.DoOnSlave("HGETALL", key))
	v, err := Values(this.DoOnSlaveArg1("HGETALL", key))
	if err != nil {
		return err
	}
	return ScanStruct(v, obj)
}

func (this *Cluster) GetObjectOnMain(key string, obj interface{}) error {
	//v, err := Values(this.DoOnMain("HGETALL", key))
	v, err := Values(this.DoOnMainArg1("HGETALL", key))
	if err != nil {
		return err
	}
	return ScanStruct(v, obj)
}

func (this *Cluster) Exist(key string) bool {
	//v, err := Bool(this.DoOnSlave("EXISTS", key))
	v, err := Bool(this.DoOnSlaveArg1("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

func (this *Cluster) ExistOnMain(key string) bool {
	//v, err := Bool(this.DoOnMain("EXISTS", key))
	v, err := Bool(this.DoOnMainArg1("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

func (this *Cluster) Incrby(key string, num int64) error {
	//return this.SendToMain(nil, "INCRBY", key, num)
	return this.SendToMainArg2(nil, "INCRBY", key, num)
}

func (this *Cluster) Send(commandName string, args ...interface{}) error {
	//return this.SendToMain(nil, commandName, args)
	return this.SendToMainArg1(nil, commandName, args)
}

func (this *Cluster) Set(key string, val interface{}) error {
	//return this.SendToMain(nil, "SET", key, val)
	return this.SendToMainArg2(nil, "SET", key, val)
}

func (this *Cluster) SetExpire(key string, second int) error {
	//return this.SendToMain(nil, "EXPIRE", key, second)
	return this.SendToMainArg2(nil, "EXPIRE", key, second)
}

func (this *Cluster) SetExpireAt(key string, timestamp int64) error {
	//return this.SendToMain(nil, "EXPIREAT", key, timestamp)
	return this.SendToMainArg2(nil, "EXPIREAT", key, timestamp)
}

func (this *Cluster) SetField(key string, fieldKey string, obj interface{}) error {
	//return this.SendToMain(nil, "HSET", Args{}.Add(key).Add(fieldKey).AddFlat(obj)...)
	return this.SendToMainArg3(nil, "HSET", key, fieldKey, obj)
}

func (this *Cluster) SetObject(key string, obj interface{}) error {
	return this.SendToMain(nil, "HMSET", Args{}.Add(key).AddFlat(obj)...)
}

func (this *Cluster) FieldIncrby(key string, fieldKey string, num int64) error {
	//return this.SendToMain(nil, "HINCRBY", key, fieldKey, num)
	return this.SendToMainArg3(nil, "HINCRBY", key, fieldKey, num)
}
