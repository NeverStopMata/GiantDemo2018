package composites

import (
	_ "fmt"

	b3 "base/behavior3go"
	_ "base/behavior3go/config"
	. "base/behavior3go/core"
)

type Sequence struct {
	Composite
}

/**
 * Tick method.
 * @method tick
 * @param {b3.Tick} tick A tick instance.
 * @return {Constant} A state constant.
**/
func (this *Sequence) OnTick(tick *Tick) b3.Status {
	//fmt.Println("tick Sequence :", this.GetTitle())
	for i := 0; i < this.GetChildCount(); i++ {
		var status = this.GetChild(i).Execute(tick)
		if status != b3.SUCCESS {
			return status
		}
	}
	return b3.SUCCESS
}
