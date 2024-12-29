package system

import (
	"log/slog"
	"math/rand"

	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/archetype"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/tags"
	"github.com/soockee/terminal-games/breakout/util"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/events"
)

func OnCollideEvent(w donburi.World, e *event.Collide) {
	ColliderType := e.Collider.Archetype().Layout()

	if ColliderType.HasComponent(tags.Ball) {
		OnBallCollisionEvent(w, e)
	}
}

func OnPickupEvent(w donburi.World, e *event.Collect) {
	switch e.Type {
	default:
		slog.Error("pickup not implemented", slog.Any("Type", e.Type))
	}
}

func OnReleaseEvent(w donburi.World, e *event.Release) {
	entry := component.Ball.MustFirst(w)
	ball := component.Ball.Get(entry)
	velocity := component.Velocity.Get(entry)

	// randomize direction
	direction := resolv.NewVector(2*rand.Float64()-1, -1)

	velocity.Velocity = velocity.Velocity.Add(direction)

	velocity.Velocity = util.LimitMagnitude(velocity.Velocity, ball.MaxSpeed)
}

func OnExplodeEvent(w donburi.World, e *event.Explode) {
	archetype.NewExplosion(w, e.Brick.Entity())
	w.Remove(e.Brick.Entity())
}

func ProcessEvents(ecs *ecs.ECS) {
	events.ProcessAllEvents(ecs.World)
}
