package system

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/event"
	dresolv "github.com/soockee/terminal-games/ldtk-snake/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/soockee/terminal-games/ldtk-snake/util"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateSnake(ecs *ecs.ECS) {
	snakeEntry, _ := component.Snake.First(ecs.World)
	// snakeData := component.Snake.Get(snakeEntry)
	snakeObject := dresolv.GetObject(snakeEntry)

	velocity := component.Velocity.Get(snakeEntry)
	snakeObject.Position = snakeObject.Position.Add(velocity.Velocity)

	if checkWallCollision(snakeObject) {
		velocity.Velocity = resolv.NewVector(0, 0)
		event.CollideEvent.Publish(ecs.World, &event.Collide{
			Type: event.CollideWall,
		})
	}

	if checkFoodCollision(snakeObject) {
		event.CollectEvent.Publish(ecs.World, &event.Collect{
			Type: component.FoodCollectable,
		})
	}

}

func DrawSnake(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Snake.Each(ecs.World, func(e *donburi.Entry) {
		velocity := component.Velocity.Get(e)
		// todo calc direction
		slog.Info("", slog.Float64("degree", util.CalculateAngle(velocity.Velocity)))
		angle := util.CalculateAngle(velocity.Velocity)
		component.DrawRotatedSprite(screen, e, angle)
	})
}

// move temporarily uses a speed of type int whiel figuring out the collision
func OnMoveEvent(w donburi.World, e *event.Move) {
	entity, _ := component.Snake.First(w)
	// snakeData := component.Snake.Get(entity)

	velocity := component.Velocity.Get(entity)
	switch e.Direction {
	case component.ActionMoveUp:
		// velocity.Velocity = resolv.NewVector(0, -1)
		velocity.Velocity = resolv.NewVector(0, -1).Add(velocity.Velocity)

	case component.ActionMoveDown:
		// velocity.Velocity = resolv.NewVector(0, 1)
		velocity.Velocity = resolv.NewVector(0, 1).Add(velocity.Velocity)

	case component.ActionMoveLeft:
		// velocity.Velocity = resolv.NewVector(-1, 0)
		velocity.Velocity = resolv.NewVector(-1, 0).Add(velocity.Velocity)

	case component.ActionMoveRight:
		// velocity.Velocity = resolv.NewVector(1, 0)
		velocity.Velocity = resolv.NewVector(1, 0).Add(velocity.Velocity)

	}
	// velocity.Velocity.X *= snakeData.Speed
	// velocity.Velocity.Y *= snakeData.Speed
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
