package ape

import (
	"base/math"
)

//粒子对象类型
type ParticleType uint8

const (
	ParticleTypeCircle ParticleType = iota + 1
	ParticleTypeRect
)

type IAbstractParticle interface {
	IAbstractItem

	SetDisplay(d *DisplayObject, offsetX float32, offsetY float32, rotation float32)
	AddForce(f *math.Vector2)
	AddMasslessForce(f *math.Vector2)
	InitDisplay()
	Update(dt2 float32)
	GetComponents(collisionNormal *math.Vector2) *Collision
	ResolveCollision(mtd *math.Vector2, vel *math.Vector2, n *math.Vector2, d float32, o int32, p IAbstractParticle)
	GetInvMass() float32
	GetFixed() bool
	SetFixed(b bool)
	GetCenter() *math.Vector2
	GetCollidable() bool
	SetCollidable(b bool)
	GetMultisample() int32
	GetSamp() *math.Vector2
	GetCurr() *math.Vector2
	GetPrev() *math.Vector2
	GetElasticity() float32
	GetFriction() float32
	GetVelocity() *math.Vector2
	SetVelocity(v *math.Vector2)
	SetPx(x float32)
	SetPy(y float32)

	GetPx() float32
	GetPy() float32
	GetParticleType() ParticleType
}
type AbstractParticle struct {
	AbstractItem

	curr *math.Vector2

	prev *math.Vector2

	samp *math.Vector2

	interval *Interval

	forces    *math.Vector2
	temp      *math.Vector2
	collision *Collision

	elasticity float32
	mass       float32
	invMass    float32
	friction   float32

	opposition float32 //阻力系数0~1

	fixed      bool
	collidable bool

	center      *math.Vector2
	multisample int32
	particType  ParticleType
}

func (this *AbstractParticle) GetElasticity() float32 {
	return this.elasticity
}
func (this *AbstractParticle) GetPrev() *math.Vector2 {
	return this.prev
}
func (this *AbstractParticle) GetCurr() *math.Vector2 {
	return this.curr
}
func (this *AbstractParticle) GetSamp() *math.Vector2 {
	return this.samp
}
func (this *AbstractParticle) GetMultisample() int32 {
	return this.multisample
}
func (this *AbstractParticle) GetFixed() bool {
	return this.fixed
}

func (this *AbstractParticle) SetFixed(b bool) {
	this.fixed = b
}
func (this *AbstractParticle) SetOpposition(oppo float32) {
	this.opposition = oppo
}
func (this *AbstractParticle) GetOpposition() float32 {
	return this.opposition
}

func (this *AbstractParticle) GetMass() float32 {
	return this.mass
}

func (this *AbstractParticle) SetMass(m float32) {
	if m <= 0 {
		panic("mass may not be set <= 0")
	}
	this.mass = m
	this.invMass = 1 / this.mass
}
func (this *AbstractParticle) GetCenter() *math.Vector2 {
	this.center.Set(this.GetPx(), this.GetPy())
	return this.center
}

func (this *AbstractParticle) GetFriction() float32 {
	return this.friction
}
func (this *AbstractParticle) SetFriction(f float32) {
	if f < 0 || f > 1 {
		panic("Legal friction must be >= 0 and <=1")
	}
	this.friction = f
}

func (this *AbstractParticle) GetPostion() *math.Vector2 {
	return &math.Vector2{X: this.curr.X, Y: this.curr.Y}
}

func (this *AbstractParticle) SetPosition(p *math.Vector2) {
	this.curr.CopyFrom(p)
	this.prev.CopyFrom(p)
}

func (this *AbstractParticle) GetPx() float32 {
	return this.curr.X
}
func (this *AbstractParticle) SetPx(x float32) {
	this.curr.X = x
	this.prev.X = x
}

func (this *AbstractParticle) GetPy() float32 {
	return this.curr.Y
}
func (this *AbstractParticle) SetPy(y float32) {
	this.curr.Y = y
	this.prev.Y = y
}
func (this *AbstractParticle) GetVelocity() *math.Vector2 {
	return this.curr.Substract(this.prev)
}
func (this *AbstractParticle) SetVelocity(v *math.Vector2) {
	this.prev = this.curr.Substract(v)
}

func (this *AbstractParticle) GetCollidable() bool {
	return this.collidable
}

func (this *AbstractParticle) SetCollidable(b bool) {
	this.collidable = b
}

func (this *AbstractParticle) SetDisplay(d *DisplayObject, offsetX float32, offsetY float32, rotation float32) {
	this.displayObject = d
	this.displayObjectRotation = rotation
	this.displayObjectOffset = &math.Vector2{
		X: offsetX,
		Y: offsetY,
	}
}
func (this *AbstractParticle) AddForce(f *math.Vector2) {
	this.forces.IncreaseBy(f.Mult(this.invMass))
}
func (this *AbstractParticle) AddMasslessForce(f *math.Vector2) {
	this.forces.IncreaseBy(f)
}

func (this *AbstractParticle) Update(dt2 float32) {
	if this.GetFixed() {
		return
	}
	//this.AddForce(APEngine_GetMe().force)
	//this.AddMasslessForce(APEngine_GetMe().masslessForce)
	this.temp.CopyFrom(this.curr)

	nv := this.GetVelocity().Add(this.forces.ScaleBy(dt2))

	if this.opposition > 0 {
		nv.ScaleBy(this.opposition)
	}

	this.curr.IncreaseBy(nv)
	this.prev.CopyFrom(this.temp)

	// clear the forces
	this.forces.Set(0, 0)
}

func (this *AbstractParticle) InitDisplay() {
	this.displayObject.x = this.displayObjectOffset.X
	this.displayObject.y = this.displayObjectOffset.Y
	this.displayObject.rotation = this.displayObjectRotation
	//|||this.sprite.addChild(this.displayObject)
}

func (this *AbstractParticle) GetComponents(collisionNormal *math.Vector2) *Collision {

	vel := this.GetVelocity()
	vdotn := collisionNormal.Dot(vel)
	this.collision.vn = collisionNormal.Mult(vdotn)
	this.collision.vt = vel.Substract(this.collision.vn)
	return this.collision
}

func (this *AbstractParticle) ResolveCollision(mtd *math.Vector2, vel *math.Vector2, n *math.Vector2, d float32, o int32, p IAbstractParticle) {
	this.curr.IncreaseBy(mtd)
	this.prev.Set(this.curr.X, this.curr.Y)
	//this.SetVelocity(vel)
}

func (this *AbstractParticle) GetInvMass() float32 {
	if this.GetFixed() {
		return 0
	} else {
		return this.invMass
	}
}
