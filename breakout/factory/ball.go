package factory

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/ldtkgo"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/archetype"
	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"

	"github.com/yohamta/donburi"
)

func CreateBall(w donburi.World, project *assets.LDtkProject, entity *ldtkgo.Entity) *donburi.Entry {
	ball := archetype.Ball.SpawnInWorld(w)

	width := float64(entity.Width)

	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])

	r := resolv.NewCircle(X, Y, width/2)
	component.Space.Get(component.Space.MustFirst(w)).Add(r)
	component.Ball.Set(ball, &component.BallData{
		Speed: 8,
		Shape: r,
	})

	sprite := project.GetSpriteByEntityInstance(entity)
	component.Sprite.SetValue(ball, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	return ball
}
