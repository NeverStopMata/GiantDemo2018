package common

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//////////////////////聊天房间信息<<<<<<<<<<<<<<<<<<<
// 聊天房间类型
const (
	CHAT_TYPE_ROOM      int32 = 1 // 临时房间
	CHAT_TYPE_TEAM      int32 = 2 // 临时队伍
	CHAT_TYPE_GROUP     int32 = 3 // 群聊
	CHAT_TYPE_MATCHTEAM int32 = 4 // 比赛组队

	CHAT_TYPE_RADAR  int32 = 6 // 雷达
	CHAT_TYPE_SYSTEM int32 = 7 // 系统群

	CHAT_TALK_0 uint32 = 0 // 互关私聊
	CHAT_TALK_1 uint32 = 1 // 快捷私聊
	CHAT_TALK_2 uint32 = 2 // 互关推送
	CHAT_TALK_3 uint32 = 3 // 未读接收消息
	CHAT_TALK_4 uint32 = 4 // 粉丝私聊

	//群设置
	CHAT_SET_TOP     int32 = 1 // 置顶
	CHAT_SET_DISTURB int32 = 2 // 消息免打扰

	CHAT_MATCHTEAM_TIME uint32 = 3 * 24 * 60 * 60 //比赛结束组队群组时效
	CHAT_ROOM_TIME      uint32 = 30 * 60          //比赛结束临时群组时效

	CHAT_RED_TYPE_1 uint32 = 1 //彩豆包

	CHAT_RED_MAX    uint32 = 1000 //最大爱心数量
	CHAT_EXCHANGE_1 uint32 = 1    //彩豆兑换爱心比
	CHAT_EXCHANGE_2 uint32 = 2    //蘑菇兑换爱心比

	CHAT_RED_END_TIME uint32 = 24 * 60 * 60 //红包过期时间

	CHAT_ROOM_ADD uint32 = 1 // 增加
	CHAT_ROOM_DEL uint32 = 2 // 删除

	CHAT_ACTIVE_TIME uint32 = 5 * 60 //群组活跃时间

	CHAT_TIME int64 = 500 //聊天频率 毫秒
)

type ChatRoom struct {
	Id       bson.ObjectId `bson:"_id"`
	RoomId   uint64        `bson:"roomid"`   //进程房间ID
	Owner    uint64        `bson:"owner"`    //房主ID
	Name     string        `bson:"name"`     //群名称
	Type     int32         `bson:"type"`     //房间类型
	UserList []uint64      `bson:"userlist"` //成员ID列表
	EndTime  uint32        `bson:"endtime"`  //房间结束时间
	TalkId   uint64        `bson:"talkid"`   //聊天内容id
}

type UserShow struct {
	Id       uint64
	Account  string // 帐号
	Sex      uint8  // 性别
	PassIcon string // 已审核头像
	Icon     uint32 // 图标
}

type ChatSynTalk struct {
	RoomId  uint64 //房间ID
	UserId  uint64 //玩家ID
	StartId uint64 //开始ID
}

type ChatSynSet struct {
	RoomId  uint64 //房间ID
	UserId  uint64 //玩家ID
	Top     bool   //置顶
	Disturb bool   //消息免打扰
}

type ChatSynPersonalSet struct {
	FriendTalk bool `bson:"friendtalk"` //是否开启互关好友才能聊天
	QuickTalk  bool `bson:"quicktalk"`  //是否开启快捷聊天
	BatchPush  bool `bson:"batchpush"`  //是否开启批量互关推送
	UnReadTalk bool `bson:"unreadtalk"` //是否开启未读接收消息
}

type ChatSyn struct {
	Id          bson.ObjectId       `bson:"_id"`
	UserId      uint64              `bson:"userid"`       // 用户ID
	Rooms       []uint64            `bson:"rooms"`        // 拥有群列表
	Talks       []*ChatSynTalk      `bson:"talks"`        // 聊天游标
	Sets        []*ChatSynSet       `bson:"sets"`         // 群设置
	PersonalSet *ChatSynPersonalSet `bson:"personalsets"` // 个人设置
	FansTalk    []uint64            `bson:"fanstalk"`     // 粉丝聊天
}

