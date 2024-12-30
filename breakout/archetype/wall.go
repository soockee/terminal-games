package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
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

func NewWall(w donburi.World, shape resolv.IShape, sprite *ebiten.Image) *donburi.Entry {
	wall := Wall.SpawnInWorld(w)

	component.Space.Get(component.Space.MustFirst(w)).Add(shape)
	component.Collidable.Set(wall, &component.CollidableData{
		Type:  tags.Wall,
		Shape: shape,
	})
	component.Sprite.SetValue(wall, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	return wall
}
