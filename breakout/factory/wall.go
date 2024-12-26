package factory

import (
	"github.com/solarlune/ldtkgo"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/archetype"
	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/yohamta/donburi"
)

func CreateWall(w donburi.World, project *assets.LDtkProject, entity *ldtkgo.Entity) *donburi.Entry {
	wall := archetype.Wall.SpawnInWorld(w)

	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])
	width := float64(entity.Width)
	height := float64(entity.Height)

	r := resolv.NewRectangleFromCorners(X, Y, X+width, Y+height)
	component.Space.Get(component.Space.MustFirst(w)).Add(r)
	component.Collidable.Set(wall, &component.CollidableData{
		Shape: r,
	})

	// sprite := project.GetSpriteByEntityInstance(entity)
	//component.Sprite.SetValue(wall, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	return wall
}
