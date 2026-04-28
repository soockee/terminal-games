package system

import (
	"math/rand/v2"

	"github.com/soockee/terminal-games/flappy-bird/component"
	"github.com/soockee/terminal-games/flappy-bird/event"
	"github.com/soockee/terminal-games/flappy-bird/physics"
	"github.com/soockee/terminal-games/flappy-bird/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// UpdateCollision checks the bird against all Collidable entities.
// A hit sets GameOver.Dead = true.
func UpdateCollision(e *ecs.ECS) {
	goEntry, ok := component.GameState.First(e.World)
	if !ok {
		return
	}
	gd := component.GameState.Get(goEntry)
	if gd.Dead || gd.Paused {
		return
	}

	birdEntry, ok := tags.Bird.First(e.World)
	if !ok {
		return
	}
	bb := component.Shape.Get(birdEntry).Shape.Bounds()

	tags.Collidable.Each(e.World, func(entry *donburi.Entry) {
		ob := component.Shape.Get(entry).Shape.Bounds()
		if physics.Overlaps(
			bb.Min.X, bb.Min.Y, bb.Max.X, bb.Max.Y,
			ob.Min.X, ob.Min.Y, ob.Max.X, ob.Max.Y,
		) {
			gd.Dead = true
			// Random death SFX.
			deathSFX := []string{component.SFXHurt, component.SFXExplosion}
			name := deathSFX[rand.IntN(len(deathSFX))]
			event.AudioEvent.Publish(e.World, event.AudioEventData{Name: name})
		}
	})
}
