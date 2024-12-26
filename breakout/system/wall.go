package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func DrawWall(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Wall.Each(ecs.World, func(e *donburi.Entry) {
		w := component.Wall.Get(e)
		component.DrawRepeatedSprite(screen, component.Sprite.Get(e).Images[0], w.Shape)
	})
}
