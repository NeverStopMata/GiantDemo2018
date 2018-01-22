package math

import (
	"math/rand"
	"time"
)

func ReSeed() {
	rand.Seed(time.Now().UnixNano())
}

func GetRandomInt(floor int, ceiling int) int {
	if floor == ceiling {
		return floor
	} else if floor > ceiling {
		floor, ceiling = ceiling, floor
	}
	return rand.Intn(ceiling-floor) + floor
}

func GetRandomF32(floor float32, ceiling float32) float32 {
	if floor == ceiling {
		return floor
	} else if floor > ceiling {
		floor, ceiling = ceiling, floor
	}
	return rand.Float32()*(ceiling-floor) + floor
}

func GetRandomPosByRectangle(x, y, width, height float32) *Vector3 {
	Vec3 := &Vector3{}
	Vec3.X = GetRandomF32(x, x+width)
	Vec3.Z = GetRandomF32(y, y+height)
	return Vec3
}
func GetRandomPosByRectangleEx(leftUp *Vector3, rightBottom *Vector3) *Vector3 {
	return GetRandomPosByRectangle(leftUp.X, leftUp.Z, rightBottom.X-leftUp.X, rightBottom.Z-leftUp.Z)
}
