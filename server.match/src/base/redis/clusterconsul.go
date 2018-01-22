package redis

const (
	NamePre                 = "dbs_r_c_"  //dbserver的name的qianzhui
	ClusterConsulServerName = "dbserver_" //service值前缀，node和address为dbserver的name。用于在dbserver启动的时候向consul注册本节点信息

	RedisGroupSlotsKey     = "redisSlotsKey_" //kv key加mainRedis的server做后缀，value为此redisGroup对应的slots
	RedisClusterStatusKey  = "RCStatusKey"    //kv 值为redisCluster的状态
	RedisClusterMovingSlot = "RedisMoingSlot" //kv 值为正在扩展的slot
	RedisClusterMapService = "RCMap"          //service值，node为mainRedis的server, tags为slaveredis的server
)
