package system

import (
	"log/slog"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/tags"
	"github.com/soockee/terminal-games/breakout/util"
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
	checkCollision(ecs.World, ball.Shape, component.Ball)
}

func DrawBall(ecs *ecs.ECS, screen *ebiten.Image) {
	e := tags.Ball.MustFirst(ecs.World)
	ball := component.Ball.Get(e)
	spriteData := component.Sprite.Get(e)
	sprite := spriteData.Images[0]
	component.DrawScaledSprite(screen, sprite, ball.Shape)
}

func OnBallCollisionEvent(w donburi.World, e *event.Collide) {
	entry := component.Ball.MustFirst(w)
	velocity := component.Velocity.Get(entry)

	CollideWithType := e.CollideWith.Archetype().Layout()
	ball := component.Ball.Get(entry)
	// wall := component.Collidable.Get(e.CollideWith).Shape
	//space := component.Space.Get(component.Space.MustFirst(w))

	if CollideWithType.HasComponent(tags.Wall) {
		velocity.Velocity = velocity.Velocity.Reflect(e.Intersection.Intersections[0].Normal)
	} else if CollideWithType.HasComponent(tags.Player) {
		playerVelocity := component.Velocity.Get(e.CollideWith)
		out := velocity.Velocity.Reflect(e.Intersection.Intersections[0].Normal)
		velocity.Velocity = out.Add(playerVelocity.Velocity.Scale(0.9))
		velocity.Velocity = util.LimitMagnitude(velocity.Velocity, ball.MaxSpeed)
	}
	slog.Debug("Ball magnitude", slog.Any("magnitude", velocity.Velocity.Magnitude()))
}

func OnReleaseEvent(w donburi.World, e *event.Release) {
	entry := component.Ball.MustFirst(w)
	ball := component.Ball.Get(entry)
	velocity := component.Velocity.Get(entry)

	// randomize direction
	direction := resolv.NewVector(2*rand.Float64()-1, -1)

	velocity.Velocity = velocity.Velocity.Add(direction)

	velocity.Velocity = util.LimitMagnitude(velocity.Velocity, ball.MaxSpeed)

}

func moveball(ecs *ecs.ECS) {
	ballEntry := component.Ball.MustFirst(ecs.World)
	ball := component.Ball.Get(ballEntry)
	velocity := component.Velocity.Get(ballEntry)
	space := component.Space.Get(component.Space.MustFirst(ecs.World))

	ball.Shape.Move(velocity.Velocity.X, velocity.Velocity.Y)
	if ball.Shape.Position().Y <= 0 {
		ball.Shape.SetPosition(ball.Shape.Position().X, ball.Shape.Position().Y+ball.Shape.Bounds().Height()+1)
	}
	if ball.Shape.Position().Y >= float64(space.Height()) {
		ball.Shape.SetPosition(ball.Shape.Position().X, ball.Shape.Position().Y-ball.Shape.Bounds().Height()-1)
	}
	if ball.Shape.Position().X <= 0 {
		ball.Shape.SetPosition(ball.Shape.Position().X+ball.Shape.Bounds().Width()+1, ball.Shape.Position().Y)
	}
	if ball.Shape.Position().X >= float64(space.Width()) {
		ball.Shape.SetPosition(ball.Shape.Position().X-ball.Shape.Bounds().Width()-1, ball.Shape.Position().Y)
	}
}
