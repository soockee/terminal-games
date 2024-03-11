package factory

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/soockee/ldtkgo"
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateSnake(ecs *ecs.ECS, sprite *ebiten.Image, entity *ldtkgo.Entity) *donburi.Entry {
	snake := archetype.Snake.Spawn(ecs)

	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])
	width := float64(entity.Width)
	height := float64(entity.Height)

	obj := resolv.NewObject(X, Y, width, height, entity.Tags...)
	component.Object.Set(snake, obj)
	component.Snake.SetValue(snake, component.SnakeData{
		Speed: 1,
	})

	component.Sprite.SetValue(snake, component.SpriteData{Image: sprite})

	obj.SetShape(resolv.NewRectangle(X, Y, width, height))

	return snake
}
