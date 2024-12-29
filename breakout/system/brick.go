package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateBrick(ecs *ecs.ECS) {
	component.Brick.Each(ecs.World, func(e *donburi.Entry) {
	})
}

func DrawBrick(ecs *ecs.ECS, screen *ebiten.Image) {

	tags.Brick.Each(ecs.World, func(e *donburi.Entry) {
		b := component.Collidable.Get(e)
		component.DrawScaledSprite(screen, component.Sprite.Get(e).Images[0], b.Shape)
	})
}
