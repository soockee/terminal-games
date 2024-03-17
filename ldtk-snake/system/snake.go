package system

import (
	"log/slog"
	"math"

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
	snakeData := component.Snake.Get(snakeEntry)
	snakeObject := dresolv.GetObject(snakeEntry)

	velocity := component.Velocity.Get(snakeEntry)

	snapshotHistory(ecs.World, snakeEntry)

	// Update the position of the snake's head
	snakeObject.Position = snakeObject.Position.Add(velocity.Velocity)

	// Stepwise movement based on velocity
	stepSize := 1.0 / math.Max(math.Abs(velocity.Velocity.X), math.Abs(velocity.Velocity.Y)) // Adjust step size based on velocity magnitude
	for step := 0.0; step <= 1.0; step += stepSize {
		snakeObject.Position = snakeObject.Position.Add(velocity.Velocity.Scale(step))
		updateSnakeBody(ecs.World, snakeData.Tail)
		checkBodyCollision(ecs.World, snakeObject)
	}

	// for each velocity ( vector of X and Y) stepwise increase the position by one, check for each step body collision

	checkWallCollision(ecs.World, snakeObject)

	checkFoodCollision(ecs.World, snakeObject)

	component.SnakeBody.Each(ecs.World, func(e *donburi.Entry) {
		dresolv.GetObject(e).AddTags(tags.Collidable.Name())
		e.AddComponent(tags.Collidable)
	})

}

func DrawSnake(ecs *ecs.ECS, screen *ebiten.Image) {
	e, ok := tags.Snake.First(ecs.World)
	if !ok {
		slog.Error("snake not found in draw")
		panic(0)
	}

	snake := component.Snake.Get(e)
	DrawSnakeBody(ecs, screen, snake.Tail)
	velocity := component.Velocity.Get(e)
	angle := util.CalculateAngle(velocity.Velocity)
	component.DrawRotatedSprite(screen, e, angle)
	component.DrawPlaceholder(screen, dresolv.GetObject(e), angle)
}

func DrawSnakeBody(ecs *ecs.ECS, screen *ebiten.Image, next *component.SnakeBodyData) {
	if next == nil {
		return
	}
	DrawSnakeBody(ecs, screen, next.Next)
	velocity := component.Velocity.Get(next.Entry)
	angle := util.CalculateAngle(velocity.Velocity)
	component.DrawRotatedSprite(screen, next.Entry, angle)
	// debug
	component.DrawPlaceholder(screen, dresolv.GetObject(next.Entry), angle)
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

func checkWallCollision(w donburi.World, snakeObject *resolv.Object) bool {
	if check := snakeObject.Check(0, 0, tags.Wall.Name()); check != nil {
		event.CollideEvent.Publish(w, &event.Collide{
			Type: event.CollideWall,
		})
	}
	return false
}

func checkFoodCollision(w donburi.World, snakeObject *resolv.Object) {
	if check := snakeObject.Check(0, 0, tags.Food.Name()); check != nil {
		event.CollectEvent.Publish(w, &event.Collect{
			Type: component.FoodCollectable,
		})
	}
}

func checkBodyCollision(w donburi.World, snakeObject *resolv.Object) {
	// if check := snakeObject.Check(0, 0, tags.Collidable.Name()); check != nil {
	// 	return true
	// }

	component.SnakeBody.Each(w, func(e *donburi.Entry) {
		obj := dresolv.GetObject(e)
		if !obj.HasTags(tags.Collidable.Name()) {
			return
		}
		if intersection := snakeObject.Shape.Intersection(0, 0, obj.Shape); intersection != nil {
			event.CollideEvent.Publish(w, &event.Collide{
				Type: event.CollideBody,
			})
		}
	})
}

func updateSnakeBody(w donburi.World, next *component.SnakeBodyData) {
	if next == nil {
		return
	}
	nextObj := dresolv.GetObject(next.Entry)

	var prev *resolv.Object
	var history []component.HistoryData

	if next.Previous == nil {
		snake, ok := component.Snake.First(w)
		if !ok {
			slog.Error("could not find snake")
			panic(0)
		}
		history = component.Snake.Get(snake).History
		prev = dresolv.GetObject(snake)
		direction := util.DirectionVector(prev.Position, nextObj.Position)
		component.Velocity.Get(next.Entry).Velocity = direction
	} else {
		history = component.SnakeBody.Get(next.Previous.Entry).History
		prev = dresolv.GetObject(next.Previous.Entry)
	}

	// keep history short
	maxlength := int(math.Min(float64(len(history)), 50))
	history = history[maxlength:]

	var vel resolv.Vector
	var pos resolv.Vector

	isInvalid := true
	for i := len(history) - 1; i > 0; i-- {
		if history[i].Velocity.X != 0 || history[i].Velocity.Y != 0 {
			vel = history[i].Velocity
			pos = history[i].Position
			if math.Abs(pos.X-prev.Position.X) > prev.Size.X || math.Abs(pos.Y-prev.Position.Y) > prev.Size.Y {
				isInvalid = false
				break
			}
		}
	}
	if isInvalid {
		return
	}
	if vel.X == 0 && vel.Y == 0 {
		return
	}
	dispossition := vel.Unit().Invert()
	nextObj.Position = pos.Add(dispossition)
	component.Velocity.Get(next.Entry).Velocity = vel
	updateSnakeBody(w, next.Next)
}

func snapshotHistory(w donburi.World, snakeEntry *donburi.Entry) {
	snakeData := component.Snake.Get(snakeEntry)
	snakeObject := dresolv.GetObject(snakeEntry)
	velocity := component.Velocity.Get(snakeEntry)
	if snakeData.HistoryTimer.IsReady() {
		snakeData.History = append(snakeData.History, component.HistoryData{
			Position: snakeObject.Position,
			Velocity: velocity.Velocity,
		})
		component.SnakeBody.Each(w, func(e *donburi.Entry) {
			obj := component.SnakeBody.Get(e)
			obj.History = append(obj.History, component.HistoryData{
				Position: dresolv.GetObject(e).Position,
				Velocity: component.Velocity.Get(e).Velocity,
			})
		})
		snakeData.HistoryTimer.Reset()
	} else {
		snakeData.HistoryTimer.Update()
	}
}
