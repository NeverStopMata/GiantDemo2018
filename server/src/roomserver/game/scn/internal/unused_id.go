package internal

// 用于获取一个 1-maxnum的ID

import (
	"base/glog"
)

type UnuseId struct {
	ids      []uint32
	ids_flag []uint32
	maxnum   uint32
}

func (this *UnuseId) Init(maxnum uint32) {
	this.maxnum = maxnum
	this.ids = make([]uint32, 0)
	this.ids_flag = make([]uint32, maxnum+1)
	for id := uint32(1); id <= maxnum; id++ {
		this.ids = append(this.ids, id)
		this.ids_flag[id] = 1
	}
}

func (this *UnuseId) GetId() (uint32, bool) {
	if len(this.ids) == 0 {
		glog.Errorln("[UnuseId]: no id can get !!!!!!!!!!!!!!!!!!!!!!!!!")
		return 0, false
	}

	id := this.ids[0]
	this.ids = this.ids[1:]

	if this.ids_flag[id] == 0 {
		return this.GetId()
	}

	this.ids_flag[id] = 0
	return id, true
}

func (this *UnuseId) ReturnId(id uint32) {
	if id <= 0 || id > this.maxnum {
		return
	}
	if this.ids_flag[id] != 0 {
		return
	}
	this.ids_flag[id] = 1
	this.ids = append(this.ids, id)
}
