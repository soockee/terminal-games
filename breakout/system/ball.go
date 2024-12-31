package system

import (
	"image/color"
	"log/slog"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
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
		collision := checkCollision(ecs.World, ball.Shape, component.Ball)
		if collision != nil {
			event.CollideEvent.Publish(ecs.World, collision)
		}
	}
}

func DrawBall(ecs *ecs.ECS, screen *ebiten.Image) {
	e := tags.Ball.MustFirst(ecs.World)
	ball := component.Ball.Get(e)

	sprite := component.Sprite.Get(e)

	// component.DrawScaledSprite(screen, sprite, ball.Shape)
	component.DrawPlaceholder(screen, sprite.Images[0], ball.Shape, 0, color.White, false)
}

func OnBallCollisionEvent(w donburi.World, e *event.Collide) {
	entry := component.Ball.MustFirst(w)
	velocity := component.Velocity.Get(entry)

	CollideWithType := e.CollideWith.Archetype().Layout()
	ball := component.Ball.Get(entry)

	var bestNormal resolv.Vector
	var bestWeight float64 = -1
	var closestIntersection float64 = math.MaxFloat64
	shortTraceDistance := 1.0

	if len(e.Intersection.Intersections) > 0 {
		for _, i := range e.Intersection.Intersections {
			normal := i.Normal // Get the surface normal
			distance := ball.Shape.Position().Distance(i.Point)

			// Check alignment with velocity (positive alignment indicates the surface is valid)
			alignment := velocity.Velocity.Dot(normal)
			if alignment >= 0 { // Only consider surfaces the ball is approaching
				continue
			}

			// Calculate weight: combine inverse distance and alignment
			weight := math.Abs(alignment) / (distance + 0.001)

			// Predict the reflected velocity
			reflectedVelocity := velocity.Velocity.Reflect(normal)

			// Raycast along the reflected trajectory to check for immediate collisions
			nextPosition := i.Point.Add(reflectedVelocity.Scale(shortTraceDistance))
			shape := resolv.NewCircle(nextPosition.X, nextPosition.Y, ball.Shape.Radius())
			if checkCollision(w, shape, component.Ball) == nil {
				weight *= 1.5 // Increase weight for open paths
			}

			// Choose the best intersection based on weight or proximity
			if weight > bestWeight || (weight == bestWeight && distance < closestIntersection) {
				bestWeight = weight
				bestNormal = normal
				closestIntersection = distance
			}
		}
	}

	if CollideWithType.HasComponent(tags.Wall) {
		velocity.Velocity = velocity.Velocity.Reflect(bestNormal)

	} else if CollideWithType.HasComponent(tags.Player) {
		playerVelocity := component.Velocity.Get(e.CollideWith)
		out := velocity.Velocity.Reflect(bestNormal)
		playerFactor := util.LimitMagnitude(playerVelocity.Velocity, 3)
		velocity.Velocity = out.Add(playerFactor)
		velocity.Velocity = util.LimitMagnitude(velocity.Velocity, ball.MaxSpeed)
		ball.CooldownTimer = *engine.NewTimer(ball.CollisionCooldownPlayer)

	} else if CollideWithType.HasComponent(tags.Brick) {
		velocity.Velocity = velocity.Velocity.Reflect(bestNormal)
		space := component.Space.Get(component.Space.MustFirst(w))

		brick := component.Brick.Get(e.CollideWith)
		brick.Health--
		if brick.Health <= 0 {
			collidable := component.Collidable.Get(e.CollideWith)
			event.CreateEntityEvent.Publish(w, &event.CreateEntityData{
				Tags:       []string{tags.Explosion.Name()},
				Identifier: "Explosion",
				X:          collidable.Shape.Bounds().Min.X,
				Y:          collidable.Shape.Bounds().Min.Y,
				W:          collidable.Shape.Bounds().Width(),
				H:          collidable.Shape.Bounds().Height(),
			})
			space.Remove(collidable.Shape)
			w.Remove(e.CollideWith.Entity())
		}
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
