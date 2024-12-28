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
	Brick = newArchetype(
		tags.Brick,

		component.Brick,
		component.Sprite,
		component.Collidable,
	)
)

func NewBrick(w donburi.World, project *assets.LDtkProject, entity *ldtkgo.Entity) *donburi.Entry {
	brick := Brick.SpawnInWorld(w)

	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])
	width := float64(entity.Width)
	height := float64(entity.Height)

	r := resolv.NewRectangleFromCorners(X, Y, X+width, Y+height)
	component.Space.Get(component.Space.MustFirst(w)).Add(r)
	component.Brick.Set(brick, &component.BrickData{})
	component.Collidable.Set(brick, &component.CollidableData{
		Type:  tags.Brick,
		Shape: r,
	})

	sprite := project.GetSpriteByEntityInstance(entity)
	component.Sprite.SetValue(brick, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	return brick
}
