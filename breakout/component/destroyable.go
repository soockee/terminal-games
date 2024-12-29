package component

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type DestroyState int

const (
	Default DestroyState = iota
	Destroyed
)

type DestroyableData struct {
	Type         *donburi.ComponentType[donburi.Tag]
	Shape        resolv.IShape
	DestroyState DestroyState
}

var Destroyable = donburi.NewComponentType[DestroyableData]()
