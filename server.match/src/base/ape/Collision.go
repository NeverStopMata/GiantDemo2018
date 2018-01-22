package ape

import (
	"base/math"
)

type Collision struct {
	vn *math.Vector2
	vt *math.Vector2
}

func NewCollision(n *math.Vector2, t *math.Vector2) *Collision {
	c := &Collision{vn: n, vt: t}
	return c
}
