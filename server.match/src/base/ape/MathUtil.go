package ape

import (
	"math"
)

const (
	ONE_EIGHTY_OVER_PI = float32(180) / float32(math.Pi)
	PI_OVER_ONE_EIGHTY = float32(math.Pi) / float32(180)
)

//超过这个距离预先剔除
var ColliderTestDistanceMax float32 = 10

type MathUtil struct {
}

func (this *MathUtil) Clamp(n, min, max float32) float32 {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

func (this *MathUtil) Sign(val float32) int32 {
	if val < 0 {
		return -1
	}
	return 1
}
