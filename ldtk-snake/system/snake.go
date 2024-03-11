package system

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/event"
	dresolv "github.com/soockee/terminal-games/ldtk-snake/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateSnake(ecs *ecs.ECS) {
	snakeEntry, _ := component.Snake.First(ecs.World)
	snakeData := component.Snake.Get(snakeEntry)
	snakeObject := dresolv.GetObject(snakeEntry)

	if checkWallCollision(snakeObject, snakeData) {
		slog.Info("Hit the Wall")
	}

}

func DrawSnake(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Snake.Each(ecs.World, func(e *donburi.Entry) {
		component.DrawSprite(screen, e)
	})
}

// move temporarily uses a speed of type int whiel figuring out the collision
func HandleMoveEvent(w donburi.World, e *event.Move) {
	entity, _ := component.Snake.First(w)
	snakeData := component.Snake.Get(entity)
	snakeObject := dresolv.GetObject(entity)

	switch e.Direction {
	case component.ActionMoveUp:
		snakeData.Direction = component.ActionMoveUp
		snakeObject.Position.Y -= snakeData.Speed
	case component.ActionMoveDown:
		snakeData.Direction = component.ActionMoveDown
		snakeObject.Position.Y += snakeData.Speed
	case component.ActionMoveLeft:
		snakeData.Direction = component.ActionMoveLeft
		snakeObject.Position.X -= snakeData.Speed
	case component.ActionMoveRight:
		snakeData.Direction = component.ActionMoveRight
		snakeObject.Position.X += snakeData.Speed
	}

}

func checkWallCollision(snakeObject *resolv.Object, snakeData *component.SnakeData) bool {
	switch snakeData.Direction {
	case component.ActionMoveUp:
		if check := snakeObject.Check(0, -snakeData.Speed, tags.Wall.Name()); check != nil {
			return true
		}
	case component.ActionMoveDown:
		if check := snakeObject.Check(0, snakeData.Speed, tags.Wall.Name()); check != nil {
			return true
		}
	case component.ActionMoveLeft:
		if check := snakeObject.Check(-snakeData.Speed, 0, tags.Wall.Name()); check != nil {
			return true
		}
	case component.ActionMoveRight:
		if check := snakeObject.Check(snakeData.Speed, 0, tags.Wall.Name()); check != nil {
			return true
		}
	}
	return false
}
