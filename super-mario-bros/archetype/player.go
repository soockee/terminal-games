package archetype

import (
	"github.com/jakecoffman/cp/v2"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/super-mario-bros/assets"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/soockee/terminal-games/super-mario-bros/system"
	"github.com/yohamta/donburi"
)

// NewPlayer creates the player ECS entity from an LDtk entity instance.
func NewPlayer(w donburi.World, ent *ldtkgo.Entity) {
	entry := w.Entry(w.Create(
		component.Position,
		component.Animation,
		component.Player,
		component.Body,
	))
	tlx, tly := ent.TopLeft()
	eW, eH := float64(ent.Width), float64(ent.Height)

	component.Position.Set(entry, &component.PositionData{
		X: float64(tlx),
		Y: float64(tly),
	})
	anim := assets.PlayerAnimations()
	component.Animation.Set(entry, &anim)

	component.Player.Set(entry, &component.PlayerData{
		MoveSpeed: 120,
		JumpForce: 350,
	})

	// Physics body — dynamic, uses real Chipmunk gravity from the space.
	spaceEntry, ok := component.PhysicsSpace.First(w)
	if !ok {
		return
	}
	space := component.PhysicsSpace.Get(spaceEntry).Space
	body := space.AddBody(cp.NewBody(1, cp.INFINITY))
	body.SetPosition(cp.Vector{X: float64(tlx) + eW/2, Y: float64(tly) + eH/2})

	// Velocity update func: inject DesiredVelX at the start of Space.Step(),
	// then let Chipmunk apply gravity for Y. The solver runs after this and
	// corrects for wall penetration, so the player cannot push through walls.
	pd := component.Player.Get(entry)
	body.SetVelocityUpdateFunc(func(b *cp.Body, gravity cp.Vector, damping, dt float64) {
		b.SetVelocityVector(cp.Vector{X: pd.DesiredVelX, Y: b.Velocity().Y})
		b.UpdateVelocity(gravity, damping, dt)
	})

	// Collision shape — bevel radius rounds corners so the player glides
	// over any remaining tile seams instead of catching on edges.
	const bevel = 1.0
	shape := space.AddShape(cp.NewBox(body, eW-2*bevel, eH-2*bevel, bevel))
	shape.SetElasticity(0)
	shape.SetFriction(0)
	shape.SetCollisionType(system.CollisionTypePlayer)

	// Foot sensor — thin rectangle at the bottom of the player, slightly
	// narrower so standing at a ledge edge doesn't count as grounded.
	const sensorH = 2.0
	const sensorInset = 2.0
	footBB := cp.NewBB(
		-(eW/2 - sensorInset), // left
		eH/2-sensorH,          // bottom - sensorH (top of sensor)
		eW/2-sensorInset,      // right
		eH/2,                  // bottom
	)
	footShape := space.AddShape(cp.NewBox2(body, footBB, 0))
	footShape.SetSensor(true)
	footShape.SetCollisionType(system.CollisionTypeFoot)

	component.Body.Set(entry, &component.BodyData{
		Body:   body,
		Shapes: []*cp.Shape{shape, footShape},
		W:      eW,
		H:      eH,
	})
}
