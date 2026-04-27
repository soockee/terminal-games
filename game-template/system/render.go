package system

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/soockee/terminal-games/pong/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// DrawEntities renders all paddles, balls, and walls.
func DrawEntities(e *ecs.ECS, screen *ebiten.Image) {

}

// DrawScore renders the score at the top center.
func DrawScore(e *ecs.ECS, screen *ebiten.Image) {
	var leftScore, rightScore int

	space := component.Space.Get(component.Space.MustFirst(e.World))
	text := fmt.Sprintf("%d - %d", leftScore, rightScore)
	ebitenutil.DebugPrintAt(screen, text, space.Width()/2-20, 10)
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
