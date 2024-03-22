package factory

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/soockee/ldtkgo"
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateWall(ecs *ecs.ECS, project *assets.LDtkProject, entity *ldtkgo.Entity) *donburi.Entry {
	wall := archetype.Wall.Spawn(ecs)

	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])
	width := float64(entity.Width)
	height := float64(entity.Height)

	obj := resolv.NewObject(X, Y, width, height, entity.Tags...)
	component.Object.Set(wall, obj)
	sprite, err := project.GetSpriteByEntityInstance(entity)
	if err != nil {
		slog.Error("Sprite not found")
		panic(0)
	}
	component.Sprite.SetValue(wall, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	obj.SetShape(resolv.NewRectangle(X, Y, width, height))

	return wall
}
