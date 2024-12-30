package component

import (
	"github.com/yohamta/donburi"
)

type BrickData struct {
	Health int
}

var Brick = donburi.NewComponentType[BrickData]()
