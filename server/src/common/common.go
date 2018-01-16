package common

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

const (
	SEASON_ROOM_DATA = 1 // 赛季房间数据
	SEASON_FANS      = 2 // 赛季新粉丝
)

const (
	BanType_LeaveMsg              = "2" //对留言进行封禁
	BanType_LeaveMsgToSuperPlayer = "3" //对大神进行封禁
	BanType_Sign                  = "4" //对签名进行封禁
	BanType_Barrage               = "5" //封禁弹幕
	BanType_All                   = "6" //全部封禁
)

// 玩家登录token数据
type UserData struct {
	Id          uint64 `redis:"Id"`          // 玩家id
	Account     string `redis:"Account"`     // 用户名(账号)
	Icon        uint32 `redis:"Icon"`        // 玩家icon
	PassIcon    string `redis:"PassIcon"`    // 已审核头像icon
	Sex         uint8  `redis:"Sex"`         // 性别
	Money       uint32 `redis:"Money"`       // 彩豆
	JiangGuo    uint32 `redis:"JiangGuo"`    // 小浆果
	Location    uint32 `redis:"Location"`    // 位置
	ShowPos     uint8  `redis:"ShowPos"`     // 显示地区
	RoomId      uint32 `redis:"RoomId"`      // 房间id
	RoomAddr    string `redis:"RoomAddr"`    // 房间地址
	Model       uint32 `redis:"Model"`       // 登录模式(观战等)
	State       string `redis:"State"`       // 动态
	CState      string `redis:"CState"`      // 自由合作状态
	TCState     string `redis:"TCState"`     // 团队合作状态
	Level       uint32 `redis:"Level"`       // 段位
	PlayNum     uint32 `redis:"PlayNum"`     // 总局数
	GameDayTime uint32 `redis:"GameDayTime"` // 当天游戏时间
	IsNewbie    bool   `redis:"IsNewbie"`    // 是否新手(团队模式)
	Robot       uint32 `redis:"Robot"`       // 机器人级别
	TeamId      uint32 `redis:"TeamId"`      // Teamid
	TeamName    uint32 `redis:"TeamName"`    // 队伍名 (组队)
	IsLeader    bool   `redis:"IsLeader"`    // 是否队长
	IsCustom    bool   `redis:"IsCustom"`    // 是否自建房间(自建房间)
	RoomName    string `redis:"RoomName"`    // 房间名(自建房间)
	WatchedId   uint64 `redis:"WatchedId"`   // 被观看者id
	UnCoop      uint32 `redis:"UnCoop"`      // 是否不合作(自由模式)
	Scores      uint32 `redis:"Scores"`      // 分数
	RegDev      string `redis:"RegDev"`
	SceneId     uint32 `redis:"SceneId"`    //地图id
	HideScore   uint32 `redis:"HideScores"` // 隐藏积分
}

// 设备数据
type RedisDevice struct {
	Id uint64 `redis:"Id"` // 玩家id
}

// ip限制数据
type RedisIpLimit struct {
	Ip    string `redis:"Ip"`
	Count uint32 `redis:"Count"`
	Time  int64  `redis:"Time"`
	Limit uint32 `redis:"Limit"`
}

// 账号数据
type RedisAccount struct {
	Id        uint64 `redis:"Id"`        // 玩家id
	Password  string `redis:"Password"`  // 密码
	MCount    uint32 `redis:"MCount"`    // 修改次数
	State     uint32 `redis:"State"`     // 账号状态
	BanTime   int64  `redis:"BanTime"`   // 解封时间
	BanReason uint32 `redis:"BanReason"` // 封禁原因
	BanDays   uint32 `redis:"BanDays"`   // 账号被封天数
}

// 注册数据
type RegGameData struct {
	Account  string `redis:"Account"`  // 账号
	DevCode  string `redis:"DevCode"`  // 设备号
	Icon     uint32 `redis:"Icon"`     // icon
	Location uint32 `redis:"Location"` // 位置
	ShowPos  uint8  `redis:"ShowPos"`  // 显示地区
	RegTime  int64  `redis:"RegTime"`  // 注册时间
	RegIp    string `redis:"RegIp"`    // 注册ip
	RegDev   string `redis:"RegDev"`   // 注册设备
	Level    uint32 `redis:"Level"`    // 段位
	Scores   uint32 `redis:"Scores"`   // 星数
}

