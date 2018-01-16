package math

import "math"

const (
	TWO_PI = math.Pi * 2
	PI     = math.Pi
	RTOD   = 180 / PI
	DTOR   = PI / 180
)

type AngleInDegrees float32

func (this AngleInDegrees) ToRadians() AngleInRadians {
	return AngleInRadians(float32(this) * DTOR)
}
func (this AngleInDegrees) ToFloat32() float32 {
	return float32(this)
}

type AngleInRadians float32

func (this AngleInRadians) ToDegrees() AngleInDegrees {
	return AngleInDegrees(float32(this) * RTOD)
}
func (this AngleInRadians) ToFloat32() float32 {
	return float32(this)
}

func AngleToRadians(angle float64) float64 {
	return angle * DTOR
}
func RadiansToAngle(v float64) float64 {
	return v * RTOD
}

const RAD0 AngleInRadians = 0
const RADHALF AngleInRadians = 0.008726646
const RAD1 AngleInRadians = 0.017453293
const RAD5 AngleInRadians = 0.087266463
const RAD10 AngleInRadians = 0.174532925
const RAD15 AngleInRadians = 0.261799388
const RAD22HALF AngleInRadians = 0.392699082
const RAD30 AngleInRadians = 0.523598776
const RAD45 AngleInRadians = 0.785398163
const RAD60 AngleInRadians = 1.047197551
const RAD90 AngleInRadians = 1.570796327
const RAD120 AngleInRadians = 2.094395102
const RAD135 AngleInRadians = 2.35619449
const RAD180 AngleInRadians = 3.141592654

//----------
func NormalizeAngleZeroToTwoPI(angle AngleInRadians) AngleInRadians {
	result := angle
	for result < 0 {
		result += TWO_PI
	}
	for result >= TWO_PI {
		result -= TWO_PI
	}
	return result
}
func NormalizeAngleNegPIToPI(angle AngleInRadians) AngleInRadians {
	result := angle
	for result < -PI {
		result += TWO_PI
	}
	for result >= PI {
		result -= TWO_PI
	}
	return result
}

//取角度的标准值(0~360)
func AbsAngle(a float64) float64 {
	if a < 0 {
		n := math.Ceil(-a / 360)
		a += n * 360
		return a
	}
	return math.Mod(a, 360)
}

//取两个角度的较小夹角 -180~180
//return: offset,flag
func AngleMinSub(a float64, b float64) (float64, float64) {
	var flag float64 = 0
	offset := AbsAngle(a - b)
	if offset == 0 {
		return 0, 0
	}
	if offset > 180 {
		offset = offset - 360
		flag = -1
	} else {
		flag = 1
	}
	return offset, flag
}

func AngleBetween(v1, v2 *Vector3) AngleInRadians {
	v1n := v1.Normalize()
	v2n := v2.Normalize()
	return AngleInRadians(ACosF32(v1n.Dot(v2n)))
}

//count-clockwise is positive
func AngleBetween2DWithSign(v1, v2 *Vector3) AngleInRadians {
	v1Side := &Vector3{v1.Z, 0, -v1.X}
	v1Side.NormalizeSelf()

	v12D := v1.Clone()
	v12D.Y = 0
	v12D.NormalizeSelf()

	v22D := v2.Clone()
	v22D.Y = 0
	v22D.NormalizeSelf()

	dotSide := v1Side.Dot(v22D)
	ret := ACosF32(v12D.Dot(v22D))
	return AngleInRadians(SelF32(dotSide, -ret, ret))
}
func AngleToVector3D(angle AngleInRadians) *Vector3 {
	return &Vector3{CosF32(angle.ToFloat32()), 0, SinF32(angle.ToFloat32())}
}
func Vector3DToAngle(v *Vector3) AngleInRadians {
	lenOfVector := SqrtF32(v.X*v.X + v.Z*v.Z)
	if IsZero(lenOfVector) {
		return 0
	}
	ret := ACosF32(v.X / lenOfVector)
	if v.Z < 0 {
		ret = TWO_PI - ret
	}
	return AngleInRadians(ret)
}
