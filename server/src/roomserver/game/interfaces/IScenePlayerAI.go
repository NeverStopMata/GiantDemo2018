package interfaces

// 玩家AI接口

import (
	"roomserver/conf"
)

type IScenePlayerAI interface {
	InitAI()
	UpdateAI(perTime float64)
	AddAICtrl()
	GetExpireTime() int64
	IsOK() bool
	OnDeadEvent(enemy interface{})
	OnBeHit(enemyBallId uint32)
	GetBevData() *conf.XmlAIBehaveData //ai配置
	FindNearAttackAnimal(disNeed float64) uint32
}