// 玩家数据
type RedisGameData struct {
	Account      string `redis:"Account"`      // 账号
	DevCode      string `redis:"DevCode"`      // 设备号
	Email        string `redis:"Email"`        // 邮箱地址
	Tel          string `redis:"Tel"`          // 手机号
	Icon         uint32 `redis:"Icon"`         // icon
	Sign         string `redis:"Sign"`         // 签名
	Sex          uint8  `redis:"Sex"`          // 性别
	Age          uint32 `redis:"Age"`          // 年龄
	Money        uint32 `redis:"Money"`        // 彩豆
	JiangGuo     uint32 `redis:"JiangGuo"`     // 小浆果
	BattTicket   uint32 `redis:"BattTicket"`   // 入场券
	Weight       uint64 `redis:"Weight"`       // 总体重
	PlayNum      uint32 `redis:"PlayNum"`      // 总局数
	EatNum       uint32 `redis:"EatNum"`       // 总吞噬球数
	FirstNum     uint32 `redis:"FirstNum"`     // 总冠军数
	MaxWeight    uint64 `redis:"MaxWeight"`    // 单局最大体重
	KillNum      uint32 `redis:"KillNum"`      // 总吞噬人数
	CState       string `redis:"CState"`       // 自由合作状态
	TCState      string `redis:"TCState"`      // 团队合作状态
	HScores      uint32 `redis:"HScores"`      // 自由隐藏分
	THScores     uint32 `redis:"THScores"`     // 团队隐藏分
	SeasonId     uint32 `redis:"SeasonId"`     // 赛季
	MaxLevel     uint32 `redis:"MaxLevel"`     // 最高段位
	MaxScores    uint32 `redis:"MaxScores"`    // 最高星
	Level        uint32 `redis:"Level"`        // 段位
	Scores       uint32 `redis:"Scores"`       // 星数
	Fans         uint32 `redis:"Fans"`         // 粉丝数
	Follows      uint32 `redis:"Follows"`      // 关注数
	Blacks       uint32 `redis:"Blacks"`       // 坏人数
	Friends      uint32 `redis:"Friends"`      // 好友数
	Location     uint32 `redis:"Location"`     // 位置
	ShowPos      uint8  `redis:"ShowPos"`      // 显示地区
	State        string `redis:"State"`        // 动态
	FansRank     string `redis:"FanRank"`      // 粉丝排行
	IconBox      uint32 `redis:"IconBox"`      // 当前头象框
	NewIcon      string `redis:"NewIcon"`      // 用户最新头像
	PassIcon     string `redis:"PassIcon"`     // 通过验证的头像
	AudienceNum  uint32 `redis:"AudienceNum"`  // 听众人数
	AudienceUrl  string `redis:"AudienceUrl"`  // 语音签名
	LoveNum      uint32 `redis:"LoveNum"`      // 被喜欢次数
	AudienceTime uint32 `redis:"AudienceTime"` // 语音签名时间
	Position     string `redis:"Position"`     // 玩家位置详情
	Robot        uint32 `redis:"Robot"`        // 对应机器人等级
	HideScores   uint32 `redis:"HideScores"`   // 无限模式与限时模式的隐藏积分

	Nick        string `redis:"Nick"`        // 昵称
	Sceneid     uint32 `redis:"Sceneid"`     // 场景ID
	ActCodeGift uint32 `redis:"ActCodeGift"` // 激活码奖励
	RegIp       string `redis:"RegIp"`
	RegTime     uint32 `redis:"RegTime"`
	RegDev      string `redis:"RegDev"`
	Speaker     uint32 `redis:"Speaker"`
	OpenId      string `redis:"OpenId"`
	LoginTime   uint32 `redis:"LoginTime"`
	IdCard      string `redis:"IdCard"`
}

// 红点数据
type RedPointData struct {
	FollowNum uint32 `redis:"FollowNum"` // 关注列表
	FansNum   uint32 `redis:"FansNum"`   // 粉丝列表
	BlackNum  uint32 `redis:"BlackNum"`  // 坏人列表
	StarNum   uint32 `redis:"StarNum"`   // 被赞次数
	MsgNum    uint32 `redis:"MsgNum"`    // 留言次数
	BansNum   uint32 `redis:"BansNum"`   // 黑名单列表
}

// 玩家在线数据
type RedisUser struct {
	RServerId uint16 `redis:"RServerId"` // 房间服id
	TServerId uint16 `redis:"TServerId"` // 组队服id
	RoomId    uint32 `redis:"RoomId"`    // 自由房间
	TRoomId   uint32 `redis:"TRoomId"`   // 队伍房间
	QRoomId   uint32 `redis:"QRoomId"`   // 闪电战房间
	Model     uint32 `redis:"Model"`     // 游戏模式
	State     uint32 `redis:"State"`     // 在线状态
	Key       string `redis:"Key"`       // 登录key
	TeamId    uint32 `redis:"TeamId"`    // 队伍id (组队)
	TeamName  uint32 `redis:"TeamName"`  // 队伍名 (组队)
	LeaderId  uint64 `redis:"LeaderId"`  // 队长id (组队)
	IsLeader  bool   `redis:"IsLeader"`  // 是否队长 (组队)
	IsNewbie  bool   `redis:"IsNewbie"`  // 是否萌新 (组队)
	RoomOwner uint64 `redis:"RoomOwner"` // 房主id (自建房间)
	Priv      uint32 `redis:"Priv"`      // 权限 (自建房间)
	ToId      uint64 `redis:"ToId"`      // 被观战者
}

// 分享类型
const (
	SHARE_TYPE_DIREADD = 1 // 直接加了
	SHARE_TYPE_PENDING = 2 // 挂起
	SHARE_TYPE_FIND    = 3 // 找回
	SHARE_TYPE_REGED   = 4 // 小号
	SHARE_TYPE_EMULAT  = 5 // 模拟器
	SHARE_TYPE_BANIP   = 6 // 非法小号
)

// 分享数据
type RedisShare struct {
	ShareId uint64 `redis:"ShareId"` // 分享者id
	SType   uint32 `redis:"SType"`   // 分享类型
	SAdded  uint32 `redis:"SAdded"`  // 分享已添加
	FindId  uint64 `redis:"FindId"`  // 找回者id
	FType   uint32 `redis:"FType"`   // 找回类型
	FTime   int64  `redis:"FTime"`   // 找回时间
	FAdded  uint32 `redis:"FAdded"`  // 找回已添加
}

