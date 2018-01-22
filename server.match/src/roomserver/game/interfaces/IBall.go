package interfaces

// 球的基本接口

import (
	"roomserver/util"
	"usercmd"
)

type IBall interface {
	GetID() uint32
	SetID(uint32)
	GetTypeId() uint16
	GetType() usercmd.BallType
	GetPos() (float64, float64)
	SetPos(float64, float64)
	SetBirthPoint(_birthPoint IBirthPoint)
	GetRect() *util.Square
	OnReset()
}

//球种类
type BallKind uint8

const (
	BallKind_None   BallKind = 0
	BallKind_Player BallKind = 1  // 玩家
	BallKind_Food   BallKind = 4  // 食物
	BallKind_Feed   BallKind = 8  // 蘑菇
	BallKind_Skill  BallKind = 16 // 技能
)

func BallTypeToKind(btype usercmd.BallType) BallKind {
	if btype == usercmd.BallType_Player {
		return BallKind_Player
	} else if btype > usercmd.BallType_FoodBegin && btype < usercmd.BallType_FoodEnd {
		return BallKind_Food
	} else if btype > usercmd.BallType_FeedBegin && btype < usercmd.BallType_FeedEnd {
		return BallKind_Feed
	} else if btype > usercmd.BallType_SkillBegin && btype < usercmd.BallType_SkillEnd {
		return BallKind_Skill
	} else {
		return BallKind_None
	}
}
