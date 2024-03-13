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

func CreateFood(ecs *ecs.ECS, sprite *ebiten.Image, entity *ldtkgo.Entity) *donburi.Entry {
	food := archetype.Food.Spawn(ecs)

	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])
	width := float64(entity.Width)
	height := float64(entity.Height)

	obj := resolv.NewObject(X, Y, width, height, entity.Tags...)
	component.Object.Set(food, obj)
	component.Sprite.SetValue(food, component.SpriteData{Image: sprite})
	component.Collectable.SetValue(food, component.CollectableData{Type: component.FoodCollectable})

	obj.SetShape(resolv.NewRectangle(X, Y, width, height))

	return food
}
