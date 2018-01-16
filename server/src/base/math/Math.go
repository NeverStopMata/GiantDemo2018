package math

import "math"

const (
	EPSILON = 0.00001
)

func IsZero(v float32) bool {
	return AbsF32(v) < EPSILON
}
func AbsInt(v int) int {
	switch {
	case v < 0:
		return -v
	case v == 0:
		return 0 // return correctly abs(-0)
	}
	return v
}
func MaxInt(v1, v2 int) int {
	if v1 > v2 {
		return v1
	}
	return v2
}
func MinInt(v1, v2 int) int {
	if v1 < v2 {
		return v1
	}
	return v2
}
func AbsF64(v float64) float64 {
	switch {
	case v < 0:
		return -v
	case v == 0:
		return 0 // return correctly abs(-0)
	}
	return v
}
func AbsF32(v float32) float32 {
	switch {
	case v < 0:
		return -v
	case v == 0:
		return 0 // return correctly abs(-0)
	}
	return v
}
func MaxF32(v1, v2 float32) float32 {
	if v1 > v2 {
		return v1
	}
	return v2
}
func MinF32(v1, v2 float32) float32 {
	if v1 < v2 {
		return v1
	}
	return v2
}
func Atan2F32(a float32, b float32) float32 {
	return float32(math.Atan2(float64(a), float64(b)))
}

func SqrtF32(v float32) float32 {
	return float32(math.Sqrt(float64(v)))
}
func SqrInt(v int) int {
	return v * v
}

func SqrF32(v float32) float32 {
	return v * v
}
func RoundF32(value float32) float32 {
	if IsZero(value) {
		return 0
	}
	if value > 0 {
		return float32(int(value + 0.5))
	}
	return float32(int(value - 0.5))
}
func RoundF32ToInt(value float32) int {
	if IsZero(value) {
		return 0
	}
	if value > 0 {
		return int(value + 0.5)
	}
	return int(value - 0.5)
}
func GetFractionF32(v float32) float32 {
	return v - float32(int(v))
}
func FloorF32(value float32) int {
	return int(value)
}
func CeilingF32(value float32) int {
	v := int(value)
	f := value - float32(int(value))
	if AbsF32(f) < EPSILON {
		return v
	}
	return v + 1
}
func ClampF32(v, min, max float32) float32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
func ClampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
func BlendF32(v1, v2, blendFactor float32) float32 {
	f := ClampF32(blendFactor, 0, 1)
	return v1*(1-f) + v2*f
}
func SignF32(v float32) float32 {
	if v > 0 {
		return 1
	} else if v < 0 {
		return -1
	}
	return 0
}
func SinF32(v float32) float32 {
	return float32(math.Sin(float64(v)))
}
func ASinF32(v float32) float32 {
	return float32(math.Asin(float64(ClampF32(v, -1, 1))))
}
func CosF32(v float32) float32 {
	return float32(math.Cos(float64(v)))
}
func ACosF32(v float32) float32 {
	return float32(math.Acos(float64(ClampF32(v, -1, 1))))
}
func TanF32(v float32) float32 {
	return float32(math.Tan(float64(v)))
}
func SelF32(selV float32, va float32, vb float32) float32 {
	if selV >= 0 {
		return va
	}
	return vb
}
func EqualF32(v1, v2 float32) bool {
	return IsZero(v1 - v2)
}
func GreaterF32(v1, v2 float32) bool {
	return v1 >= (v2 + EPSILON)
}
func GreaterEqualF32(v1, v2 float32) bool {
	return v1 > (v2 - EPSILON)
}
func LessF32(v1, v2 float32) bool {
	return v1 <= (v2 - EPSILON)
}
func LessEqualF32(v1, v2 float32) bool {
	return v1 < (v2 + EPSILON)
}
func IsNaNF32(v float32) bool {
	return math.IsNaN(float64(v))
}

// 三元运算函数
// for example var a,b int = 2,3  max := If(a > b,a,b).(int)
func If(condition bool, trueValue, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}

func Pow10F32(n int) float32 {
	return float32(math.Pow10(n))
}

func PowF32(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}
