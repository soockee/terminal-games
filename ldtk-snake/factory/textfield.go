package factory

import (
	"log/slog"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/soockee/ldtkgo"
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateTextField(ecs *ecs.ECS, project *assets.LDtkProject, entity *ldtkgo.Entity) *donburi.Entry {
	textfield := archetype.TextField.Spawn(ecs)

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
	component.Text.Set(textfield, &component.TextData{
		Identifier:          entity.Identifier,
		Text:                text,
		IsAnimated:          false,
		CurrentCharPosition: 0,
	})

	obj := resolv.NewObject(X, Y, width, height, entity.Tags...)
	component.Object.Set(textfield, obj)

	sprite, err := project.GetSpriteByEntityInstance(entity)
	if err != nil {
		slog.Error("Sprite not found")
		panic(0)
	}
	component.Sprite.SetValue(textfield, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	obj.SetShape(resolv.NewRectangle(X, Y, width, height))

	return textfield
}
