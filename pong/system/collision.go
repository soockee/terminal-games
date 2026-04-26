package system

import (
	"math"

	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/pong/component"
	"github.com/soockee/terminal-games/pong/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

const maxBounceAngle = math.Pi / 3 // 60 degrees

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
			resolveWallBounce(ballShape, ballVel, inter)
		}
	})

	// Ball vs paddles
	tags.Paddle.Each(e.World, func(entry *donburi.Entry) {
		paddleShape := component.Shape.Get(entry).Shape
		if inter := ballShape.Intersection(paddleShape); !inter.IsEmpty() {
			resolvePaddleBounce(ballShape, ballVel, paddleShape, inter)
		}
	})
}

func resolveWallBounce(ballShape resolv.IShape, vel *component.VelocityData, inter resolv.IntersectionSet) {
	// Use the first intersection's normal to reflect
	if len(inter.Intersections) == 0 {
		return
	}
	normal := inter.Intersections[0].Normal
	v := resolv.NewVector(vel.X, vel.Y)
	reflected := v.Reflect(normal)
	vel.X = reflected.X
	vel.Y = reflected.Y

	// Push ball out of wall
	mtv := inter.MTV
	ballShape.Move(mtv.X, mtv.Y)
}

func resolvePaddleBounce(ballShape resolv.IShape, vel *component.VelocityData, paddleShape resolv.IShape, inter resolv.IntersectionSet) {
	// Push ball out of paddle
	mtv := inter.MTV
	ballShape.Move(mtv.X, mtv.Y)

	// Angle-based deflection: where on the paddle face did the ball hit?
	paddleBounds := paddleShape.Bounds()
	paddleCenterY := (paddleBounds.Min.Y + paddleBounds.Max.Y) / 2
	paddleHalfH := paddleBounds.Height() / 2

	ballBounds := ballShape.Bounds()
	ballCenterY := (ballBounds.Min.Y + ballBounds.Max.Y) / 2

	relativeIntersect := (paddleCenterY - ballCenterY) / paddleHalfH
	angle := relativeIntersect * maxBounceAngle

	speed := math.Sqrt(vel.X*vel.X + vel.Y*vel.Y)
	dir := 1.0
	if vel.X > 0 {
		dir = -1.0
	}
	vel.X = dir * speed * math.Cos(angle)
	vel.Y = -speed * math.Sin(angle)
}
