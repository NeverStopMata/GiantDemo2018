package skill

import (
	b3 "base/behavior3go"
	b3config "base/behavior3go/config"
	b3core "base/behavior3go/core"
	_ "base/glog"
	"roomserver/game/scn/plr"
)

type ActionHammerTryHit struct {
	b3core.Action
	scale  float64
	gethit uint32
}

func (this *ActionHammerTryHit) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.scale = setting.GetProperty("scale")
	this.gethit = uint32(setting.GetPropertyAsInt("gethit"))
}

func (this *ActionHammerTryHit) OnTick(tick *b3core.Tick) b3.Status {
	ballskill := tick.Blackboard.Get("ballskill", "", "").(*SkillBall).ball
	player := tick.Blackboard.Get("player", "", "").(*plr.ScenePlayer)
	hits := tick.Blackboard.Get("hits", "", "").(map[uint32]int)
	scene := player.GetScene()

	attckRect := ballskill.GetRect()
	attckRect.SetRadius(ballskill.GetRadius() + 0.5)
	cells := scene.GetAreaCells(attckRect)

	for _, other := range scene.GetPlayers() {
		if other.GetId() == player.GetId() {
			continue
		}
		if _, ok := hits[other.SelfAnimal.GetID()]; ok {
			continue
		}
		if BallSkillAttack(tick, player, ballskill, this.scale, other.SelfAnimal) {
			hits[other.SelfAnimal.GetID()] = 1
			x, y := ballskill.GetPos()
			other.Skill.GetHit2(x, y, this.gethit)
		}
	}

	for _, cell := range cells {
		for _, feed := range cell.Feeds {
			distance := feed.GetPosV().SqrMagnitudeTo(ballskill.GetPosV())
			tmp := feed.GetRadius() + ballskill.GetRadius()
			if distance <= tmp*tmp {
				return b3.FAILURE
			}
		}
	}

	return b3.SUCCESS
}
