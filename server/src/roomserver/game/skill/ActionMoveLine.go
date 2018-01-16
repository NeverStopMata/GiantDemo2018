package skill

import (
	b3 "base/behavior3go"
	b3config "base/behavior3go/config"
	b3core "base/behavior3go/core"
	_ "base/glog"
	"roomserver/game/scn/plr"
	"roomserver/util"
)

type ActionMoveLine struct {
	b3core.Action
	d        float64
	n        uint64
	dir_type int
}

func (this *ActionMoveLine) Initialize(setting *b3config.BTNodeCfg) {
	this.Action.Initialize(setting)
	this.d = setting.GetProperty("d")
	this.n = uint64(setting.GetPropertyAsInt("n"))
	this.dir_type = setting.GetPropertyAsInt("dir_type")
}

func (this *ActionMoveLine) OnOpen(tick *b3core.Tick) {
	player := tick.Blackboard.Get("player", "", "").(*plr.ScenePlayer)

	var speed util.Vector2

	if this.dir_type == 2 {
		x0, y0 := player.SelfAnimal.GetPos()
		x1 := tick.Blackboard.GetFloat64("source_pos_x", "", "")
		y1 := tick.Blackboard.GetFloat64("source_pos_y", "", "")

		v := &util.Vector2{float64(x0), float64(y0)}
		hv := v.SubMethod(&util.Vector2{x1, y1})
		speed = hv.Normalize()
	} else if this.dir_type == 1 {
		speed = util.Vector2{float64(player.SelfAnimal.GetAngleVel().X), float64(player.SelfAnimal.GetAngleVel().Y)}
	} else {
		panic("error dir_type!")
	}

	player.SelfAnimal.ClearForce()

	force1 := speed
	force1.ScaleBy(this.d / float64(this.n) * 2)
	player.SelfAnimal.AddForce(force1, this.n)
}

func (this *ActionMoveLine) OnTick(tick *b3core.Tick) b3.Status {
	player := tick.Blackboard.Get("player", "", "").(*plr.ScenePlayer)
	if player.SelfAnimal.HasForce() == true {
		return b3.RUNNING
	}
	player.SelfAnimal.ClearForce()
	return b3.SUCCESS
}
