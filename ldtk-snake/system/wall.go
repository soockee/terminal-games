package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func DrawWall(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Wall.Each(ecs.World, func(e *donburi.Entry) {
		component.DrawRepeatedSprite(screen, e)
	})
}
