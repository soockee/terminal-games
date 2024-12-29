package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/ganim8/v2"
)

func UpdateExplosion(ecs *ecs.ECS) {
	tags.Explosion.Each(ecs.World, func(e *donburi.Entry) {
		animations := component.Animations.Get(e)
		if (*animations)[component.BrickExplosion].Animation.IsEnd() {
			ecs.World.Remove(e.Entity())
		}
	})
}

func DrawExplosion(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Explosion.Each(ecs.World, func(e *donburi.Entry) {
		animations := component.Animations.Get(e)
		if (*animations)[component.BrickExplosion].Animation.Status() != ganim8.Playing {
			drawOptions := ganim8.DrawOpts((*animations)[component.BrickExplosion].Position.X, (*animations)[component.BrickExplosion].Position.Y, 0, 2, 2)

			(*animations)[component.BrickExplosion].Animation.Draw(screen, drawOptions)
			(*animations)[component.BrickExplosion].Animation.Update()
		}
	})
}
