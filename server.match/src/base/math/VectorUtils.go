package math

import (
	"strconv"
	"strings"
)

func Vector3FromString(v string) *Vector3 {
	vArray := strings.Split(v, ",")
	ret := new(Vector3).Create()
	lenOfV := len(vArray)
	if lenOfV > 0 {
		ret.X, _ = StringToF32(vArray[0])
	}
	if lenOfV > 1 {
		ret.Y, _ = StringToF32(vArray[1])
	}
	if lenOfV > 2 {
		ret.Z, _ = StringToF32(vArray[2])
	}
	return ret
}
func Vector3ZeroY(v *Vector3) *Vector3 {
	return &Vector3{v.X, 0, v.Z}
}
func Vector3FromVector2(v *Vector2) *Vector3 {
	return &Vector3{v.X, 0, v.Y}
}
func Vector2FromVector3(v *Vector3) *Vector2 {
	return &Vector2{v.X, v.Y}
}
func GetDirection2D(to *Vector3, from *Vector3) *Vector3 {
	tmp := to.Clone()
	tmp.DecreaseBy(from)
	tmp.Y = 0
	if IsZero(tmp.SqrMagnitude()) {
		tmp.Set(0, 0, 0)
	} else {
		tmp.NormalizeSelf()
	}
	return tmp
}
func GetDistance(v1 *Vector3, v2 *Vector3) float32 {
	tmp := v1.Clone()
	tmp.DecreaseBy(v2)
	return tmp.Magnitude()
}
func GetDistanceSqr(v1 *Vector3, v2 *Vector3) float32 {
	tmp := v1.Clone()
	tmp.DecreaseBy(v2)
	return tmp.SqrMagnitude()
}
func GetDistance2D(v1 *Vector3, v2 *Vector3) float32 {
	tmp := v1.Clone()
	tmp.DecreaseBy(v2)
	tmp.Y = 0
	return tmp.Magnitude()
}
func GetDistance2DSqr(v1 *Vector3, v2 *Vector3) float32 {
	tmp := v1.Clone()
	tmp.DecreaseBy(v2)
	tmp.Y = 0
	return tmp.SqrMagnitude()
}
func Blend3D(v1 *Vector3, v2 *Vector3, blendFactor float32) *Vector3 {
	bf := ClampF32(blendFactor, 0, 1)
	tmp1 := v1.Clone()
	tmp1.ScaleBy(1 - bf)
	tmp2 := v2.Clone()
	tmp2.ScaleBy(bf)
	tmp1.IncreaseBy(tmp2)
	return tmp1
}
func RotateVectorAroundY(v *Vector3, angle AngleInRadians) *Vector3 {
	q := new(Quaternion).FromAxisAngle(Vector3_Y, angle)
	return q.MultiplyVector(v)
}
func GetPendicular2DRight(v *Vector3) *Vector3 {
	ret := &Vector3{v.Z, 0, -v.X}
	ret.NormalizeSelf()
	return ret
}
func GetPendicular2DLeft(v *Vector3) *Vector3 {
	ret := &Vector3{-v.Z, 0, v.X}
	ret.NormalizeSelf()
	return ret
}
func GetPositionByFacingAndDistance(startPos, facing *Vector3, distance float32) *Vector3 {
	temp := facing.Normalize()
	temp.ScaleBy(distance)
	temp.IncreaseBy(startPos)
	temp.Y = 0
	return temp
}
func GetPositionByPointAndDistance(startPos, endPos *Vector3, distance float32) *Vector3 {
	temp := GetDirection2D(endPos, startPos)
	return GetPositionByFacingAndDistance(startPos, temp, distance)
}
func GetPositionByPointAndDistanceWithMin(startPos, endPos *Vector3, distance float32) *Vector3 {
	temp := GetDirection2D(endPos, startPos)
	maxDist := GetDistance(startPos, endPos)
	realDis := MinF32(maxDist, distance)
	return GetPositionByFacingAndDistance(startPos, temp, realDis)
}
func GetPositionByPointAndDistanceWithMax(startPos, endPos *Vector3, distance float32) *Vector3 {
	temp := GetDirection2D(endPos, startPos)
	maxDist := GetDistance(startPos, endPos)
	realDis := MaxF32(maxDist, distance)
	return GetPositionByFacingAndDistance(startPos, temp, realDis)
}
func GetReflectedVector(v *Vector3, n *Vector3) *Vector3 {
	nnorm := n.Normalize()
	return v.Substract(nnorm.ScaleBy(nnorm.Dot(v) * 2))
}
func GetMirrorVector(v *Vector3, center *Vector3) *Vector3 {
	ret := v.Clone()
	return ret.DecreaseBy(center).ScaleBy(-1).IncreaseBy(center)
}
func GetRandomVector2D() *Vector3 {
	ret := &Vector3{GetRandomF32(-1, 1), 0, GetRandomF32(-1, 1)}
	ret.NormalizeSelf()
	return ret
}

func StringToF32(v string) (float32, bool) {
	f, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return 0, false
	}
	return float32(f), true
}

func F32ToString(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', 6, 32)
}
