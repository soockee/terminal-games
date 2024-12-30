package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/tags"
	"github.com/yohamta/donburi"
)

var (
	TextField = newArchetype(
		tags.TextField,
		component.Text,
		component.Sprite,
	)
)

func NewTextField(w donburi.World, shape resolv.IShape, sprite *ebiten.Image) *donburi.Entry {
	textfield := TextField.SpawnInWorld(w)

	component.Space.Get(component.Space.MustFirst(w)).Add(shape)

	component.Text.Set(textfield, &component.TextData{
		Shape: shape,
	})

	component.Sprite.SetValue(textfield, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	return textfield
}
