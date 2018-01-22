package ape

import (
	"base/math"
)

func NewAPEngine() *APEngine {

	apengine := &APEngine{}

	return apengine
}

type APEngine struct {
	force *math.Vector2

	masslessForce *math.Vector2

	groups    []*Group
	numGroups int32
	timeStep  float32

	damping float32
	//_container:DisplayObjectContainer

	constraintCycles          int32
	constraintCollisionCycles int32
}

func (this *APEngine) Init(t float32) {
	this.timeStep = t * t

	this.numGroups = 0
	this.groups = make([]*Group, 0)

	this.force = &math.Vector2{}
	this.masslessForce = &math.Vector2{}

	this.damping = 1

	this.constraintCycles = 0
	this.constraintCollisionCycles = 1
}

func (this *APEngine) AddForce(v *math.Vector2) {
	this.force.IncreaseBy(v)
}

func (this *APEngine) AddMasslessForce(v *math.Vector2) {
	this.masslessForce.IncreaseBy(v)
}

func (this *APEngine) AddGroup(g *Group) {
	this.groups = append(this.groups, g)
	g.IsParented = true
	this.numGroups++
	g.Init()
}

/**
 * @private
 */
func (this *APEngine) RemoveGroup(g *Group) {

	// var gpos:int = groups.indexOf(g);
	// if (gpos == -1) return;

	// groups.splice(gpos, 1);
	// g.isParented = false;
	// numGroups--;
	// g.cleanup();

	for index, value := range this.groups {
		if value == g {
			this.groups = append(this.groups[:index], this.groups[index+1:]...)
			g.IsParented = false
			this.numGroups--
			g.Cleanup()
		}
	}
}

/**
 * The main step function of the engine. This method should be called
 * continously to advance the simulation. The faster this method is
 * called, the faster the simulation will run. Usually you would call
 * this in your main program loop.
 */
func (this *APEngine) Step() {
	this.Integrate()

	var j int32
	for j = 0; j < this.constraintCycles; j++ {
		this.SatisfyConstraints()
	}
	var i int32
	for i = 0; i < this.constraintCollisionCycles; i++ {
		this.SatisfyConstraints()
		this.CheckCollisions()
	}
}

/**
 * Calling this method will in turn call each particle and constraint's paint method.
 * Generally you would call this method after stepping the engine in the main program
 * cycle.
 */
func (this *APEngine) Paint() {
	var j int32
	for j = 0; j < this.numGroups; j++ {
		g := this.groups[j]
		g.Paint()
	}
}

func (this *APEngine) Integrate() {
	var j int32
	for j = 0; j < this.numGroups; j++ {
		g := this.groups[j]
		g.Integrate(this.timeStep)
	}
}

func (this *APEngine) SatisfyConstraints() {
	var j int32
	for j = 0; j < this.numGroups; j++ {
		g := this.groups[j]
		g.SatisfyConstraints()
	}
}

func (this *APEngine) CheckCollisions() {
	var j int32
	for j = 0; j < this.numGroups; j++ {
		g := this.groups[j]
		g.CheckCollisions()
	}
}
