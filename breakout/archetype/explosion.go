package archetype

import (
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/ganim8/v2"
)

var (
	Explosion = newArchetype(
		tags.Animation,

		component.Animation,
	)
)

func NewExplosion(w donburi.World, shape resolv.IShape, sprite ganim8.Animation) *donburi.Entry {
	explosion := Explosion.SpawnInWorld(w)

	component.Space.Get(component.Space.MustFirst(w)).Add(shape)

	sprite.SetOnLoop(ganim8.PauseAtEnd)
	data := component.AnimationsData{
		Animation: &sprite,
		Shape:     shape,
		Loop:      false,
	}

	component.Animation.SetValue(explosion, data)

	return explosion
}
