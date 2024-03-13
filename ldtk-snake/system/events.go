package system

import (
	"log/slog"

	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/event"
	"github.com/soockee/terminal-games/ldtk-snake/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

func OnCollideEvent(w donburi.World, e *event.Collide) {
	slog.Info("Process Wall Event")
}

func OnPickupEvent(w donburi.World, e *event.Collect) {
	slog.Info("Process Food Event")

	switch e.Type {
	case component.FoodCollectable:
		// got only one collectable right now
		food, ok := query.NewQuery(filter.Contains(component.Collectable)).First(w)
		if !ok {
			slog.Error("food not found")
		}
		foodObj := component.Collectable.Get(food)
		space, ok := component.Space.First(w)
		if !ok {
			slog.Error("space not found")
		}
		resolv.Remove(space, food)

		slog.Error("", slog.Any("obj", foodObj))
		w.Remove(food.Entity())
		// food = factory.CreateFood()

	default:
	}
}

func ProcessEvents(ecs *ecs.ECS) {
	events.ProcessAllEvents(ecs.World)
}
