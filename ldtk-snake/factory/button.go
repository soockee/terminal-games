package factory

import (
	"log/slog"
	"math/rand"

	"github.com/solarlune/resolv"
	"github.com/soockee/ldtkgo"
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/event"
	"github.com/soockee/terminal-games/ldtk-snake/util"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var buttonhandlerMapping = map[string]func(w donburi.World){
	"StartButton": func(w donburi.World) {
		randomLevel := rand.Intn(len(component.SnakeLevels))
		slog.Info("starting level", slog.Any("Level", component.SnakeLevels[randomLevel]))
		event.SceneStateEvent.Publish(w, &event.SceneStateData{
			CurrentScene: component.SnakeLevels[randomLevel],
		})
	},
	"GithubButton": func(w donburi.World) { util.OpenUrl("https://github.com/soockee") },
	"ResetButton": func(w donburi.World) {
		event.SceneStateEvent.Publish(w, &event.SceneStateData{
			CurrentScene: component.SnakeScene,
		})
	},
}

func CreateButton(ecs *ecs.ECS, project *assets.LDtkProject, entity *ldtkgo.Entity) *donburi.Entry {
	button := archetype.Button.Spawn(ecs)

	width := float64(entity.Width)
	height := float64(entity.Height)
	// Calculate adjusted position based on pivot
	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])

	obj := resolv.NewObject(X, Y, width, height, entity.Tags...)
	component.Object.Set(button, obj)

	component.Button.SetValue(button, component.ButtonData{
		Clicked:     false,
		HandlerFunc: buttonhandlerMapping[entity.Identifier],
	})

	sprite, err := project.GetSpriteByEntityInstance(entity)
	if err != nil {
		slog.Error("Sprite not found")
		panic(0)
	}
	component.Sprite.SetValue(button, component.SpriteData{Image: sprite})

	obj.SetShape(resolv.NewRectangle(X, Y, width, height))

	return button
}
