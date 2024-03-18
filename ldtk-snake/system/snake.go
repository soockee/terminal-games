package system

import (
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
	snakeEntry := component.Snake.MustFirst(ecs.World)
	snakeObject := dresolv.GetObject(snakeEntry)

	snapshotHistory(ecs.World, snakeEntry)

	moveSnake(ecs)

	checkWallCollision(ecs.World, snakeObject)

	checkFoodCollision(ecs.World, snakeObject)

	// component.SnakeBody.Each(ecs.World, func(e *donburi.Entry) {
	// 	dresolv.GetObject(e).AddTags(tags.Collidable.Name())
	// 	e.AddComponent(tags.Collidable)
	// })

}

func DrawSnake(ecs *ecs.ECS, screen *ebiten.Image) {
	e := tags.Snake.MustFirst(ecs.World)
	snake := component.Snake.Get(e)
	DrawSnakeBody(ecs, screen, snake.Tail)
	velocity := component.Velocity.Get(e)
	angle := util.CalculateAngle(velocity.Velocity)
	component.DrawRotatedSprite(screen, e, angle)
	// component.DrawPlaceholder(screen, dresolv.GetObject(e), angle)
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
	// component.DrawPlaceholder(screen, dresolv.GetObject(next.Entry), angle)
}

// move temporarily uses a speed of type int whiel figuring out the collision
func OnMoveEvent(w donburi.World, e *event.Move) {
	entity := component.Snake.MustFirst(w)
	// snakeData := component.Snake.Get(entity)

	velocity := component.Velocity.Get(entity)
	switch e.Action {
	case component.ActionMoveUp:
		velocity.Velocity = resolv.NewVector(0, -1).Add(velocity.Velocity)

	case component.ActionMoveDown:
		velocity.Velocity = resolv.NewVector(0, 1).Add(velocity.Velocity)

	case component.ActionMoveLeft:
		velocity.Velocity = resolv.NewVector(-1, 0).Add(velocity.Velocity)

	case component.ActionMoveRight:
		velocity.Velocity = resolv.NewVector(1, 0).Add(velocity.Velocity)

	case component.ActionMoveHalt:
		velocity.Velocity = velocity.Velocity.Unit()
	}
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

func moveSnake(ecs *ecs.ECS) {
	snakeEntry := component.Snake.MustFirst(ecs.World)
	snakeData := component.Snake.Get(snakeEntry)
	snakeObject := dresolv.GetObject(snakeEntry)

	velocity := component.Velocity.Get(snakeEntry)

	// Stepwise movement based on velocity
	stepSize := 1.0 / math.Max(math.Abs(velocity.Velocity.X), math.Abs(velocity.Velocity.Y)) // Adjust step size based on velocity magnitude
	for step := 0.0; step <= 1.0; step += stepSize {
		updateSnakeBody(ecs.World, snakeData.Tail)
		// check if out of level bounds and teleport to opposite side
		pos := snakeObject.Position.Add(velocity.Velocity.Scale(step))
		space := component.Space.MustFirst(ecs.World)
		spaceObj := component.Space.Get(space)

		maxX := float64(spaceObj.Width() * spaceObj.CellWidth)
		maxY := float64(spaceObj.Height() * spaceObj.CellHeight)

		if pos.X > maxX {
			pos.X = 0
		} else if pos.X < 0 {
			pos.X = maxX
		}
		if pos.Y > maxY {
			pos.Y = 0
		} else if pos.Y < 0 {
			pos.Y = maxY
		}

		snakeObject.Position = pos
		checkBodyCollision(ecs.World, snakeObject)
	}
}

func updateSnakeBody(w donburi.World, next *component.SnakeBodyData) {
	if next == nil {
		return
	}
	nextObj := dresolv.GetObject(next.Entry)

	var prev *resolv.Object
	var history []component.HistoryData

	if next.Previous == nil {
		snake := component.Snake.MustFirst(w)

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
