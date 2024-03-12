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

	velocity := component.Velocity.Get(snakeEntry)
	snakeObject.Position = snakeObject.Position.Add(velocity.Velocity)

	if checkWallCollision(snakeObject) {
		velocity.Velocity = resolv.NewVector(0, 0)
		slog.Info("Hit the Wall")
	}

	if checkFoodCollision(snakeObject) {
		snakeData.Speed += 1
		slog.Info("Hit Food")
	}

}

func DrawSnake(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Snake.Each(ecs.World, func(e *donburi.Entry) {
		velocity := component.Velocity.Get(e)

		angle := 0.0
		if velocity.Velocity.X == 1 {
			angle = 90.0
		}
		if velocity.Velocity.Y == 1 {
			angle = 180
		}
		if velocity.Velocity.X == -1 {
			angle = 270.0
		}

		// todo calc direction
		component.DrawRotatedSprite(screen, e, angle)
	})
}

// move temporarily uses a speed of type int whiel figuring out the collision
func HandleMoveEvent(w donburi.World, e *event.Move) {
	entity, _ := component.Snake.First(w)
	snakeData := component.Snake.Get(entity)

	velocity := component.Velocity.Get(entity)
	switch e.Direction {
	case component.ActionMoveUp:
		velocity.Velocity = resolv.NewVector(0, -1)
	case component.ActionMoveDown:
		velocity.Velocity = resolv.NewVector(0, 1)
	case component.ActionMoveLeft:
		velocity.Velocity = resolv.NewVector(-1, 0)
	case component.ActionMoveRight:
		velocity.Velocity = resolv.NewVector(1, 0)
	}
	velocity.Velocity.X *= snakeData.Speed
	velocity.Velocity.Y *= snakeData.Speed

}

func checkWallCollision(snakeObject *resolv.Object) bool {
	if check := snakeObject.Check(0, 0, tags.Wall.Name()); check != nil {
		return true
	}
	return false
}

func checkFoodCollision(snakeObject *resolv.Object) bool {
	if check := snakeObject.Check(0, 0, tags.Food.Name()); check != nil {
		return true
	}
	return false
}
