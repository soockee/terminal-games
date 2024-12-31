package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/ganim8/v2"
)

func UpdateAnimation(ecs *ecs.ECS) {

}

func DrawAnimation(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Animation.Each(ecs.World, func(e *donburi.Entry) {
		animation := component.Animation.Get(e)

		drawOptions := ganim8.DrawOpts(
			animation.Shape.Bounds().Min.X, animation.Shape.Bounds().Min.Y, 0, 2, 2)

		animation.Animation.Draw(screen, drawOptions)
		animation.Animation.Update()
		if animation.Animation.IsEnd() && !animation.Loop {
			ecs.World.Remove(e.Entity())
			space := component.Space.Get(component.Space.MustFirst(ecs.World))
			space.Remove(animation.Shape)
		}
	})
}
