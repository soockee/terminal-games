package systems

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/components"
	dresolv "github.com/soockee/terminal-games/ldtk-snake/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateSnake(ecs *ecs.ECS) {
	snakeEntry, _ := components.Snake.First(ecs.World)
	snakeData := components.Snake.Get(snakeEntry)
	snakeObject := dresolv.GetObject(snakeEntry)
	control := components.Control.Get(snakeEntry)

	move(control.InputHandler, snakeObject, snakeData)
	if checkWallCollision(snakeObject, snakeData) {
		slog.Info("Hit the Wall")
	}
}

func DrawSnake(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Snake.Each(ecs.World, func(e *donburi.Entry) {
		o := dresolv.GetObject(e)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(o.Position.X), float64(o.Position.Y))
		sprite := components.Sprite.Get(e)
		screen.DrawImage(sprite.Image, op)
	})
}

// move temporarily uses a speed of type int whiel figuring out the collision
func move(inputHandler *input.Handler, snakeObject *resolv.Object, snakeData *components.SnakeData) {
	if inputHandler.ActionIsPressed(components.ActionMoveUp) {
		snakeData.Direction = components.ActionMoveUp
		snakeObject.Position.Y -= snakeData.Speed
	}
	if inputHandler.ActionIsPressed(components.ActionMoveDown) {
		snakeData.Direction = components.ActionMoveDown
		snakeObject.Position.Y += snakeData.Speed
	}
	if inputHandler.ActionIsPressed(components.ActionMoveLeft) {
		snakeData.Direction = components.ActionMoveLeft
		snakeObject.Position.X -= snakeData.Speed
	}
	if inputHandler.ActionIsPressed(components.ActionMoveRight) {
		snakeData.Direction = components.ActionMoveRight
		snakeObject.Position.X += snakeData.Speed
	}
}

func checkWallCollision(snakeObject *resolv.Object, snakeData *components.SnakeData) bool {
	switch snakeData.Direction {
	case components.ActionMoveUp:
		if check := snakeObject.Check(0, -snakeData.Speed, tags.Wall.Name()); check != nil {
			slog.Info("Check Up", slog.Any("Check Info", check))
			return true
		}
	case components.ActionMoveDown:
		if check := snakeObject.Check(0, snakeData.Speed, tags.Wall.Name()); check != nil {
			slog.Info("Check Down", slog.Any("Check Info", check))

			return true
		}
	case components.ActionMoveLeft:
		if check := snakeObject.Check(-snakeData.Speed, 0, tags.Wall.Name()); check != nil {
			slog.Info("Check Left", slog.Any("Check Info", check))

			return true
		}
	case components.ActionMoveRight:
		if check := snakeObject.Check(snakeData.Speed, 0, tags.Wall.Name()); check != nil {
			slog.Info("Check Right", slog.Any("Check Info", check))

			return true
		}
	}
	return false
}
