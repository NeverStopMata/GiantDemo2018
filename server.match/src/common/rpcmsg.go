package common

import (
	"gopkg.in/mgo.v2/bson"
)

type ReqGetUserSumInfo struct {
	UID uint64
}

type RspGetUserSumInfo struct {
	UID    uint64
	EMoney uint32
	Money  uint32
	Ticket uint32
}

// 创建设备
type ReqNewDevice struct {
	DevCode   string // 设备号
	Device    string // 设备名
	LoginType uint32 // 登录类型1为新方式登录
	Location  uint32 // 位置
	Emulator  int    // 是否模拟器
	RegIp     string // 玩家ip
	//Channl    int    // 渠道号
}
type RetNewDevice struct {
	Id       uint64 // 账号id
	Account  string // 账号
	IsNewbie bool   // 是否新玩家
	Password string // 密码
}

// 基本数据
type ReqBaseInfo struct {
	UserId uint64
}
type RetBaseInfo struct {
	Udata *RedisBaseInfo
}

// 根据设备号获取数据
type ReqUserByDev struct {
	DevCode string
}

// 根据账号获取数据
type ReqUserByAcc struct {
	Account string
}

// 根据手机号获取数据
type ReqUserByTel struct {
	Tel string
}

// 根据id获取数据
type ReqUserById struct {
	Id uint64
}

// 详细数据
type ReqDetailData struct {
	ToId   uint64
	FromId uint64
}
type RetDetailData struct {
	RetUserInfo
}

type ReqSnsDetailData struct {
	FromId uint64
	ToId   uint64
}
type RetSnsDetailData struct {
	RetSnsUserInfo
}

// 喜欢事件
type ReqLoveState struct {
	FromId uint64 //用户ID
	ToId   uint64 //被喜欢用户ID
}

type RetLoveState struct {
	ErrCode uint32 //状态码
	LoveNum uint32 //喜欢次数
	ShowPos uint8
}

// 设置语音签名
type ReqSetAudience struct {
	Uid          uint64 //用户ID
	Url          string //地址
	AudienceTime int    //语音时间
	Time         int64  //上传时间
	Kind         int    // 类型
}
type RetSetAudience struct {
	Ret    bool
	OldUrl string
}

type BanMsgInfo struct {
	BanType string
	BanTime int64
}

// 返回数据
type RetUserInfo struct {
	Id           uint64 // 玩家id
	Account      string // 用户名(账号)
	DevCode      string // 设备号
	Email        string // 邮箱
	Tel          string // 手机
	Password     string // 密码
	AState       uint32 // 账号状态
	Icon         uint32 // 玩家icon
	NewIcon      string // 头象
	Sign         string // 签名
	PassIcon     string
	AudienceNum  uint32 // 听众人数
	AudienceTime uint32 // 语音时间
	AudienceUrl  string // 语音签名
	Money        uint32 // 游戏币
	JiangGuo     uint32 // 小浆果
	Ticket       uint32 // 入场券
	Sex          uint8  // 性别
	Age          uint32 // 年龄
	PlayNum      uint32 // 总局数
	Location     uint32 // 位置
	ShowPos      uint8  // 显示地区
	State        string // 动态
	Position     string // 玩家位置详情
	OnlineTime   int64
	Level        uint32 // 段位
	Scores       uint32 // 星数
	MaxLevel     uint32 // 历史最大段位
	MaxScores    uint32 // 历史最大星数
	HideScores   uint32 // 隐藏经验分数
	Nick         string // 渠道昵称
	Sceneid      uint32
	ActCodeGift  uint32 // 微信公众号激活码奖励
	RegIp        string
	RegTime      uint32
	RegDev       string
	Speaker      uint32 // 小喇叭
	OpenId       string // 渠道帐号
	LoginTime    uint32 // 每天第一次登录时间
	IdCard       string
}

// 返回社交数据
type RetSnsUserInfo struct {
	Id           uint64 // 玩家id
	Account      string // 用户名(账号)
	NickName     string // 渠道昵称
	Icon         uint32 // 玩家icon
	Sign         string // 签名
	Sex          uint8  // 性别
	Age          uint32 // 年龄
	Level        uint32 // 段位
	Scores       uint32 // 星数
	PassIcon     string // 上一次头像
	AudienceNum  uint32 // 听众人数
	AudienceTime uint32 // 语音时间
	AudienceUrl  string // 语音签名
	IsBanList    bool   // 是否黑名单
	Location     uint32
	Position     string // 玩家位置详情
	IsFollow     bool   // 是否关注
	Follows      uint32 // 关注数
	Fans         uint32 // 粉丝数
	Friends      uint32 // 好友数
}

// 根据账号获取数据
type ReqAccInfo struct {
	Account string
}
type RetAccInfo struct {
	Ret      bool
	Account  string
	Location uint32
	Adata    *RedisAccount
}

