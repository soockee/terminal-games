package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/soockee/terminal-games/flappy-bird/component"
	"github.com/soockee/terminal-games/flappy-bird/event"
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

	// Jump on Space, mouse click, or touch.
	goEntry, goOK := component.GameOver.First(e.World)
	jumpPressed := inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		len(inpututil.AppendJustPressedTouchIDs(nil)) > 0

	if goOK && jumpPressed {
		go_ := component.GameOver.Get(goEntry)
		if go_.Dead || go_.Won {
			// Signal restart — handled by the game loop.
			go_.Restart = true
			return
		}
		if !go_.Started {
			go_.Started = true
		}
	}

	if entry, ok := component.Player.First(e.World); ok && jumpPressed {
		p := component.Player.Get(entry)
		p.VelY = -p.Jump
		event.AudioEvent.Publish(e.World, event.AudioEventData{Name: component.SFXJump})
	}
}
