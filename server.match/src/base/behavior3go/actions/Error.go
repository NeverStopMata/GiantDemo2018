package actions

import (
	b3 "base/behavior3go"
	. "base/behavior3go/core"
)

type Error struct {
	Action
}

func (this *Error) OnTick(tick *Tick) b3.Status {
	return b3.ERROR
}
