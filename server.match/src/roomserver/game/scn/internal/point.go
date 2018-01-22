// Copyright 2012 Daniel Connelly.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"
	"math"
)

// Point represents a point in n-dimensional Euclidean space.
type Point struct {
	X float64
	Y float64
}

func (r Point) Equal(other Point) bool {
	return r.X == other.X && r.Y == other.Y
}

func (r Point) String() string {
	return fmt.Sprintf("[%.2f, %.2f]", r.X, r.Y)
}

// Dist computes the Euclidean distance between two points p and q.
func (p Point) dist(q Point) float64 {
	dx := p.X - q.X
	sum := dx * dx
	dx = p.Y - q.Y
	sum = sum + dx*dx
	return math.Sqrt(sum)
}