type RetStatus struct {
	Status uint32
}

type RetStatuss struct {
	Id     []uint64
	Status []uint32
}

// 根据设备号获取数据
type ReqDevInfo struct {
	DevCode string
}
type RetDevInfo struct {
	UserId uint64
	Ret    bool
}

// 检查账号
type ReqCheckIp struct {
	Ip string
}

// 检查账号
type RetCheckIp struct {
	Count uint32
	Ret   bool
}

// 修改账号
type ReqSetAccount struct {
	UserId  uint64
	Account string
}
type RetSetAccount struct {
}

// 检查是否绑定
type ReqIsBindId struct {
	UserId uint64
}
type RetIsBindId struct {
	Ret bool
}

// 检查玩家账号
type ReqIsExistAcc struct {
	UID     uint64
	Account string
}
type RetIsExistAcc struct {
	Ret    bool
	OkName []string
}

// 检查玩家id
type ReqIsExistId struct {
	UserId uint64
}
type RetIsExistId struct {
	Ret bool
}

// 检查手机号
type ReqIsExistTel struct {
	Tel string
}
type RetIsExistTel struct {
	Ret bool
}

// 修改手机号
type ReqSetTel struct {
	UserId uint64
	Tel    string
	PassWd string
}
type RetSetTel struct {
}

// 检查密码
type ReqCheckPasswd struct {
	UserId uint64
	Passwd string
}
type RetCheckPasswd struct {
	ErrorCode uint32
}

// 修改手机号
type ReqGetTel struct {
	UserId uint64
}
type RetGetTel struct {
	Tel string
}

// 修改密码
type ReqSetPasswd struct {
	UserId uint64
	Passwd string
	OldPw  string
	NewPw  string
}
type RetSetPasswd struct {
	ErrorCode uint32
}

type ReqSetPosition struct {
	UserId   uint64
	Position string
}
type RetSetPosition struct {
}

type ReqGetPositions struct {
	UIDs []uint64
}

type UserPostionDesc struct {
	UID  uint64
	Desc string
}

type RetGetPostions struct {
	PostionsDesc []*UserPostionDesc
}

// 修改头象
type ReqSetIcon struct {
	UserId  uint64
	Icon    uint32
	PhotoId string
}
type RetSetIcon struct {
	PassIcon string
}

// 修改签名
type ReqSetSign struct {
	UserId uint64
	Sign   string
}
type RetSetSign struct {
	ErrorCode uint32
	BanTTL    *BanLeaveMsgTTL
}

// 修改邮箱
type ReqSetEmail struct {
	UserId uint64
	Email  string
}
type RetSetEmail struct {
}

// 修改年龄
type ReqSetAge struct {
	UserId uint64
	Age    uint32
}
type RetSetAge struct {
}

// 修改性别
type ReqSetSex struct {
	UserId uint64
	Sex    uint8
}
type RetSetSex struct {
}

// 添加关注
type ReqAddFollow struct {
	FromId uint64
	ToId   uint64
	FType  string
}
type RetAddFollow struct {
	MaxFollow uint32
	RetCode   uint32
	IsFriend  bool
	Account   string
	Near      bool
	Icon      uint32
	PassIcon  string
}

// 添加关注
type ReqBatchAddFollow struct {
	FromId uint64
	UIDs   []uint64
	TypeId uint32
}
type RetBatchAddFollow struct {
	UIDs     []uint64
	RelType  []uint32
	Account  string
	Icon     uint32
	PassIcon string
}

// 取消关注
type ReqUnFollow struct {
	FromId uint64
	ToId   uint64
}
type RetUnFollow struct {
}

// 请求点赞
type ReqAddStars struct {
	ToId uint64
}
type RetAddStars struct {
}

// 修改账号状态
type ReqAccState struct {
	Account string
	State   uint32
}
type RetAccState struct {
}

// 请求比赛列表
type ReqMatchList struct {
	UserId uint64
}
type RetMatchList struct {
	MDatas []*MatchUser
}

// 请求关注列表
type ReqFollowList struct {
	UserId uint64
	Page   string
}
type RetFollowList struct {
	FDatas []*FollowUser
}

// 请求粉丝列表
type ReqFanList struct {
	UserId uint64
	Page   string
	TypeId uint32
}
type RetFanList struct {
	FDatas []*FanUser
}

// 请求荣誉列表
type ReqHonorList struct {
	UserId uint64
}
type RetHonorList struct {
	HDatas []*HonorUser
}

// 统计数据
type ReqDynamicMsg struct {
	UserId uint64
}
type RetDynamicMsg struct {
	RelationData  *RelationData
	RedPointData  *RedPointData
	SnsDynamicMsg *SnsDynamicMsg
	Reward        uint32
}

