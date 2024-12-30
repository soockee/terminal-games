package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/tags"

	"github.com/yohamta/donburi"
)

var (
	Player = newArchetype(
		tags.Player,

		component.Player,
		component.Collidable,
		component.Sprite,
		component.Velocity,
	)
)

func NewPlayer(w donburi.World, shape resolv.IShape, sprite *ebiten.Image) *donburi.Entry {
	player := Player.SpawnInWorld(w)

	component.Space.Get(component.Space.MustFirst(w)).Add(shape)
	component.Player.Set(player, &component.PlayerData{
		Speed:             8,
		SpeedAcceleration: 1.05,
		SpeedFriction:     0.94,
	})

	component.Collidable.Set(player, &component.CollidableData{
		Type:  tags.Player,
		Shape: shape,
	})

	component.Sprite.SetValue(player, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	return player
}