// 玩家基本数据
type RedisBaseInfo struct {
	Id           uint64
	Account      string // 账号
	Icon         uint32 // icon
	NewIcon      string // 头象
	PassIcon     string // 通过验证的头像
	Sex          uint8  // 性别
	Location     uint32 // 位置
	LoveLocation uint32
	MaxLevel     uint32 // 最高段位
	Level        uint32 // 当前段位
	Scores       uint32 // 当前星数
	State        string // 动态
	AudienceNum  uint32 // 听众人数
	LoveNum      uint32 // 被喜欢次数
	PlayNum      uint32 // 游戏次数
	THScores     uint32 // 团队隐藏分
	ShowPos      uint8
	Age          uint32
	Sign         string
	AudienceUrl  string
	AudienceTime uint32
	ValidUser    bool //是否有效用户
	Nickname     string
}

// 玩家排行数据
type RankNode struct {
	Value uint32
	Time  uint32
}

type FansRanks struct {
	DayFans   RankNode
	WeekFans  RankNode
	TotalFans RankNode
}
type StarsRanks struct {
	DayStars   RankNode
	WeekStars  RankNode
	TotalStars RankNode
}

const (
	RankDay      = 1 // 日
	RankLocation = 2 // 月/地区
	RankTotal    = 3 // 总
)

const (
	//EatRank   = 1 // 吞噬球数
	FanRank    = 2 // 粉丝
	StarRank   = 3 // 赞
	KillRank   = 4 // 吞噬人数
	LevelRank  = 5 // 段位排名
	LoveRank   = 6 // 喜欢排名
	FlowerRank = 7 // 小红花排行

	MedalRank     = 9  // 勋章排名
	MvpRank       = 10 // MVP排行榜
	ZhenRank      = 11 // 棒棒糖价值排名
	FTopScoreRank = 12 // 无限最高分排行榜
	QTopScoreRank = 13 // 限时模式最高分排行榜
)

type ReportBugInfo struct {
	Id          bson.ObjectId `bson:"_id"` //bugID
	Uid         uint64        `bson:"uid"`
	Type        int32         `bson:"type"`
	BugType     int32         `bson:"bug_type"` // 1:登录, 2:程序, 3:美术, 4:文字, 5:卡顿, 6:其它
	BugName     string        `bson:"bug_name"`
	BugDesc     string        `bson:"bug_desc"`
	SubmitUser  string        `bson:"submit_user"`
	DeviceType  string        `bson:"device_type"`
	GameVersion string        `bson:"game_version"`
	BugPicture  string        `bson:"bug_picture"`
	BugFrom     string        `bson:"bug_from"`
	Contact     string        `bson:"contact"`
	SubmitTime  int64         `bson:"submit_time"`
	BugState    int32         `bson:"bug_state"`
}

type PhotoInfo struct {
	Id        bson.ObjectId     `bson:"_id"` //图片ID
	Uid       uint64            //玩家ID
	FileSize  uint32            //文件大小
	ShowNum   uint32            //浏览次数
	LaudNum   uint32            //赞次数
	UploadUrl string            //图片地址
	Thumbnail map[string]string //缩略图
	Time      int64             //上传时间
	State     uint32            //状态 (0 未审核  1 审核通过  2审核未通过  3删除)
	Account   string
}

type PhotoExtInfo struct {
	*PhotoInfo
	Account  string
	IsPraise bool // 是否点赞
	IsShow   bool // 是否查看
}

// 关注数据
type FollowUser struct {
	Id         uint64 // 玩家id
	Account    string // 帐号
	Icon       uint32 // 图标
	Sex        uint8  // 性别
	State      string // 状态
	PassIcon   string // 头像
	IsNew      uint32
	RelType    uint32
	Level      uint32 //段位
	Scores     uint32 //星数
	Follows    uint32 //关注数
	Fans       uint32 //粉丝数
	Remark     string //备注
	IsBan      uint32 //是否是黑名单
	FollowTime uint32 // 关注时间
}

// 好友列表
type FriendUser struct {
	Id       uint64 // 玩家id
	Account  string // 帐号
	Icon     uint32 // 图标
	Sex      uint8  // 性别
	State    string // 状态
	PassIcon string // 头像
	NickName string // 昵称
	RelType  uint32
	Level    uint32 //段位
	Scores   uint32 //星数
}

// 新手数据
type NewbieUser struct {
	Id       uint64 // 玩家id
	Account  string // 玩家账号
	Icon     uint32 // 玩家图标
	Sex      uint8  // 性别
	PassIcon string // 头象
	Location uint32 // 位置
}

// 粉丝数据
type FanUser struct {
	Id         uint64 // 玩家id
	Account    string // 帐号
	Icon       uint32 // 图标
	Sex        uint8  // 性别
	PassIcon   string //头像
	IsNew      uint32
	RelType    uint32
	Level      uint32 //段位
	Scores     uint32 //星数
	Follows    uint32 //关注数
	Fans       uint32 //粉丝数
	Sign       string //签名
	Remark     string // 备注
	Distance   uint32 // 距离
	FollowTime uint32 // 关注时间
}

