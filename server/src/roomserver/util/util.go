package util

const (
	EPSILON = 0.0000001
)

func IsZero(v float64) bool {
	return Abs32(v) < EPSILON
}

func AbsInt(X int) int {
	switch {
	case X < 0:
		return -X
	case X == 0:
		return 0 // return correctly abs(-0)
	}
	return X
}
func MaxInt(X, Y int) int {
	if X > Y {
		return X
	}
	return Y
}
func MinInt(X, Y int) int {
	if X < Y {
		return X
	}
	return Y
}
func Abs32(X float64) float64 {
	switch {
	case X < 0:
		return -X
	case X == 0:
		return 0 // return correctly abs(-0)
	}
	return X
}
func Max32(X, Y float64) float64 {
	if X > Y {
		return X
	}
	return Y
}
func Min32(X, Y float64) float64 {
	if X < Y {
		return X
	}
	return Y
}

//限制value在min~max之间
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}
