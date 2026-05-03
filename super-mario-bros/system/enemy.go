package system

import (
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

// enemyCleanupQuery matches all enemy entities.
var enemyCleanupQuery = donburi.NewQuery(
	filter.Contains(component.Enemy, component.Position),
)

// UpdateEnemyCleanup ticks cleanup timers on dead enemies and removes their
// ECS entities once the timer expires.
func UpdateEnemyCleanup(e *ecs.ECS) {
	const dt = 1.0 / 60.0 // fixed timestep matching physics

	var toRemove []*donburi.Entry
	enemyCleanupQuery.Each(e.World, func(entry *donburi.Entry) {
		ed := component.Enemy.Get(entry)
		if ed.State != component.EnemyDead || ed.CleanupTimer <= 0 {
			return
		}
		ed.CleanupTimer -= dt
		if ed.CleanupTimer <= 0 {
			toRemove = append(toRemove, entry)
		}
	})

	for _, entry := range toRemove {
		entry.Remove()
	}
}
