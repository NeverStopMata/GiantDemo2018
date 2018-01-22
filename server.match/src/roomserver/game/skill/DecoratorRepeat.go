package skill

import (
	b3 "base/behavior3go"
	. "base/behavior3go/config"
	. "base/behavior3go/core"
)

type DecoratorRepeater struct {
	Decorator
	maxLoop int
}

func (this *DecoratorRepeater) Initialize(setting *BTNodeCfg) {
	this.Decorator.Initialize(setting)
	this.maxLoop = setting.GetPropertyAsInt("maxLoop")
	if this.maxLoop < 1 {
		panic("maxLoop parameter in MaxTime decorator is an obligatory parameter")
	}
}

func (this *DecoratorRepeater) OnOpen(tick *Tick) {
	tick.Blackboard.Set("i", 0, tick.GetTree().GetID(), this.GetID())
}

func (this *DecoratorRepeater) OnTick(tick *Tick) b3.Status {
	if this.GetChild() == nil {
		return b3.ERROR
	}
	var i = tick.Blackboard.GetInt("i", tick.GetTree().GetID(), this.GetID())
	if i < this.maxLoop {
		i = i + 1
		status := this.GetChild().Execute(tick)

		if status == b3.FAILURE {
			return b3.SUCCESS
		}

		if status != b3.RUNNING {
			tick.Blackboard.Set("i", i, tick.GetTree().GetID(), this.GetID())

		}
		return b3.RUNNING
	} else {
		return b3.SUCCESS
	}
}
