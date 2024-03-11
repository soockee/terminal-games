package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	dresolv "github.com/soockee/terminal-games/ldtk-snake/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func DrawWall(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Wall.Each(ecs.World, func(e *donburi.Entry) {
		o := dresolv.GetObject(e)
		sprite := component.Sprite.Get(e)

		dx := 0.0
		for dx < o.Size.X {
			dy := 0.0
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(o.Position.X)+dx, float64(o.Position.Y)+dy)
			screen.DrawImage(sprite.Image, op)
			for dy < o.Size.Y {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(o.Position.X)+dx, float64(o.Position.Y)+dy)
				screen.DrawImage(sprite.Image, op)
				dy += float64(sprite.Image.Bounds().Dx())
			}
			dx += float64(sprite.Image.Bounds().Dx())
		}
	})
}
