package system

import (
	"github.com/soockee/terminal-games/pong/component"
	"github.com/soockee/terminal-games/pong/physics"
	"github.com/soockee/terminal-games/pong/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// UpdateCollision checks ball vs walls and paddles, resolves overlaps, and reflects velocity.
func UpdateCollision(e *ecs.ECS) {
	ballEntry, ok := tags.Ball.First(e.World)
	if !ok {
		return
	}

	ballShape := component.Shape.Get(ballEntry).Shape
	ballVel := component.Velocity.Get(ballEntry)

	// Ball vs walls
	tags.Wall.Each(e.World, func(entry *donburi.Entry) {
		wallShape := component.Shape.Get(entry).Shape
		if inter := ballShape.Intersection(wallShape); !inter.IsEmpty() {
			if len(inter.Intersections) == 0 {
				return
			}
			normal := inter.Intersections[0].Normal
			r := physics.WallBounce(ballVel.X, ballVel.Y, normal.X, normal.Y, inter.MTV.X, inter.MTV.Y)
			ballVel.X = r.VelX
			ballVel.Y = r.VelY
			ballShape.Move(r.PushX, r.PushY)
		}
	})

	// Ball vs paddles
	tags.Paddle.Each(e.World, func(entry *donburi.Entry) {
		paddleShape := component.Shape.Get(entry).Shape
		if inter := ballShape.Intersection(paddleShape); !inter.IsEmpty() {
			paddleBounds := paddleShape.Bounds()
			paddleCenterY := (paddleBounds.Min.Y + paddleBounds.Max.Y) / 2
			paddleHalfH := paddleBounds.Height() / 2

			ballBounds := ballShape.Bounds()
			ballCenterY := (ballBounds.Min.Y + ballBounds.Max.Y) / 2

			r := physics.PaddleBounce(ballVel.X, ballVel.Y, ballCenterY, paddleCenterY, paddleHalfH, inter.MTV.X, inter.MTV.Y)
			ballVel.X = r.VelX
			ballVel.Y = r.VelY
			ballShape.Move(r.PushX, r.PushY)
		}
	})
}