// 统计数据
type ReqRelationStat struct {
	UserId uint64
}
type RetRelationStat struct {
	RData *RelationData
}

type SetLoginTime struct {
	UserId    uint64
	LoginTime uint32
}

type RetSetLoginTime struct {
}

// 红点数据
type ReqRedPoints struct {
	UserId uint64
}
type RetRedPoints struct {
	RData *RedPointData
}

// 检查关注
type ReqCheckFollows struct {
	FromId uint64
	ToIds  []uint64
}
type FollowFlag struct {
	Id       uint64
	IsFollow bool
	FType    int
}
type RetCheckFollows struct {
	Follows []*FollowFlag
}

// 添加好友
type ReqAddFriend struct {
	FromId uint64
	ToId   uint64
}
type RetAddFriend struct {
}

// 好友id列表
type ReqFriendIdList struct {
	UserId uint64
}
type RetFriendIdList struct {
	Friends []uint64
}

// 好友列表
type ReqFriendList struct {
	UserId uint64
	//Page   string
}
type RetFriendList struct {
	Friends []*FriendUser
}

// 检查好友
type ReqCheckFriends struct {
	FromId uint64
	ToIds  []uint64
}
type FriendFlag struct {
	Id       uint64
	IsFriend bool
}
type RetCheckFriends struct {
	Friends []*FriendFlag
}

// 检查关注
type ReqCheckFollow struct {
	FromId uint64
	ToId   uint64
}

type RetCheckFollow struct {
	IsFriend bool
}

type RoomIncData struct {
	Id        uint64 // 玩家id
	IsEndRoom bool   // 是否房间结束结算
	EndEatNum uint32 // 吞噬球数 (结算)
	EndWeight uint64 // 结束体重 (结算)
	KillNum   uint32 // 吞噬人数 (结算)
	CState    string // 自由合作状态 (结算)
	TCState   string // 团队合作状态 (结算)
	HScores   uint32 // 自由隐藏分 (结算)
	THScores  uint32 // 团队隐藏分 (结算)
	SeasonId  uint32 // 赛季 (结算)
	Level     uint32 // 段位 (结算)
	Scores    uint32 // 分数 (结算)
	IsTMvp    bool   // 是否为mvp
	Robot     uint32 // 机器人级别 (结算)
	NickName  string // 昵称 (比赛列表)
	Icon      uint32 // 图标 (比赛列表)
	Rank      uint32 // 排名 (比赛列表)
	JiangGuo  uint32 // 添加的小浆果(比赛列表)
	SceneId   uint32 // 游戏模式(比赛列表)
	AnimalID  uint32 // 对应的动物(比赛列表)
	TAnimal   uint32 // 达到的最高动物
	TopScore  uint32 // 达到的最高分数
	Location  uint32 // 地理位置
	AddExp    uint32 // 结算经验值(比赛列表)
	GameTime  uint32 // 房间游戏时间
}

type ReqRoomInc struct {
	UData    []*RoomIncData
	TData    [][]uint64 // 队伍成员
	RoomType uint32     // 房间类型
	RoomId   uint32     // 房间id
}
type RetRoomInc struct {
}

// 闪电战结束
type ReqEndQRoom struct {
	UserId     uint64
	RoomType   uint32
	RoomId     uint32
	EatNum     uint32
	KillNum    uint32
	DropId     uint32
	Relivesec  uint32
	AddExp     uint32 // 增加经验值
	AddExpText string // 增加经验值理由
	GameTime   uint32 // 游戏时间
}
type RetEndQRoom struct {
}

// 排名列表
type ReqRankList struct {
	FromId uint64 // 请求者id
	RType  int    // 类型(吞噬，粉丝,赞)
	RSort  int    // 日，周，月
	Page   string
}
type RetRankList struct {
	RDatas []*RankUser
}

// 排行数据
type LRankUser struct {
	Id       uint64 // 玩家id
	Rank     uint32 // 排名
	Account  string // 帐号
	Icon     uint32 // 图标
	Sex      uint8  // 性别
	PassIcon string //头像
	Level    uint32 // 段位
	Stars    uint32 // 星数
	MaxLevel uint32 // 最高段位
	MaxStars uint32 // 最高分数
	IsFollow bool
}

type ReqLRankList struct {
	FromId uint64 //玩家id
}
type RetLRankList struct {
	RDatas []*LRankUser
}

// 扣除数值
type ReqSubMoney struct {
	Id    uint64 // 玩家id
	MType uint32 // 类型
	MNum  uint32 // 数量
	Text  string // 注释
	Event string // 事件
}
type RetSubMoney struct {
}

// 添加数值
type ReqAddMoney struct {
	Id    uint64 // 玩家id
	MType uint32 // 类型
	MNum  uint32 // 数量
	Text  string // 注释
	Event string // 事件
}
type RetAddMoney struct {
}

