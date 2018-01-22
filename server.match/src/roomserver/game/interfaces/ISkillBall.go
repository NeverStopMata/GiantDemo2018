package interfaces

// 技能球 的 技能 接口

type ISkillBall interface {
	CastSkill(skillid uint32)
	Update()
	IsFinish() bool
	GetBeginFrame() uint32
}
