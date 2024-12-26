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
	velocity := component.Velocity.Get(entry)
	velocity.Velocity = velocity.Velocity.Add(e.Direction)
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
	maxX := float64(space.Width())
	if player.Shape.Bounds().Min.X <= 0 {
		// allows player movement to the right
		player.Shape.SetX(player.Shape.Center().X + 1)
		velocity.Velocity.X = 0
		return
	} else if player.Shape.Bounds().Max.X >= maxX {
		// allows player movement to the left
		player.Shape.SetX(player.Shape.Center().X - 1)
		velocity.Velocity.X = 0
		return
	}

	player.Shape.Move(velocity.Velocity.X, velocity.Velocity.Y)

	velocity.Velocity = velocity.Velocity.Mult(resolv.NewVector(player.SpeedFriction, player.SpeedFriction))
}
