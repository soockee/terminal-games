package systems

import (
	"github.com/soockee/terminal-games/ldtk-snake/components"
	dresolv "github.com/soockee/terminal-games/ldtk-snake/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateObjects(ecs *ecs.ECS) {
	components.Object.Each(ecs.World, func(e *donburi.Entry) {
		obj := dresolv.GetObject(e)
		obj.Update()
	})
}
