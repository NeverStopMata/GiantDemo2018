package ape

import (
	"base/glog"
	"base/math"
	_ "reflect"
)

var collisionDetector *CollisionDetector

func CollisionDetector_GetMe() *CollisionDetector {
	if collisionDetector == nil {
		collisionDetector = &CollisionDetector{}
	}
	return collisionDetector
}

type CollisionDetector struct {
}

func (this *CollisionDetector) Test(objA IAbstractParticle, objB IAbstractParticle) {

	if objA.GetFixed() && objB.GetFixed() {

		return
	}

	//优化长距离的碰撞忽略,只处理圆对圆
	a := objA.GetParticleType()
	b := objB.GetParticleType()
	if a == ParticleTypeCircle && b == ParticleTypeCircle {
		if math.AbsF32(objA.GetPx()-objB.GetPx()) > ColliderTestDistanceMax ||
			math.AbsF32(objA.GetPy()-objB.GetPy()) > ColliderTestDistanceMax {
			return
		}
	}

	if objA.GetMultisample() == 0 && objB.GetMultisample() == 0 {
		this.NormVsNorm(objA, objB)

	} else if objA.GetMultisample() > 0 && objB.GetMultisample() == 0 {
		this.SampVsNorm(objA, objB)

	} else if objB.GetMultisample() > 0 && objA.GetMultisample() == 0 {
		this.SampVsNorm(objB, objA)

	} else if objA.GetMultisample() == objB.GetMultisample() {
		this.SampVsSamp(objA, objB)

	} else {
		this.NormVsNorm(objA, objB)
	}
}

func (this *CollisionDetector) NormVsNorm(objA IAbstractParticle, objB IAbstractParticle) {
	objA.GetSamp().CopyFrom(objA.GetCurr())
	objB.GetSamp().CopyFrom(objB.GetCurr())
	this.TestTypes(objA, objB)
}

func (this *CollisionDetector) SampVsNorm(objA IAbstractParticle, objB IAbstractParticle) {

	s := float32(1 / (objA.GetMultisample() + 1))
	t := s

	objB.GetSamp().CopyFrom(objB.GetCurr())

	var i int32
	for i = 0; i <= objA.GetMultisample(); i++ {
		objA.GetSamp().Set(objA.GetPrev().X+t*(objA.GetCurr().X-objA.GetPrev().X),
			objA.GetPrev().Y+t*(objA.GetCurr().Y-objA.GetPrev().Y))

		if this.TestTypes(objA, objB) {
			return
		}
		t += s
	}
}
func (this *CollisionDetector) SampVsSamp(objA IAbstractParticle, objB IAbstractParticle) {

	s := float32(1 / (objA.GetMultisample() + 1))
	t := s
	var i int32
	for i = 0; i <= objA.GetMultisample(); i++ {

		objA.GetSamp().Set(objA.GetPrev().X+t*(objA.GetCurr().X-objA.GetPrev().X),
			objA.GetPrev().Y+t*(objA.GetCurr().Y-objA.GetPrev().Y))

		objB.GetSamp().Set(objB.GetPrev().X+t*(objB.GetCurr().X-objB.GetPrev().X),
			objB.GetPrev().Y+t*(objB.GetCurr().Y-objB.GetPrev().Y))

		if this.TestTypes(objA, objB) {
			return
		}
		t += s
	}
}
func (this *CollisionDetector) TestTypes(objA IAbstractParticle, objB IAbstractParticle) bool {

	a := objA.GetParticleType()
	b := objB.GetParticleType()

	if a == ParticleTypeRect && b == ParticleTypeRect {
		return this.TestOBBvsOBB(objA.(*RectangleParticle), objB.(*RectangleParticle))
	} else if a == ParticleTypeCircle && b == ParticleTypeCircle {
		return this.TestCirclevsCircle(objA.(*CircleParticle), objB.(*CircleParticle))
	} else if a == ParticleTypeRect && b == ParticleTypeCircle {
		return this.TestOBBvsCircle(objA.(*RectangleParticle), objB.(*CircleParticle))
	} else if a == ParticleTypeCircle && b == ParticleTypeRect {
		return this.TestOBBvsCircle(objB.(*RectangleParticle), objA.(*CircleParticle))
	}
	/*
		a := reflect.TypeOf(objA).String()
		b := reflect.TypeOf(objB).String()

		if a == "*ape.RectangleParticle" && b == "*ape.RectangleParticle" {
			return this.TestOBBvsOBB(objA.(*RectangleParticle), objB.(*RectangleParticle))
		} else if a == "*ape.CircleParticle" && b == "*ape.CircleParticle" {
			return this.TestCirclevsCircle(objA.(*CircleParticle), objB.(*CircleParticle))
		} else if a == "*ape.RectangleParticle" && b == "*ape.CircleParticle" {
			return this.TestOBBvsCircle(objA.(*RectangleParticle), objB.(*CircleParticle))

		} else if a == "*ape.CircleParticle" && b == "*ape.RectangleParticle" {
			return this.TestOBBvsCircle(objB.(*RectangleParticle), objA.(*CircleParticle))
		}
	*/

	glog.Error("[ape.TestTypes] error:", a, " ", b)
	return false
}

