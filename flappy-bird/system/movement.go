package system

import (
	"github.com/soockee/terminal-games/flappy-bird/component"
	"github.com/soockee/terminal-games/flappy-bird/physics"
	"github.com/soockee/terminal-games/flappy-bird/tags"
	"github.com/yohamta/donburi/ecs"
)

// UpdateMovement applies gravity and velocity to the bird each tick.
func UpdateMovement(e *ecs.ECS) {
	// No movement until the player has started (first jump), or when dead/won.
	if goEntry, ok := component.GameOver.First(e.World); ok {
		go_ := component.GameOver.Get(goEntry)
		if !go_.Started || go_.Dead || go_.Won {
			return
		}
	}

	entry, ok := tags.Bird.First(e.World)
	if !ok {
		return
	}

	space := component.Space.Get(component.Space.MustFirst(e.World))
	screenH := float64(space.Height())

	p := component.Player.Get(entry)
	s := component.Shape.Get(entry)

	p.VelY = physics.ApplyGravity(p.VelY, p.Gravity)
	s.Shape.Move(p.VelX, p.VelY)

	bounds := s.Shape.Bounds()

	// Clamp: prevent flying above the screen top.
	if pushY, clamped := physics.ClampTop(bounds.Min.Y); clamped {
		s.Shape.Move(0, pushY)
		p.VelY = 0
	}

	// Clamp: stop at ground top.
	groundY := screenH
	if gEntry, ok := tags.Ground.First(e.World); ok {
		groundY = component.Ground.Get(gEntry).Y
	}
	if pushY, clamped := physics.ClampGround(bounds.Max.Y, groundY); clamped {
		s.Shape.Move(0, pushY)
		p.VelY = 0
	}
}
