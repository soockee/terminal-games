package factory

import (
	"log/slog"

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
		event.SceneStateEvent.Publish(w, &event.SceneStateData{
			CurrentScene: component.SnakeScene,
		})
	},
	"GithubButton": func(w donburi.World) { util.OpenUrl("https://github.com/soockee") },
}

func CreateButton(ecs *ecs.ECS, project *assets.LDtkProject, entity *ldtkgo.Entity) *donburi.Entry {
	button := archetype.Button.Spawn(ecs)

	width := float64(entity.Width)
	height := float64(entity.Height)
	pivotX := float64(entity.Pivot[0]) * width  // Calculate pivot offset for X
	pivotY := float64(entity.Pivot[1]) * height // Calculate pivot offset for Y
	// Calculate adjusted position based on pivot
	X := float64(entity.Position[0]) - pivotX
	Y := float64(entity.Position[1]) - pivotY

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
