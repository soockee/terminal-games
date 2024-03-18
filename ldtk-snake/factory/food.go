package factory

import (
	"log/slog"

	"github.com/solarlune/resolv"
	"github.com/soockee/ldtkgo"
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	dresolv "github.com/soockee/terminal-games/ldtk-snake/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/soockee/terminal-games/ldtk-snake/util"
	"github.com/yohamta/donburi"
)

func CreateFood(w donburi.World, project *assets.LDtkProject, entity *ldtkgo.EntityDefinition) *donburi.Entry {
	food := archetype.Food.SpawnInWorld(w)

	sprite, err := project.GetSpriteByDefinition(entity)
	if err != nil {
		slog.Error("Sprite not found")
		panic(0)
	}
	component.Sprite.SetValue(food, component.SpriteData{Image: sprite})
	component.Collectable.SetValue(food, component.CollectableData{Type: component.FoodCollectable})

	width := float64(entity.Width)
	height := float64(entity.Height)

	xBound := project.Project.Levels[component.Level_0].Width
	yBound := project.Project.Levels[component.Level_0].Height
	collidableTags := []string{tags.Wall.Name(), tags.Snake.Name(), tags.SnakeBody.Name()}

	space := component.Space.MustFirst(w)

	maxAttempts := 1000 // Adjust this value as needed

	var attempt int
	for attempt = 0; attempt < maxAttempts; attempt++ {
		x, y := util.RandomPointInBounds(entity.Width, entity.Height, xBound-entity.Width, yBound-entity.Width)
		obj := resolv.NewObject(float64(x), float64(y), width, height, entity.Tags...)
		component.Object.Set(food, obj)
		obj.SetShape(resolv.NewRectangle(float64(x), float64(y), width, height))
		dresolv.Add(space, food)
		check := obj.Check(0, 0, collidableTags...)
		if check == nil {
			break // No collision, exit the loop
		}
		dresolv.Remove(space, food) // If collision, try a different position
	}

	if attempt == maxAttempts {
		// Handle the case where all attempts failed to find a non-colliding position
		slog.Error("could not fine a space for food")
		panic(0)
	}

	return food
}
