package internal

// 附加压迫力

import (
	"base/glog"
	"roomserver/util"
)

//附加力
type AddonForceData struct {
	leftTime uint64
	force    util.Vector2
}

type Force struct {
	addonForceDatas []*AddonForceData //附加力数据
	currForce       util.Vector2      //附加力
}

func (this *Force) ClearForce() {
	this.currForce.X = 0
	this.currForce.Y = 0
	this.addonForceDatas = this.addonForceDatas[:0]
}

func (this *Force) HasForce() bool {
	return len(this.addonForceDatas) != 0
}

func (this *Force) AddForce(force util.Vector2, time uint64) {
	data := &AddonForceData{force: force, leftTime: time}
	this.addonForceDatas = append(this.addonForceDatas, data)
}

func (this *Force) UpdateForce(detaTime float64) {
	this.currForce.X = 0
	this.currForce.Y = 0

	if len(this.addonForceDatas) < 1 {
		return
	}

	var tempList []*AddonForceData

	for _, data := range this.addonForceDatas {
		if data.leftTime > 0 {
			data.leftTime -= 1
			this.currForce.IncreaseBy(&data.force)
			tempList = append(tempList, data)
		}
	}

	if len(tempList) != len(this.addonForceDatas) {
		this.addonForceDatas = tempList
	}

	if len(this.addonForceDatas) > 100 {
		glog.Error("BIG ERROR,addonForceDatas overflow:", len(this.addonForceDatas))
		this.ClearForce()
	}
}
func (this *Force) GetForce() *util.Vector2 {
	return &this.currForce
}
