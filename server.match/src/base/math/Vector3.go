package math

import ()

type Vector3 struct {
	X, Y, Z float32
}

func (this *Vector3) Create() *Vector3 {
	this.X, this.Y, this.Z = 0, 0, 0
	return this
}
func (this *Vector3) Set(x, y, z float32) *Vector3 {
	this.X = x
	this.Y = y
	this.Z = z
	return this
}
func (this *Vector3) CopyFrom(v *Vector3) *Vector3 {
	this.X = v.X
	this.Y = v.Y
	this.Z = v.Z
	return this
}
func (this *Vector3) Clone() *Vector3 {
	return &Vector3{this.X, this.Y, this.Z}
}
func (this *Vector3) Add(v *Vector3) *Vector3 {
	return &Vector3{this.X + v.X, this.Y + v.Y, this.Z + v.Z}
}
func (this *Vector3) Substract(v *Vector3) *Vector3 {
	return &Vector3{this.X - v.X, this.Y - v.Y, this.Z - v.Z}
}
func (this *Vector3) IncreaseBy(v *Vector3) *Vector3 {
	this.X += v.X
	this.Y += v.Y
	this.Z += v.Z
	return this
}
func (this *Vector3) IncreaseByF32(v float32) *Vector3 {
	this.X += v
	this.Y += v
	this.Z += v
	return this
}
func (this *Vector3) DecreaseBy(v *Vector3) *Vector3 {
	this.X -= v.X
	this.Y -= v.Y
	this.Z -= v.Z
	return this
}
func (this *Vector3) DecreaseByF32(v float32) *Vector3 {
	this.X -= v
	this.Y -= v
	this.Z -= v
	return this
}
func (this *Vector3) ScaleBy(v float32) *Vector3 {
	this.X *= v
	this.Y *= v
	this.Z *= v
	return this
}
func (this *Vector3) Magnitude() float32 {
	return SqrtF32(this.X*this.X + this.Y*this.Y + this.Z*this.Z)
}
func (this *Vector3) SqrMagnitude() float32 {
	return this.X*this.X + this.Y*this.Y + this.Z*this.Z
}
func (this *Vector3) Normalize() *Vector3 {
	temp := this.Clone()
	temp.NormalizeSelf()
	return temp
}
func (this *Vector3) NormalizeSelf() float32 {
	var magn = this.Magnitude()
	if IsZero(magn) {
		return 0
	}
	this.X = this.X / magn
	this.Y = this.Y / magn
	this.Z = this.Z / magn
	return magn
}
func (this *Vector3) IsZero() bool {
	return IsZero(this.SqrMagnitude())
}
func (this *Vector3) Dot(v *Vector3) float32 {
	return this.X*v.X + this.Y*v.Y + this.Z*v.Z
}

func (v *Vector3) MulC(u *Vector3) *Vector3 {
	return &Vector3{v.X * u.X, v.Y * u.Y, v.Z * u.Z}
}

func (this *Vector3) Cross(v *Vector3) *Vector3 {
	temp := &Vector3{}
	temp.X = this.Y*v.Z - this.Z*v.Y
	temp.Y = this.Z*v.X - this.X*v.Z
	temp.Z = this.X*v.Y - this.Y*v.X
	return temp
}
func (this *Vector3) Equal(v *Vector3) bool {
	return EqualF32(this.X, v.X) && EqualF32(this.Y, v.Y) && EqualF32(this.Z, v.Z)
}
func (this *Vector3) IsValid() bool {
	return !(IsNaNF32(this.X) || IsNaNF32(this.Y) || IsNaNF32(this.Z))
}

var Vector3_0 *Vector3 = &Vector3{0, 0, 0}
var Vector3_X *Vector3 = &Vector3{1, 0, 0}
var Vector3_Y *Vector3 = &Vector3{0, 1, 0}
var Vector3_Z *Vector3 = &Vector3{0, 0, 1}
