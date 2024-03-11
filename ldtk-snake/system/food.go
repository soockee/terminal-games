package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	dresolv "github.com/soockee/terminal-games/ldtk-snake/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateFood(ecs *ecs.ECS) {

}

func DrawFood(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Food.Each(ecs.World, func(e *donburi.Entry) {
		o := dresolv.GetObject(e)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(o.Position.X), float64(o.Position.Y))
		sprite := component.Sprite.Get(e)
		screen.DrawImage(sprite.Image, op)
	})
}