// 检查数值
type ReqCheckMoney struct {
	Id    uint64 // 玩家id
	MType uint32 // 类型
	MNum  uint32 // 数量
}
type RetCheckMoney struct {
}

type ReqAddHideExp struct {
	UserId     uint64
	AddHideExp uint32
}
type RetAddHideExp struct {
}

type RetUpdateSeasonCourse struct {
}
type ReqUpdateSeasonCourse struct {
	UserId   uint64
	CType    uint32
	KillNum  uint32
	MaxExp   uint32
	SceneId  uint32
	Animalid uint32
	Rank     uint32
	Fans     uint32
}

// 检查赛季
type ReqCheckSeason struct {
	Id uint64 // 玩家id
}
type RetCheckSeason struct {
	Ret bool
}

type ReqGetRegDev struct {
	Id uint64 // 玩家id
}
type RetGetRegDev struct {
	RegDev string
}

// 赛季结算
type GReward struct {
	OldSeason uint32 // 赛季编号
	OldLevel  uint32 // 老段位
	OldScores uint32 // 老星数
	NewLevel  uint32 // 新段位
	NewScores uint32 // 新星数
	EMoney    uint32 // 奖励棒棒糖数量
	Medal     uint32 // 荣誉
	IconBox   uint32 // 头象框
	Objs      []AddObjNode
}

type ReqUpdateSeason struct {
	Id uint64 // 玩家id
}
type RetUpdateSeason struct {
	Reward *GReward
	Ret    bool
}

type AddObjNode struct {
	Id    uint32
	Num   uint32
	Money uint32
}

// 设置在线状态
type ReqSetOnlineData struct {
	UserId   uint64
	ClientId string
	Device   string
	MKey     string
	Emulator int
	Location uint32
	CityId   uint32
}
type RetSetOnlineData struct {
	TimeNow int64
}

// 获取在线状态
type ReqGetOnlineData struct {
	UserId uint64
}
type RetGetOnlineData struct {
	ClientId string
	Device   string
}

type RetGameDayTime struct {
	GameDayTime uint32
}

type SetCrossUser struct {
	ID  uint64
	Val string
}

// 添加备注
type ReqAddRemark struct {
	FromId uint64
	ToId   uint64
	Mark   string
}
type RetAddRemark struct {
}

// 结算时同步玩家数据
type SyncData struct {
	UserId   uint64
	SeasonId uint32
	Level    uint32
	Scores   uint32
	IsValid  bool
}
type ReqSyncDatas struct {
	UserIds []uint64
}
type RetSyncDatas struct {
	SData []SyncData
}

// 获取玩的次数
type ReqPlayNum struct {
	UserId uint64
}
type RetPlayNum struct {
	PlayNum uint32
}

// 获取注册ip
type ReqRegIp struct {
	UserId uint64
}
type RetRegIp struct {
	RegIp string
}

// 队友数据
type TeamMateUser struct {
	Id       uint64 // 玩家id
	Account  string // 帐号
	Icon     uint32 // 图标
	Sex      uint8  // 性别
	Level    uint32 // 段位
	Scores   uint32 // 星数
	PassIcon string // 头像
}
type ReqTeamMates struct {
	UserId uint64
}
type RetTeamMates struct {
	TDatas []*TeamMateUser
}

// vip玩家信息
type VipUser struct {
	Id       uint64
	Account  string
	Icon     uint32
	PassIcon string
	Sex      uint8
	Level    uint32
	Scores   uint32
}
type ReqVipUser struct {
	UserId uint64
}
type RetVipUser struct {
	Info *VipUser
}

/////////////////// 语音相关 /////////////////////
// 创建房间
type ReqNewRoom struct {
	RoomId   uint32 // 房间id
	LastTime uint32 // 持续时间
}
type RetNewRoom struct {
}

// 当前发言人
type ReqUserSpeak struct {
	RoomId   uint32 // 房间id
	UserId   uint64 // 玩家id
	TeamId   uint32 // 队伍id
	LastTm   uint32 // 持续时间
	BanVoice bool
	NewLogin map[uint64]bool
}
type RetUserSpeak struct {
}

/////////////////// 房间分配相关 开始 /////////////////////
// 获取房间地址
type ReqByRoomId struct {
	RoomType uint32
	RoomId   uint32
}
type RetByRoomId struct {
	WAddress string //观看服务器地址，在观看请求中返回
	Address  string
	UserNum  int32
	EndTime  uint32
	UnCoop   uint32
	SceneId  uint32
}

// 分配自由房间
type ReqFreeRoom struct {
	UserId  uint64 // 玩家id
	IsNew   bool   // 是否新手
	IsCoop  bool   // 是否合作
	HScores uint32 // 隐藏分
	UnCoop  uint32 // 是否非合作
	SceneId uint32 // 场景id
	Level   uint32 // 段位
}
type RetFreeRoom struct {
	ServerId uint16
	Address  string
	RoomId   uint32
}

