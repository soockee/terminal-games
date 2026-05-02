package archetype

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp/v2"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/soockee/terminal-games/super-mario-bros/system"
	"github.com/yohamta/donburi"
)

// tileRect is a helper for the merge pass.
type tileRect struct {
	X, Y, W, H float64
}

// addMergedStaticCollision groups horizontally adjacent collidable entities
// on the same row (same Y and same height) into wider rectangles, then creates
// one static physics shape per merged run. This eliminates ghost collisions
// at tile seams.
func addMergedStaticCollision(w donburi.World, entities []*ldtkgo.Entity) {
	if len(entities) == 0 {
		return
	}

	// Collect rects.
	rects := make([]tileRect, 0, len(entities))
	for _, ent := range entities {
		tlx, tly := ent.TopLeft()
		rects = append(rects, tileRect{
			X: float64(tlx), Y: float64(tly),
			W: float64(ent.Width), H: float64(ent.Height),
		})
	}

	// Sort by Y, then X for stable row-wise merging.
	sort.Slice(rects, func(i, j int) bool {
		if rects[i].Y != rects[j].Y {
			return rects[i].Y < rects[j].Y
		}
		return rects[i].X < rects[j].X
	})

	// Merge horizontally adjacent tiles on the same row with the same height.
	var merged []tileRect
	cur := rects[0]
	for i := 1; i < len(rects); i++ {
		r := rects[i]
		if r.Y == cur.Y && r.H == cur.H && r.X == cur.X+cur.W {
			// Extend current run.
			cur.W += r.W
		} else {
			merged = append(merged, cur)
			cur = r
		}
	}
	merged = append(merged, cur)

	// Create one static shape per merged rectangle.
	for _, r := range merged {
		addStaticCollision(w, r.X, r.Y, r.W, r.H)
	}
}

// addStaticCollision creates a static physics shape from world-space coordinates.
func addStaticCollision(w donburi.World, tlx, tly, eW, eH float64) {
	spaceEntry, ok := component.PhysicsSpace.First(w)
	if !ok {
		return
	}
	space := component.PhysicsSpace.Get(spaceEntry).Space
	body := space.StaticBody

	verts := []cp.Vector{
		{X: tlx, Y: tly},
		{X: tlx + eW, Y: tly},
		{X: tlx + eW, Y: tly + eH},
		{X: tlx, Y: tly + eH},
	}
	shape := space.AddShape(cp.NewPolyShapeRaw(body, len(verts), verts, 0))
	shape.SetElasticity(0)
	shape.SetFriction(1.0)
	shape.SetCollisionType(system.CollisionTypeGround)
}

// NewBlock creates a static block entity (visual only; collision is tag-driven).
func NewBlock(w donburi.World, ent *ldtkgo.Entity) {
	entry := w.Entry(w.Create(
		component.Position,
		component.Sprite,
	))
	tlx, tly := ent.TopLeft()
	component.Position.Set(entry, &component.PositionData{
		X: float64(tlx),
		Y: float64(tly),
	})
	if img := ent.SubImage(); img != nil {
		component.Sprite.Set(entry, &component.SpriteData{
			Image: ebiten.NewImageFromImage(img),
		})
	}
}

// NewJumpBlock creates a question-mark block entity (visual only; collision is tag-driven).
func NewJumpBlock(w donburi.World, ent *ldtkgo.Entity) {
	entry := w.Entry(w.Create(
		component.Position,
		component.Sprite,
	))
	tlx, tly := ent.TopLeft()
	component.Position.Set(entry, &component.PositionData{
		X: float64(tlx),
		Y: float64(tly),
	})
	if img := ent.SubImage(); img != nil {
		component.Sprite.Set(entry, &component.SpriteData{
			Image: ebiten.NewImageFromImage(img),
		})
	}
	// TODO: read Items field and attach item component
}

// NewGround creates a static ground tile entity (visual only; collision is tag-driven).
func NewGround(w donburi.World, ent *ldtkgo.Entity) {
	entry := w.Entry(w.Create(
		component.Position,
		component.Sprite,
	))
	tlx, tly := ent.TopLeft()
	component.Position.Set(entry, &component.PositionData{
		X: float64(tlx),
		Y: float64(tly),
	})
	if img := ent.SubImage(); img != nil {
		component.Sprite.Set(entry, &component.SpriteData{
			Image: ebiten.NewImageFromImage(img),
		})
	}
}
