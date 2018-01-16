package math

type Quaternion struct {
	X, Y, Z, W float32
}

func (this *Quaternion) Create() *Quaternion {
	this.X, this.Y, this.Z, this.W = 0, 0, 0, 0
	return this
}

func (this *Quaternion) Clone() *Quaternion {
	return &Quaternion{this.X, this.Y, this.Z, this.W}
}

func (this *Quaternion) FromAxisAngle(axis *Vector3, angle AngleInRadians) *Quaternion {
	halfAngle := 0.5 * angle.ToFloat32()
	sinValue := SinF32(halfAngle)
	this.W = CosF32(halfAngle)
	axisN := axis.Normalize()
	this.X = sinValue * axisN.X
	this.Y = sinValue * axisN.Y
	this.Z = sinValue * axisN.Z
	return this
}

func (this *Quaternion) FromYawPitchRoll(yaw AngleInRadians, pitch AngleInRadians, roll AngleInRadians) *Quaternion {
	qp := new(Quaternion).FromAxisAngle(&Vector3{1, 0, 0}, pitch)
	qr := new(Quaternion).FromAxisAngle(&Vector3{0, 0, 1}, roll)
	this.FromAxisAngle(&Vector3{0, 1, 0}, yaw)
	this.MultiplyBy(qp).MultiplyBy(qr)
	return this
}

func (this *Quaternion) MultiplyBy(q *Quaternion) *Quaternion {
	this.W = this.W*q.W - this.X*q.X - this.Y*q.Y - this.Z*q.Z
	this.X = this.W*q.X + this.X*q.W + this.Y*q.Z - this.Z*q.Y
	this.Y = this.W*q.Y + this.Y*q.W + this.Z*q.X - this.X*q.Z
	this.Z = this.W*q.Z + this.Z*q.W + this.X*q.Y - this.Y*q.X
	return this
}

func (this *Quaternion) Conjugate() *Quaternion {
	this.X, this.Y, this.Z = -this.X, -this.Y, -this.Z
	return this
}

func (this *Quaternion) MultiplyVector(v *Vector3) *Vector3 {
	qvec := &Vector3{this.X, this.Y, this.Z}
	uv := v.Cross(qvec)
	uuv := uv.Cross(qvec)
	uv.ScaleBy(2 * this.W)
	uuv.ScaleBy(2)
	return v.Add(uv).IncreaseBy(uuv)
}
