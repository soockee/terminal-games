package system

import (
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

// physicsBodyQuery matches entities that have both a physics Body and a Position.
var physicsBodyQuery = donburi.NewQuery(
	filter.Contains(component.Body, component.Position),
)

// UpdatePhysics steps the cp.Space and syncs body positions back to the
// ECS Position component. Call once per Update tick.
//
// cp body position is the center of the entity.
// ECS Position is top-left. The conversion uses BodyData.W and .H.
func UpdatePhysics(e *ecs.ECS) {
	entry, ok := component.PhysicsSpace.First(e.World)
	if !ok {
		return
	}
	space := component.PhysicsSpace.Get(entry).Space

	// Step at a fixed 1/60s timestep (matches Ebitengine default TPS).
	space.Step(1.0 / 60.0)

	// Sync cp body positions → ECS Position component (top-left).
	physicsBodyQuery.Each(e.World, func(entry *donburi.Entry) {
		bd := component.Body.Get(entry)
		if bd.Body == nil {
			return
		}
		pos := component.Position.Get(entry)
		p := bd.Body.Position()
		// cp position = center → ECS top-left
		pos.X = p.X - bd.W/2
		pos.Y = p.Y - bd.H/2
	})
}
