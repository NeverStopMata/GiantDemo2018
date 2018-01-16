package ape

import (
	"base/math"
)

type CircleParticle struct {
	AbstractParticle

	radius float32
}

func NewCircleParticle(x float32, y float32, radius float32) *CircleParticle {

	c := &CircleParticle{}

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

	c.SetMass(1)
	c.elasticity = 0
	c.SetFriction(0)
	c.center = vect.Clone()
	c.multisample = 0

	c.radius = radius
	return c
}

func (this *CircleParticle) Init() {
	this.CleanUp()
	if this.displayObject != nil {
		this.InitDisplay()
	} else {

	}
	this.Paint()
}
func (this *CircleParticle) Paint() {
	//println(this.curr.X)
	this.GetSprite().X = this.curr.X
	this.GetSprite().Y = this.curr.Y
}
func (this *CircleParticle) GetProjection(axis *math.Vector2) *Interval {
	c := this.samp.Dot(axis)
	this.interval.min = c - this.radius
	this.interval.max = c + this.radius
	return this.interval
}
func (this *CircleParticle) GetIntervalX() *Interval {
	this.interval.min = this.curr.X - this.radius
	this.interval.max = this.curr.X + this.radius
	return this.interval
}
func (this *CircleParticle) GetIntervalY() *Interval {
	this.interval.min = this.curr.Y - this.radius
	this.interval.max = this.curr.Y + this.radius
	return this.interval
}
func (this *CircleParticle) GetRadius() float32 {
	return this.radius
}
func (this *CircleParticle) SetRadius(r float32) {
	this.radius = r
}
func (this *CircleParticle) GetParticleType() ParticleType {
	return ParticleTypeCircle
}