func (this *CollisionDetector) TestOBBvsOBB(ra *RectangleParticle, rb *RectangleParticle) bool {

	var collisionNormal *math.Vector2
	var collisionDepth float32 = 99999999

	for i := 0; i < 2; i++ {

		axisA := ra.axes[i]
		depthA := this.TestIntervals(ra.GetProjection(axisA), rb.GetProjection(axisA))
		if depthA == 0 {
			return false
		}

		axisB := rb.axes[i]
		depthB := this.TestIntervals(ra.GetProjection(axisB), rb.GetProjection(axisB))
		if depthB == 0 {
			return false
		}

		absA := math.AbsF32(depthA)
		absB := math.AbsF32(depthB)

		if absA < math.AbsF32(collisionDepth) || absB < math.AbsF32(collisionDepth) {
			var altb bool = absA < absB
			if altb {
				collisionNormal = axisA
				collisionDepth = depthA
			} else {
				collisionNormal = axisB
				collisionDepth = depthB
			}

		}
	}

	CollisionResolver_GetMe().ResolveParticleParticle(ra, rb, collisionNormal, collisionDepth)
	return true
}

func (this *CollisionDetector) TestOBBvsCircle(ra *RectangleParticle, ca *CircleParticle) bool {

	var collisionNormal *math.Vector2
	var collisionDepth float32 = 99999999
	depths := make([]float32, 2)

	// first go through the axes of the rectangle
	for i := 0; i < 2; i++ {

		boxAxis := ra.axes[i]
		depth := this.TestIntervals(ra.GetProjection(boxAxis), ca.GetProjection(boxAxis))

		if depth == 0 {
			return false
		}

		if math.AbsF32(depth) < math.AbsF32(collisionDepth) {
			collisionNormal = boxAxis
			collisionDepth = depth
		}
		depths[i] = depth
	}

	// determine if the circle's center is in a vertex region
	r := ca.radius
	if math.AbsF32(depths[0]) < r && math.AbsF32(depths[1]) < r {

		vertex := this.ClosestVertexOnOBB(ca.samp, ra)

		// get the distance from the closest vertex on rect to circle center
		collisionNormal = vertex.Substract(ca.samp)
		mag := collisionNormal.Magnitude()
		collisionDepth = r - mag

		if collisionDepth > 0 {
			// there is a collision in one of the vertex regions
			collisionNormal.DivEquals(mag)

		} else {
			// ra is in vertex region, but is not colliding
			return false
		}
	}
	CollisionResolver_GetMe().ResolveParticleParticle(ra, ca, collisionNormal, collisionDepth)
	return true
}
func (this *CollisionDetector) TestCirclevsCircle(ca *CircleParticle, cb *CircleParticle) bool {

	depthX := this.TestIntervals(ca.GetIntervalX(), cb.GetIntervalX())

	if depthX == 0 {
		return false
	}

	depthY := this.TestIntervals(ca.GetIntervalY(), cb.GetIntervalY())

	if depthY == 0 {
		return false
	}

	collisionNormal := ca.samp.Substract(cb.samp)
	mag := collisionNormal.Magnitude()

	collisionDepth := (ca.radius + cb.radius) - mag

	if collisionDepth > 0 {
		collisionNormal.DivEquals(mag)
		CollisionResolver_GetMe().ResolveParticleParticle(ca, cb, collisionNormal, collisionDepth)
		//glog.Info("[TestCirclevsCircle]碰撞了,", ca, " ", cb)
		return true
	}
	return false
}

/**
 * Returns 0 if intervals do not overlap. Returns smallest depth if they do.
 */
func (this *CollisionDetector) TestIntervals(intervalA *Interval, intervalB *Interval) float32 {
	// println(intervalA.max)
	// println(intervalA.min)
	if intervalA.max < intervalB.min {
		return 0
	}
	if intervalB.max < intervalA.min {
		return 0
	}

	lenA := intervalB.max - intervalA.min
	lenB := intervalB.min - intervalA.max

	if math.AbsF32(lenA) < math.AbsF32(lenB) {
		return lenA
	} else {
		return lenB
	}

}

/**
 * Returns the location of the closest vertex on r to point p
 */
func (this *CollisionDetector) ClosestVertexOnOBB(p *math.Vector2, r *RectangleParticle) *math.Vector2 {

	d := p.Substract(r.samp)
	q := &math.Vector2{r.samp.X, r.samp.Y}

	for i := 0; i < 2; i++ {
		dist := d.Dot(r.axes[i])

		if dist >= 0 {
			dist = r.extents[i]
		} else if dist < 0 {
			dist = -r.extents[i]
		}

		q.IncreaseBy(r.axes[i].Mult(dist))
	}
	return q
}
