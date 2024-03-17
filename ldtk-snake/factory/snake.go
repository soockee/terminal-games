package factory

import (
	"log/slog"
	"time"

	"github.com/solarlune/resolv"
	"github.com/soockee/ldtkgo"
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/engine"
	dresolv "github.com/soockee/terminal-games/ldtk-snake/resolv"
	dtags "github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateSnake(ecs *ecs.ECS, project *assets.LDtkProject, entity *ldtkgo.Entity) *donburi.Entry {
	snake := archetype.Snake.Spawn(ecs)

	width := float64(entity.Width)
	height := float64(entity.Height)
	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])

	obj := resolv.NewObject(X, Y, width, height, entity.Tags...)
	obj.SetShape(resolv.NewRectangle(X, Y, width, height))
	component.Object.Set(snake, obj)

	component.Snake.SetValue(snake, component.SnakeData{
		Speed:        5,
		Tail:         nil,
		History:      []component.HistoryData{},
		HistoryTimer: engine.NewTimer(time.Millisecond * 16),
	})

	sprite, err := project.GetSpriteByEntityInstance(entity)
	if err != nil {
		slog.Error("Sprite not found")
		panic(0)
	}
	component.Sprite.SetValue(snake, component.SpriteData{Image: sprite})

	// no initial body
	// CreateBodyPart(ecs.World, project, snake, project.Project.EntityDefinitionByIdentifier(dtags.SnakeBody.Name()), dtags.SnakeBody.Name())

	return snake
}

func CreateBodyPart(world donburi.World, project *assets.LDtkProject, snakeEntry *donburi.Entry, entity *ldtkgo.EntityDefinition, tags ...string) {
	part := archetype.SnakeBody.SpawnInWorld(world)

	boundsObj := dresolv.GetObject(snakeEntry)

	x, y := boundsObj.Shape.Bounds()
	width := y.X - x.X
	height := y.X - x.X
	X := float64(boundsObj.Position.X)
	Y := float64(boundsObj.Position.Y)

	obj := resolv.NewObject(X, Y, width, height, tags...)
	center := obj.Center()
	obj.SetShape(resolv.NewCircle(center.X, center.Y, width/height))
	component.Object.Set(part, obj)

	component.SnakeBody.SetValue(part, component.SnakeBodyData{
		Entry:    part,
		Next:     nil,
		Previous: nil,
	})

	sprite, err := project.GetSpriteByIdentifier(dtags.SnakeBody.Name())
	if err != nil {
		slog.Error("Sprite not found")
		panic(0)
	}

	component.Sprite.SetValue(part, component.SpriteData{Image: sprite})

	spaceEntry, ok := component.Space.First(world)
	if !ok {
		slog.Error("Sprite not found")
		panic(0)
	}
	dresolv.Add(spaceEntry, part)

	snakehead := component.Snake.Get(snakeEntry)
	partData := component.SnakeBody.Get(part)
	prev, _ := component.GetTail(snakehead)
	snakehead.SetTail(partData)
	partData.Previous = prev
}