// 比赛列表
type MatchUser struct {
	Icon     uint32
	NickName string
	Time     uint32
	Rank     uint32
	SceneID  uint32
	JiangGuo uint32
	AnimalID uint32
	Score    uint32
	PassIcon string //头像
}

// 荣誉列表
type HonorUser struct {
	Season uint32 // 赛季
	Honor  uint32 // 荣誉
	Level  uint32 // 段位
	Scores uint32 // 星数
	Rank   uint32 // 排名
	ZpRank int32  // 排名
}

//动态消息
type DynamicMsg struct {
	RelationData  *RelationData
	RedPointData  *RedPointData
	SnsDynamicMsg *SnsDynamicMsg
	Reward        uint32
}

//sns相关的动态
type SnsDynamicMsg struct {
	LoveNum      uint32 //喜欢数量
	FollowNum    uint32 //关注数量
	LeaveMsgNum  uint32 //主页留言
	StarPhoneNum uint32 //照片点赞
	MsgReplyNum  uint32 //我的留言的回复数量
}

// 统计数据
type RelationData struct {
	FollowNum  uint32 // 关注数量
	FansNum    uint32 // 粉丝数量
	BlackNum   uint32 // 坏人数量
	StarNum    uint32 // 被赞次数
	BanListNum uint32 // 黑名单数量
	FriendNum  uint32 // 好友数量
}

// 关注动态
type UserState struct {
	Ftype int    // 类型
	Time  uint32 // 触发时间
	Param uint64 // 参数
}

// 留言数据
type RMessageData struct {
	Id      bson.ObjectId `bson:"_id"` // 留言id
	FromId  uint64        // 留言者
	ToId    uint64        // 被留言者
	Time    uint32        // 留言时间
	Stars   []uint64      // 点赞玩家列表
	StarNum uint32        // 点赞数量
	Text    string        // 留言内容
	ReplyId string        // 回复id
	RAcc    string        // 回复账号
	RText   string        // 回复内容
	FromAcc string
	ToAcc   string
	RId     uint64 //被回复者
	SortId  uint64
}

// 举报数据
type ReportUserMsg struct {
	Id     bson.ObjectId `bson:"_id"` // id
	Uid    uint64        // 玩家id
	RClass uint32        //举报大类
	RType  uint32        // 举报类型
	RId    string        // 举报资源ID
	Reason string        // 举报理由
	Time   int64         // 结算时间
}

// 在线状态
const (
	UserOnline  = 1 // 上线
	UserOffline = 2 // 下线
)

// 排行数据
type RankUser struct {
	Rank       uint32 // 排名
	Id         uint64 // 玩家id
	Account    string // 帐号
	Icon       uint32 // 图标
	Sex        uint8  // 性别
	IncNum     uint32 // 增加赞数
	AllNum     uint32 // 总吞噬数
	PassIcon   string //头像
	RankChange int32
	IsFollow   bool // 是否是关注
	Rtroom     uint32
}

// 自由模式排行榜
type FreeRank struct {
	Id       uint64
	Score    float64
	LastRank int32
}
type FreeRankList []FreeRank

func (self FreeRankList) Len() int {
	return len(self)
}

func (self FreeRankList) Less(i, j int) bool {
	return self[i].Score > self[j].Score
}

func (self FreeRankList) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

// 组队模式排行榜
type TeamRank struct {
	TeamId   uint32
	TeamName uint32
	UScores  map[uint64]float64
	LastRank int32
}

func (this *TeamRank) Scores() uint64 {
	var score float64
	for _, s := range this.UScores {
		score += s
	}
	return uint64(score)
}

type TeamRankList []*TeamRank

func (self TeamRankList) Len() int {
	return len(self)
}

func (self TeamRankList) Less(i, j int) bool {
	var score1, score2 float64
	for _, s := range self[i].UScores {
		score1 += s
	}
	for _, s := range self[j].UScores {
		score2 += s
	}
	return score1 > score2
}

func (self TeamRankList) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

// 语音类型
const (
	Audio_Kind_Default = 0 // 主页语音
	Audio_Kind_RedEnv  = 1 // 红包语音
	Audio_Kind_Chat    = 2 // 聊天语音
)

// 组队房间模式
const (
	TEAM_MODEL_NONE = 0 // 团队
)

// 登录模式
const (
	UserModelJoin   = 1 // 加入模式(自由/组队)
	UserModelWatch  = 2 // 观战模式
	UserModelRand   = 3 // 随机模式(无效)
	UserModelTeam   = 4 // 组队模式
	UserModelQuick  = 5 // 限时模式
	UserModelCustom = 6 // 自建房间模式
	UserModelOb     = 7 // OB模式
)

// 在线状态
const (
	UserStateOffline  = 0  // 离线
	UserStateOnline   = 1  // 在线
	UserStateFPlaying = 2  // 游戏中(无尽)
	UserStateQPlaying = 3  // 游戏中(限时)
	UserStateTeam     = 4  // 组队中
	UserStateTeamInv  = 5  // 邀请中
	UserStateTPlaying = 7  // 游戏中(团战)
	UserStateCustom   = 9  // 自建房间中
	UserStateWatching = 10 // 观战中
)

// 服务器类型
const (
	ServerTypeRoom  = 1 // 房间服务器
	ServerTypeTeam  = 2 // 队伍服务器
	ServerTypeLogin = 3 // 登录服务器
	ServerTypeChat  = 4 // 聊天服务器
)

