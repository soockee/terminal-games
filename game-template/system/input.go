package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/soockee/terminal-games/game-template/component"
	"github.com/yohamta/donburi/ecs"
)

// UpdateInput reads keyboard/mouse/touch input.
// F3 toggles the debug overlay. Action input handles game flow
// (start, restart) and can be extended with game-specific logic.
func UpdateInput(e *ecs.ECS) {
	// Toggle debug on F3 (single press via inpututil)
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		if entry, ok := component.Debug.First(e.World); ok {
			d := component.Debug.Get(entry)
			d.Enabled = !d.Enabled
		}
	}

	// Action input: Space, mouse click, or touch.
	goEntry, goOK := component.GameState.First(e.World)
	actionPressed := inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		len(inpututil.AppendJustPressedTouchIDs(nil)) > 0

	if goOK && actionPressed {
		gs := component.GameState.Get(goEntry)
		if gs.Paused {
			gs.Paused = false
			return
		}
		if gs.Dead || gs.Won {
			gs.Restart = true
			return
		}
		if !gs.Started {
			gs.Started = true
		}
	}

	// TODO: add game-specific input handling here.
	// Example: publish audio event on action.
	// if actionPressed {
	//     event.AudioEvent.Publish(e.World, event.AudioEventData{Name: "action"})
	// }
}
