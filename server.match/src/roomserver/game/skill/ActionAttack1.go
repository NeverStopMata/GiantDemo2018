package skill

import (
	b3 "base/behavior3go"
	b3config "base/behavior3go/config"
	b3core "base/behavior3go/core"
	_ "base/glog"
	"roomserver/game/scn/plr"
)

type ActionAttack1 struct {
	b3core.Action
	skillid     int
	attackType  int
	attackRange float64
}

func (this *ActionAttack1) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.skillid = setting.GetPropertyAsInt("skillid")
	this.attackType = setting.GetPropertyAsInt("attack_type")
	this.attackRange = setting.GetProperty("attack_range")
}

func (this *ActionAttack1) OnOpen(tick *b3core.Tick) {
	tick.Blackboard.Set("attackType", this.attackType, "", "")
	tick.Blackboard.Set("attackRange", this.attackRange, "", "")
}

func (this *ActionAttack1) OnTick(tick *b3core.Tick) b3.Status {
	player := tick.Blackboard.Get("player", "", "").(*plr.ScenePlayer)
	if player.IsLive == false {
		return b3.FAILURE
	}
	SkillAttack(tick, player, this.attackType)
	return b3.SUCCESS
}
