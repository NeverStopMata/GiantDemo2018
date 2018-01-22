// 包 playermgr 处理玩家管理类.
package playermgr

import (
	"common"
	"roomserver/game/scn/plr"
	"sync"
)

type ScenePlayerMgr struct {
	mutex   sync.RWMutex
	players map[string]*plr.ScenePlayer
}

var (
	sptaskm             *ScenePlayerMgr
	sceneplayermgr_once sync.Once
)

func GetMe() *ScenePlayerMgr {
	if sptaskm == nil {
		sceneplayermgr_once.Do(func() {
			sptaskm = &ScenePlayerMgr{
				players: make(map[string]*plr.ScenePlayer),
			}
		})
	}
	return sptaskm
}

func (this *ScenePlayerMgr) GetNum() int {
	this.mutex.Lock()
	num := len(this.players)
	this.mutex.Unlock()
	return num
}

func (this *ScenePlayerMgr) Add(task *plr.ScenePlayer) {
	this.mutex.Lock()
	this.players[task.Key] = task
	this.mutex.Unlock()
}

func (this *ScenePlayerMgr) Remove(player *plr.ScenePlayer) {
	if player.Key == "" {
		return
	}
	this.mutex.Lock()
	delete(this.players, player.Key)
	this.mutex.Unlock()
}

func (this *ScenePlayerMgr) Removes(players map[uint64]*plr.ScenePlayer) {
	this.mutex.Lock()
	for _, player := range players {
		if _, ok := this.players[player.Key]; ok {
			delete(this.players, player.Key)
		}
	}
	this.mutex.Unlock()
}

func (this *ScenePlayerMgr) GetPlayer(key string) (player *plr.ScenePlayer) {
	this.mutex.RLock()
	player, _ = this.players[key]
	this.mutex.RUnlock()
	return
}

func (this *ScenePlayerMgr) GetUDataFromKey(key string) *common.UserData {
	player := this.GetPlayer(key)
	if player == nil {
		return nil
	}
	return player.UData()
}
