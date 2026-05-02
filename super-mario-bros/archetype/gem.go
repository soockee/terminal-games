package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
)

// NewGem creates a collectible gem entity.
func NewGem(w donburi.World, ent *ldtkgo.Entity) {
	entry := w.Entry(w.Create(
		component.Position,
		component.Sprite,
	))
	tlx, tly := ent.TopLeft()
	component.Position.Set(entry, &component.PositionData{
		X: float64(tlx),
		Y: float64(tly),
	})
	if img := ent.SubImage(); img != nil {
		component.Sprite.Set(entry, &component.SpriteData{
			Image: ebiten.NewImageFromImage(img),
		})
	}
}
