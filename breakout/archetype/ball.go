package archetype

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/engine"
	"github.com/soockee/terminal-games/breakout/tags"

	"github.com/yohamta/donburi"
)

var (
	Ball = newArchetype(
		tags.Ball,

		component.Ball,
		component.Sprite,
		component.Velocity,
	)
)

func NewBall(w donburi.World, shape *resolv.Circle, sprite *ebiten.Image) *donburi.Entry {
	ball := Ball.SpawnInWorld(w)

	component.Space.Get(component.Space.MustFirst(w)).Add(shape)
	component.Ball.Set(ball, &component.BallData{
		Speed:                   8,
		Shape:                   shape,
		MaxSpeed:                10,
		CollisionCooldownBlock:  time.Duration(time.Millisecond * 17),  // 1 frame
		CollisionCooldownPlayer: time.Duration(time.Millisecond * 134), // 8 frames
		CooldownTimer:           engine.Timer{},
	})

	component.Sprite.SetValue(ball, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	return ball
}