// 分配组队房间
type ReqTeamRoom struct {
	TServerId uint16
	RLevel    uint32
	TModel    uint32
}
type RetTeamRoom struct {
	Address string
	RoomId  uint32
}

// 分配闪电战房间
type ReqQuickRoom struct {
	UserId  uint64
	IsNew   bool // 是否新手
	SceneId uint32
	Level   uint32 // 段位
	HScores uint32 // 隐藏分
}

type RetQuickRoom struct {
	ServerId uint16
	Address  string
	RoomId   uint32
}

// 获得在线人数
type ReqOnlineNum struct {
}

type RetOnlineNum struct {
	OnlineNum uint32
}

// 分配自建房间
type ReqCustomRoom struct {
}
type RetCustomRoom struct {
	ServerId uint16
	RAddress string
	WAddress string
	RoomId   uint32
}

// 分配OB房间
type ReqObRoom struct {
}
type RetObRoom struct {
	ServerId uint16
	RAddress string
	WAddress string
	RoomId   uint32
}

// 获取剩余时间(组队)
type ReqRoomTTL struct {
	RoomType uint32
	RoomId   uint32
}
type RetRoomTTL struct {
	TTL uint32
}

// 获取组队服(组队)
type ReqTServerId struct {
	RoomId uint32
}
type RetTServerId struct {
	TServerId uint16
}

/////////////////// 房间分配相关 结束 /////////////////////

/////////////////// 组队分配相关 开始 /////////////////////
// 请求网关地址
type ReqGateAddr struct {
	ServerId uint16
}
type RetGateAddr struct {
	ServerId uint16
	Address  string
}

type ReqGatewayList struct {
	Type    uint32
	Address string
}
type RetGatewayList struct {
	Servers []string
	Id      uint32
}

// 获取组队列表
type ReqTServerList struct {
}
type RetTServerList struct {
	TotalNum  uint32
	ServerIds []uint16
}

// 转发给其它服
type ReqS2SType struct {
	Type uint32 // 服务器类型
	Flag byte   // 压缩标识
	Data []byte // 数据
}
type RetS2SType struct {
}

// 转发给其它服的客户端
type ReqS2SClient struct {
	ServerId uint16 // 服务器id
	UserId   uint64 // 玩家id
	Flag     byte   // 压缩标识
	Data     []byte // 数据
}
type RetS2SClient struct {
}

// 批量转发给其它服的客户端
type ReqS2SClients struct {
	UserIds []uint64 // 玩家id
	Flag    byte     // 压缩标识
	Data    []byte   // 数据
}
type RetS2SClients struct {
}

// 批量转发邮件信息给其它服的客户端
type ReqS2SMailClients struct {
	Flag byte   // 压缩标识
	Data []byte // 数据
}
type RetS2SMailClients struct {
}

/////////////////// 组队分配相关 结束 /////////////////////

type ReqPhotoInfoList struct {
	UserId uint64
	FormId uint64
}

type RetPhotoInfoList struct {
	ErrorCode int
	Data      []*PhotoExtInfo
}

type ReqPhotoData struct {
	UserId  uint64
	PhotoId string
}

type RetPhotoData struct {
	ErrorCode uint32
	Data      uint32
	UID       uint64
}

type RetMailData struct {
	ErrorCode uint32
	Ids       []*Mailids
}

type ReqAddUserExp struct {
	UserId uint64
	Exp    uint32
	Text   string
}

type ReqAddReport struct {
	Id    bson.ObjectId `bson:"_id"` //ID
	Uid   uint64        // 举报人
	Class uint32        // 举报大类
	RId   string        // 举报ID
	RType uint32        // 举报类型
	Time  int64         // 时间
}

type RetAddReport struct {
	ErrorCode uint32
}

type ReqIsReport struct {
	Class uint32 // 举报大类
	RId   string // 举报ID
	Uid   uint64 // 用户ID
}
type RetIsReport struct {
	Ret bool
}

type ReqGetReportMsg struct {
	UserId uint64
}

type RetGetReportMsg struct {
	Datas []*ReportUserMsg
}

type ReqBanInfoList struct {
	UserId uint64
}

type RetBanInfoList struct {
	Data []*BanInfoExt
}

type ReqDelBanInfo struct {
	Uid   uint64 //用户ID
	ToId  uint64 //黑名单用户ID
	Admin uint64
}

type ReqUserId struct {
	Uid uint64 //用户ID
}

type ReqUserIdList struct {
	UsersId []uint64 //用户ID
}

type ReqGetSnsInfo struct {
	FromId uint64
	ToId   uint64
}

type ReqIsBanInfo struct {
	UserIds []uint64 //用户ID
	ToId    uint64   //黑名单ID
}

