package archetype

import (
	"log/slog"

	"github.com/solarlune/ldtkgo"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

var (
	Explosion = newArchetype(
		tags.Explosion,

		component.Animations,
		component.Position,
	)
)

func NewExplosion(w donburi.World, project *assets.LDtkProject, entity *ldtkgo.Entity) *donburi.Entry {
	explosion := Explosion.SpawnInWorld(w)

	X := float64(entity.Position[0])
	Y := float64(entity.Position[1])
	width := float64(entity.Width)
	height := float64(entity.Height)

	r := resolv.NewRectangleFromCorners(X, Y, X+width, Y+height)
	component.Space.Get(component.Space.MustFirst(w)).Add(r)

	explosionFX := project.Project.EntityDefinitionByIdentifier("Explosion")
	animationData, err := project.GetAnimatedSpriteByDefinition(explosionFX)
	if err != nil {
		slog.Error("explosion animation not found")
		panic(0)
	}

	animation := component.Animation{
		Animation: animationData,
		Position:  math.NewVec2(X, Y),
	}
	animations := component.AnimationsData{}
	animations[component.BrickExplosion] = animation

	component.Animations.SetValue(explosion, animations)

	return explosion
}

