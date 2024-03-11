package component

import (
	"github.com/soockee/terminal-games/ldtk-snake/engine"
	"github.com/yohamta/donburi"
)

type CameraData struct {
	Moving    bool
	MoveTimer *engine.Timer
}

var Camera = donburi.NewComponentType[CameraData]()
