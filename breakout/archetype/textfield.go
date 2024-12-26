package archetype

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/ldtkgo"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/assets"
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

func NewTextField(w donburi.World, project *assets.LDtkProject, entity *ldtkgo.Entity) *donburi.Entry {
	textfield := TextField.SpawnInWorld(w)

	width := float64(entity.Width)
	height := float64(entity.Height)
	// Calculate adjusted position based on pivot
	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])

	textPropertiy := entity.PropertyByIdentifier("text")

	text := []string{}
	if len(entity.Properties) > 0 && !textPropertiy.IsNull() {
		text = strings.Split(textPropertiy.AsString(), "\n")
	}

	r := resolv.NewRectangleFromCorners(X, Y, X+width, Y+height)
	component.Space.Get(component.Space.MustFirst(w)).Add(r)

	component.Text.Set(textfield, &component.TextData{
		Identifier:          entity.Identifier,
		Text:                text,
		IsAnimated:          false,
		CurrentCharPosition: 0,
		Shape:               r,
	})

	sprite := project.GetSpriteByEntityInstance(entity)
	component.Sprite.SetValue(textfield, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	return textfield
}
