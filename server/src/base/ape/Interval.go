package ape

type Interval struct {
	min float32
	max float32
}

func NewInterval(min float32, max float32) *Interval {
	i := &Interval{min, max}
	return i
}
