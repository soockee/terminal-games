package component

import (
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type SnakeData struct {
	Speed     float64
	OnWall    *resolv.Object
	Direction input.Action
}

var Snake = donburi.NewComponentType[SnakeData]()
