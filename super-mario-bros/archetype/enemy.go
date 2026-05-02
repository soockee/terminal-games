package archetype

import (
	"log"

	"github.com/jakecoffman/cp/v2"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/super-mario-bros/assets"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/soockee/terminal-games/super-mario-bros/system"
	"github.com/yohamta/donburi"
)

// NewKooper creates a Kooper enemy entity from its LDtk instance.
func NewKooper(w donburi.World, ent *ldtkgo.Entity, iidMap map[string]*ldtkgo.Entity) {
	spawnEnemy(w, ent, iidMap, component.EnemyKooper)
}

// NewGoomba creates a Goomba enemy entity from its LDtk instance.
func NewGoomba(w donburi.World, ent *ldtkgo.Entity, iidMap map[string]*ldtkgo.Entity) {
	spawnEnemy(w, ent, iidMap, component.EnemyGoomba)
}

// spawnEnemy is the shared logic for creating any enemy type.
func spawnEnemy(w donburi.World, ent *ldtkgo.Entity, iidMap map[string]*ldtkgo.Entity, enemyType component.EnemyType) {
	loopPath := false
	if raw, ok := ent.CustomFields["Loop"]; ok {
		if b, ok := raw.(bool); ok {
			loopPath = b
		}
	}

	waypoints := resolvePatrolPath(ent, iidMap)

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

	tlx, tly := ent.TopLeft()
	eW, eH := float64(ent.Width), float64(ent.Height)

	component.Position.Set(entry, &component.PositionData{
		X: float64(tlx),
		Y: float64(tly),
	})
	component.Enemy.Set(entry, &component.EnemyData{
		Type: enemyType,
	})
	anim := assets.EnemyAnimations(enemyType)
	component.Animation.Set(entry, &anim)

	// Physics body -- dynamic, no Chipmunk gravity (game systems handle it).
	// cp position = center of entity.
	spaceEntry, ok := component.PhysicsSpace.First(w)
	if ok {
		space := component.PhysicsSpace.Get(spaceEntry).Space
		body := space.AddBody(cp.NewBody(1, cp.INFINITY)) // mass=1, no rotation
		body.SetPosition(cp.Vector{X: float64(tlx) + eW/2, Y: float64(tly) + eH/2})

		shape := space.AddShape(cp.NewBox(body, eW, eH, 0))
		shape.SetElasticity(0)
		shape.SetFriction(0.5)
		shape.SetCollisionType(system.CollisionTypeEnemy)
		// Same group = enemies ignore each other.
		shape.SetFilter(cp.NewShapeFilter(1, cp.ALL_CATEGORIES, cp.ALL_CATEGORIES))

		component.Body.Set(entry, &component.BodyData{
			Body:   body,
			Shapes: []*cp.Shape{shape},
			W:      eW,
			H:      eH,
		})
	}

	if len(waypoints) > 1 {
		component.Patrol.Set(entry, &component.PatrolData{
			Waypoints: waypoints,
			Current:   1,
			Speed:     30,
			Forward:   true,
			Loop:      loopPath,
		})
	}

	log.Printf("spawned %s at (%d,%d) with %d waypoints (loop=%v)",
		enemyType, tlx, tly, len(waypoints), loopPath)
}

// resolvePatrolPath walks the To → To linked list starting from the enemy
// entity, collecting waypoint positions. All positions use TopLeft coordinates.
func resolvePatrolPath(start *ldtkgo.Entity, iidMap map[string]*ldtkgo.Entity) []component.Vec2 {
	tlx, tly := start.TopLeft()
	waypoints := []component.Vec2{
		{X: float64(tlx), Y: float64(tly)},
	}

	visited := map[string]bool{start.IID: true}
	current := start

	for {
		ref := entityRefIID(current)
		if ref == "" {
			break
		}
		if visited[ref] {
			break // prevent infinite loops
		}

		next, ok := iidMap[ref]
		if !ok {
			log.Printf("warning: entity ref %s not found in level", ref)
			break
		}

		visited[ref] = true
		ntlx, ntly := next.TopLeft()
		waypoints = append(waypoints, component.Vec2{
			X: float64(ntlx),
			Y: float64(ntly),
		})
		current = next
	}

	return waypoints
}

// entityRefIID extracts the target entity IID from a "To" EntityRef custom field.
// Returns "" if the field is absent or null.
func entityRefIID(ent *ldtkgo.Entity) string {
	raw, ok := ent.CustomFields["To"]
	if !ok || raw == nil {
		return ""
	}
	// In the simplified export, EntityRef is an object with "entityIid".
	m, ok := raw.(map[string]interface{})
	if !ok {
		return ""
	}
	iid, _ := m["entityIid"].(string)
	return iid
}
