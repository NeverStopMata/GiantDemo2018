package ape

import (
	"base/math"
)

var collisionResolver *CollisionResolver

func CollisionResolver_GetMe() *CollisionResolver {
	if collisionResolver == nil {
		collisionResolver = &CollisionResolver{}
	}
	return collisionResolver
}

type CollisionResolver struct {
}

func (this *CollisionResolver) name() {

}

func (this *CollisionResolver) ResolveParticleParticle(pa IAbstractParticle, pb IAbstractParticle, normal *math.Vector2, depth float32) {
	if normal == nil {
		return
	}
	// a collision has occured. set the current positions to sample locations
	pa.GetCurr().CopyFrom(pa.GetSamp())
	pb.GetCurr().CopyFrom(pb.GetSamp())

	// //圆形碰撞同一直线上方向相反的碰撞

	if normal.X == 0 {

		velDirA := pa.GetVelocity().Normalize()
		velDirB := pb.GetVelocity().Normalize()

		if velDirA.Y == -velDirB.Y || pa.GetFixed() || pb.GetFixed() {

			if pa.GetCurr().Y > pb.GetCurr().Y {
				if !pa.GetFixed() {
					pa.GetCurr().Set(pa.GetCurr().X-depth, pa.GetCurr().Y)
				}
				if !pb.GetFixed() {
					pb.GetCurr().Set(pb.GetCurr().X+depth, pb.GetCurr().Y)
				}
			} else {
				if !pa.GetFixed() {
					pa.GetCurr().Set(pa.GetCurr().X+depth, pa.GetCurr().Y)
				}
				if !pb.GetFixed() {
					pb.GetCurr().Set(pb.GetCurr().X-depth, pb.GetCurr().Y)
				}
			}
		}

	} else if normal.Y == 0 {

		velDirA := pa.GetVelocity().Normalize()
		velDirB := pb.GetVelocity().Normalize()

		if velDirA.X == -velDirB.X || pa.GetFixed() || pb.GetFixed() {

			if pa.GetCurr().X > pb.GetCurr().X {

				if !pa.GetFixed() {
					pa.GetCurr().Set(pa.GetCurr().X, pa.GetCurr().Y-depth)
				}
				if !pb.GetFixed() {
					pb.GetCurr().Set(pb.GetCurr().X, pb.GetCurr().Y+depth)
				}
			} else {

				if !pa.GetFixed() {
					pa.GetCurr().Set(pa.GetCurr().X, pa.GetCurr().Y+depth)
				}
				if !pb.GetFixed() {
					pb.GetCurr().Set(pb.GetCurr().X, pb.GetCurr().Y-depth)
				}
			}
		}

	} else if pa.GetFixed() {

		this.ResolveWithFixedBody(pa, pb, depth*2)

	} else if pb.GetFixed() {

		this.ResolveWithFixedBody(pb, pa, depth*2)
	}

	mtd := normal.Mult(depth)
	te := pa.GetElasticity() + pb.GetElasticity()
	sumInvMass := pa.GetInvMass() + pb.GetInvMass()

	// the total friction in a collision is combined but clamped to [0,1]
	tf := this.Clamp(1-(pa.GetFriction()+pb.GetFriction()), 0, 1)

	// get the collision components, vn and vt
	ca := pa.GetComponents(normal)
	cb := pb.GetComponents(normal)

	// calculate the coefficient of restitution based on the mass, as the normal component
	vnA := (cb.vn.Mult((te + 1) * pa.GetInvMass()).Add(
		ca.vn.Mult(pb.GetInvMass() - te*pa.GetInvMass()))).DivEquals(sumInvMass)
	vnB := (ca.vn.Mult((te + 1) * pb.GetInvMass()).Add(
		cb.vn.Mult(pa.GetInvMass() - te*pb.GetInvMass()))).DivEquals(sumInvMass)

	// apply friction to the tangental component
	ca.vt.ScaleBy(tf)
	cb.vt.ScaleBy(tf)

	// scale the mtd by the ratio of the masses. heavier particles move less
	mtdA := mtd.Mult(pa.GetInvMass() / sumInvMass)
	mtdB := mtd.Mult(-pb.GetInvMass() / sumInvMass)

	// add the tangental component to the normal component for the new velocity
	vnA.IncreaseBy(ca.vt)
	vnB.IncreaseBy(cb.vt)

	if !pa.GetFixed() {
		pa.ResolveCollision(mtdA, vnA, normal, depth, -1, pb)
	}
	if !pb.GetFixed() {
		pb.ResolveCollision(mtdB, vnB, normal, depth, 1, pa)
	}
}

func (this *CollisionResolver) Clamp(input float32, min float32, max float32) float32 {
	if input > max {
		return max
	}
	if input < min {
		return min
	}
	return input
}

func (this *CollisionResolver) ResolveWithFixedBody(fixedBody IAbstractParticle, body IAbstractParticle, shiftDis float32) {
	normalizeVel := body.GetVelocity().Normalize()
	if math.AbsF32(normalizeVel.Y) > math.AbsF32(normalizeVel.X) {
		if body.GetCurr().Y > fixedBody.GetCurr().Y {

			if body.GetVelocity().Y < 0 {

				if body.GetCurr().X > fixedBody.GetCurr().X {

					body.GetCurr().Set(body.GetCurr().X+shiftDis, body.GetCurr().Y)

				} else {

					body.GetCurr().Set(body.GetCurr().X-shiftDis, body.GetCurr().Y)
				}
			}

		} else if body.GetVelocity().Y > 0 {

			if body.GetCurr().X > fixedBody.GetCurr().X {

				body.GetCurr().Set(body.GetCurr().X+shiftDis, body.GetCurr().Y)

			} else {

				body.GetCurr().Set(body.GetCurr().X-shiftDis, body.GetCurr().Y)
			}
		}
	} else {
		if body.GetCurr().X > fixedBody.GetCurr().X {

			if body.GetVelocity().X < 0 {

				if body.GetCurr().Y > fixedBody.GetCurr().Y {

					body.GetCurr().Set(body.GetCurr().X, body.GetCurr().Y+shiftDis)

				} else {

					body.GetCurr().Set(body.GetCurr().X, body.GetCurr().Y-shiftDis)
				}
			}

		} else if body.GetVelocity().X > 0 {

			if body.GetCurr().Y > fixedBody.GetCurr().Y {

				body.GetCurr().Set(body.GetCurr().X, body.GetCurr().Y+shiftDis)

			} else {

				body.GetCurr().Set(body.GetCurr().X, body.GetCurr().Y-shiftDis)
			}
		}
	}

}
