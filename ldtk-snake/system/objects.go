package system

import (
	"github.com/soockee/terminal-games/ldtk-snake/component"
	resolv "github.com/soockee/terminal-games/ldtk-snake/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateObjects(ecs *ecs.ECS) {
	component.Object.Each(ecs.World, func(e *donburi.Entry) {
		obj := resolv.GetObject(e)
		obj.Update()
	})
}
