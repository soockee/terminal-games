package factory

import (
	"log/slog"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
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
	center := obj.Center()
	radius := math.Abs(obj.Position.X - center.X)
	obj.SetShape(resolv.NewCircle(center.X, center.Y, radius))
	component.Object.Set(snake, obj)

	component.Snake.SetValue(snake, component.SnakeData{
		Speed:             8,
		SpeedAcceleration: 1.05,
		SpeedFriction:     0.94,
		Tail:              nil,
		History:           []component.HistoryData{},
		HistoryTimer:      engine.NewTimer(time.Millisecond * 16),
	})

	sprite, err := project.GetSpriteByEntityInstance(entity)
	if err != nil {
		slog.Error("Sprite not found")
		panic(0)
	}
	component.Sprite.SetValue(snake, component.SpriteData{Images: map[int]*ebiten.Image{0: sprite}})

	CreateBodyPart(ecs.World, project, snake, project.Project.EntityDefinitionByIdentifier(dtags.SnakeBody.Name()), component.SnakeBodyTypeFood, dtags.SnakeBody.Name())

	return snake
}

func CreateBodyPart(world donburi.World, project *assets.LDtkProject, snakeEntry *donburi.Entry, entity *ldtkgo.EntityDefinition, sbt component.SnakeBodyType, tags ...string) {
	part := archetype.SnakeBody.SpawnInWorld(world)
	snakehead := component.Snake.Get(snakeEntry)
	partData := component.SnakeBody.Get(part)
	prev, _ := component.GetTail(snakehead)

	var boundsObj *resolv.Object
	if prev == nil {
		boundsObj = dresolv.GetObject(snakeEntry)
	} else {
		boundsObj = dresolv.GetObject(prev.Entry)
	}

	x, y := boundsObj.Shape.Bounds()
	width := y.X - x.X
	height := y.X - x.X
	X := float64(boundsObj.Position.X)
	Y := float64(boundsObj.Position.Y)

	obj := resolv.NewObject(X, Y, width, height, tags...)
	center := obj.Center()
	radius := math.Abs(obj.Position.X - center.X)
	obj.SetShape(resolv.NewCircle(center.X, center.Y, radius))
	component.Object.Set(part, obj)

	component.SnakeBody.SetValue(part, component.SnakeBodyData{
		Entry:         part,
		SnakeBodyType: sbt,
		Next:          nil,
		Previous:      nil,
	})

	spritesMap, err := project.GetSpritesByTag(dtags.SnakeBody.Name())

	sprites := map[int]*ebiten.Image{}
	for identifier, sprite := range spritesMap {
		switch identifier {
		case "SnakeBodyFood":
			sprites[int(component.SnakeBodyTypeFood)] = sprite
		case "SnakeBodyMouse":
			sprites[int(component.SnakeBodyTypeMouse)] = sprite
		}
	}

	if err != nil {
		slog.Error("Error", err)
		panic(0)
	}

	component.Sprite.SetValue(part, component.SpriteData{Images: sprites})

	spaceEntry := component.Space.MustFirst(world)
	dresolv.Add(spaceEntry, part)

	snakehead.SetTail(partData)
	partData.Previous = prev
}
