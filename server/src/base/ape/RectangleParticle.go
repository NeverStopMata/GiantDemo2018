package ape

import (
	"base/math"
)

type RectangleParticle struct {
	AbstractParticle

	extents []float32
	axes    []*math.Vector2
	radian  float32
}

func NewRectangleParticle(x float32, y float32, width float32, height float32) *RectangleParticle {

	c := &RectangleParticle{
		extents: make([]float32, 0),
		axes:    make([]*math.Vector2, 0),
	}

	vect := &math.Vector2{}
	vect.Create()

	c.interval = &Interval{0, 0}
	c.curr = vect.Clone()
	c.curr.Set(x, y)
	c.prev = vect.Clone()
	c.prev.Set(x, y)
	c.samp = vect.Clone()
	c.temp = vect.Clone()
	c.fixed = false
	c.forces = vect.Clone()
	c.collision = NewCollision(vect.Clone(), vect.Clone())
	c.collidable = true

	c.mass = 1
	c.elasticity = 0
	c.friction = 0
	c.center = vect.Clone()
	c.multisample = 0

	c.extents = append(c.extents, width/2, height/2)
	c.axes = append(c.axes, vect.Clone(), vect.Clone())
	c.SetRadian(0)

	return c
}

func (this *RectangleParticle) GetRadian() float32 {
	return this.radian
}

func (this *RectangleParticle) SetRadian(t float32) {
	this.radian = t
	this.SetAxes(t)
}

func (this *RectangleParticle) GetAngle() float32 {
	return this.radian * ONE_EIGHTY_OVER_PI
}
func (this *RectangleParticle) SetAngle(a float32) {
	this.radian = a * PI_OVER_ONE_EIGHTY
}

func (this *RectangleParticle) Init() {
	this.CleanUp()
	if this.displayObject != nil {
		this.InitDisplay()
	} else {

		// w := this.extents[0] * 2
		// h := this.extents[1] * 2
	}
	this.Paint()
}
func (this *RectangleParticle) Paint() {
	this.GetSprite().X = this.curr.X
	this.GetSprite().Y = this.curr.Y
	this.GetSprite().rotation = this.GetAngle()
}

func (this *RectangleParticle) SetWidth(w float32) {
	this.extents[0] = w / float32(2)
}
func (this *RectangleParticle) GetWidth() float32 {
	return this.extents[0] * 2
}

func (this *RectangleParticle) SetHeight(h float32) {
	this.extents[1] = h / float32(2)
}
func (this *RectangleParticle) GetHeight() float32 {
	return this.extents[1] * 2
}

func (this *RectangleParticle) GetProjection(axis *math.Vector2) *Interval {

	radius := this.extents[0]*math.AbsF32(axis.Dot(this.axes[0])) + this.extents[1]*math.AbsF32(axis.Dot(this.axes[1]))
	c := this.samp.Dot(axis)

	this.interval.min = c - radius
	this.interval.max = c + radius
	return this.interval
}

func (this *RectangleParticle) SetAxes(t float32) {

	s := math.SinF32(t)
	c := math.CosF32(t)

	this.axes[0].X = c
	this.axes[0].Y = s
	this.axes[1].X = -s
	this.axes[1].Y = c
}
func (this *RectangleParticle) GetParticleType() ParticleType {
	return ParticleTypeRect
}
