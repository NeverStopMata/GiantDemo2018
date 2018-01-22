//ai全局管理
package ai

import (
	b3 "base/behavior3go"
	b3config "base/behavior3go/config"
	b3core "base/behavior3go/core"
	b3loader "base/behavior3go/loader"
	"base/env"
	"fmt"
	"roomserver/conf"
)

//公用树
var mopeBevTrees map[string]*b3core.BehaviorTree

func GetBevTree(name string) *b3core.BehaviorTree {
	tree, ok := mopeBevTrees[name]
	if ok {
		return tree
	}
	CreateBevTree(name)
	return mopeBevTrees[name]
}

//创建黑板
func CreateBlackBoard() *b3core.Blackboard {
	return b3core.NewBlackboard()
}

//
func CreateBevTree(name string) {
	_, ok := mopeBevTrees[name]
	if ok {
		return
	}
	fmt.Println("create tree:", name)
	config, ok := b3config.LoadTreeCfg(env.Get("global", "aicfg") + name)
	if !ok {
		panic("LoadTreeCfg fail:" + name)
	}
	extMaps := createExtStructMaps()
	mopeBevTree := b3loader.CreateBevTreeFromConfig(config, extMaps)
	//mopeBevTree.Print()
	mopeBevTrees[name] = mopeBevTree
}
func CreateBevAIMgr() bool {
	fmt.Println("load tree...", b3.VERSION)
	//b3core.IsPrintLog = true //打印ai日志
	mopeBevTrees = make(map[string]*b3core.BehaviorTree)
	for _, data := range conf.ConfigMgr_GetMe().AIDatas.AIData.Behave.Datas {
		CreateBevTree(data.Aifile)
	}
	return true
}

func createExtStructMaps() *b3.RegisterStructMaps {
	st := b3.NewRegisterStructMaps()
	//actions
	st.Register("Rand", &RandAction{})
	st.Register("TurnTarget", &TurnTarget{})
	st.Register("TurnIndex", &TurnIndex{})
	st.Register("FindNearUnit", &FindNearUnit{})
	st.Register("TurnTargetPlayer", &TurnTargetPlayer{})
	st.Register("TurnAwayTarget", &TurnAwayTargetPlayer{})
	st.Register("FindAttackTarget", &FindAttackTarget{})
	st.Register("SubTree", &SubTreeNode{})
	st.Register("CheckDis", &CheckDisNode{})
	st.Register("BBCastSkill", &BBCastSkillNode{})
	st.Register("EnemyToAttackTarget", &EnemyToAttackTarget{})
	st.Register("MoveCtrl", &MoveCtrl{})
	st.Register("CheckDis2", &CheckDis2Node{})
	st.Register("MoveBack", &MoveBack{})
	st.Register("FindNearUnit2", &FindNearUnit2{})
	st.Register("WaitSkillIdle", &WaitSkillIdle{})

	//conditions
	st.Register("CheckBall", &CheckBall{})
	st.Register("CheakBall", &CheckBall{}) //fake name
	st.Register("AttrLimit", &AttrLimit{})
	st.Register("TargetAttrLess", &TargetAttrLess{})
	st.Register("CheckBool", &CheckBool{})
	st.Register("CheckNearPlayer", &CheckNearPlayer{})
	st.Register("CheckNearAttackPlayer", &CheckNearAttackPlayer{})
	st.Register("HpMoreThan", &HpMoreThan{})

	//composite
	st.Register("Random", &RandomComposite{})
	st.Register("Parallel", &ParallelComposite{})
	return st
}
