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
		fallthrough

	case event.CollideWall:
		event.GameStateEvent.Publish(w, &event.GameStateData{
			IsGameOver: true,
		})
	}
}

func OnPickupEvent(w donburi.World, e *event.Collect) {

	switch e.Type {
	case component.FoodCollectable:
		slog.Info("Process Food Event")
		// got only one collectable right now
		food, ok := query.NewQuery(filter.Contains(component.Collectable)).First(w)
		if !ok {
			slog.Error("food not found")
		}
		foodObj := component.Collectable.Get(food)

		resolv.Remove(component.Space.MustFirst(w), food)

		slog.Debug("Food Object OnPickupEvent", slog.Any("obj", foodObj))
		w.Remove(food.Entity())

		sceneObj := component.SceneState.Get(component.SceneState.MustFirst(w))
		snakeEntity := component.Snake.MustFirst(w)
		snakeData := component.Snake.Get(snakeEntity)
		snakeData.Speed *= snakeData.SpeedAcceleration

		if !ok {
			slog.Error("snakeEntity not found OnCollideEvent")
			panic(0)
		}

		sceneData := component.SceneState.Get(component.SceneState.MustFirst(w))

		factory.CreateBodyPart(w, sceneObj.Project, snakeEntity, sceneData.Project.Project.EntityDefinitionByIdentifier(tags.SnakeBody.Name()), tags.SnakeBody.Name(), tags.Collidable.String())
		factory.CreateFood(w, sceneObj.Project, sceneData.Project.Project.EntityDefinitionByIdentifier(tags.Food.Name()))

		gameStateDate := component.GameState.Get(component.GameState.MustFirst(w))
		gameStateDate.Score++
		if gameStateDate.Score == 15 {
			if next, ok := component.GetNextLevel(sceneData.CurrentScene); ok {
				event.SceneStateEvent.Publish(w, &event.SceneStateData{
					CurrentScene: component.LevelClearScene,
					NextScene:    next,
				})
			} else {
				event.SceneStateEvent.Publish(w, &event.SceneStateData{
					CurrentScene: component.GameOverScene,
				})
			}
		}

	default:
		slog.Error("unknown collectable")
		panic(0)
	}
}

func ProcessEvents(ecs *ecs.ECS) {
	events.ProcessAllEvents(ecs.World)
}
