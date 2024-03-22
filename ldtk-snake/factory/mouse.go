package factory

import (
	"log/slog"
	"time"

	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/engine"
	"github.com/soockee/terminal-games/ldtk-snake/tags"

	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/yohamta/donburi"
	"github.com/yohamta/ganim8/v2"
)

func CreateMouse(w donburi.World, project *assets.LDtkProject) *donburi.Entry {
	mouse := archetype.Mouse.SpawnInWorld(w)

	component.Mouse.SetValue(mouse, component.MouseData{
		IsHidden:   false,
		Invincible: engine.NewTimer(time.Second * 2),
	})

	entityDefinition := project.Project.EntityDefinitionByIdentifier("MouseMoving")
	mouseHiddenDefinition := project.Project.EntityDefinitionByIdentifier("MouseHidden")

	animationMoving, err := project.GetAnimatedSpriteByDefinition(entityDefinition)
	if err != nil {
		slog.Error("mouse moving animation not found")
		panic(0)
	}
	animationHidden, err := project.GetAnimatedSpriteByDefinition(mouseHiddenDefinition)
	if err != nil {
		slog.Error("mouse hidden animation not found")
		panic(0)
	}
	animation := map[int]*ganim8.Animation{
		int(component.MouseMoving): animationMoving,
		int(component.MouseHidden): animationHidden,
	}
	component.Animation.SetValue(mouse, component.AnimationData{Animations: animation})
	component.Collectable.SetValue(mouse, component.CollectableData{Type: component.MouseCollectable})

	width := float64(entityDefinition.Width)
	height := float64(entityDefinition.Height)

	entityDefinition.Tags = append(entityDefinition.Tags, tags.Collectable.Name())

	obj := resolv.NewObject(0, 0, width, height, entityDefinition.Tags...)

	center := obj.Center()
	offsetX := obj.Size.X / 2
	offsetY := obj.Size.Y / 2
	center = center.Sub(resolv.NewVector(offsetX, offsetY))
	radius := offsetX
	obj.SetShape(resolv.NewCircle(center.X-width/2, center.Y-height/2, radius))
	component.Object.Set(mouse, obj)

	return mouse
}
