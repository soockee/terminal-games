package component

import (
	"github.com/yohamta/donburi"
)

type Scene int

const (
	StartScreen = iota
	SnakeScene
)

type GamestateData struct {
	CurrentScene Scene
}

var Gamestate = donburi.NewComponentType[GamestateData]()
