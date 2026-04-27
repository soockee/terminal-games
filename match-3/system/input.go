package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/soockee/terminal-games/pong/component"
	"github.com/yohamta/donburi/ecs"
)

// UpdateInput reads keyboard input and sets paddle velocities.
// F3 toggles the debug overlay.
func UpdateInput(e *ecs.ECS) {
	// Toggle debug on F3 (single press via inpututil)
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		if entry, ok := component.Debug.First(e.World); ok {
			d := component.Debug.Get(entry)
			d.Enabled = !d.Enabled
		}
	}

}
