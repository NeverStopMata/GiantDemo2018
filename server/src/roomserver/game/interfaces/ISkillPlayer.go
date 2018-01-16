package interfaces

// 玩家球 的 技能 接口

type ISkillPlayer interface {
	Update()
	CastSkill(skillid uint32, targetId uint32) bool
	GetCurSkillId() uint32
	GetBlackboard(p1, p2, p3 string) interface{}
	GetHit2(sourceX, sourceY float64, skillid uint32)
	TryTurn(rawAngle *float64, rawFace *uint32)
}
