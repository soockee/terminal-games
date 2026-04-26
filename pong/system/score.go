package system

import (
	"github.com/soockee/terminal-games/pong/component"
	"github.com/soockee/terminal-games/pong/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// UpdateScore checks if the ball has left the screen horizontally and awards a point.
// Returns true if a reset occurred.
func UpdateScore(e *ecs.ECS) {
	ballEntry, ok := tags.Ball.First(e.World)
	if !ok {
		return
	}

	space := component.Space.Get(component.Space.MustFirst(e.World))
	screenW := float64(space.Width())

	ballBounds := component.Shape.Get(ballEntry).Shape.Bounds()
	ballVel := component.Velocity.Get(ballEntry)
	ball := component.Ball.Get(ballEntry)

	scored := false

	// Ball off left edge → right scores
	if ballBounds.Max.X < 0 {
		scorePaddle(e.World, component.SideRight)
		scored = true
	}
	// Ball off right edge → left scores
	if ballBounds.Min.X > screenW {
		scorePaddle(e.World, component.SideLeft)
		scored = true
	}

	if scored {
		resetBall(ballEntry, ballVel, ball)
	}
}

func scorePaddle(w donburi.World, side component.PaddleSide) {
	tags.Paddle.Each(w, func(entry *donburi.Entry) {
		p := component.Paddle.Get(entry)
		if p.Side == side {
			p.Score++
		}
	})
}

func resetBall(entry *donburi.Entry, vel *component.VelocityData, ball *component.BallData) {
	ref := component.EntityRef.Get(entry)
	cx, cy := ref.Entity.Center()

	// Move shape so its center aligns with the LDtk entity center.
	shape := component.Shape.Get(entry).Shape
	pos := shape.Position()
	shape.Move(float64(cx)-pos.X, float64(cy)-pos.Y)

	// Reset velocity toward the side that just scored
	dir := 1.0
	if vel.X > 0 {
		dir = -1.0
	}
	vel.X = dir * ball.Speed
	vel.Y = ball.Speed
}
