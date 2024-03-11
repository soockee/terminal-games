package system

import (
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi"
)

func UpdateCamera(w donburi.World) {
	camera := archetype.MustFindCamera(w)
	cam := component.Camera.Get(camera)

	if !cam.Moving {
		cam.MoveTimer.Update()
		if cam.MoveTimer.IsReady() {
			cam.Moving = true
			component.Velocity.Get(camera).Velocity.Y = -0.5
		}
	}
}
