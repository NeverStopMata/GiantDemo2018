package internal

// 用于获取地图上，一个可以放置的点

import (
	"base/glog"
	"math"
	"math/rand"
)

type UnusedPos struct {
	ids      []int
	ids_flag map[int]int
}

func (this *UnusedPos) Init(roomSize float64) {
	this.ids = make([]int, int(roomSize*roomSize))
	this.ids_flag = make(map[int]int)

	tmplist := rand.Perm(int(roomSize * roomSize))
	index := 0
	for i := 0; i < int(roomSize); i++ {
		for j := 0; j < int(roomSize); j++ {
			id := i*10000 + j

			if _, ok := this.ids_flag[id]; ok {
				panic("UnusedPos::Init fail. #1")
			}

			if this.ids[tmplist[index]] != 0 {
				panic("UnusedPos::Init fail. #2")
			}

			this.ids[tmplist[index]] = id
			index++

			this.ids_flag[id] = 1
		}
	}

	if len(this.ids) != len(this.ids_flag) {
		panic("UnusedPos::Init fail. #3")
	}
}

func (this *UnusedPos) AppendFixedPos(x, y int) {
	id := x*10000 + y
	this.ids_flag[id] = 0
}

func (this *UnusedPos) GetPos() (x, y float64, op bool) {

	if len(this.ids) == 0 {
		glog.Errorln("no pos can get !!!!!!!!!!!!!!!!!!!!!!!!!")
		return float64(0), float64(0), false
	}

	id := this.ids[0]
	this.ids = this.ids[1:]

	if value, ok := this.ids_flag[id]; (!ok) || value == 0 {
		return this.GetPos()
	}

	x0 := int(math.Floor(float64(id) / 10000))
	y0 := id % 10000

	this.ids_flag[id] = 0
	return float64(x0), float64(y0), true
}

func (this *UnusedPos) ReturnPos(x, y float64) {
	id := int(x)*10000 + int(y)
	if value, ok := this.ids_flag[id]; (!ok) || value != 0 {
		return
	}
	this.ids_flag[id] = 1
	this.ids = append(this.ids, id)
}

func (this *UnusedPos) GetRandomPos() (x, y float64, op bool) {
	x, y, op = this.GetPos()
	if op {
		this.ReturnPos(x, y)
	}
	return
}
