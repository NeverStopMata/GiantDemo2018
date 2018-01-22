package skill

import (
	b3 "base/behavior3go"
	b3config "base/behavior3go/config"
	b3core "base/behavior3go/core"
	_ "base/glog"
	"math"
	"roomserver/game/consts"
	"roomserver/game/scn/plr"
)

type ActionNextSceneRander5 struct {
	b3core.Action
	n         uint64
	canAttack uint64
}

func (this *ActionNextSceneRander5) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.n = uint64(setting.GetPropertyAsInt("n"))
	this.canAttack = uint64(setting.GetPropertyAsInt("canAttack"))
}

func (this *ActionNextSceneRander5) OnOpen(tick *b3core.Tick) {
	player := tick.Blackboard.Get("player", "", "").(*plr.ScenePlayer)
	endframe := GetEndFrame(player, this.n)
	tick.Blackboard.Set("endframe", endframe, tick.GetTree().GetID(), this.GetID())

}

func (this *ActionNextSceneRander5) OnTick(tick *b3core.Tick) b3.Status {
	player := tick.Blackboard.Get("player", "", "").(*plr.ScenePlayer)
	endframe := tick.Blackboard.Get("endframe", tick.GetTree().GetID(), this.GetID()).(uint64)
	if endframe <= uint64(player.GetFrame()) {
		return b3.SUCCESS
	} else {
		if this.canAttack != 0 {
			attackType := tick.Blackboard.GetInt("attackType", "", "")
			hits := tick.Blackboard.Get("hits", "", "").(map[uint32]int)
			if len(hits) == 0 {
				SkillAttack(tick, player, attackType)
			}
		}
		return b3.RUNNING
	}
}

func GetEndFrame(player *plr.ScenePlayer, n uint64) uint64 {
	return uint64(math.Floor(float64(player.GetFrame()-1)/float64(consts.FrameCountBy100MS)))*consts.FrameCountBy100MS + consts.FrameCountBy100MS*n + 1
}
