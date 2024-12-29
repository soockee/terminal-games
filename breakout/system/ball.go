package system

import (
	"image/color"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/engine"
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

	ball.CooldownTimer.Update()
	moveball(ecs.World)
	if ball.CooldownTimer.IsReady() {
		checkCollision(ecs.World, ball.Shape, component.Ball)
	}
}

func DrawBall(ecs *ecs.ECS, screen *ebiten.Image) {
	e := tags.Ball.MustFirst(ecs.World)
	ball := component.Ball.Get(e)
	spriteData := component.Sprite.Get(e)
	sprite := spriteData.Images[0]
	// component.DrawScaledSprite(screen, sprite, ball.Shape)
	component.DrawPlaceholder(screen, sprite, ball.Shape, 0, color.White, false)
}

func OnBallCollisionEvent(w donburi.World, e *event.Collide) {
	entry := component.Ball.MustFirst(w)
	velocity := component.Velocity.Get(entry)

	CollideWithType := e.CollideWith.Archetype().Layout()
	ball := component.Ball.Get(entry)

	if CollideWithType.HasComponent(tags.Wall) {
		velocity.Velocity = velocity.Velocity.Reflect(e.Intersection.Intersections[0].Normal)

	} else if CollideWithType.HasComponent(tags.Player) {
		playerVelocity := component.Velocity.Get(e.CollideWith)
		out := velocity.Velocity.Reflect(e.Intersection.Intersections[0].Normal)
		playerFactor := util.LimitMagnitude(playerVelocity.Velocity, 3)
		velocity.Velocity = out.Add(playerFactor)
		velocity.Velocity = util.LimitMagnitude(velocity.Velocity, ball.MaxSpeed)
		ball.CooldownTimer = *engine.NewTimer(ball.CollisionCooldownPlayer)

	} else if CollideWithType.HasComponent(tags.Brick) {
		velocity.Velocity = velocity.Velocity.Reflect(e.Intersection.Intersections[0].Normal)

		w.Remove(e.CollideWith.Entity())
	}

	moveball(w)

	slog.Debug("ball.CooldownTimer", slog.Any("cooldown", ball.CooldownTimer))
}

func moveball(w donburi.World) {
	ballEntry := component.Ball.MustFirst(w)
	ball := component.Ball.Get(ballEntry)
	velocity := component.Velocity.Get(ballEntry)
	space := component.Space.Get(component.Space.MustFirst(w))

	ball.Shape.Move(velocity.Velocity.X, velocity.Velocity.Y)
	if ball.Shape.Position().Y <= 0 {
		ball.Shape.SetPosition(ball.Shape.Position().X, ball.Shape.Position().Y+ball.Shape.Bounds().Height()+2)
	}
	if ball.Shape.Position().Y >= float64(space.Height()) {
		ball.Shape.SetPosition(ball.Shape.Position().X, ball.Shape.Position().Y-ball.Shape.Bounds().Height()-2)
	}
	if ball.Shape.Position().X <= 0 {
		ball.Shape.SetPosition(ball.Shape.Position().X+ball.Shape.Bounds().Width()+2, ball.Shape.Position().Y)
	}
	if ball.Shape.Position().X >= float64(space.Width()) {
		ball.Shape.SetPosition(ball.Shape.Position().X-ball.Shape.Bounds().Width()-2, ball.Shape.Position().Y)
	}
}
