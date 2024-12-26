package system

import (
	"log/slog"

	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/event"
	"github.com/yohamta/donburi"
)

func checkCollision[T any](w donburi.World, shape resolv.IShape, c *donburi.ComponentType[T]) {
	c.Each(w, func(e *donburi.Entry) {
		component.Collidable.Each(w, func(e *donburi.Entry) {
			collidable := component.Collidable.Get(e)

			if intersection := shape.Intersection(collidable.Shape); !intersection.IsEmpty() {
				slog.Debug("Collision", slog.Any("intersection", intersection))
				event.CollideEvent.Publish(w, &event.Collide{
					Type: collidable.Type,
				})
			}
		})
	})
}
