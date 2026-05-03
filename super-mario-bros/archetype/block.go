package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
)

// AddStaticCollision creates a static physics shape from world-space top-left coordinates.
func AddStaticCollision(w donburi.World, x, y, width, height float64) {
	spaceEntry, ok := component.PhysicsSpace.First(w)
	if !ok {
		return
	}
	component.PhysicsSpace.Get(spaceEntry).Space.AddStaticRect(x, y, width, height)
}

// NewBlock creates a static block entity (visual only; collision is tag-driven).
func NewBlock(w donburi.World, x, y float64, img *ebiten.Image) {
	entry := w.Entry(w.Create(component.Position, component.Sprite))
	component.Position.Set(entry, &component.PositionData{X: x, Y: y})
	if img != nil {
		component.Sprite.Set(entry, &component.SpriteData{Image: img})
	}
}

// NewJumpBlock creates a question-mark block entity (visual only; collision is tag-driven).
func NewJumpBlock(w donburi.World, x, y float64, img *ebiten.Image) {
	entry := w.Entry(w.Create(component.Position, component.Sprite))
	component.Position.Set(entry, &component.PositionData{X: x, Y: y})
	if img != nil {
		component.Sprite.Set(entry, &component.SpriteData{Image: img})
	}
	// TODO: read Items field and attach item component
}

// NewGround creates a static ground tile entity (visual only; collision is tag-driven).
func NewGround(w donburi.World, x, y float64, img *ebiten.Image) {
	entry := w.Entry(w.Create(component.Position, component.Sprite))
	component.Position.Set(entry, &component.PositionData{X: x, Y: y})
	if img != nil {
		component.Sprite.Set(entry, &component.SpriteData{Image: img})
	}
}
