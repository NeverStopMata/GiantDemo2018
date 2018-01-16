package ape

import (
	"base/math"
)

type IAbstractItem interface {
	CleanUp()
	Init()
	Paint()
	GetSprite() *Sprite
	GetAlwaysRepaint() bool
}

type AbstractItem struct {
	sprite *Sprite
	//visible       bool
	alwaysRepaint bool

	displayObject *DisplayObject

	displayObjectOffset *math.Vector2

	displayObjectRotation float32
}

func (this *AbstractItem) Init() {

}
func (this *AbstractItem) Paint() {

}
func (this *AbstractItem) CleanUp() {
	// sprite.graphics.clear();
	// 			for (var i:int = 0; i < sprite.numChildren; i++) {
	// 				sprite.removeChildAt(i);
	// 			}
}

func (this *AbstractItem) GetSprite() *Sprite {

	if this.sprite != nil {
		return this.sprite
	}

	// |||if APEngine.container == nil {
	// 	panic("The container property of the APEngine class has not been set")
	// }

	this.sprite = &Sprite{}
	//|||APEngine.container.addChild(this.sprite)
	return this.sprite
}
func (this *AbstractItem) GetAlwaysRepaint() bool {
	return this.alwaysRepaint
}
