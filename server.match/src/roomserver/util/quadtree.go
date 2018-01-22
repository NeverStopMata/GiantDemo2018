/**
 * \brief 四叉数
 */
package util

import (
	"fmt"
	"math"
)

type Square struct {
	X         float64
	Y         float64
	Left      float64
	Right     float64
	Top       float64 /** 上下左右象限*/
	Bottom    float64
	Radius    float64
	indexX    int
	indexY    int
	oldIndexX int
	oldIndexY int
}

func CirleToSquare(px, py, Radius float64) *Square {
	s := Square{X: px, Y: py}
	s.SetRadius(Radius)
	return &s
}

func NewSquare(px, py, pw, ph float64) *Square {
	sideX := pw * 0.5
	sideY := ph * 0.5
	return &Square{X: px, Y: py, Left: px - sideX, Right: px + sideX, Top: py + sideY, Bottom: py - sideY}
}

func (s *Square) Reset() {
	s.X = 0
	s.Y = 0
	s.Left = 0
	s.Right = 0
	s.Top = 0
	s.Bottom = 0
	s.Radius = 0
}

func (s *Square) MDistance(s2 *Square) float64 {
	return math.Abs(s.X-s2.X) + math.Abs(s.Y-s2.Y)
}

func (s *Square) UpdatePos(px, py float64) {
	s.X = px
	s.Y = py
	s.Left = s.X - s.Radius
	s.Right = s.X + s.Radius
	s.Top = s.Y + s.Radius
	s.Bottom = s.Y - s.Radius
}

func (s *Square) SetRadius(Radius float64) {
	s.Radius = Radius
	s.Left = s.X - Radius
	s.Right = s.X + Radius
	s.Top = s.Y + Radius
	s.Bottom = s.Y - Radius
}

func (s *Square) GetRadius() float64 {
	return s.Radius
}

func (s *Square) ContainsRect(rect *Square) bool {
	return s.Left < rect.Left && s.Right > rect.Right && s.Top > rect.Top && s.Bottom < rect.Bottom
}

//测试是否包含，如果碰撞返回补正向量
func (s *Square) ContainsCircle(px, py, Radius float64) (float64, float64, bool) {
	//Left
	if px-Radius < s.Left {
		return 1, 0, false
	}
	//Right
	if px+Radius > s.Right {
		return -1, 0, false
	}

	//Top
	if py+Radius > s.Top {
		return 0, -1, false
	}
	//Bottom
	if py-Radius < s.Bottom {
		return 0, 1, false
	}
	return 0, 0, true
}

func (s *Square) Union(s2 *Square) *Square {
	if s2.IsEmpty() {
		return s
	}
	if s.IsEmpty() {
		return s2
	}
	Left := math.Min(s.Left, s2.Left)
	Right := math.Max(s.Right, s2.Right)
	Top := math.Max(s.Top, s2.Top)
	Bottom := math.Min(s.Bottom, s2.Bottom)

	px := (s.Left + s.Right) * 0.5
	py := (s.Top + s.Bottom) * 0.5
	return &Square{X: px, Y: py, Left: Left, Right: Right, Top: Top, Bottom: Bottom}
}

func (s *Square) UnionFrom(s2 *Square) {
	//glog.Info("[Square] source:", s, ", union:", s2)
	if s2.IsEmpty() {
		return
	}
	if s.IsEmpty() {
		s.CopyFrom(s2)
		return
	}
	s.Left = math.Min(s.Left, s2.Left)
	s.Bottom = math.Min(s.Bottom, s2.Bottom)
	s.Right = math.Max(s.Right, s2.Right)
	s.Top = math.Max(s.Top, s2.Top)
	s.X = (s.Left + s.Right) * 0.5
	s.Y = (s.Top + s.Bottom) * 0.5
}

func (s *Square) Scale(rate float64) {
	w := (s.Right - s.Left) * rate * 0.5
	h := (s.Top - s.Bottom) * rate * 0.5
	s.Left = s.X - w
	s.Right = s.X + w
	s.Top = s.Y + h
	s.Bottom = s.Y - h
}

//缩放到，这里是长宽的一半
func (s *Square) ScaleTo(pw, ph float64) {
	s.Left = s.X - pw
	s.Right = s.X + pw
	s.Top = s.Y + ph
	s.Bottom = s.Y - ph
}

func (r1 *Square) Intersects(r2 *Square) bool {
	//	glog.Info("r2.Left > r1.Right :", r2.Left > r1.Right)
	//	glog.Info("r2.Right < r1.Left :", r2.Right < r1.Left)
	//	glog.Info("r2.Top > r1.Bottom :", r2.Top > r1.Bottom)
	//	glog.Info("r2.Bottom < r1.Top :", r2.Bottom < r1.Top)
	return !(r2.Left > r1.Right || r2.Right < r1.Left || r2.Top < r1.Bottom || r2.Bottom > r1.Top)
}

//正方形和圆形相交
func (r1 *Square) IntersectsCircle(X, Y, r float64) bool {
	addR := r + r1.Radius
	out := X-r1.X > addR || r1.X-X > addR || Y-r1.Y > addR || r1.Y-Y > addR
	return !out
}

func (s *Square) IsEmpty() bool {
	return s.Right-s.Left < 0.001 || s.Top-s.Bottom < 0.001
}

func (s *Square) CopyFrom(s2 *Square) {
	s.X = s2.X
	s.Y = s2.Y
	s.Left = s2.Left
	s.Right = s2.Right
	s.Top = s2.Top
	s.Bottom = s2.Bottom
	s.Radius = s2.Radius
}

func (s *Square) String() string {
	return fmt.Sprintf("(%.2f:%.2f,%.2f,%.2f,%.2f,%.2f)", s.X, s.Y, s.Left, s.Right, s.Top, s.Bottom)
}
