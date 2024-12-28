package system

import (
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/event"
	"github.com/yohamta/donburi"
)

func checkCollision[T any](w donburi.World, shape resolv.IShape, c *donburi.ComponentType[T]) {
	c.Each(w, func(collider *donburi.Entry) {
		component.Collidable.Each(w, func(collideWith *donburi.Entry) {
			collidable := component.Collidable.Get(collideWith)

			if intersection := shape.Intersection(collidable.Shape); !intersection.IsEmpty() {
				event.CollideEvent.Publish(w, &event.Collide{
					CollideWith:  collideWith,
					Collider:     collider,
					Intersection: intersection,
				})
			}
		})
	})
}
