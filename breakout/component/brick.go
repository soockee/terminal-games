package component

import (
	"github.com/yohamta/donburi"
)

type BrickData struct {
}

var Brick = donburi.NewComponentType[BrickData]()
