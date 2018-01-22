package scn

// 场景中 球出生点 辅助类

import (
	bmath "base/math"
	"math"
	"math/rand"
	"roomserver/conf"
	"usercmd"
)

type SceneBirthPointHelper struct {
	birthPoints []*BirthPoint // 球出生点
}

func (this *SceneBirthPointHelper) AddBirthPoint(point *BirthPoint) {
	this.birthPoints = append(this.birthPoints, point)
}

func (this *SceneBirthPointHelper) RefreshBirthPoint(d float64, scene *Scene) {
	for _, birth := range this.birthPoints {
		birth.Refresh(d, scene)
	}
}

// 生成食物 、 动态障碍物
func (this *SceneBirthPointHelper) CreateAllBirthPoint(scene *Scene) {
	var xlist, ylist []float64
	items := conf.ConfigMgr_GetMe().GetXmlFoodItems(scene.SceneID())
	for _, item := range items.Items {
		ftype := item.FoodType
		fid := item.FoodId
		birthTime := item.BirthTime
		foodnum := conf.ConfigMgr_GetMe().GetFoodMapNum(scene.SceneID(), fid)
		size := float64(conf.ConfigMgr_GetMe().GetFoodSize(scene.SceneID(), fid))
		if foodnum > 0 {
			for i := 0; i < int(foodnum); i++ {
				x, y, op := scene.GetPos() //在场景中随机的分配出一个空闲的位置
				if op {
					point := NewBirthPoint(scene.NewBallId(), float32(x), float32(y), float32(size), float32(size), fid, ftype, birthTime, 1, scene)
					this.AddBirthPoint(point)
				}
				if ftype == uint16(usercmd.BallType_FoodNormal) {
					xlist = append(xlist, x)
					ylist = append(ylist, y)
				}
			}
		}
	}
	for i := 0; i < len(xlist); i++ {
		scene.ReturnPos(xlist[i], ylist[i])
	}
}

// 不同的食物，初始位置会做调整。如 食物(普通) 根据输入x,y 随机附近的值； 如食物(锤子) 根据输入x,y 对齐到地图上对应格子的中心 等等
func BallFood_InitPos(pos *bmath.Vector2, t usercmd.BallType, birthRadiusMin, birthRadiusMax float32) *bmath.Vector2 {
	switch t {
	case usercmd.BallType_FoodNormal:
		x := math.Floor(float64(pos.X)) + rand.Float64()*0.5
		y := math.Floor(float64(pos.Y)) + rand.Float64()*0.5
		posNew := &bmath.Vector2{float32(x), float32(y)}
		return posNew

	default:
		x := math.Floor(float64(pos.X)) + 0.25
		y := math.Floor(float64(pos.Y)) + 0.25
		posNew := &bmath.Vector2{float32(x), float32(y)}
		return posNew
	}
}
