package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateBall(ecs *ecs.ECS) {
	ballEntry := component.Ball.MustFirst(ecs.World)
	ball := component.Ball.Get(ballEntry)
	if ball == nil {
		return
	}

	moveball(ecs)
	// slog.Debug("Ball", slog.Any("Ball", ball))
	checkCollision(ecs.World, ball.Shape, component.Ball)
}

func DrawBall(ecs *ecs.ECS, screen *ebiten.Image) {
	e := tags.Ball.MustFirst(ecs.World)
	ball := component.Ball.Get(e)
	spriteData := component.Sprite.Get(e)
	sprite := spriteData.Images[0]
	component.DrawScaledSprite(screen, sprite, ball.Shape)
}

func OnCollisionEvent(w donburi.World, e *event.Move) {
	entry := component.Ball.MustFirst(w)
	velocity := component.Velocity.Get(entry)
	velocity.Velocity = velocity.Velocity.Add(e.Direction)
}

func OnReleaseEvent(w donburi.World, e *event.Release) {
	entry := component.Ball.MustFirst(w)
	//ball := component.Ball.Get(entry)
	velocity := component.Velocity.Get(entry)
	// TODO: smarter release direction
	direction := resolv.NewVector(0, -1)
	velocity.Velocity = velocity.Velocity.Add(direction)
}

func moveball(ecs *ecs.ECS) {
	ballEntry := component.Ball.MustFirst(ecs.World)
	ball := component.Ball.Get(ballEntry)
	velocity := component.Velocity.Get(ballEntry)
	//space := component.Space.Get(component.Space.MustFirst(ecs.World))

	ball.Shape.Move(velocity.Velocity.X, velocity.Velocity.Y)
}