type RetIsBanInfo struct {
	Data map[uint64]int //状态码
}

type ReqGetUserTalkInfo struct {
	UserId   uint64 //用户ID
	ToUserId uint64 //黑名单ID
	TypeId   uint32 //私聊类型
}

type RetGetUserTalkInfo struct {
	IsBlack  bool //是否拉黑
	IsFriend bool //是否通过互关好友才能聊天
	IsQuick  bool //是否通过快捷才能聊天
	IsUnread bool //是否通过未读聊天
	IsFans   int8 //是否通过粉丝聊天 0通过 1被关注方关闭聊天
	IsBan    bool //是否被封禁
}

type ReqCheckPrivatePower struct {
	UserId   uint64 //用户ID
	ToUserId uint64 //对方ID
}

type RetCheckPrivatePower struct {
	Succ  bool //是否可以聊天
	IsBan bool //是否被封禁
}

type ReqFansTalkClose struct {
	UserId   uint64 //用户ID
	ToUserId uint64 //对方ID
}
type RetFansTalkClose struct {
	Succ bool
}

//GEO相关//////////////////////////
type ReqUpdateLocationInfo struct {
	UID          uint64 //账号id
	Account      string //账号
	Icon         uint32
	Sex          uint8
	City         string
	Desc         string
	PassIcon     string
	Latitude     float64
	Longitude    float64
	State        uint32
	Level        uint32 // 当前段位
	Scores       uint32 // 当前星数
	NewPlayer    int
	Age          uint32
	Sign         string
	AudienceUrl  string
	AudienceTime uint32
	UpdateType   int
}

type RspUpdateLocationInfo struct {
}

type ReqGetUserNearby struct {
	City      string
	UID       uint64
	Latitude  float64
	Longitude float64
	Skip      int
	Limit     int
	Time      uint64
	Del       bool
}

type RspGetUserNearby struct {
	ErrCode  uint32
	MyLongi  float64
	MyLati   float64
	UserList []*GeoUser
}

type ReqSetLastPlayTime struct {
	UID uint64
}

type RspSetLastPlayTime struct {
}

type ReqAddMessage struct {
	UID     uint64
	Itemid  uint32
	ItemNum uint32
}

type ReqMessage struct {
	UID uint64
}

type RspMessage struct {
	Itemid  uint32
	ItemNum uint32
}

type ReqAddIcon struct {
	UID    uint64
	Iconid uint32
}

type ReqIcons struct {
	UID uint64
}

type RspIcons struct {
	Icons []uint32
}

type ReqQuest struct {
	UID uint64
}

type RspQuset struct {
	Flag     bool
	Complete bool
}

type ReqSetNick struct {
	UID  uint64
	Nick string
}

type RspSetNick struct {
	Flag bool
}

type ReqSetChanllAcc struct {
	UID    uint64
	Acc    string
	Openid string
	Channl uint32
}

type RspSetChanllAcc struct {
	Flag bool
}

type ReqAddWhiteList struct {
	Acc string
}

type ReqWhiteList struct {
}

type RspWhitList struct {
	Accs []string
}

type RspSceneId struct {
}

type ReqSceneId struct {
	ID      uint64
	Sceneid uint32
}

type ReqCharBaseEvent struct {
	ID uint64
}
type RspCharBaseEvent struct {
}

type ReqSetActCodeGift struct {
	ID uint64
}
type RspSetActCodeGift struct {
}

//日志
type RspWriteULog struct {
}

type ReqWriteULog struct {
	Info string
}

type NickName struct {
	Id   uint64 // ID
	Name string //玩家昵称
}

type ReqChatCreate struct {
	Owner        uint64      //房主ID
	Name         string      //群名称
	Type         int32       //群组类型
	UserList     []uint64    //成员列表 玩家ID
	LastTime     uint32      //时效
	NickNameList []*NickName //昵称
	OnlyStr      string
}

type RetChatCreate struct {
	RoomId uint64
}

//解散聊天房间
type ReqChatDelete struct {
	RoomId     uint64
	UserIdList []uint64
}

type RetChatDelete struct {
	Succ bool
}

//加入聊天房间
type ReqChatJoin struct {
	RoomId     uint64
	TypeId     int32
	UserIdList []uint64
}

type RetChatJoin struct {
	Succ bool
}

//离开聊天房间
type ReqChatLeave struct {
	RoomId   uint64
	UserId   uint64
	NewOwner uint64
}

type RetChatLeave struct {
	Succ bool
}

//踢人
type ReqChatKick struct {
	RoomId     uint64
	UserIdList []uint64
}

type RetChatKick struct {
	Succ bool
}

//聊天
type ReqChatTalk struct {
	UserId uint64
	RoomId uint64
	Text   string
}

type RetChatTalk struct {
	Succ bool
}

