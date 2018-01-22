package plr

// 移动处理辅助类

import (
	"math"
	"time"
)

const (
	MAX_OPS_MOVE = 12 // 移动操作，最大每秒处理请求次数
)

type MoveHelper struct {
	cacheAngle   float64 // 缓存上次移动数据
	cacheangle   float64 // 缓存上次移动数据
	cacheFace    uint32  // 缓存上次移动数据
	lastMoveTime int64   // 用于控制移动发包频率
	lastMoveops  int     // 用于控制移动发包频率
	Angle        float64 // 当前移动方向
	Face         uint32  // 当前朝向的目标，值为0表示没有朝向目标
	Power        float64 // 当前摇杆力度（目前恒为0或者1，来简化同步计算）
}

func (this *MoveHelper) CheckMoveMsg(power, angle float64, face uint32) (float64, float64, uint32, bool) {
	if this.lastMoveTime == 0 {
		this.lastMoveTime = time.Now().Add(time.Second).Unix()
	}
	if this.lastMoveops > MAX_OPS_MOVE {
		now := time.Now()
		if this.lastMoveTime > now.Unix() {
			return 0, 0, 0, false
		}
		this.lastMoveTime = now.Add(time.Second).Unix()
		this.lastMoveops = 0
	}
	this.lastMoveops++

	power = math.Min(math.Max(0, float64(power)), 100) * 0.01
	angle = math.Min(math.Max(0, float64(angle)), 360)
	if this.cacheAngle != power || this.cacheangle != angle || this.cacheFace != face {
		this.cacheAngle = power
		this.cacheangle = angle
		this.cacheFace = face
		return power, angle, face, true
	}
	return 0, 0, 0, false
}
