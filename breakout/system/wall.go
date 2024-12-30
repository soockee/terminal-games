package system

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func DrawWall(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Wall.Each(ecs.World, func(e *donburi.Entry) {
		w := component.Collidable.Get(e)
		spriteData := component.Sprite.Get(e)

		sprite := spriteData.Images[0]
		component.DrawPlaceholder(screen, sprite, w.Shape, 0, color.RGBA{255, 255, 255, 0}, true)
	})
}
