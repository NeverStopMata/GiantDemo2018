package ape

import ()

type IAbstractConstraint interface {
	IAbstractItem
	Resolve()
}
type AbstractConstraint struct {
	AbstractItem

	Stiffness float32
}

func (this *AbstractConstraint) Resolve() {

}
