package ape

import (
	"base/math"
)

func NewComposite() *Composite {
	c := &Composite{
		delta: &math.Vector2{},
	}
	return c
}

type Composite struct {
	AbstractCollection
	delta *math.Vector2
}

func (this *Composite) RotateByRadian(angleRadians float32, center *math.Vector2) {
	var p IAbstractParticle
	pa := this.Particles
	// var len:int = pa.length;
	// for (var i:int = 0; i < len; i++) {
	// 	p = pa[i];
	// 	var radius:Number = p.center.distance(center);
	// 	var angle:Number = getRelativeAngle(center, p.center) + angleRadians;
	// 	p.px = (Math.cos(angle) * radius) + center.x;
	// 	p.py = (Math.sin(angle) * radius) + center.y;
	// }

	for _, v := range pa {
		p = v
		radius := p.GetCenter().Distance(center)
		angle := this.GetRelativeAngle(center, p.GetCenter()) + angleRadians

		p.SetPx(math.CosF32(angle)*radius + center.X)
		p.SetPy(math.SinF32(angle)*radius + center.Y)
	}
}

func (this *Composite) RotateByAngle(angleDegrees float32, center *math.Vector2) {
	angleRadians := angleDegrees * PI_OVER_ONE_EIGHTY
	this.RotateByRadian(angleRadians, center)
}

func (this *Composite) GetFixed() bool {

	for _, v := range this.Particles {
		if v.GetFixed() == false {
			return false
		}
	}
	return true
}

/**
 * @private
 */
func (this *Composite) Setfixed(b bool) {

	for _, v := range this.Particles {
		v.SetFixed(b)
	}
}

func (this *Composite) GetRelativeAngle(center *math.Vector2, p *math.Vector2) float32 {
	this.delta.Set(p.X-center.X, p.Y-center.Y)
	return math.Atan2F32(this.delta.Y, this.delta.X)
}
