package util

import (
	"fmt"
	"math"
)

func Dot2D(v1, v2 *Vector2) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

/**
*2D点
 */
type Vector2 struct {
	X float64
	Y float64
}

func (v *Vector2) SubMethod(value *Vector2) *Vector2 {
	return &Vector2{v.X - value.X, v.Y - value.Y}
}

func (v *Vector2) AddMethod(value *Vector2) *Vector2 {
	return &Vector2{v.X + value.X, v.Y + value.Y}
}

func (v *Vector2) MultiMethod(value float64) *Vector2 {
	return &Vector2{v.X * value, v.Y * value}
}

func (v *Vector2) SetMethod(value *Vector2) {
	v.X = value.X
	v.Y = value.Y
}
func (v *Vector2) IncreaseBy(value *Vector2) {
	v.X += value.X
	v.Y += value.Y
}
func (v *Vector2) DecreaseBy(value *Vector2) {
	v.X -= value.X
	v.Y -= value.Y
}
func (v *Vector2) ScaleBy(value float64) {
	v.X *= value
	v.Y *= value
}
func (v1 *Vector2) CopyFrom(value *Vector2) (v *Vector2) {
	v.X = value.X
	v.Y = value.Y
	return v
}

func (v *Vector2) Clone() *Vector2 {
	return &Vector2{v.X, v.Y}
}

func (v *Vector2) Magnitude() float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y))
}

func (v *Vector2) DistanceTo(v2 *Vector2) float64 {
	return math.Sqrt((v.X-v2.X)*(v.X-v2.X) + (v.Y-v2.Y)*(v.Y-v2.Y))
}

func (v *Vector2) SqrMagnitudeTo(v2 *Vector2) float64 {
	return (v.X-v2.X)*(v.X-v2.X) + (v.Y-v2.Y)*(v.Y-v2.Y)
}
func (v *Vector2) SqrMagnitude() float64 {

	return (v.X*v.X + v.Y*v.Y)
}
func (v *Vector2) Normalize() Vector2 {

	var temp Vector2
	var magn = v.Magnitude()
	if IsZero(magn) {
		return temp
	}
	temp.X = v.X / magn
	temp.Y = v.Y / magn
	return temp
}

func (v *Vector2) NormalizeSelf() {
	var magn = v.Magnitude()
	if IsZero(magn) {
		return
	}
	v.X = v.X / magn
	v.Y = v.Y / magn
}

func (v *Vector2) Set(px, py float64) {
	v.X = px
	v.Y = py
}

func (v *Vector2) SetVector(value *Vector2) {
	v.X = value.X
	v.Y = value.Y
}

func (v *Vector2) String() string {
	return fmt.Sprintf("(%.2f,%.2f)", v.X, v.Y)
}

func (v *Vector2) IsEmpty() bool {
	return v.X == 0 && v.Y == 0
}

// 方向是否一致
func IsSameDir(dir, pos1, pos2 *Vector2) bool {
	if dir != nil {
		if Dot2D(dir, &Vector2{pos1.X - pos2.X, pos1.Y - pos2.Y}) > 0 {
			return true
		}
	}
	return false
}