//私聊
type ReqChatPrivate struct {
	UserId   uint64
	ToUserId uint64
	Text     string
}

type RetChatPrivate struct {
	Succ bool
}

//成为互关好友
type ReqChatFriend struct {
	UserId   uint64
	ToUserId uint64
}

type RetChatFriend struct {
	Succ bool
}

//获取聊天房间列表
type ReqChatStart struct {
	ChatId uint32
}

type RetChatStart struct {
	ChatReds []*ChatRedEnd
}

//获取聊天房间信息
type ReqChatGetRoom struct {
	RoomId uint64
}

type RetChatGetRoom struct {
	ChatRoom *ChatRoom
	RedList  map[uint64]*ChatRed
}

//获取玩家聊天展示信息列表
type ReqUserShows struct {
	UserList []uint64
}

type RetUserShows struct {
	UserList []*UserShow
}

//获取玩家聊天展示信息
type ReqUserShow struct {
	UserId uint64
}

type RetUserShow struct {
	UserShow *UserShow
}

//设置聊天同步数据
type ReqChatSynTalk struct {
	UserIdList []uint64
	TalkId     uint64 //聊天ID
	UserId     uint64 //玩家ID
	RoomId     uint64 //房间ID
	Text       string //内容
	Time       uint32 //时间
}

type RetChatSynTalk struct {
	Succ bool
}

//个推
type ReqSendTalkNotice struct {
	UserIdList []uint64
	UserId     uint64 //玩家ID
	RoomId     uint64 //房间ID
	Text       string //内容
}

type RetSendTalkNotice struct {
	Succ bool
}

//聊天设置
type ReqChatSynSet struct {
	UserId   uint64 //玩家ID
	DoUserId uint64 //操作玩家ID
	RoomId   uint64 //房间ID
	TypeId   int32  //类型
	IsOpen   bool   //是否开启
}

type RetChatSynSet struct {
	Succ bool
}

//聊天设置
type ReqChatPersonalSet struct {
	UserId uint64 //玩家ID
	TypeId int32  //类型
	IsOpen bool   //是否开启
}

type RetChatPersonalSet struct {
	Succ bool
}

//获取离线同步数据
type ReqChatSynGet struct {
	UserId uint64 //玩家ID
}

type RetChatSynGet struct {
	Rooms       []uint64            // 拥有房间ID列表
	Talks       []*RetTalk          // 聊天内容
	Sets        []*ChatSynSet       // 设置
	PersonalSet *ChatSynPersonalSet // 设置
}

//获取离线同步数据
type TmpChatRoomTimeSyn struct {
	RoomId  uint64
	TalkId  uint64
	EndTime uint32
}

type ReqChatRoomTimeSyn struct {
	RoomList []*TmpChatRoomTimeSyn
}

type RetChatRoomTimeSyn struct {
	Succ bool
}

//修改房间名称
type ReqChatChangeName struct {
	RoomId uint64
	Name   string
}

type RetChatChangeName struct {
	Succ bool
}

//发送爱心包
type ReqChatRedSend struct {
	ChatRed *ChatRed
}

type RetChatRedSend struct {
	Succ bool
}

// 领取爱心包
type ReqChatRedReceive struct {
	SendId uint64 //用户ID
	UserId uint64 //被喜欢用户ID
	TalkId uint64
	RoomId uint64
}

type RetChatRedReceive struct {
	ErrId     int
	LessMoney uint32
	AddLoves  uint32
	EndTime   uint32
}

// 爱心包过期
type ReqChatRedEnd struct {
	RoomId uint64 //房间ID
	TalkId uint64 //爱心包ID
}

type RetChatRedEnd struct {
	Succ bool
}

////////////////mgr服务器相关
//获取指定类型节点连接地址
type ReqMgrNodeList struct {
	SerType      int
	Address      string
	NodeTypeList []int
}

type RetMgrNodeList struct {
	NodeTypeList map[int]string
}

////////////////mgr服务器相关 end

//----------------------------------------------
type ReqGeoUserList struct {
	UserId uint64
}

//雷达数据
type GeoUser struct {
	Id        uint64 //玩家id
	Account   string //账号
	Icon      uint32
	Sex       uint32
	Age       uint32
	IsFriend  bool
	PassIcon  string
	Distance  uint32  //m
	Longitude float64 //经度
	Latitude  float64 //纬度
	Level     uint32
	Scores    uint32
	UserState uint32
}
type RetGeoUserList struct {
	Users []*GeoUser
}

//开启雷达
type ReqOpenRadar struct {
	UserId    uint64
	Longitude float64
	Latitude  float64
}

type RetOpenRadar struct {
}

//关闭雷达
type ReqCloseRadar struct {
	UserId uint64
}

type RetCloseRadar struct {
}

