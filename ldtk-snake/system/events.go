package system

import (
	"log/slog"

	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/event"
	"github.com/soockee/terminal-games/ldtk-snake/factory"
	"github.com/soockee/terminal-games/ldtk-snake/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/events"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

func OnCollideEvent(w donburi.World, e *event.Collide) {
	switch e.Type {
	case event.CollideBody:
		slog.Info("Process Collide Body Event")

	case event.CollideWall:
		slog.Info("Process Collide Wall Event")
	}
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

		slog.Debug("Food Object OnPickupEvent", slog.Any("obj", foodObj))
		w.Remove(food.Entity())

		sceneStateEntity, ok := component.SceneState.First(w)
		if !ok {
			slog.Error("sceneStateEntity not found OnCollideEvent")
			panic(0)
		}
		sceneObj := component.SceneState.Get(sceneStateEntity)
		snakeEntity, ok := component.Snake.First(w)

		if !ok {
			slog.Error("snakeEntity not found OnCollideEvent")
			panic(0)
		}
		scene, _ := component.SceneState.First(w)
		sceneData := component.SceneState.Get(scene)

		factory.CreateBodyPart(w, sceneObj.Project, snakeEntity, sceneData.Project.Project.EntityDefinitionByIdentifier(tags.SnakeBody.Name()), tags.SnakeBody.Name())
		factory.CreateFood(w, sceneObj.Project, sceneData.Project.Project.EntityDefinitionByIdentifier(tags.Food.Name()))

	default:
	}
}

func ProcessEvents(ecs *ecs.ECS) {
	events.ProcessAllEvents(ecs.World)
}
