package ape

import ()

type IAbstractCollection interface {
	AddConstraint(c IAbstractConstraint)
	RemoveConstraint(c IAbstractConstraint)
	AddParticle(p IAbstractParticle)
	RemoveParticle(p IAbstractParticle)

	Init()
	Paint()
	Cleanup()
	GetParticles() []IAbstractParticle
}

type Sprite struct {
	X        float32
	Y        float32
	rotation float32
}

type AbstractCollection struct {
	sprite      *Sprite
	Particles   []IAbstractParticle
	Constraints []IAbstractConstraint
	IsParented  bool
}

func (this *AbstractCollection) GetParticles() []IAbstractParticle {
	return this.Particles
}

func (this *AbstractCollection) AddParticle(p IAbstractParticle) {
	this.Particles = append(this.Particles, p)
	if this.IsParented {
		p.Init()
	}
}

/**
 * Removes an AbstractParticle from the AbstractCollection.
 *
 * @param p The particle to be removed.
 */
func (this *AbstractCollection) RemoveParticle(p IAbstractParticle) {
	// var ppos:int = particles.indexOf(p);
	// if (ppos == -1) return;
	// particles.splice(ppos, 1);
	// p.cleanup();

	for index, value := range this.Particles {
		if value == p {
			this.Particles = append(this.Particles[:index], this.Particles[index+1:]...)
			break
		}
	}
	p.CleanUp()
}

/**
 * Adds a constraint to the Collection.
 *
 * @param c The constraint to be added.
 */
func (this *AbstractCollection) AddConstraint(c IAbstractConstraint) {
	this.Constraints = append(this.Constraints, c)

	if this.GetIsParented() {
		c.Init()
	}
}

/**
 * Removes a constraint from the Collection.
 *
 * @param c The constraint to be removed.
 */
func (this *AbstractCollection) RemoveConstraint(c IAbstractConstraint) {

	for index, value := range this.Constraints {
		if value == c {
			this.Constraints = append(this.Constraints[:index], this.Constraints[index+1:]...)
			break
		}
	}
	c.CleanUp()
}

/**
 * Initializes every member of this AbstractCollection by in turn calling
 * each members <code>init()</code> method.
 */
func (this *AbstractCollection) Init() {

	for i := 0; i < len(this.Particles); i++ {
		this.Particles[i].Init()
	}
	for i := 0; i < len(this.Constraints); i++ {
		this.Constraints[i].Init()
	}
}

/**
 * paints every member of this AbstractCollection by calling each members
 * <code>paint()</code> method.
 */
func (this *AbstractCollection) Paint() {

	var p IAbstractParticle

	for i := 0; i < len(this.Particles); i++ {
		p = this.Particles[i]
		if (!p.GetFixed()) || p.GetAlwaysRepaint() {
			p.Paint()
		}
	}
	//|||以后添加
	// var c:SpringConstraint;
	// len = _constraints.length;
	// for (i = 0; i < len; i++) {
	// 	c = _constraints[i];
	// 	if ((! c.fixed) || c.alwaysRepaint) c.paint();
	// }
}

/**
 * Calls the <code>cleanup()</code> method of every member of this AbstractCollection.
 * The cleanup() method is called automatically when an AbstractCollection is removed
 * from its parent.
 */
func (this *AbstractCollection) Cleanup() {

	for _, v := range this.Particles {
		v.CleanUp()
	}

	for _, v := range this.Constraints {
		v.CleanUp()
	}
}

/**
 * Provides a Sprite to use as a container for drawing or adding children. When the
 * sprite is requested for the first time it is automatically added to the global
 * container in the APEngine class.
 */
func (this *AbstractCollection) GetSprite() *Sprite {

	if this.sprite != nil {
		return this.sprite
	}

	//|||if APEngine.container == nil {
	// 	panic("The container property of the APEngine class has not been set")
	// }

	this.sprite = &Sprite{}
	//|||APEngine.container.addChild(_sprite)
	return this.sprite
}

/**
 * Returns an array of every particle and constraint added to the AbstractCollection.
 */
// func (this *AbstractCollection) GetAll():Array {
// 	return particles.concat(constraints);
// }

/**
 * @private
 */
func (this *AbstractCollection) GetIsParented() bool {
	return this.IsParented
}

/**
 * @private
 */
func (this *AbstractCollection) SetIsParented(b bool) {
	this.IsParented = b
}

/**
 * @private
 */
func (this *AbstractCollection) Integrate(dt2 float32) {
	// var len:int = _particles.length;
	// for (var i:int = 0; i < len; i++) {
	// 	var p:AbstractParticle = _particles[i];
	// 	p.update(dt2);
	// }
	for _, v := range this.Particles {
		v.Update(dt2)
	}
}

/**
 * @private
 */
