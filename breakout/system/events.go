package system

import (
	"log/slog"

	"github.com/soockee/terminal-games/breakout/event"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/events"
)

func OnCollideEvent(w donburi.World, e *event.Collide) {
	switch e.Type {
	case event.CollideBody:
		fallthrough

	case event.CollideWall:
		event.GameStateEvent.Publish(w, &event.GameStateData{
			IsGameOver: true,
		})
	}
}

func OnPickupEvent(w donburi.World, e *event.Collect) {

	switch e.Type {
	default:
		slog.Error("unknown collectable")
		panic(0)
	}
}

func ProcessEvents(ecs *ecs.ECS) {
	events.ProcessAllEvents(ecs.World)
}
