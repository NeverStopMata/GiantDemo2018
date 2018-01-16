package loader

import (
	_ "fmt"
	_ "reflect"

	b3 "base/behavior3go"
	. "base/behavior3go/actions"
	. "base/behavior3go/composites"
	. "base/behavior3go/config"
	. "base/behavior3go/core"
	. "base/behavior3go/decorators"
)

func createBaseStructMaps() *b3.RegisterStructMaps {
	st := b3.NewRegisterStructMaps()
	//actions
	st.Register("Error", &Error{})
	st.Register("Failer", &Failer{})
	st.Register("Runner", &Runner{})
	st.Register("Succeeder", &Succeeder{})
	st.Register("Wait", &Wait{})
	st.Register("Log", &Log{})
	//composites
	st.Register("MemPriority", &MemPriority{})
	st.Register("MemSequence", &MemSequence{})
	st.Register("Priority", &Priority{})
	st.Register("Sequence", &Sequence{})

	//decorators
	st.Register("Inverter", &Inverter{})
	st.Register("Limiter", &Limiter{})
	st.Register("MaxTime", &MaxTime{})
	st.Register("Repeater", &Repeater{})
	st.Register("RepeatUntilFailure", &RepeatUntilFailure{})
	st.Register("RepeatUntilSuccess", &RepeatUntilSuccess{})
	return st
}

func CreateBevTreeFromConfig(config *BTTreeCfg, extMap *b3.RegisterStructMaps) *BehaviorTree {
	baseMaps := createBaseStructMaps()
	tree := NewBeTree()
	tree.Load(config, baseMaps, extMap)
	return tree
}
