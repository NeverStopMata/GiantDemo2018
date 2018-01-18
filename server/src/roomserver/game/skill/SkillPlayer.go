package skill

// 用于玩家释放技能

import (
	b3core "base/behavior3go/core"
	_ "base/glog"
	"roomserver/game/bll"
	"roomserver/game/interfaces"
	"roomserver/game/scn/plr"
)

type SkillPlayer struct {
	player            *plr.ScenePlayer
	blackboard        *b3core.Blackboard
	bevTree           *b3core.BehaviorTree
	nextSkillId       uint32
	nextSkillTargetId uint32
}

func NewISkillPlayer(player *plr.ScenePlayer) interfaces.ISkillPlayer {
	return NewSkillPlayer(player)
}

func NewSkillPlayer(player *plr.ScenePlayer) *SkillPlayer {
	s := &SkillPlayer{player: player, blackboard: nil, bevTree: nil, nextSkillId: 0}
	s.blackboard = b3core.NewBlackboard()
	s.blackboard.Set("castskill", float64(0), "", "")
	s.blackboard.Set("player", player, "", "")
	s.blackboard.Set("playerskill", s, "", "")

	return s
}

func (this *SkillPlayer) Update() {

	// 上个技能执行完毕
	if this.bevTree != nil {
		skillid := int(this.blackboard.Get("castskill", "", "").(float64))
		if skillid == 0 {
			this.bevTree = nil
		}
	}

	// 执行下个技能
	if this.bevTree == nil && this.nextSkillId != 0 {
		this.bevTree = GetSkillBevTree(this.nextSkillId)
		this.blackboard.Set("skillTargetId", this.nextSkillTargetId, "", "")
		this.nextSkillId = 0
		this.nextSkillTargetId = 0
	}

	// 执行当前技能
	if this.bevTree != nil {
		this.bevTree.Tick(this, this.blackboard)
	}
}

func (this *SkillPlayer) CastSkill(skillid uint32, targetId uint32) bool {
	if nil == this.player {
		return false
	}

	s := this.GetCurSkillId()
	if s >= SKILL_GETHIT_MIN {
		// 受击中无法接受下个技能释放
		return false
	}

	if skillid == SKILL_ID_BOMB {
		if this.player.SelfAnimal.GetAttr(bll.AttrBombNum) == 0 {
			//return false  //mata: no limit of bombNum
		}
		//this.player.SelfAnimal.SetAttr(bll.AttrBombNum, 0)//mata: never loose ur bomb
	} else if skillid == SKILL_ID_HAMMER {
		if this.player.SelfAnimal.GetAttr(bll.AttrHammerNum) == 0 {
			return false
		}
		this.player.SelfAnimal.SetAttr(bll.AttrHammerNum, 0)
	}

	this.nextSkillId = skillid
	this.nextSkillTargetId = targetId
	//this.player.SetIsRunning(false)
	return true
}

func (this *SkillPlayer) GetHit(source *plr.ScenePlayer, skillid uint32) {
	pos := source.SelfAnimal.GetPosV()
	this.blackboard.Set("source_pos_x", pos.X, "", "")
	this.blackboard.Set("source_pos_y", pos.Y, "", "")

	if this.bevTree != nil {
		this.bevTree.Close(nil, this.blackboard)
	}
	this.bevTree = GetGetHitBevTree(skillid)
	this.bevTree.Tick(nil, this.blackboard)
}

func (this *SkillPlayer) GetHit2(sourceX, sourceY float64, skillid uint32) {
	this.blackboard.Set("source_pos_x", sourceX, "", "")
	this.blackboard.Set("source_pos_y", sourceY, "", "")

	if this.bevTree != nil {
		this.bevTree.Close(nil, this.blackboard)
	}
	this.bevTree = GetGetHitBevTree(skillid)
	this.bevTree.Tick(nil, this.blackboard)
}

func (this *SkillPlayer) TryTurn(rawAngle *float64, rawFace *uint32) {
	skillid := this.GetCurSkillId()
	if skillid != 0 {
		angle := this.GetBlackboard("playerangle", "", "")
		if angle != nil {
			*rawAngle = angle.(float64)
		}
		face := this.GetBlackboard("playerface", "", "")
		if face != nil {
			*rawFace = face.(uint32)
		}
	}
}

func (this *SkillPlayer) GetCurSkillId() uint32 {
	return uint32(this.blackboard.Get("castskill", "", "").(float64))
}

func (this *SkillPlayer) GetBlackboard(p1, p2, p3 string) interface{} {
	return this.blackboard.Get(p1, p2, p3)
}