func (this *AbstractCollection) SatisfyConstraints() {
	// var len:int = _constraints.length;
	// for (var i:int = 0; i < len; i++) {
	// 	var c:AbstractConstraint = _constraints[i];
	// 	c.resolve();
	// }
	for _, v := range this.Constraints {
		v.Resolve()
	}
}

/**
 * @private
 */
func (this *AbstractCollection) CheckInternalCollisions() {

	// var plen:int = _particles.length;
	// for (var j:int = 0; j < plen; j++) {

	// 	var pa:AbstractParticle = _particles[j];
	// 	if (! pa.collidable) continue;

	// 	// ...vs every other particle in this AbstractCollection
	// 	for (var i:int = j + 1; i < plen; i++) {
	// 		var pb:AbstractParticle = _particles[i];
	// 		if (pb.collidable) CollisionDetector.test(pa, pb);
	// 	}

	// 	// ...vs every other constraint in this AbstractCollection
	// 	var clen:int = _constraints.length;
	// 	for (var n:int = 0; n < clen; n++) {
	// 		var c:SpringConstraint = _constraints[n];
	// 		if (c.collidable && ! c.isConnectedTo(pa)) {
	// 			c.scp.updatePosition();
	// 			CollisionDetector.test(pa, c.scp);
	// 		}
	// 	}
	// }

	plen := len(this.Particles)
	for j := 0; j < plen; j++ {
		pa := this.Particles[j]
		if !pa.GetCollidable() {
			continue
		}
		for i := j + 1; i < plen; i++ {
			pb := this.Particles[i]
			if pb.GetCollidable() {
				CollisionDetector_GetMe().Test(pa, pb)
			}
		}
		//|||后续添加
		// var clen:int = _constraints.length;
		// for (var n:int = 0; n < clen; n++) {
		// 	var c:SpringConstraint = _constraints[n];
		// 	if (c.collidable && ! c.isConnectedTo(pa)) {
		// 		c.scp.updatePosition();
		// 		CollisionDetector.test(pa, c.scp);
		// 	}
		// }
	}
}

/**
 * @private
 */
func (this *AbstractCollection) CheckCollisionsVsCollection(ac IAbstractCollection) {

	// every particle in this collection...
	// var plen:int = _particles.length;
	// for (var j:int = 0; j < plen; j++) {

	// 	var pga:AbstractParticle = _particles[j];
	// 	if (! pga.collidable) continue;

	// 	// ...vs every particle in the other collection
	// 	var acplen:int = ac.particles.length;
	// 	for (var x:int = 0; x < acplen; x++) {
	// 		var pgb:AbstractParticle = ac.particles[x];
	// 		if (pgb.collidable) CollisionDetector.test(pga, pgb);
	// 	}
	// 	// ...vs every constraint in the other collection
	// 	var acclen:int = ac.constraints.length;
	// 	for (x = 0; x < acclen; x++) {
	// 		var cgb:SpringConstraint = ac.constraints[x];
	// 		if (cgb.collidable && ! cgb.isConnectedTo(pga)) {
	// 			cgb.scp.updatePosition();
	// 			CollisionDetector.test(pga, cgb.scp);
	// 		}
	// 	}
	// }

	// // every constraint in this collection...
	// var clen:int = _constraints.length;
	// for (j = 0; j < clen; j++) {
	// 	var cga:SpringConstraint = _constraints[j];
	// 	if (! cga.collidable) continue;

	// 	// ...vs every particle in the other collection
	// 	acplen = ac.particles.length;
	// 	for (var n:int = 0; n < acplen; n++) {
	// 		pgb = ac.particles[n];
	// 		if (pgb.collidable && ! cga.isConnectedTo(pgb)) {
	// 			cga.scp.updatePosition();
	// 			CollisionDetector.test(pgb, cga.scp);
	// 		}
	// 	}
	// }

	//
	//
	plen := len(this.Particles)
	for j := 0; j < plen; j++ {
		pga := this.GetParticles()[j]
		if !pga.GetCollidable() {
			continue
		}

		// ...vs every particle in the other collection
		acplen := len(ac.GetParticles())
		for x := 0; x < acplen; x++ {
			pgb := ac.GetParticles()[x]
			if pgb.GetCollidable() {
				CollisionDetector_GetMe().Test(pga, pgb)
			}
		}
		// ...vs every constraint in the other collection
		// |||以后再添加
		// var acclen:int = ac.constraints.length;
		// for (x = 0; x < acclen; x++) {
		// 	var cgb:SpringConstraint = ac.constraints[x];
		// 	if (cgb.collidable && ! cgb.isConnectedTo(pga)) {
		// 		cgb.scp.updatePosition();
		// 		CollisionDetector.test(pga, cgb.scp);
		// 	}
		// }
	}
}
