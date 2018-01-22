package math

type Vector2 struct {
	X, Y float32
}

//角度转向量
func AngleToVector2D(angle float32) *Vector2 {
	return &Vector2{CosF32(angle), SinF32(angle)}
}

//2d向量转角度
func Vector2DToAngle(v *Vector2) AngleInDegrees {
	lenOfVector := SqrtF32(v.X*v.X + v.Y*v.Y)
	if IsZero(lenOfVector) {
		return 0
	}
	ret := ACosF32(v.X / lenOfVector)
	if v.Y < 0 {
		ret = TWO_PI - ret
	}
	return AngleInRadians(ret).ToDegrees()
}

//2d向量转角度
func Vector2DToAngleXY(vx, vy float32) AngleInDegrees {
	lenOfVector := SqrtF32(vx*vx + vy*vy)
	if IsZero(lenOfVector) {
		return 0
	}
	ret := ACosF32(vx / lenOfVector)
	if vy < 0 {
		ret = TWO_PI - ret
	}
	return AngleInRadians(ret).ToDegrees()
}

func (this *Vector2) Create() *Vector2 {
	this.X, this.Y = 0, 0
	return this
}
func (this *Vector2) Set(X, Y float32) *Vector2 {
	this.X = X
	this.Y = Y
	return this
}
func (this *Vector2) CopyFrom(v *Vector2) *Vector2 {
	this.X = v.X
	this.Y = v.Y
	return this
}
func (this *Vector2) Clone() *Vector2 {
	return &Vector2{this.X, this.Y}
}
func (this *Vector2) Add(v *Vector2) *Vector2 {
	return &Vector2{this.X + v.X, this.Y + v.Y}
}
func (this *Vector2) Substract(v *Vector2) *Vector2 {
	return &Vector2{this.X - v.X, this.Y - v.Y}
}
func (this *Vector2) IncreaseBy(v *Vector2) *Vector2 {
	this.X += v.X
	this.Y += v.Y
	return this
}
func (this *Vector2) DecreaseBy(v *Vector2) *Vector2 {
	this.X -= v.X
	this.Y -= v.Y
	return this
}
func (this *Vector2) ScaleBy(v float32) *Vector2 {
	this.X *= v
	this.Y *= v
	return this
}
func (this *Vector2) Magnitude() float32 {
	return SqrtF32(this.X*this.X + this.Y*this.Y)
}
func (this *Vector2) SqrMagnitude() float32 {
	return this.X*this.X + this.Y*this.Y
}
func (this *Vector2) Normalize() *Vector2 {
	temp := this.Clone()
	temp.NormalizeSelf()
	return temp
}
func (this *Vector2) NormalizeSelf() float32 {
	var magn = this.Magnitude()
	if IsZero(magn) {
		return 0
	}
	this.X = this.X / magn
	this.Y = this.Y / magn
	return magn
}
func (this *Vector2) IsZero() bool {
	return IsZero(this.SqrMagnitude())
}
func (this *Vector2) Dot(v *Vector2) float32 {
	return this.X*v.X + this.Y*v.Y
}
func (this *Vector2) Mult(s float32) *Vector2 {
	return &Vector2{this.X * s, this.Y * s}
}

func (this *Vector2) DivEquals(s float32) *Vector2 {
	if s == 0 {
		s = 0.0001
	}
	this.X /= s
	this.Y /= s
	return this
}
func (this *Vector2) Distance(v *Vector2) float32 {
	delta := this.Substract(v)
	return delta.Magnitude()
}
