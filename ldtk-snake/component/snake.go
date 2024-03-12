package component

import (
	"github.com/yohamta/donburi"
)

type SnakeData struct {
	Speed float64
}

var Snake = donburi.NewComponentType[SnakeData]()