// 房间类型
const (
	RoomTypeTeam  = 1 // 组队模式
	RoomTypeQuick = 2 // 闪电战模式
)

const (
	//总类型
	RedisType  int = 1 << 24
	MongoType  int = 2 << 24
	ServerType int = 3 << 24
	SsdbType   int = 4 << 24

	//子类型
	DefaultSubType   int = 0
	OutServerSubType int = 1 << 16
	InServerSubType  int = 2 << 16

	//redis类型
	WAccRedis   int = RedisType | DefaultSubType | 1  // 账号写
	RAccRedis   int = RedisType | DefaultSubType | 2  // 账号读
	WLstRedis   int = RedisType | DefaultSubType | 3  // 排行写
	RLstRedis   int = RedisType | DefaultSubType | 4  // 排行读
	WRelRedis   int = RedisType | DefaultSubType | 5  // 关系写
	RRelRedis   int = RedisType | DefaultSubType | 6  // 关系读
	StatRedis   int = RedisType | DefaultSubType | 7  // 临时数据
	WActRedis   int = RedisType | DefaultSubType | 8  // 活动数据
	GeoRedis    int = RedisType | DefaultSubType | 9  // lbs位置
	LoveRedis   int = RedisType | DefaultSubType | 10 // 喜欢数据
	TokenRedis  int = RedisType | DefaultSubType | 11 // 登录认证
	PlayerRedis int = RedisType | DefaultSubType | 12 // 玩家数据
	CacheRedis  int = RedisType | DefaultSubType | 13 // 缓存数据
	FansRedis   int = RedisType | DefaultSubType | 14 // 粉丝排序数据

	ChatRedis     int = RedisType | DefaultSubType | 17 // 聊天房间路由数据
	TCenterRedis  int = RedisType | DefaultSubType | 18 // 聊天房间路由数据
	PayOrderRedis int = RedisType | DefaultSubType | 19 // 充值订单数据

	PbeRedis int = RedisType | DefaultSubType | 100 // 测试redis

	//mongo类型
	Mongodb     int = MongoType | DefaultSubType | 1 // 默认数据库
	FfMongodb   int = MongoType | DefaultSubType | 2 // lbs位置数据
	ReadMongodb int = MongoType | DefaultSubType | 3 // 找朋友数据
	ChatMongodb int = MongoType | DefaultSubType | 4 // 聊天

	//服务器类型
	ChatDbServer    int = ServerType | InServerSubType | 1  // 聊天DB服务
	ChatServer      int = ServerType | InServerSubType | 2  // 聊天服务
	GeoDbserver     int = ServerType | InServerSubType | 3  // LBS DB服务
	GeoServer       int = ServerType | InServerSubType | 4  // LBS服务
	GmServer        int = ServerType | InServerSubType | 5  // GM服务
	LoginDbserver   int = ServerType | InServerSubType | 6  // 登陆DB服务
	LoginServer     int = ServerType | InServerSubType | 7  // 登陆服务
	MgrServer       int = ServerType | InServerSubType | 8  // 配置管理服务
	PresenceServer  int = ServerType | InServerSubType | 9  //
	RCenterServer   int = ServerType | InServerSubType | 10 // 房间中心服务
	RoomDbServer    int = ServerType | InServerSubType | 11 // 房间DB服务
	RoomServer      int = ServerType | InServerSubType | 12 // 房间服务
	ShareServer     int = ServerType | InServerSubType | 13 // 分享服务
	TCenterServer   int = ServerType | InServerSubType | 14 // 组队中心服务
	TeamDbServer    int = ServerType | InServerSubType | 15 // 组队DB服务
	TeamServer      int = ServerType | InServerSubType | 16 // 组队服务
	TimerServer     int = ServerType | InServerSubType | 17 //
	UploadDbServer  int = ServerType | InServerSubType | 18 // 上传DB服务
	UploadServer    int = ServerType | InServerSubType | 19 // 上传服务
	VoiceServer     int = ServerType | InServerSubType | 20 // 语音服务
	WatcherDbServer int = ServerType | InServerSubType | 21 // 观战DB服务
	WatcherServer   int = ServerType | InServerSubType | 22 // 观战服务
	WCenterServer   int = ServerType | InServerSubType | 23 // 观战中心服务
	DbServer        int = ServerType | InServerSubType | 24 // DB服务
	GatewayServer   int = ServerType | InServerSubType | 25 // 网关服务

	MLstSsdb int = SsdbType | DefaultSubType | 1 // 房间比赛记录
)

// 机器人最小id
const (
	MINROBOTID uint64 = 900000000
)

// 货币类型
const (
	MONEY_ID_1 = 1001 // 彩豆
	MONEY_ID_2 = 1002 // 棒棒糖

	MONEY_ID_TICKET  = 190001 // 入场券
	MONEY_ID_SPEAKER = 195001 //喇叭
)

// 账号状态
const (
	ACCOUNT_STATE_NORMAL = 0 // 正常
	ACCOUNT_STATE_BAN    = 1 // 封号
	ACCOUNT_STATE_OB     = 2 // OB
)

