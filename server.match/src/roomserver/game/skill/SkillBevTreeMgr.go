package skill

// 用于加载技能行为树、获取某技能行为树

import (
	b3 "base/behavior3go"
	b3config "base/behavior3go/config"
	b3core "base/behavior3go/core"
	b3loader "base/behavior3go/loader"
	"base/env"
	"base/glog"
	"fmt"
	"io/ioutil"
)

const (
	SKILL_ID_MIN     = 50
	SKILL_ID_HAMMER  = 103
	SKILL_ID_BOMB    = 104
	SKILL_ID_MAX     = 250
	SKILL_GETHIT_MIN = 151
)

//公用树
var mopeSkillBevTrees map[uint32]*b3core.BehaviorTree
var mopGetHitBevTrees map[uint32]*b3core.BehaviorTree

func GetSkillBevTree(skillid uint32) *b3core.BehaviorTree {
	if _, ok := mopeSkillBevTrees[skillid]; ok {
		return mopeSkillBevTrees[skillid]
	}

	return nil
}

func GetGetHitBevTree(skillid uint32) *b3core.BehaviorTree {
	id := skillid + 100
	if _, ok := mopGetHitBevTrees[id]; ok {
		return mopGetHitBevTrees[id]
	}
	return mopGetHitBevTrees[200]
}

func LoadSkillBevTree() bool {
	glog.Infoln("load skill tree...", b3.VERSION)
	mopeSkillBevTrees = make(map[uint32]*b3core.BehaviorTree)
	for i := SKILL_ID_MIN; i <= SKILL_ID_MAX; i++ {
		tree := createBevTree(uint32(i), "skill")
		if tree != nil {
			mopeSkillBevTrees[uint32(i)] = tree
		}
	}

	glog.Infoln("load gethit tree...", b3.VERSION)
	mopGetHitBevTrees = make(map[uint32]*b3core.BehaviorTree)
	for i := SKILL_ID_MIN + 100; i <= SKILL_ID_MAX+100; i++ {
		tree := createBevTree(uint32(i), "gethit")
		if tree != nil {
			mopGetHitBevTrees[uint32(i)] = tree
		}
	}

	return true
}

func createBevTree(skillid uint32, prefix string) *b3core.BehaviorTree {

	path := fmt.Sprintf("%s/%s_%d.json", env.Get("global", "skillcfg"), prefix, skillid)

	_, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	config, ok := b3config.LoadTreeCfg(path)
	if !ok {
		return nil
	}
	glog.Infof("create %s tree: %s", prefix, path)
	extMaps := createSkillExtStructMaps()
	return b3loader.CreateBevTreeFromConfig(config, extMaps)
}

func createSkillExtStructMaps() *b3.RegisterStructMaps {
	st := b3.NewRegisterStructMaps()

	//actions
	st.Register("ActionAttack1", &ActionAttack1{})
	st.Register("ActionAttack2", &ActionAttack2{})
	st.Register("ActionAttack3", &ActionAttack3{})
	st.Register("ActionNextSceneRander5", &ActionNextSceneRander5{})
	st.Register("ActionMoveLine", &ActionMoveLine{})
	st.Register("ActionThrowBall", &ActionThrowBall{})
	st.Register("ActionHammerTryHit", &ActionHammerTryHit{})
	st.Register("ActionBallSkillStopMove", &ActionBallSkillStopMove{})
	st.Register("ActionBombTryHit", &ActionBombTryHit{})

	//composite
	st.Register("CompositeSkillMemSeq", &CompositeSkillMemSeq{})

	//decorator
	st.Register("DecoratorRepeater", &DecoratorRepeater{})

	return st
}