//房间离线消息
type ChatTalkRoom struct {
	Id     bson.ObjectId `bson:"_id"`
	RoomId uint64        `bson:"roomid"` // 房间ID
	TalkId uint64        `bson:"talkid"` // 聊天ID
	UserId uint64        `bson:"userid"` // 说话玩家ID
	Text   string        `bson:"text"`   // 内容
	Time   uint32        `bson:"time"`   // 时间
	Date   time.Time     `bson:"date"`   // 日期
}

//玩家离线消息
type ChatTalkUser struct {
	Id     uint64    `bson:"_id"`
	SelfId uint64    `bson:"selfid"` // 玩家自身ID
	UserId uint64    `bson:"userid"` // 说话玩家ID
	Text   string    `bson:"text"`   // 内容
	Time   uint32    `bson:"time"`   // 时间
	Date   time.Time `bson:"date"`   // 日期
}

type RetTalk struct {
	RoomId string         //房间ID
	UserId uint64         //玩家ID
	Talks  []*RetTalkTalk //内容
}

type RetTalkTalk struct {
	UserId uint64 //玩家ID
	Text   string //内容
	Time   uint32 //时间
	TalkId uint64 //聊天ID
}

type ChatRed struct {
	TalkId   uint64                     //聊天ID
	SendId   uint64                     //发送者ID
	RoomId   uint64                     //房间ID
	TypeId   uint32                     //爱心包类型 1彩豆包 2蘑菇包
	Nums     uint32                     //数量
	Total    uint32                     //总额
	Less     uint32                     //剩余额度
	SendTime uint32                     //发送时间
	UserList map[uint64]*ChatRedReceive //已领取红包列表
	EndTime  uint32                     //领完时间
}

type ChatRedMongo struct {
	Id       bson.ObjectId     `bson:"_id"`
	TalkId   uint64            `bson:"talkid"`   //聊天ID
	SendId   uint64            `bson:"sendid"`   //发送者ID
	RoomId   uint64            `bson:"roomid"`   //房间ID
	TypeId   uint32            `bson:"typeid"`   //爱心包类型 1彩豆包 2蘑菇包
	Nums     uint32            `bson:"nums"`     //数量
	Total    uint32            `bson:"total"`    //总额
	Less     uint32            `bson:"less"`     //剩余额度
	SendTime uint32            `bson:"sendtime"` //发送时间
	UserList []*ChatRedReceive `bson:"userlist"` //已领取红包列表
	EndTime  uint32            `bson:"endtime"`  //领完时间
}

type ChatRedReceive struct {
	UserId uint64 `json:"userid"` //玩家ID
	Loves  uint32 `json:"loves"`  //爱心数
	Time   uint32 `json:"time"`   //时间
}

type ChatRedEnd struct {
	RoomId  uint64 //群组ID
	TalkId  uint64 //聊天ID
	EndTime uint32 //过期时间
}

const (
	NOTICE_TYPE_1 uint32 = 1 //滚屏公告
	NOTICE_TYPE_2 uint32 = 2 //系统消息
	NOTICE_TYPE_3 uint32 = 3 //弹出公告
)

type NoticeInfo struct {
	Id        uint64 `redis:"Id"`
	TypeId    uint32 `redis:"TypeId"`    //类型1滚屏公告 2系统消息
	Text      string `redis:"Text"`      //内容
	SysType   string `redis:"SysType"`   //系统类型 android ios
	ClientVer string `redis:"ClientVer"` //客户端版本
	BetUserid string `redis:"BetUserid"` //玩家ID区间
	Awards    string `redis:"Awards"`    //奖励
	STime     uint32 `redis:"STime"`     //开始时间
	ETime     uint32 `redis:"ETime"`     //结束时间
	RTime     uint32 `redis:"RTime"`     //间隔频率
	CanSend   bool   `redis:"CanSend"`   //是否可以发送
}
