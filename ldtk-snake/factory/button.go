package factory

import (
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/config"
	"github.com/soockee/terminal-games/ldtk-snake/system"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var buttonhandlerMapping = map[string]func(){
	"StartButton": system.Start,
	"GithubButton": system.OpenGithub,
}

func CreateButton(ecs *ecs.ECS, idd string) *donburi.Entry {
	button := archetype.Button.Spawn(ecs)

	entity := config.C.GetEntityByIID(idd, config.C.CurrentLevel)

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

	component.Sprite.SetValue(button, component.SpriteData{Image: config.C.GetSprite(entity)})

	obj.SetShape(resolv.NewRectangle(X, Y, width, height))

	return button
}
