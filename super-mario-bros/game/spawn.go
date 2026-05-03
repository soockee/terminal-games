package game

import (
	"log"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/super-mario-bros/archetype"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
)

// spawnEntities reads all entity instances from the LDtk level and creates
// corresponding ECS entities. All LDtk-format knowledge is contained here;
// archetype functions receive only game-owned types.
func spawnEntities(w donburi.World, level *ldtkgo.Level) {
	iidMap := make(map[string]*ldtkgo.Entity)
	for _, ent := range level.AllEntities() {
		iidMap[ent.IID] = ent
	}

	if p := level.Entity("Player"); p != nil {
		tlx, tly := p.TopLeft()
		archetype.NewPlayer(w, float64(tlx), float64(tly), float64(p.Width), float64(p.Height))
	}
	for _, ent := range level.EntitiesByID("Kooper") {
		tlx, tly := ent.TopLeft()
		archetype.NewKooper(w,
			float64(tlx), float64(tly), float64(ent.Width), float64(ent.Height),
			resolvePatrolPath(ent, iidMap), boolField(ent, "Loop"),
		)
	}
	for _, ent := range level.EntitiesByID("Goomba") {
		tlx, tly := ent.TopLeft()
		archetype.NewGoomba(w,
			float64(tlx), float64(tly), float64(ent.Width), float64(ent.Height),
			resolvePatrolPath(ent, iidMap), boolField(ent, "Loop"),
		)
	}
	for _, ent := range level.EntitiesByID("Block") {
		tlx, tly := ent.TopLeft()
		archetype.NewBlock(w, float64(tlx), float64(tly), toEbitenImage(ent))
	}
	for _, ent := range level.EntitiesByID("JumpBlock") {
		tlx, tly := ent.TopLeft()
		archetype.NewJumpBlock(w, float64(tlx), float64(tly), toEbitenImage(ent))
	}
	for _, ent := range level.EntitiesByID("Ground") {
		tlx, tly := ent.TopLeft()
		archetype.NewGround(w, float64(tlx), float64(tly), toEbitenImage(ent))
	}
	for _, ent := range level.EntitiesByID("Coin") {
		tlx, tly := ent.TopLeft()
		archetype.NewCoin(w, float64(tlx), float64(tly), toEbitenImage(ent))
	}

	addMergedStaticCollision(w, level.EntitiesByTag("collidable"))
}

// addMergedStaticCollision groups horizontally adjacent collidable entities on
// the same row into wider rectangles, then registers one static shape per run.
// This eliminates ghost collisions at tile seams.
func addMergedStaticCollision(w donburi.World, entities []*ldtkgo.Entity) {
	if len(entities) == 0 {
		return
	}

	type rect struct{ X, Y, W, H float64 }
	rects := make([]rect, 0, len(entities))
	for _, ent := range entities {
		tlx, tly := ent.TopLeft()
		rects = append(rects, rect{
			X: float64(tlx), Y: float64(tly),
			W: float64(ent.Width), H: float64(ent.Height),
		})
	}

	sort.Slice(rects, func(i, j int) bool {
		if rects[i].Y != rects[j].Y {
			return rects[i].Y < rects[j].Y
		}
		return rects[i].X < rects[j].X
	})

	var merged []rect
	cur := rects[0]
	for i := 1; i < len(rects); i++ {
		r := rects[i]
		if r.Y == cur.Y && r.H == cur.H && r.X == cur.X+cur.W {
			cur.W += r.W
		} else {
			merged = append(merged, cur)
			cur = r
		}
	}
	merged = append(merged, cur)

	for _, r := range merged {
		archetype.AddStaticCollision(w, r.X, r.Y, r.W, r.H)
	}
}

// resolvePatrolPath walks the To → To linked list starting from the enemy
// entity, collecting waypoint positions in top-left world coordinates.
func resolvePatrolPath(start *ldtkgo.Entity, iidMap map[string]*ldtkgo.Entity) []component.Vec2 {
	tlx, tly := start.TopLeft()
	waypoints := []component.Vec2{{X: float64(tlx), Y: float64(tly)}}

	visited := map[string]bool{start.IID: true}
	current := start

	for {
		ref := entityRefIID(current)
		if ref == "" {
			break
		}
		if visited[ref] {
			break
		}
		next, ok := iidMap[ref]
		if !ok {
			log.Printf("warning: entity ref %s not found in level", ref)
			break
		}
		visited[ref] = true
		ntlx, ntly := next.TopLeft()
		waypoints = append(waypoints, component.Vec2{X: float64(ntlx), Y: float64(ntly)})
		current = next
	}

	return waypoints
}

// entityRefIID extracts the target entity IID from a "To" EntityRef custom field.
func entityRefIID(ent *ldtkgo.Entity) string {
	raw, ok := ent.CustomFields["To"]
	if !ok || raw == nil {
		return ""
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return ""
	}
	iid, _ := m["entityIid"].(string)
	return iid
}

// boolField reads a boolean custom field from an LDtk entity, defaulting to false.
func boolField(ent *ldtkgo.Entity, key string) bool {
	if raw, ok := ent.CustomFields[key]; ok {
		if b, ok := raw.(bool); ok {
			return b
		}
	}
	return false
}

// toEbitenImage converts the LDtk entity's sub-image to an *ebiten.Image.
// Returns nil if the entity has no sub-image.
func toEbitenImage(ent *ldtkgo.Entity) *ebiten.Image {
	img := ent.SubImage()
	if img == nil {
		return nil
	}
	return ebiten.NewImageFromImage(img)
}
