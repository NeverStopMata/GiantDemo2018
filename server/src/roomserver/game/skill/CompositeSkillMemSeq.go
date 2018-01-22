package skill

import (
	b3 "base/behavior3go"
	. "base/behavior3go/core"
	_ "base/glog"
	"roomserver/game/scn/plr"
)

type CompositeSkillMemSeq struct {
	Composite
}

func (this *CompositeSkillMemSeq) OnOpen(tick *Tick) {
	tick.Blackboard.Set("runningChild", 0, tick.GetTree().GetID(), this.GetID())

	cfg := tick.GetTree().Dump()
	if v, ok := cfg.Properties["skillid"]; ok {
		tick.Blackboard.Set("castskill", v, "", "")
	} else {
		panic("error skill bev tree cfg. can't get skillid")
	}
	tick.Blackboard.Set("endframe", uint64(0), tick.GetTree().GetID(), this.GetID())
	tick.Blackboard.Set("hits", make(map[uint32]int), "", "")
	tick.Blackboard.Set("condfail", int(0), "", "")

	// 扣蓝
	player := tick.Blackboard.Get("player", "", "").(*plr.ScenePlayer)
	if _, ok := cfg.Properties["skillcost"]; ok {
		//d := int32(v.(float64))
		d := int32(0)
		curmp := int32(player.SelfAnimal.GetMP()) - d
		if d > 0 && curmp >= 0 {
			player.SelfAnimal.SetMP(float64(curmp))
		} else if d > 0 && curmp < 0 {
			tick.Blackboard.Set("condfail", int(1), "", "")
		}
	}

	tick.Blackboard.Set("playerangle", nil, "", "")
	tick.Blackboard.Set("playerface", nil, "", "")
	tick.Blackboard.Set("playerpower", nil, "", "")
	// 停止移动
	if false {
		playerskill := tick.Blackboard.Get("playerskill", "", "")
		if playerskill != nil {
			tick.Blackboard.Set("playerpower", float64(player.Power), "", "")
			player.Power = 0
		}
	}
}

func (this *CompositeSkillMemSeq) OnTick(tick *Tick) b3.Status {
	condfail := tick.Blackboard.Get("condfail", "", "").(int)
	if condfail == 1 {
		return b3.FAILURE
	}

	var child = tick.Blackboard.GetInt("runningChild", tick.GetTree().GetID(), this.GetID())
	for i := child; i < this.GetChildCount(); i++ {
		var status = this.GetChild(i).Execute(tick)

		if status != b3.SUCCESS {
			if status == b3.RUNNING {
				tick.Blackboard.Set("runningChild", i, tick.GetTree().GetID(), this.GetID())
			} else {
			}
			return status
		}
	}
	if child != this.GetChildCount() {
		tick.Blackboard.Set("runningChild", this.GetChildCount(), tick.GetTree().GetID(), this.GetID())
	}

	// 释放完毕后，补完这个100毫秒
	player := tick.Blackboard.Get("player", "", "").(*plr.ScenePlayer)
	endframe := tick.Blackboard.Get("endframe", tick.GetTree().GetID(), this.GetID()).(uint64)
	if endframe == 0 {
		endframe = GetEndFrame(player, 1)
		tick.Blackboard.Set("endframe", endframe, tick.GetTree().GetID(), this.GetID())
	}

	if endframe+1 > uint64(player.GetFrame()) {
		return b3.RUNNING
	}

	return b3.SUCCESS
}

func (this *CompositeSkillMemSeq) OnClose(tick *Tick) {
	skillid := uint32(tick.Blackboard.Get("castskill", "", "").(float64))
	tick.Blackboard.Set("castskill", float64(0), "", "")
	tick.Blackboard.Set("condfail", int(0), "", "")
	tick.Blackboard.Set("hits", make(map[uint32]int), "", "")

	playerskill := tick.Blackboard.Get("playerskill", "", "")
	if playerskill != nil && skillid >= SKILL_GETHIT_MIN {
		playerskill.(*SkillPlayer).nextSkillId = 0
	}

	// 技能结束，恢复速度
	if false {
		if playerskill != nil {
			power := tick.Blackboard.Get("playerpower", "", "")
			if power != nil {
				player := tick.Blackboard.Get("player", "", "").(*plr.ScenePlayer)
				player.Power = power.(float64)
			}
		}
	}
	tick.Blackboard.Set("playerangle", nil, "", "")
	tick.Blackboard.Set("playerface", nil, "", "")
	tick.Blackboard.Set("playerpower", nil, "", "")
}
