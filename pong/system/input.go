package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/soockee/terminal-games/pong/component"
	"github.com/soockee/terminal-games/pong/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

const paddleSpeed = 5.0

// UpdateInput reads keyboard input and sets paddle velocities.
// F3 toggles the debug overlay.
func UpdateInput(e *ecs.ECS) {
	// Toggle debug on F3 (single press via inpututil)
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		if entry, ok := component.Debug.First(e.World); ok {
			d := component.Debug.Get(entry)
			d.Enabled = !d.Enabled
		}
	}

	tags.Paddle.Each(e.World, func(entry *donburi.Entry) {
		paddle := component.Paddle.Get(entry)
		vel := component.Velocity.Get(entry)

		vel.Y = 0
		switch paddle.Side {
		case component.SideLeft:
			if ebiten.IsKeyPressed(ebiten.KeyW) {
				vel.Y = -paddle.Speed
			}
			if ebiten.IsKeyPressed(ebiten.KeyS) {
				vel.Y = paddle.Speed
			}
		case component.SideRight:
			if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
				vel.Y = -paddle.Speed
			}
			if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
				vel.Y = paddle.Speed
			}
		}
	})
}
