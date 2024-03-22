package component

import (
	"github.com/yohamta/donburi"
)

type CollectableType int

const (
	FoodCollectable CollectableType = iota
	MouseCollectable
)

type CollectableData struct {
	Type CollectableType
}

var Collectable = donburi.NewComponentType[CollectableData]()
