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
	// slog.Info("t", snakeObject)
	control := components.Control.Get(snakeEntry)
	move(control.InputHandler, snakeObject, snakeData)
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
func move(inputHandler *input.Handler, snake *resolv.Object, snakeData *components.SnakeData) {
	if inputHandler.ActionIsPressed(components.ActionMoveUp) {
		slog.Info("Up")
		snake.Position.Y -= snakeData.Speed
	}
	if inputHandler.ActionIsPressed(components.ActionMoveDown) {
		slog.Info("Down")
		snake.Position.Y += snakeData.Speed
	}
	if inputHandler.ActionIsPressed(components.ActionMoveLeft) {
		slog.Info("Left")
		snake.Position.X -= snakeData.Speed
	}
	if inputHandler.ActionIsPressed(components.ActionMoveRight) {
		slog.Info("Right")
		snake.Position.X += snakeData.Speed
	}
}
