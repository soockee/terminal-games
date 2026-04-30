package system

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/soockee/terminal-games/game-template/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// DrawEntities renders all game entities using the camera projection.
// Extend this with your entity-specific draw logic.
func DrawEntities(e *ecs.ECS, screen *ebiten.Image) {
	_ = NewProjection(e.World)
	// TODO: iterate over tagged entities and draw them using proj.WorldToScreen*
}

// DrawScore renders the current score at the top center of the screen.
func DrawScore(e *ecs.ECS, screen *ebiten.Image) {
	entry, ok := component.Score.First(e.World)
	if !ok {
		return
	}
	score := component.Score.Get(entry)
	bounds := screen.Bounds()
	text := fmt.Sprintf("%d", score.Value)
	ebitenutil.DebugPrintAt(screen, text, bounds.Dx()/2-10, 10)
}

// DrawHUD renders overlays for game state: start prompt, pause, game over, or win.
func DrawHUD(e *ecs.ECS, screen *ebiten.Image) {
	entry, ok := component.GameState.First(e.World)
	if !ok {
		return
	}
	gs := component.GameState.Get(entry)
	bounds := screen.Bounds()
	cx := bounds.Dx() / 2
	cy := bounds.Dy() / 2

	if gs.Paused {
		ebitenutil.DebugPrintAt(screen, "PAUSED", cx-20, cy-10)
		ebitenutil.DebugPrintAt(screen, "Press P or tap to resume", cx-75, cy+10)
		return
	}
	if !gs.Started {
		ebitenutil.DebugPrintAt(screen, "Tap or press Space to start", cx-80, cy)
		return
	}
	if gs.Won {
		ebitenutil.DebugPrintAt(screen, "YOU WIN!", cx-30, cy-20)
		ebitenutil.DebugPrintAt(screen, "Tap or press Space to restart", cx-90, cy+10)
		return
	}
	if gs.Dead {
		ebitenutil.DebugPrintAt(screen, "GAME OVER", cx-35, cy-20)
		ebitenutil.DebugPrintAt(screen, "Tap or press Space to restart", cx-90, cy+10)
		return
	}
}

func getSpriteImage(entry *donburi.Entry) (*ebiten.Image, bool) {
	if !entry.HasComponent(component.Sprite) {
		return nil, false
	}
	s := component.Sprite.Get(entry)
	if s.Image == nil {
		return nil, false
	}
	return s.Image, true
}
