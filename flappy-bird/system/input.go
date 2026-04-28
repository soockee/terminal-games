package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/soockee/terminal-games/flappy-bird/component"
	"github.com/soockee/terminal-games/flappy-bird/event"
	"github.com/yohamta/donburi/ecs"
)

// pauseButtonTapped returns true if any just-pressed touch was inside the pause button.
// We use JustPressedTouchIDs (not JustReleasedTouchIDs) because TouchPosition
// returns (0,0) for released touches — the position is only available while active.
func pauseButtonTapped(screenW int) bool {
	bx, by, bw, bh := pauseButtonBounds(screenW)
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if x >= bx && x <= bx+bw && y >= by && y <= by+bh {
			return true
		}
	}
	return false
}

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

	goEntry, goOK := component.GameState.First(e.World)

	// Derive virtual screen width for the pause button hit-test.
	virtualW := 0
	if camEntry, ok := component.Camera.First(e.World); ok {
		cam := component.Camera.Get(camEntry)
		space := component.Space.Get(component.Space.MustFirst(e.World))
		virtualW = int(float64(space.Width()) * cam.ScaleX)
	}

	// Toggle pause on Escape/P (keyboard) or pause button tap (mobile).
	pauseTriggered := inpututil.IsKeyJustPressed(ebiten.KeyEscape) ||
		inpututil.IsKeyJustPressed(ebiten.KeyP) ||
		(virtualW > 0 && pauseButtonTapped(virtualW))
	if goOK && pauseTriggered {
		go_ := component.GameState.Get(goEntry)
		if go_.Started && !go_.Dead && !go_.Won {
			go_.Paused = !go_.Paused
			return
		}
	}

	// Jump on Space, mouse click, or touch.
	jumpPressed := inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		len(inpututil.AppendJustPressedTouchIDs(nil)) > 0

	if goOK && jumpPressed {
		go_ := component.GameState.Get(goEntry)
		if go_.Paused {
			go_.Paused = false
			return
		}
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
		if goOK {
			go_ := component.GameState.Get(goEntry)
			if go_.Paused {
				return
			}
		}
		p := component.Player.Get(entry)
		p.VelY = -p.Jump
		event.AudioEvent.Publish(e.World, event.AudioEventData{Name: component.SFXJump})
	}
}
