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
	Ball = newArchetype(
		tags.Ball,

		component.Ball,
		component.Sprite,
		component.Velocity,
	)
)

func NewBall(w donburi.World, project *assets.LDtkProject, entity *ldtkgo.Entity) *donburi.Entry {
	ball := Ball.SpawnInWorld(w)

	width := float64(entity.Width)

	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])

	r := resolv.NewCircle(X, Y, width/2)
	component.Space.Get(component.Space.MustFirst(w)).Add(r)
	component.Ball.Set(ball, &component.BallData{
		Speed:    8,
		Shape:    r,
		MaxSpeed: 30,
	})

	sprite := project.GetSpriteByEntityInstance(entity)
	component.Sprite.SetValue(ball, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	return ball
}
