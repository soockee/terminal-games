package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

// playerInputQuery matches the player entity.
var playerInputQuery = donburi.NewQuery(
	filter.Contains(component.Player, component.Body),
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

	// Player movement input — only while game is active.
	if !goOK {
		return
	}
	gs := component.GameState.Get(goEntry)
	if !gs.Started || gs.Dead || gs.Won || gs.Paused {
		return
	}

	playerInputQuery.Each(e.World, func(entry *donburi.Entry) {
		pd := component.Player.Get(entry)

		// Horizontal: Arrow keys or A/D
		pd.MoveDir = 0
		if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			pd.MoveDir = -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			pd.MoveDir = 1
		}

		// Jump: Space or Up or W (single press)
		pd.JumpInput = inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			inpututil.IsKeyJustPressed(ebiten.KeyUp) ||
			inpututil.IsKeyJustPressed(ebiten.KeyW)
	})
}
