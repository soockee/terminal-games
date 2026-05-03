package archetype

import (
	"log"

	"github.com/soockee/terminal-games/super-mario-bros/assets"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
)

// NewKooper creates a Kooper enemy entity.
func NewKooper(w donburi.World, x, y, width, height float64, waypoints []component.Vec2, loop bool) {
	spawnEnemy(w, x, y, width, height, waypoints, loop, component.EnemyKooper)
}

// NewGoomba creates a Goomba enemy entity.
func NewGoomba(w donburi.World, x, y, width, height float64, waypoints []component.Vec2, loop bool) {
	spawnEnemy(w, x, y, width, height, waypoints, loop, component.EnemyGoomba)
}

func spawnEnemy(w donburi.World, x, y, width, height float64, waypoints []component.Vec2, loop bool, enemyType component.EnemyType) {
	components := []donburi.IComponentType{
		component.Position,
		component.Enemy,
		component.Animation,
		component.Body,
	}
	if len(waypoints) > 1 {
		components = append(components, component.Patrol)
	}

	entry := w.Entry(w.Create(components...))

	component.Position.Set(entry, &component.PositionData{X: x, Y: y})
	component.Enemy.Set(entry, &component.EnemyData{Type: enemyType})
	anim := assets.EnemyAnimations(enemyType)
	component.Animation.Set(entry, &anim)

	spaceEntry, ok := component.PhysicsSpace.First(w)
	if ok {
		space := component.PhysicsSpace.Get(spaceEntry).Space
		body := space.AddEnemyBody(x, y, width, height)
		component.Body.Set(entry, &component.BodyData{Body: body})
	}

	if len(waypoints) > 1 {
		component.Patrol.Set(entry, &component.PatrolData{
			Waypoints: waypoints,
			Current:   1,
			Speed:     30,
			Forward:   true,
			Loop:      loop,
		})
	}

	log.Printf("spawned %s at (%.0f,%.0f) with %d waypoints (loop=%v)", enemyType, x, y, len(waypoints), loop)
}
