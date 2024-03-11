package component

import (
	"github.com/yohamta/donburi"
)

type GameStateData struct {
}

var Gamestate = donburi.NewComponentType[GameStateData]()
