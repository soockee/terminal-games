package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
)

// NewCoin creates a collectible coin entity.
func NewCoin(w donburi.World, x, y float64, img *ebiten.Image) {
	entry := w.Entry(w.Create(
		component.Position,
		component.Sprite,
		component.Collectable,
	))
	component.Position.Set(entry, &component.PositionData{X: x, Y: y})
	if img != nil {
		component.Sprite.Set(entry, &component.SpriteData{Image: img})
	}
}
