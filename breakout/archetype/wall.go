package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/ldtkgo"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/tags"
	"github.com/yohamta/donburi"
)

var (
	Wall = newArchetype(
		tags.Wall,
		component.Collidable,
		component.Sprite,
	)
)

func NewWall(w donburi.World, project *assets.LDtkProject, entity *ldtkgo.Entity) *donburi.Entry {
	wall := Wall.SpawnInWorld(w)

	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])
	width := float64(entity.Width)
	height := float64(entity.Height)

	r := resolv.NewRectangleFromCorners(X, Y, X+width, Y+height)
	component.Space.Get(component.Space.MustFirst(w)).Add(r)
	component.Collidable.Set(wall, &component.CollidableData{
		Type:  tags.Wall,
		Shape: r,
	})

	// sprite := project.GetSpriteByEntityInstance(entity)
	sprite := ebiten.NewImage(int(r.Bounds().Width()), int(r.Bounds().Height()))
	component.Sprite.SetValue(wall, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	return wall
}