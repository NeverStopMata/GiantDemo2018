package scn

// 场景中分配球相关的辅助类。包括分配球位置、分配球ID、机器人UID等

import (
	"common"
	"math"
	"roomserver/game/scn/internal"
	"sync/atomic"
)

var (
	NewBallBeginId uint32 = 251 // 除玩家球、技能球以外的球ID
)

type SceneBallHelper struct {
	internal.UnusedPos        // 可用的位置
	internal.UnuseId          // 可用的ID，用于玩家球、技能球
	ballID             uint32 // 分配球使用的球ID
	robotUID           uint64 // 分配机器人UID
}

func (this *SceneBallHelper) Init(roomSize float64) {
	this.UnusedPos.Init(roomSize)
	this.UnuseId.Init(NewBallBeginId - 1)
	this.ballID = NewBallBeginId
	this.robotUID = common.MINROBOTID
}

// 获取一个新的球ID
func (this *SceneBallHelper) NewBallId() uint32 {
	return atomic.AddUint32(&this.ballID, 1)
}

// 获取一个新的玩家球ID
func (this *SceneBallHelper) NewBallPlayerId() uint32 {
	ballid, ok := this.GetId()
	if !ok {
		ballid = this.NewBallId()
	}
	return ballid
}

// 获取一个新的技能球ID
func (this *SceneBallHelper) NewBallSkillId() uint32 {
	ballid, ok := this.GetId()
	if !ok {
		ballid = this.NewBallId()
	}
	return ballid
}

// 获取一个新的机器人UID
func (this *SceneBallHelper) NewRobotUID() uint64 {
	return atomic.AddUint64(&this.robotUID, 1)
}

// 获取一个圆所占用的 x、y范围
func (this *SceneBallHelper) GetSquare(x, y, radius float64) []*internal.Point {
	retvalue := make([]*internal.Point, 0)
	leftup_x := int32(x - radius)
	leftup_y := int32(math.Ceil(y + radius))
	rightup_x := int32(math.Ceil(x + radius))
	leftdown_y := int32(y - radius)
	for i := leftdown_y - 1; i <= leftup_y+1; i++ {
		for j := leftup_x - 1; j <= rightup_x+1; j++ {
			retvalue = append(retvalue, &internal.Point{float64(j), float64(i)})
		}
	}
	return retvalue
}
