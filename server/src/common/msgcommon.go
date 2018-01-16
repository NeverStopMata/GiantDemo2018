package common

const (
	MsgType_Team        uint8 = 0  // 组队消息
	MsgType_Chat        uint8 = 2  // 聊天
	MsgType_ChatForward uint8 = 3  // 群聊天
	MsgType_WorldChat   uint8 = 4  // 世界聊天
	MsgType_TeamChat    uint8 = 5  // 团战聊天
	MsgType_Server      uint8 = 16 // 服务器内部转发
)
