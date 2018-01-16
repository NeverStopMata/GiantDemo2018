package scn

// 场景类-玩家离线逻辑

import (
	"roomserver/game/scn/plr"
	"time"
)

type PlayerOffline struct {
	acc     string
	offtime int64  //断线时间
	exp     uint32 //断线时的经验值
	roomid  uint32
	x       float64
	y       float64
}

type ScenePlayerOffline struct {
	Offlineinfos []PlayerOffline // 离线玩家列表
}

func (this *Scene) CheckOffline(acc string) (uint16, bool) {
	for index, _ := range this.Offlineinfos {
		if this.Offlineinfos[index].acc == acc {
			return uint16(index), true
		}
	}
	return 0, false
}

func (this *Scene) GetOfflineExp(acc string) uint32 {
	for index, _ := range this.Offlineinfos {
		if this.Offlineinfos[index].acc == acc {
			return this.Offlineinfos[index].exp
		}
	}
	return 0
}

func (this *Scene) DeleteOffline() {
	for index, _ := range this.Offlineinfos {
		if this.Offlineinfos[index].offtime+60 < time.Now().Unix() {
			this.Offlineinfos = append(this.Offlineinfos[:index], this.Offlineinfos[index+1:]...)
			return
		}
	}
}

func (this *Scene) DeleteActOffline(acc string) {
	for index, _ := range this.Offlineinfos {
		if this.Offlineinfos[index].acc == acc {
			this.Offlineinfos = append(this.Offlineinfos[:index], this.Offlineinfos[index+1:]...)
			return
		}
	}
}

func (this *Scene) GetOfflinePos(acc string) (float64, float64, bool) {
	for index, _ := range this.Offlineinfos {
		if this.Offlineinfos[index].acc == acc {
			return this.Offlineinfos[index].x, this.Offlineinfos[index].y, true
		}
	}
	return 0, 0, false
}

func (this *Scene) AddOffline(player *plr.ScenePlayer) {
	index, ex := this.CheckOffline(player.UData().Account)
	if true == ex /*&& true == player.IsLive*/ {
		this.Offlineinfos[index].exp = player.GetExp()
		this.Offlineinfos[index].offtime = time.Now().Unix()
		this.Offlineinfos[index].x = player.SelfAnimal.Pos.X
		this.Offlineinfos[index].y = player.SelfAnimal.Pos.Y
		this.Offlineinfos[index].roomid = this.room.ID()
	} else /*if true == player.IsLive*/ {
		var tmp PlayerOffline
		tmp.acc = player.UData().Account
		tmp.exp = player.GetExp()
		tmp.offtime = time.Now().Unix()
		tmp.roomid = this.room.ID()
		tmp.x = player.SelfAnimal.Pos.X
		tmp.y = player.SelfAnimal.Pos.Y
		this.Offlineinfos = append(this.Offlineinfos, tmp)
	}
}
