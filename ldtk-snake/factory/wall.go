package factory

import (
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/archetypes"
	"github.com/soockee/terminal-games/ldtk-snake/components"
	"github.com/soockee/terminal-games/ldtk-snake/config"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateWall(ecs *ecs.ECS, idd string) *donburi.Entry {
	wall := archetypes.Wall.Spawn(ecs)

	entity := config.C.GetEntityByIID(idd, config.C.CurrentLevel)
	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])
	width := float64(entity.Width)
	height := float64(entity.Height)

	obj := resolv.NewObject(X, Y, width, height, entity.Tags...)
	components.Object.Set(wall, obj)
	components.Sprite.SetValue(wall, components.SpriteData{Image: config.C.GetSprite(entity)})

	obj.SetShape(resolv.NewRectangle(X, Y, width, height))

	return wall
}
