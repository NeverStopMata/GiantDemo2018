package math

const DIM_4 = 4

type Matrix4 struct {
	data [DIM_4][DIM_4]float32
}

func (this *Matrix4) Create() *Matrix4 {
	for i := 0; i < DIM_4; i++ {
		for j := 0; j < DIM_4; j++ {
			this.data[i][j] = 0
		}
	}
	return this
}

func (this *Matrix4) Clone() *Matrix4 {
	m := new(Matrix4)
	for i := 0; i < DIM_4; i++ {
		for j := 0; j < DIM_4; j++ {
			m.data[i][j] = this.data[i][j]
		}
	}
	return m
}

func (this *Matrix4) GetRawData() [][DIM_4]float32 {
	return this.data[:][:]
}

func (this *Matrix4) SetColumn(cIndex int, v *Vector3, w float32) *Matrix4 {
	if cIndex < 0 || cIndex >= DIM_4 {
		return this
	}
	this.data[0][cIndex] = v.X
	this.data[1][cIndex] = v.Y
	this.data[2][cIndex] = v.Z
	this.data[3][cIndex] = w
	return this
}

func (this *Matrix4) SetRow(rIndex int, v *Vector3, w float32) *Matrix4 {
	if rIndex < 0 || rIndex >= DIM_4 {
		return this
	}
	this.data[rIndex][0] = v.X
	this.data[rIndex][1] = v.Y
	this.data[rIndex][2] = v.Z
	this.data[rIndex][3] = w
	return this
}

func (this *Matrix4) Multiply(m *Matrix4) *Matrix4 {
	ret := new(Matrix4)
	for i := 0; i < DIM_4; i++ {
		for j := 0; j < DIM_4; j++ {
			ret.data[i][j] = this.data[i][0]*m.data[0][j] + this.data[i][1]*m.data[1][j] + this.data[i][2]*m.data[2][j] + this.data[i][3]*m.data[3][j]
		}
	}
	return ret
}

func (this *Matrix4) MultiplyBy(m *Matrix4) *Matrix4 {
	calM := this.Clone()
	for i := 0; i < DIM_4; i++ {
		for j := 0; j < DIM_4; j++ {
			this.data[i][j] = calM.data[i][0]*m.data[0][j] + calM.data[i][1]*m.data[1][j] + calM.data[i][2]*m.data[2][j] + calM.data[i][3]*m.data[3][j]
		}
	}
	return this
}

func (this *Matrix4) Identity() *Matrix4 {
	this.data[0][0] = 1
	this.data[0][1] = 0
	this.data[0][2] = 0
	this.data[0][3] = 0
	this.data[1][0] = 0
	this.data[1][1] = 1
	this.data[1][2] = 0
	this.data[1][3] = 0
	this.data[2][0] = 0
	this.data[2][1] = 0
	this.data[2][2] = 1
	this.data[2][3] = 0
	this.data[3][0] = 0
	this.data[3][1] = 0
	this.data[3][2] = 0
	this.data[3][3] = 1
	return this
}

func (this *Matrix4) MultiplyPoint3x4(v *Vector3) *Vector3 {
	ret := new(Vector3).Create()
	ret.X = this.data[0][0]*v.X + this.data[0][1]*v.Y + this.data[0][2]*v.Z + this.data[0][3]*1
	ret.Y = this.data[1][0]*v.X + this.data[1][1]*v.Y + this.data[1][2]*v.Z + this.data[1][3]*1
	ret.Z = this.data[2][0]*v.X + this.data[2][1]*v.Y + this.data[2][2]*v.Z + this.data[2][3]*1
	return ret
}

func (this *Matrix4) MultiplyVector(v *Vector3) *Vector3 {
	ret := new(Vector3).Create()
	ret.X = this.data[0][0]*v.X + this.data[0][1]*v.Y + this.data[0][2]*v.Z
	ret.Y = this.data[1][0]*v.X + this.data[1][1]*v.Y + this.data[1][2]*v.Z
	ret.Z = this.data[2][0]*v.X + this.data[2][1]*v.Y + this.data[2][2]*v.Z
	return ret
}
