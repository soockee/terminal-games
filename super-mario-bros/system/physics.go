package system

import (
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var physicsBodyQuery = donburi.NewQuery(
	filter.Contains(component.Body, component.Position),
)

// UpdatePhysics steps the physics space and syncs body positions back to the
// ECS Position component. Call once per Update tick.
func UpdatePhysics(e *ecs.ECS) {
	entry, ok := component.PhysicsSpace.First(e.World)
	if !ok {
		return
	}
	component.PhysicsSpace.Get(entry).Space.Step()

	physicsBodyQuery.Each(e.World, func(entry *donburi.Entry) {
		bd := component.Body.Get(entry)
		if !bd.Body.IsAlive() {
			return
		}
		pos := component.Position.Get(entry)
		pos.X, pos.Y = bd.Body.Position()
	})
}
