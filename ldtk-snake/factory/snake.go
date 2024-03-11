package factory

import (
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/config"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateSnake(ecs *ecs.ECS, iid string) *donburi.Entry {
	snake := archetype.Snake.Spawn(ecs)

	entity := config.C.GetEntityByIID(iid, config.C.CurrentLevel)
	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])
	width := float64(entity.Width)
	height := float64(entity.Height)

	obj := resolv.NewObject(X, Y, width, height, entity.Tags...)
	component.Object.Set(snake, obj)
	component.Snake.SetValue(snake, component.SnakeData{
		Speed: 1,
	})

	component.Sprite.SetValue(snake, component.SpriteData{Image: config.C.GetSprite(entity)})

	obj.SetShape(resolv.NewRectangle(X, Y, width, height))

	return snake
}
