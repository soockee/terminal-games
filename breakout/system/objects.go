package system

import (
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateObjects(ecs *ecs.ECS) {
	component.ConvexPolygon.Each(ecs.World, func(e *donburi.Entry) {

	})
}