//-------------------------------------------------
//摇一摇
//开启
type ReqOpenShake struct {
	UserId    uint64
	Longitude float64
	Latitude  float64
}

type RetOpenShake struct {
}

//搜索
type ReqShakeUserList struct {
	UserId uint64
}

type RetShakeUserList struct {
	Users []*GeoUser
}

//更新
type ReqUpdateShake struct {
	UserId uint64
}

type RetUpdateShake struct {
}

//删除
type ReqDelShake struct {
	UserId uint64
}

type RetDelShake struct {
}

//请求资源信息
type ReqResInfo struct {
	UserId uint64
}

//返回资源信息
type RetResInfo struct {
	ResList []*RedisResInfo
}

//添加资源
type ReqAddRes struct {
	UserId uint64
	Resid  uint32
	Num    uint32
}

//添加资源返回
type RetAddRes struct {
	Ret uint32
}

//请求设置资源信息
type ReqSetResInfo struct {
	UserId  uint64
	ResList []*RedisResInfo
}

//返回设置资源信息
type RetSetResInfo struct {
	Ret uint32
}

//请求解锁列表
type ReqUnlockList struct {
	UserId uint64
}

//返回解锁列表
type RetUnlockList struct {
	UnList []uint32
}

type ReqSetIdCard struct {
	UID    uint64
	IdCard string
	Name   string
}

type RspSetIdCard struct {
	Status bool
}

//------------------------------------------邮件相关---------------
//增加系统邮件
type ReqAddSystemMail struct {
	MailInfo *SystemMailInfo
}

type RetAddSystemMail struct {
	Ret uint32
}

//增加个人邮件
type ReqAddSingleMail struct {
	MailInfo *SingleMailInfo
}

type RetAddSingleMail struct {
	Ret uint32
}

//增加索要中心
type ReqAddRequestMsg struct {
	MailInfo *RequestMessage
}

type RetAddRequestMsg struct {
	Ret uint32
}

//获取系统邮件列表
type ReqSysMailList struct {
	UserId uint64
}

//返回系统邮件列表
type RetSysMailList struct {
	Ret      uint32
	MailList []*UserSysMailInfo
}

//获取个人邮件列表
type ReqUserMailList struct {
	UserId uint64
}

//返回玩家好友邮件列表
type RetUserMailList struct {
	Ret      uint32
	MailList []*FriendMailItem
}

//获取索取列表
type ReqRequestMessageList struct {
	UserId uint64
}

type RetRequestMessageList struct {
	Ret      uint32
	MailList []*ReqMail
}

type ReqDelAllReqMsg struct {
	UserId uint64
}

type RetDelAllReqMsg struct {
	Ret uint32
}

type ReqUpdateReqMsg struct {
	MailId string
	State  uint8
}

type RetUpdateReqMsg struct {
	Ret uint32
}

//请求删除个人邮件
type ReqDelSingleMail struct {
	UserId uint64
	MailId string
}

//返回删除个人邮件
type RetDelSingleMail struct {
	Ret uint32
}

//请求读取个人邮件
type ReqReadSingleMail struct {
	UserId uint64
	MailId string
}

//返回读取个人邮件
type RetReadSingleMail struct {
	Ret uint32
}

//请求更新 UpdateSysMailUserList
type ReqUserSysMailUpdate struct {
	UserId uint64
	MailId string
	State  uint8
}

//返回更新
type RetUserSysMailUpdate struct {
	Ret uint32
}

//删除所有系统邮件
type ReqDelAllSysMail struct {
}

type RetDelAllSysMail struct {
	Ret uint32
}

//获取不读的状态
type ReqNonReadState struct {
	UserId uint64
}

type RetNonReadState struct {
	Info *NonReadMailInfo
}

//领取一个邮件的附件
type ReqPickMail struct {
	UserId uint64
	MailId string
	Type   uint8
}

type RetPickMail struct {
	Ret uint32
}

//领取所有附件
type ReqRecvAllMail struct {
	UserId uint64
	Type   uint8
}

type Mailids struct {
	MailId string
}

type RetRecvAllMail struct {
	Ret uint32
	Ids []*Mailids
}

//请求删除所有已读无附件的系统邮件
type ReqDelUserNonAttrSysMail struct {
	UserId uint64
}

type RetDelUserNonAttrSysMail struct {
	Ret uint32
}

type ReqRespReqMessage struct {
	MailId string
	State  uint8
}

type RetRespReqMessage struct {
	Ret uint32
}

type ReqMailById struct {
	MailId string
}

type RetMailById struct {
	MailInfo *RequestMessage
}

type ReqRecvApplyState struct {
	UserId uint64
	State  uint8
}

type RetRecvApplyState struct {
	Ret uint32
}

type RetGetRegTime struct {
	RegTime uint32
}

type ReqGetRegTime struct {
	Id uint64 // 玩家id
}
