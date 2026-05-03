package physics

import "github.com/jakecoffman/cp/v2"

// Collision types — internal to the package; callers use typed factory methods.
const (
	collisionPlayer cp.CollisionType = 1
	collisionEnemy  cp.CollisionType = 2
	collisionGround cp.CollisionType = 3
	collisionFoot   cp.CollisionType = 5
)

// Body wraps a cp.Body and its shapes.
// All positions are in top-left world-space; the center-origin conversion
// required by Chipmunk is handled internally.
type Body struct {
	inner  *cp.Body
	shapes []*cp.Shape
	w, h   float64
}

// Position returns the entity's top-left world position.
func (b *Body) Position() (x, y float64) {
	p := b.inner.Position()
	return p.X - b.w/2, p.Y - b.h/2
}

// SetVelocity sets the body's velocity.
func (b *Body) SetVelocity(vx, vy float64) {
	b.inner.SetVelocityVector(cp.Vector{X: vx, Y: vy})
}

// Velocity returns the current velocity.
func (b *Body) Velocity() (vx, vy float64) {
	v := b.inner.Velocity()
	return v.X, v.Y
}

// Size returns the entity's pixel dimensions.
func (b *Body) Size() (w, h float64) {
	return b.w, b.h
}

// IsAlive reports whether the body is still active in the physics space.
// Returns false after RemoveBody has completed.
func (b *Body) IsAlive() bool {
	return b != nil && b.inner != nil
}

// Space wraps a cp.Space, hiding Chipmunk coordinate conventions and setup
// details behind a game-friendly interface.
type Space struct {
	inner             *cp.Space
	bodies            map[*cp.Body]*Body // maps cp body → our wrapper for callbacks
	footGroundHandler *cp.CollisionHandler
}

// New creates a Space with the given downward gravity (pixels/s²).
func New(gravity float64) *Space {
	inner := cp.NewSpace()
	inner.SetGravity(cp.Vector{X: 0, Y: gravity})
	inner.SetDamping(1.0)
	inner.Iterations = 20
	return &Space{
		inner:  inner,
		bodies: make(map[*cp.Body]*Body),
	}
}

func (s *Space) register(b *Body) {
	s.bodies[b.inner] = b
}

func (s *Space) lookup(cpBody *cp.Body) *Body {
	return s.bodies[cpBody]
}

// footGround returns (creating if needed) the collision handler for
// foot sensor vs ground contacts so Begin and Separate can share it.
func (s *Space) footGround() *cp.CollisionHandler {
	if s.footGroundHandler == nil {
		s.footGroundHandler = s.inner.NewCollisionHandler(collisionFoot, collisionGround)
	}
	return s.footGroundHandler
}

// AddPlayerBody creates a dynamic body for the player at top-left position
// (tlx, tly) with pixel dimensions (w, h).
//
// desiredVelX is read each physics step to inject horizontal movement intent
// before Chipmunk's solver runs. Pass a pointer to PlayerData.DesiredVelX.
//
// The body has a beveled box shape (prevents catching on tile seams) and a
// thin foot-sensor below it for grounded detection.
func (s *Space) AddPlayerBody(tlx, tly, w, h float64, desiredVelX *float64) *Body {
	body := s.inner.AddBody(cp.NewBody(1, cp.INFINITY))
	body.SetPosition(cp.Vector{X: tlx + w/2, Y: tly + h/2})
	body.SetVelocityUpdateFunc(func(b *cp.Body, gravity cp.Vector, damping, dt float64) {
		b.SetVelocityVector(cp.Vector{X: *desiredVelX, Y: b.Velocity().Y})
		b.UpdateVelocity(gravity, damping, dt)
	})

	const bevel = 1.0
	shape := s.inner.AddShape(cp.NewBox(body, w-2*bevel, h-2*bevel, bevel))
	shape.SetElasticity(0)
	shape.SetFriction(0)
	shape.SetCollisionType(collisionPlayer)

	const sensorH = 2.0
	const sensorInset = 2.0
	footBB := cp.NewBB(
		-(w/2 - sensorInset),
		h/2-sensorH,
		w/2-sensorInset,
		h/2,
	)
	footShape := s.inner.AddShape(cp.NewBox2(body, footBB, 0))
	footShape.SetSensor(true)
	footShape.SetCollisionType(collisionFoot)

	b := &Body{inner: body, shapes: []*cp.Shape{shape, footShape}, w: w, h: h}
	s.register(b)
	return b
}

