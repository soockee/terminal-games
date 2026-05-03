package system

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

// spriteQuery matches entities with a Sprite and a Position (static entities).
var spriteQuery = donburi.NewQuery(
	filter.Contains(component.Sprite, component.Position),
)

// DrawEntities renders static sprite entities (blocks, gems, etc.)
// using the camera projection. Animated entities are handled by DrawAnimation.
func DrawEntities(e *ecs.ECS, screen *ebiten.Image) {
	proj := NewProjection(e.World)

	spriteQuery.Each(e.World, func(entry *donburi.Entry) {
		// Skip entities that have an Animation — those are drawn by DrawAnimation.
		if entry.HasComponent(component.Animation) {
			return
		}
		img, ok := getSpriteImage(entry)
		if !ok {
			return
		}
		pos := component.Position.Get(entry)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(proj.WorldToScreenX(pos.X), proj.WorldToScreenY(pos.Y))
		screen.DrawImage(img, op)
	})
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

	switch gs.Phase {
	case component.PhasePaused:
		ebitenutil.DebugPrintAt(screen, "PAUSED", cx-20, cy-10)
		ebitenutil.DebugPrintAt(screen, "Press P or tap to resume", cx-75, cy+10)
	case component.PhaseIdle:
		ebitenutil.DebugPrintAt(screen, "Tap or press Space to start", cx-80, cy)
	case component.PhaseWon:
		ebitenutil.DebugPrintAt(screen, "YOU WIN!", cx-30, cy-20)
		ebitenutil.DebugPrintAt(screen, "Tap or press Space to restart", cx-90, cy+10)
	case component.PhaseDead:
		ebitenutil.DebugPrintAt(screen, "GAME OVER", cx-35, cy-20)
		ebitenutil.DebugPrintAt(screen, "Tap or press Space to restart", cx-90, cy+10)
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