// GM操作方式
const (
	GmExecQuery int = 1 // 查询

	GmExecEMoney       int = 3  // 货币
	GmExecObject       int = 4  // 道具
	GmExecBanAcc       int = 5  // 封号
	GmExecTDrop        int = 6  // 礼包
	GmExecChPass       int = 7  // 重置密码
	GmExecReloadCfg    int = 8  // 动态加载配置文件
	GmExecChkWhitelist int = 9  //查检白名单
	GmExecAddWhite     int = 10 // 添加白名单
)

// 错误码
const (
	ErrorCodeOkay      = 0  // 成功
	ErrorCodeName      = 1  // 名字非法
	ErrorCodeVersion   = 2  // 版本号错误
	ErrorCodeRoom      = 3  // 未找到房间
	ErrorCodeDb        = 4  // 数据库错误
	ErrorCodePass      = 5  // 账号或密码错误
	ErrorCodeParam     = 6  // 参数错误
	ErrorCodeSession   = 7  // session错误
	ErrorCodeDecode    = 8  // 解码错误
	ErrorCodeCmd       = 9  // 未知指令
	ErrorCodeExist     = 10 // 账号名已存在
	ErrorCodeNotExist  = 11 // 账号不存在
	ErrorCodeOffline   = 12 // 玩家离线
	ErrorCodeIsFull    = 13 // 目标房间已满
	ErrorCodeUnPlay    = 14 // 玩家不在游戏
	ErrorCodeMaxFollow = 15 // 关注人数满
	ErrorCodeSelf      = 16 // 不能关注自己
	ErrorCodeVerify    = 17 // 验证失败
	ErrorCodeTooQuick  = 18 // 操作太快
	ErrorCodeBeBinded  = 19 // 邮箱已绑定
	ErrorCodeInValid   = 20 // 邮箱不合法
	ErrorCodeReLogin   = 21 // 其它地方登录
	ErrorCodeOutTeam   = 31 // 不在组队状态

	ErrorCodeMicBusy = 23 // 麦上有人
	ErrorCodeVoice   = 24 // 语音服务器出错
	ErrorCodeTeam    = 25 // 未找到组队服务器

	ErrorCodeInTeam   = 27 // 已在组队游戏中
	ErrorCodeTeamFull = 28 // 队伍已满
	ErrorCodeInvFail  = 29 // 组队邀请失效
	ErrorCodeLeader   = 30 // 只有队长才能点开始游戏

	ErrorCodeAuthErr = 35 // 验证码错误
	ErrorCodeMsgErr  = 36 // 短信服务器出错
	ErrorCodeUnBind  = 37 // 未绑定手机

	ErrorCodeTelErr = 38 // 手机号错误
	ErrorCodeIsBind = 41 // 手机号已被绑定
	ErrorCodeMaxMsg = 43 // 今天短信次数已用完
	ErrorCodeMaxErr = 44 // 密码错误已超过最大限制

	ErrorCodeNoWidget = 45 // 数量不足
	ErrorCodeBanAcc   = 53 // 账号被封
	ErrorCodeLevel    = 54 // 段位不够

	ErrorCodeMaxFile    = 57 // 文件太大
	ErrorCodeUploadFile = 58 // 文件上传失败
	ErrorCodeLauding    = 59 // 已经点赞
	ErrorCodeNotLaud    = 60 // 没有点赞

	ErrorCodeFast = 66 // 操作太快

	ErrorCodeFollowed   = 68 // 已关注
	ErrorCodeMaxInvalid = 69 // 小号关注太多

	ErrorCodeBlack  = 72 // 你已经被增加到黑名单
	ErrorCodeSeason = 73 // 赛季切换

	ErrorCodeBanListMaxErr = 76 // 增加黑名单数量太多
	ErrorCodeNewbieMemInv  = 77 // 萌新不能邀请
	ErrorCodeSetAccountTm  = 78 // 修改用户名

	ErrorCodeUnInTRoom     = 82  // 不在自建房间中
	ErrorCodeNotRoomOwner  = 83  // 自建房间不是房主
	ErrorCodeLittleUserNum = 84  // 自建房间人数不够
	ErrorCodeLittleRPriv   = 85  // 自建房间权限不够
	ErrorCodeTRNotEnd      = 86  // 自建房间还未结束
	ErrorCodeTJoinRoom     = 89  // 自建房间模式不能加入
	ErrorCodeTPassword     = 98  // 自建房间密码错误
	ErrorCodeTRoomPriv     = 99  // 自建房间权限不够加入失败
	ErrorCodeTExpQrCode    = 100 // 二维码过期
	ErrorCodeTTeamUser     = 101 // 自建房间人数或队伍数量错误
	ErrorCodeTKickUser     = 102 // 自建房间已被踢出房间

	ErrorCodeIllegalChars    = 133 // 存在非法字符
	ErrorCodeUserBeBanned    = 137 // 自己玩家被封
	ErrorCodeOtherBeBanned   = 138 // 该玩家被封
	ErrorCodeIllegalPosition = 139 // 位置信息非法

	ErrorCodeCloseLeaveMsg      = 150 // 封禁留言
	ErrorCodeCloseSuperLeaveMsg = 151 // 封禁大神留言
	ErrorCodeCloseSign          = 152 // 封禁签名
	ErrorCodeCloseBarrage       = 153 // 封禁弹幕
	ErrorCodeCloseForever       = 154 // 永久封禁

	ErrorCodeCharTooLong = 166 // 字符过长

	ErrorCodeAccountAwarded = 201 // 账号已经领过奖励

	ErrorCodeAddStarNum = 207 // 防止小号刷赞数量

	ErrorCodeChkWhiteList = 215 // 检查白名单
	ErrorCodeActiExpire   = 216 // 活动已结束
	ErrorCodeHaveReceive  = 217 // 已领取
	ErrorCodeCondErr      = 218 // 条件不满足
	ErrorCodeExistBind    = 219 // 已绑定帐号
	ErrorCodeExistChannl  = 220 // 渠道号已绑定

	//聊天相关错误码 start>>>>>>>>>>>>>>>>>>>>>>>>>>>
	ErrorCodeChatNotOwner       = 301 // 不是群主
	ErrorCodeChatNoRoom         = 302 // 群组不存在或已关闭
	ErrorCodeChatIsJoin         = 303 // 你已加入该群组
	ErrorCodeChatCannotCreate   = 304 // 你无权限创建此类房间
	ErrorCodeChatNumsMax        = 305 // 群组人数已达上限
	ErrorCodeChatKickSelf       = 306 // 不能将自己踢出群组
	ErrorCodeChatNotGroup       = 307 // 只有群组才能修改群名称
	ErrorCodeChatNeedFreind     = 308 // 对方已设置互关好友才能聊天
	ErrorCodeChatRedMaxLove     = 309 // 爱心数量填写不合法
	ErrorCodeChatRedReceive     = 310 // 你已领取该爱心包
	ErrorCodeChatJoinNeedFreind = 311 // 只能拉互关好友进群
	ErrorCodeChatCloseQuick     = 312 // 对方已关闭快捷聊天
	ErrorCodeChatRedRecSelf     = 313 // 不能领取自己的爱心包
	ErrorCodeChatCreateSucc     = 314 // 团战群组创建成功，可以去聊天了
	ErrorCodeChatTempNotSend    = 315 // 临时群组无法发放爱心包
	ErrorCodeChatUnreadMuch     = 316 // 对方未读信息过多，暂时无法发送
	ErrorCodeChatAlreadyCreate  = 317 // 队友已创建
	ErrorCodeChatFast           = 318 // 聊天发送频率过快
	ErrorCodeChatFollowClose    = 319 // 对方删除会话无法继续聊
	ErrorCodeChatClosePrivate   = 320 // 对方关闭私聊权限
	//聊天相关错误码 end>>>>>>>>>>>>>>>>>>>>>>>>>>>

	ErrorCodeNotLocation = 407 //没有位置

	//邮件的相关错误码
	ErrorCodeMailId     = 600 //邮件不存在
	ErrorCodeMailInfo   = 601 //邮件信息错误
	ErrorCodeMailDB     = 602 //DB出错相关
	ErrorCodeAddMail    = 603 //添加邮件失败
	ErrorCodeGetMail    = 604 //获取邮件失败
	ErrorCodeForbitRecv = 605 //禁止接收索要邮件
	ErrorCodeInsertDB   = 606 //插入DB失败
	ErrorCodeDeleteMail = 607 //删除出错

)

