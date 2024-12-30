package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/tags"

	"github.com/yohamta/donburi"
)

var (
	Brick = newArchetype(
		tags.Brick,

		component.Brick,
		component.Sprite,
		component.Collidable,
	)
)

func NewBrick(w donburi.World, shape resolv.IShape, sprite *ebiten.Image) *donburi.Entry {
	brick := Brick.SpawnInWorld(w)

	component.Space.Get(component.Space.MustFirst(w)).Add(shape)
	component.Brick.Set(brick, &component.BrickData{
		Health: 1,
	})
	component.Collidable.Set(brick, &component.CollidableData{
		Type:  tags.Brick,
		Shape: shape,
	})
	component.Sprite.SetValue(brick, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	return brick
}
