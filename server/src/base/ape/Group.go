package ape

func NewGroup(_collideInternal bool) *Group {

	g := &Group{
		composites:      make([]*Composite, 0),
		collisionList:   make([]*Group, 0),
		collideInternal: _collideInternal,
	}
	return g
}

type Group struct {
	AbstractCollection

	composites      []*Composite
	collisionList   []*Group
	collideInternal bool
}

func (this *Group) Init() {
	this.AbstractCollection.Init()

	for _, v := range this.composites {
		v.Init()
	}
}

func (this *Group) AddComposite(c *Composite) {
	this.composites = append(this.composites, c)
	c.IsParented = true
	if this.IsParented {
		c.Init()
	}
}

func (this *Group) RemoveComposite(c *Composite) {

	for index, v := range this.composites {
		if v == c {
			this.composites = append(this.composites[:index], this.composites[index+1:]...)
			c.IsParented = false
			c.Cleanup()
		}
	}
}

func (this *Group) Paint() {

	this.AbstractCollection.Paint()

	for _, v := range this.composites {
		v.Paint()
	}
}

func (this *Group) AddCollidable(g *Group) {
	this.collisionList = append(this.collisionList, g)
}

func (this *Group) RemoveCollidable(g *Group) {

	for index, v := range this.collisionList {
		if v == g {
			this.collisionList = append(this.collisionList[:index], this.collisionList[index+1:]...)
		}

	}
}
func (this *Group) Cleanup() {
	this.AbstractCollection.Cleanup()
	for _, v := range this.composites {
		v.Cleanup()
	}
}

func (this *Group) Integrate(dt2 float32) {

	this.AbstractCollection.Integrate(dt2)

	for _, v := range this.composites {
		v.Integrate(dt2)
	}
}

func (this *Group) SatisfyConstraints() {

	this.AbstractCollection.SatisfyConstraints()
	for _, v := range this.composites {
		v.SatisfyConstraints()
	}
}

func (this *Group) CheckCollisions() {

	if this.collideInternal {
		this.CheckCollisionGroupInternal()
	}

	for _, v := range this.collisionList {
		this.CheckCollisionVsGroup(v)
	}
}
func (this *Group) CheckCollisionGroupInternal() {

	// check collisions not in composites
	this.CheckInternalCollisions()

	// for every composite in this Group..
	clen := len(this.composites)
	for j := 0; j < clen; j++ {

		ca := this.composites[j]

		// .. vs non composite particles and constraints in this group
		ca.CheckCollisionsVsCollection(this)

		// ...vs every other composite in this Group
		for i := j + 1; i < clen; i++ {
			cb := this.composites[i]
			ca.CheckCollisionsVsCollection(cb)
		}
	}

}

func (this *Group) CheckCollisionVsGroup(g *Group) {

	// check particles and constraints not in composites of either group
	this.CheckCollisionsVsCollection(g)

	clen := len(this.composites)
	gclen := len(g.composites)

	// for every composite in this group..
	for i := 0; i < clen; i++ {

		// check vs the particles and constraints of g
		c := this.composites[i]
		c.CheckCollisionsVsCollection(g)

		// check vs composites of g
		for j := 0; j < gclen; j++ {
			gc := g.composites[j]
			c.CheckCollisionsVsCollection(gc)
		}
	}

	// check particles and constraints of this group vs the composites of g
	for j := 0; j < gclen; j++ {
		gc := g.composites[j]
		this.CheckCollisionsVsCollection(gc)
	}
}
