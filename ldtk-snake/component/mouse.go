package component

import (
	"github.com/soockee/terminal-games/ldtk-snake/engine"
	"github.com/yohamta/donburi"
)

type MouseState int

const (
	MouseMoving MouseState = iota
	MouseHidden
)

type MouseData struct {
	IsHidden   bool
	Invincible *engine.Timer
}

var Mouse = donburi.NewComponentType[MouseData]()
