package system

import (
	"github.com/soockee/terminal-games/pong/component"
	"github.com/soockee/terminal-games/pong/physics"
	"github.com/soockee/terminal-games/pong/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// UpdateScore checks if the ball has left the screen horizontally and awards a point.
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

	result := physics.CheckScore(ballBounds.Min.X, ballBounds.Max.X, screenW)
	if !result.Scored {
		return
	}

	if result.RightScored {
		scorePaddle(e.World, component.SideRight)
	}
	if result.LeftScored {
		scorePaddle(e.World, component.SideLeft)
	}

	resetBall(ballEntry, ballVel, ball)
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
	spawn := component.SpawnPos.Get(entry)

	// Move shape so its center aligns with the spawn position.
	shape := component.Shape.Get(entry).Shape
	pos := shape.Position()
	shape.Move(spawn.X-pos.X, spawn.Y-pos.Y)

	// Reset velocity toward the side that just scored.
	vel.X, vel.Y = physics.ResetVelocity(vel.X, ball.Speed)
}
