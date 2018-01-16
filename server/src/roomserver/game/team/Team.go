// 包team定义队伍
package team

import "sync"

// 队伍数据
type Team struct {
	Id              uint64          // 战队id
	Name            string          // 战队名
	ClassicLogo     uint32          // 默认战旗
	Logo            string          // 自定义战队Logo
	Level           uint32          // 星级
	CupNum          int32           // 杯数
	MatchNum        int32           // 匹配值
	Rank            uint32          // 排名
	AddMatchNum     int32           // 增加的匹配值
	AddCupNum       int32           // 增加的杯数
	LeaderID        uint64          // 队长ID
	IsNewbie        bool            // 是否萌新
	MemList         map[uint64]bool // 队员列表
	NoticeTime      int64           // 发标记时间
	PlayerListMutex sync.RWMutex
	PlayerList      []byte
	PlayerListFlag  byte
}
