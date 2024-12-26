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
	collidable := component.Collidable.Get(e)
	spriteData := component.Sprite.Get(e)
	sprite := spriteData.Images[0]
	component.DrawRepeatedSprite(screen, sprite, collidable.Shape)
}

func OnMoveEvent(w donburi.World, e *event.Move) {
	entry := component.Player.MustFirst(w)
	velocity := component.Velocity.Get(entry)
	velocity.Velocity = velocity.Velocity.Add(e.Direction)
}

func moveplayer(ecs *ecs.ECS) {
	playerEntry := component.Player.MustFirst(ecs.World)
	player := component.Player.Get(playerEntry)
	collidable := component.Collidable.Get(playerEntry)

	shape := collidable.Shape.(*resolv.ConvexPolygon)

	velocity := component.Velocity.Get(playerEntry)
	space := component.Space.Get(component.Space.MustFirst(ecs.World))
	maxX := float64(space.Width())
	if shape.Bounds().Min.X <= 0 {
		// allows player movement to the right
		shape.SetX(shape.Center().X + 1)
		velocity.Velocity.X = 0
		return
	} else if shape.Bounds().Max.X >= maxX {
		// allows player movement to the left
		shape.SetX(shape.Center().X - 1)
		velocity.Velocity.X = 0
		return
	}

	shape.Move(velocity.Velocity.X, velocity.Velocity.Y)

	velocity.Velocity = velocity.Velocity.Mult(resolv.NewVector(player.SpeedFriction, player.SpeedFriction))
}