var (
	PhotoNumError = errors.New("photo num error") //已经照片数量太多
)

const (
	MIN_ICON_ID = 1000
	MAX_ICON_ID = 1011
)

// 照片状态
const (
	PHOTO_STATE_INIT    = 0 // 未审核
	PHOTO_STATE_PASS    = 1 // 通过审核
	PHOTO_STATE_NOPASS  = 2 // 未通过审核
	PHOTO_STATE_DEL     = 3 // 照片删除
	PHOTO_STATE_PENDING = 4 // 挂起状态
)

// 关注类型
const (
	FOLLOW_TYPE_NONE   = 1 // 互不关注
	FOLLOW_TYPE_FOLLOW = 2 // 关注
	FOLLOW_TYPE_FAN    = 3 // 粉丝
	FOLLOW_TYPE_FRIEND = 4 // 好友
)

// 玩家社区数据
type SnsUserData struct {
	Uid           uint64 `redis:"Uid"`           // 用户ID
	Level         uint32 `redis:"Level"`         // 等级
	Exp           uint32 `redis:"Exp"`           // 经验
	ReportNum     uint32 `redis:"ReportNum"`     // 被成功举报次数
	IsReportMsg   uint32 `redis:"IsReportMsg"`   // 是否有举报消息
	AudienceNum   uint32 `redis:"-"`             // 听众人数
	AudienceUrl   string `redis:"-"`             // 语音签名
	LoveNum       uint32 `redis:"-"`             // 被喜欢次数
	AudienceTime  uint32 `redis:"-"`             // 语音签名时间
	LoveUseNum    uint32 `redis:"-"`             // 能喜欢他的次数
	SoundId       uint32 `redis:"SoundId"`       // 社区声音
	IsForceFollow uint32 `redis:"IsForceFollow"` // 是否强制推荐0没有 1有
}

//玩家经纬度信息
type UserLbs struct {
	Longi float64 `redis:"Longi"` //经度
	Lati  float64 `redis:"Lati"`  //纬度
	City  string  `redis:"City"`
}

// 语音签名信息
type AudienceMsg struct {
	Uid          uint64 `bson:"_id"` //用户ID
	Url          string //地址
	AudienceTime int    //语音时间
	Time         int64  `bson:"uptime"` //上传时间
	Account      string `bson:"-"`      //帐号
}

