package system

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/tags"
	"github.com/soockee/terminal-games/breakout/util"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdatePlayer(ecs *ecs.ECS) {
	// playerEntry := component.Player.MustFirst(ecs.World)
	// player := component.Player.Get(playerEntry)

	moveplayer(ecs)
	// slog.Debug("Player", slog.Any("Player", player))
	//checkWallCollision(ecs.World, snakeObject)
}

func DrawPlayer(ecs *ecs.ECS, screen *ebiten.Image) {
	e := tags.Player.MustFirst(ecs.World)
	player := component.Player.Get(e)
	spriteData := component.Sprite.Get(e)
	sprite := spriteData.Images[0]
	component.DrawRepeatedSprite(screen, sprite, player.Shape)
}

func OnMoveEvent(w donburi.World, e *event.Move) {
	entry := component.Player.MustFirst(w)
	player := component.Player.Get(entry)

	velocity := component.Velocity.Get(entry)
	slog.Debug("Button pressed", slog.Any("Button", e.Action))
	direction := util.DirectionVector(player.Shape.Position(), e.Direction)
	speed := player.Speed

	if e.Boost {
		speed *= 2
	}

	velocity.Velocity = direction.ClampMagnitude(speed)
}

func checkCollision(w donburi.World, playerObject *resolv.ConvexPolygon) {
	component.Player.Each(w, func(e *donburi.Entry) {
		tags.Collidable.Each(w, func(e *donburi.Entry) {
			collidableObject := component.Space.Get(e)
			if intersection := playerObject.Intersection(collidableObject.FilterShapes().First()); !intersection.IsEmpty() {
				event.CollideEvent.Publish(w, &event.Collide{
					Type: event.CollideBody,
				})
			}
		})
	})
}

func moveplayer(ecs *ecs.ECS) {
	playerEntry := component.Player.MustFirst(ecs.World)
	player := component.Player.Get(playerEntry)

	velocity := component.Velocity.Get(playerEntry)

	space := component.Space.Get(component.Space.MustFirst(ecs.World))

	maxX := float64(space.Width() * space.CellWidth())

	if player.Shape.Bounds().Min.X <= 0 {
		player.Shape.SetX(0)
		velocity.Velocity.X = 0
		return
	} else if player.Shape.Bounds().Max.X >= maxX {
		player.Shape.SetX(maxX)
		velocity.Velocity.X = 0
		return
	}

	player.Shape.Move(velocity.Velocity.X, velocity.Velocity.Y)
}
