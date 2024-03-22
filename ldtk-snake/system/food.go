package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateFood(ecs *ecs.ECS) {

}

func DrawFood(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Food.Each(ecs.World, func(e *donburi.Entry) {
		component.DrawScaledSprite(screen, component.Sprite.Get(e).Images[0], e)
	})
}