// AddEnemyBody creates a dynamic body for an enemy at top-left position
// (tlx, tly) with pixel dimensions (w, h).
// Enemies share a collision group so they do not collide with each other.
func (s *Space) AddEnemyBody(tlx, tly, w, h float64) *Body {
	body := s.inner.AddBody(cp.NewBody(1, cp.INFINITY))
	body.SetPosition(cp.Vector{X: tlx + w/2, Y: tly + h/2})

	shape := s.inner.AddShape(cp.NewBox(body, w, h, 0))
	shape.SetElasticity(0)
	shape.SetFriction(0.5)
	shape.SetCollisionType(collisionEnemy)
	shape.SetFilter(cp.NewShapeFilter(1, cp.ALL_CATEGORIES, cp.ALL_CATEGORIES))

	b := &Body{inner: body, shapes: []*cp.Shape{shape}, w: w, h: h}
	s.register(b)
	return b
}

// AddStaticRect adds a non-moving ground shape at top-left world position
// (tlx, tly) with pixel dimensions (w, h).
func (s *Space) AddStaticRect(tlx, tly, w, h float64) {
	body := s.inner.StaticBody
	verts := []cp.Vector{
		{X: tlx, Y: tly},
		{X: tlx + w, Y: tly},
		{X: tlx + w, Y: tly + h},
		{X: tlx, Y: tly + h},
	}
	shape := s.inner.AddShape(cp.NewPolyShapeRaw(body, len(verts), verts, 0))
	shape.SetElasticity(0)
	shape.SetFriction(1.0)
	shape.SetCollisionType(collisionGround)
}

// Step advances the simulation by one fixed timestep (1/60 s).
func (s *Space) Step() {
	s.inner.Step(1.0 / 60.0)
}

// OnFootGroundBegin registers a callback fired when the player's foot sensor
// first touches a ground shape.
func (s *Space) OnFootGroundBegin(fn func(player *Body)) {
	h := s.footGround()
	h.BeginFunc = func(arb *cp.Arbiter, _ *cp.Space, _ any) bool {
		cpBodyA, _ := arb.Bodies()
		if b := s.lookup(cpBodyA); b != nil {
			fn(b)
		}
		return true // sensors must return true
	}
}

// OnFootGroundSeparate registers a callback fired when the foot sensor leaves
// a ground shape.
func (s *Space) OnFootGroundSeparate(fn func(player *Body)) {
	h := s.footGround()
	h.SeparateFunc = func(arb *cp.Arbiter, _ *cp.Space, _ any) {
		cpBodyA, _ := arb.Bodies()
		if b := s.lookup(cpBodyA); b != nil {
			fn(b)
		}
	}
}

// OnPlayerEnemyContact registers a callback fired when the player body first
// contacts an enemy body. normalY is the Y component of the collision normal
// (positive means the player is above the enemy — a stomp).
// Return false to suppress the physics contact impulse (usual for gameplay hits).
func (s *Space) OnPlayerEnemyContact(fn func(player, enemy *Body, normalY float64) bool) {
	h := s.inner.NewCollisionHandler(collisionPlayer, collisionEnemy)
	h.BeginFunc = func(arb *cp.Arbiter, _ *cp.Space, _ any) bool {
		cpBodyA, cpBodyB := arb.Bodies()
		playerBody := s.lookup(cpBodyA)
		enemyBody := s.lookup(cpBodyB)
		if playerBody == nil || enemyBody == nil {
			return true
		}
		return fn(playerBody, enemyBody, arb.Normal().Y)
	}
}

// RemoveBody safely defers removal of a body and its shapes via a post-step
// callback (Chipmunk forbids modifying the space during a collision callback).
// onRemoved is called after removal completes; use it to nil out ECS pointers.
func (s *Space) RemoveBody(b *Body, onRemoved func()) {
	if !b.IsAlive() {
		return
	}
	delete(s.bodies, b.inner)
	s.inner.AddPostStepCallback(func(space *cp.Space, _ any, _ any) {
		for _, sh := range b.shapes {
			space.RemoveShape(sh)
		}
		space.RemoveBody(b.inner)
		b.inner = nil
		b.shapes = nil
		if onRemoved != nil {
			onRemoved()
		}
	}, b.inner, nil)
}
