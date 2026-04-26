package system

import (
	"github.com/soockee/terminal-games/pong/component"
	"github.com/soockee/terminal-games/pong/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// UpdateMovement applies velocity to all shapes (paddles + ball).
func UpdateMovement(e *ecs.ECS) {
	space := component.Space.Get(component.Space.MustFirst(e.World))
	screenH := float64(space.Height())

	// Move paddles (clamped to screen)
	tags.Paddle.Each(e.World, func(entry *donburi.Entry) {
		vel := component.Velocity.Get(entry)
		s := component.Shape.Get(entry)
		shape := s.Shape

		shape.Move(0, vel.Y)

		// Clamp to screen bounds
		bounds := shape.Bounds()
		if bounds.Min.Y < 0 {
			shape.Move(0, -bounds.Min.Y)
		}
		if bounds.Max.Y > screenH {
			shape.Move(0, screenH-bounds.Max.Y)
		}
	})

	// Move ball
	if ballEntry, ok := tags.Ball.First(e.World); ok {
		vel := component.Velocity.Get(ballEntry)
		s := component.Shape.Get(ballEntry)
		s.Shape.Move(vel.X, vel.Y)
	}
}