// 举报数据
type ReportMsg struct {
	Id          bson.ObjectId     `bson:"_id"`    //ID
	ReportClass int32             `bson:"class"`  //举报大类
	ReportId    string            `bson:"rid"`    //举报ID
	ReportNum   uint32            `bson:"num"`    //举报数量
	ReportInfo  map[string]uint32 `bson:"info"`   //举报每种类型数量
	ReportTime  int64             `bson:"time"`   //举报时间
	Status      uint32            `bson:"status"` //举报状态 (未处理、已删除、白名单)
}

type ReportExtInfo struct {
	*ReportMsg
	Account      string      //被举报帐号
	Uid          uint64      //被举报用户ID
	AccReportNum uint32      //被举报帐号的被成功举报次数
	ReportData   interface{} //被举报资源信息
}

// 黑名单
type BanInfoMsg struct {
	Id    bson.ObjectId `bson:"_id"` //ID
	Uid   uint64        //用户ID
	ToId  uint64        //黑名单用户ID
	Time  int64         //时间
	Admin uint64
}

type BanInfoExt struct {
	*BanInfoMsg
	Account  string // 帐号
	Sex      uint32 // 性别
	Icon     uint32 // 图标
	PassIcon string // 已审核头像
	Remark   string // 备注
}

// 渠道类型
const (
	PAY_CHANNEL_ANDROID = 1
	PAY_CHANNEL_IOS     = 2
	PAY_CHANNEL_WP      = 3
)

// 订单信息
type RedisOrderData struct {
	UserId    uint64 `redis:"UserId"`
	Amout     uint32 `redis:"Amount"`
	Channel   uint32 `redis:"Channel"`
	ProductId string `redis:"ProductId"`
	State     uint32 `redis:"State"`
	Time      string `redis:"Time"`
}

const (
	SubType_Login   = "1" //订阅登陆类型
	SubType_Playing = "2" //订阅玩游戏类型
)

type BanLeaveMsgTTL struct {
	Reason uint32 //封禁原因
	Days   uint32 //封禁天数
	TTL    int64  //解封倒计时
}

//资源信息
type RedisResInfo struct {
	Resid uint32 `redis:"Resid"`
	Num   uint32 `redis:"Num"`
}

//-------------------------------------
type GiftItem struct {
	GiftId  uint32 //礼品ID
	GiftNum uint32 //礼品数量
}

//系统邮件
type SystemMailInfo struct {
	Id         string `bson:"_id"` //objID -- mailid
	UserId     uint64 `bson:"uid"` //userid
	Type       uint8  `bson:"type"`
	Title      string `bson:"title"`
	Content    string `bson:"content"`
	HasAttr    bool   `bson:"hasAttr"`
	Attr       string `bson:"attr"` //格式 10001,5|10002,5|#
	CreateTime int64  `bson:"createTime"`
}

//好友邮件
type SingleMailInfo struct {
	Id         string `bson:"_id"` //objID -- mailid
	Type       uint8  `bson:"type"`
	Title      string `bson:"title"`
	Content    string `bson:"content"`
	FromId     uint64 `bson:"fromId"`
	ToId       uint64 `bson:"toId"`
	HasAttr    bool   `bson:"hasAttr"`
	Attr       string `bson:"attr"`
	State      uint8  `bson:"state"`
	CreateTime int64  `bson:"createTime"`
}

type FriendMailItem struct {
	SInfo    *SingleMailInfo
	Account  string
	Icon     uint32
	PassIcon string
	Sex      uint8
}

//索要消息
type RequestMessage struct {
	Id   string `bson:"_id"` //objID -- mailid
	Type uint8  `bson:"type"`
	//Title      string `bson:"title"`
	//Content    string `bson:"content"`
	FromId     uint64 `bson:"fromId"`     //索要人
	ToId       uint64 `bson:"toId"`       //被索要人
	State      uint8  `bson:"state"`      //是否读取
	ItemId     uint32 `bson:"itemId"`     //道具的编号
	RespState  uint8  `bson:"respState"`  //响应状态
	CreateTime int64  `bson:"createTime"` //创建时间
}

type RecvApplyState struct {
	Uid uint64 `bson:"uid"` //
}

type ReqMail struct {
	Msg      *RequestMessage
	Account  string
	Icon     uint32
	PassIcon string
	Sex      uint8
}

//单一系统邮件信息
type SysInfo struct {
	MailId string //邮件id
	State  uint8  //状态 是否领取或者阅读
}

//玩家系统邮件
type UserSysMailInfo struct {
	Id         string
	Type       uint8
	UserId     uint64
	Title      string
	Content    string
	HasAttr    bool
	Attr       string
	CreateTime int64
	State      uint8
}

//未读邮件数量
type NonReadMailInfo struct {
	Friend_UnRead     uint8 //好友邮件未读
	SysMail_UnRead    uint8 //系统邮件未读
	Request_UnRead    uint8 //索要邮件未读
	ForbidApply_State uint8 //是否禁止索要邮件
}

//玩家的系统邮件信息
type SingleSysInfo struct {
	Id    string     `bson:"_id"`   //objID -- mailid
	Uid   uint64     `bson:"uid"`   //userid
	Infos []*SysInfo `bson:"infos"` //系统邮件信息
}
