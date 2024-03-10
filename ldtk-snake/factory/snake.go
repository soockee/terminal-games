package factory

import (
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/archetypes"
	"github.com/soockee/terminal-games/ldtk-snake/components"
	"github.com/soockee/terminal-games/ldtk-snake/config"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateSnake(ecs *ecs.ECS, iid string) *donburi.Entry {
	snake := archetypes.Snake.Spawn(ecs)

	entity := config.C.GetEntityByIID(iid, config.C.CurrentLevel)
	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])
	width := float64(entity.Width)
	height := float64(entity.Height)

	obj := resolv.NewObject(X, Y, width, height, entity.Tags...)
	components.Object.Set(snake, obj)
	components.Snake.SetValue(snake, components.SnakeData{
		Speed: 1,
	})

	components.Sprite.SetValue(snake, components.SpriteData{Image: config.C.GetSprite(entity)})
	components.Control.SetValue(snake, components.ControlData{
		InputHandler: components.InputSytem.NewHandler(components.SnakeHandler, input.Keymap{
			components.ActionMoveUp:    {input.KeyGamepadUp, input.KeyUp, input.KeyW},
			components.ActionMoveDown:  {input.KeyGamepadDown, input.KeyDown, input.KeyS},
			components.ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
			components.ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
			components.ActionClick:     {input.KeyTouchTap, input.KeyMouseLeft},
		}),
	})

	obj.SetShape(resolv.NewRectangle(X, Y, width, height))

	return snake
}
