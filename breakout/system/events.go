package system

import (
	"log/slog"

	"github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/events"
)

func OnCollideEvent(w donburi.World, e *event.Collide) {
	switch e.Type {
	case tags.Wall:
		slog.Debug("collide with Wall")

	case tags.Player:
		slog.Debug("collide with Player")
	}
}

func OnPickupEvent(w donburi.World, e *event.Collect) {
	switch e.Type {
	default:
		slog.Error("pickup not implemented", slog.Any("Type", e.Type))
	}
}

func ProcessEvents(ecs *ecs.ECS) {
	events.ProcessAllEvents(ecs.World)
}
