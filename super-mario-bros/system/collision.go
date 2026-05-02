package system

import (
	"github.com/jakecoffman/cp/v2"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

// Collision types for cp collision handlers.
const (
	CollisionTypePlayer cp.CollisionType = 1
	CollisionTypeEnemy  cp.CollisionType = 2
	CollisionTypeGround cp.CollisionType = 3
	CollisionTypeItem   cp.CollisionType = 4
	CollisionTypeFoot   cp.CollisionType = 5
)

// Queries used by collision callbacks to find ECS entries from physics bodies.
var (
	collisionPlayerQuery = donburi.NewQuery(
		filter.Contains(component.Player, component.Body, component.Animation),
	)
	collisionEnemyQuery = donburi.NewQuery(
		filter.Contains(component.Enemy, component.Body, component.Animation),
	)
)

// SetupCollisionHandlers registers collision callbacks on the physics space.
// Call once after the space is created.
func SetupCollisionHandlers(w donburi.World) {
	entry, ok := component.PhysicsSpace.First(w)
	if !ok {
		return
	}
	space := component.PhysicsSpace.Get(entry).Space

	// Foot sensor vs Ground — increment/decrement ground contact counter.
	footGround := space.NewCollisionHandler(CollisionTypeFoot, CollisionTypeGround)
	footGround.BeginFunc = makeFootGroundBeginHandler(w)
	footGround.SeparateFunc = makeFootGroundSeparateHandler(w)

	// Player vs Enemy — closure captures the world for ECS lookups.
	playerEnemy := space.NewCollisionHandler(CollisionTypePlayer, CollisionTypeEnemy)
	playerEnemy.BeginFunc = makePlayerEnemyBeginHandler(w)
}

// makeFootGroundBeginHandler returns a callback that increments the player's
// ground contact counter when the foot sensor touches a ground shape.
func makeFootGroundBeginHandler(w donburi.World) func(*cp.Arbiter, *cp.Space, any) bool {
	return func(arb *cp.Arbiter, space *cp.Space, _ any) bool {
		bodyA, _ := arb.Bodies()
		collisionPlayerQuery.Each(w, func(entry *donburi.Entry) {
			if component.Body.Get(entry).Body == bodyA {
				component.Player.Get(entry).GroundContacts++
			}
		})
		return true // sensors always return true
	}
}

// makeFootGroundSeparateHandler returns a callback that decrements the player's
// ground contact counter when the foot sensor leaves a ground shape.
func makeFootGroundSeparateHandler(w donburi.World) func(*cp.Arbiter, *cp.Space, any) {
	return func(arb *cp.Arbiter, space *cp.Space, _ any) {
		bodyA, _ := arb.Bodies()
		collisionPlayerQuery.Each(w, func(entry *donburi.Entry) {
			if component.Body.Get(entry).Body == bodyA {
				pd := component.Player.Get(entry)
				pd.GroundContacts--
				if pd.GroundContacts < 0 {
					pd.GroundContacts = 0
				}
			}
		})
	}
}

// makePlayerEnemyBeginHandler returns a collision begin callback that resolves
// player/enemy ECS entries from the physics bodies and handles stomp vs side hit.
func makePlayerEnemyBeginHandler(w donburi.World) func(*cp.Arbiter, *cp.Space, any) bool {
	return func(arb *cp.Arbiter, space *cp.Space, _ any) bool {
		// A = Player (TypeA), B = Enemy (TypeB).
		bodyA, bodyB := arb.Bodies()

		// Resolve ECS entries by matching body pointers.
		var playerEntry *donburi.Entry
		collisionPlayerQuery.Each(w, func(entry *donburi.Entry) {
			if component.Body.Get(entry).Body == bodyA {
				playerEntry = entry
			}
		})
		var enemyEntry *donburi.Entry
		collisionEnemyQuery.Each(w, func(entry *donburi.Entry) {
			if component.Body.Get(entry).Body == bodyB {
				enemyEntry = entry
			}
		})
		if playerEntry == nil || enemyEntry == nil {
			return true
		}

		enemyData := component.Enemy.Get(enemyEntry)

		// Ignore enemies that are already dead (cleanup pending).
		if enemyData.State == component.EnemyDead {
			return false
		}

		n := arb.Normal()
		if n.Y > 0 {
			// Stomp — player is above the enemy.
			handleStomp(w, space, playerEntry, enemyEntry)
		} else {
			// Side or below hit.
			// Kooper shells are harmless when idle — only side-hits from
			// alive enemies kill the player.
			if enemyData.State == component.EnemyShell {
				// Treat a side hit on a shell as a second stomp → kill.
				handleStomp(w, space, playerEntry, enemyEntry)
			} else {
				if gsEntry, ok := component.GameState.First(w); ok {
					component.GameState.Get(gsEntry).Dead = true
				}
			}
		}
		// Reject the contact so the solver doesn't apply normal impulses.
		return false
	}
}

// cleanupDelay is how long (seconds) a dead enemy stays visible before removal.
const cleanupDelay = 0.5

// handleStomp advances the enemy state machine and bounces the player.
//
// Goomba:  Alive → Dead  (squashed anim, remove body, start cleanup timer)
// Kooper:  Alive → Shell (shell anim, stop patrol, keep body)
//
//	Shell → Dead  (remove body, start cleanup timer)
func handleStomp(w donburi.World, space *cp.Space, playerEntry, enemyEntry *donburi.Entry) {
	enemyData := component.Enemy.Get(enemyEntry)
	enemyAnim := component.Animation.Get(enemyEntry)

	// Bounce the player upward (60 % of normal jump force).
	pd := component.Player.Get(playerEntry)
	playerBody := component.Body.Get(playerEntry).Body
	playerBody.SetVelocityVector(cp.Vector{
		X: playerBody.Velocity().X,
		Y: -pd.JumpForce * 0.6,
	})
	component.Animation.Get(playerEntry).Play("stomp")

	// Increment score.
	if scoreEntry, ok := component.Score.First(w); ok {
		component.Score.Get(scoreEntry).Value += 100
	}

	switch enemyData.Type {
	case component.EnemyGoomba:
		// One-hit kill: Alive → Dead.
		enemyAnim.Play("squashed")
		enemyData.State = component.EnemyDead
		enemyData.CleanupTimer = cleanupDelay
		removeEnemyBody(space, enemyEntry)
		stopPatrol(enemyEntry)

	case component.EnemyKooper:
		switch enemyData.State {
		case component.EnemyAlive:
			// First stomp: Alive → Shell.
			enemyAnim.Play("shell")
			enemyData.State = component.EnemyShell
			stopPatrol(enemyEntry)
			// Stop horizontal movement but keep body for next collision.
			enemyBd := component.Body.Get(enemyEntry)
			if enemyBd.Body != nil {
				enemyBd.Body.SetVelocityVector(cp.Vector{X: 0, Y: enemyBd.Body.Velocity().Y})
			}
		case component.EnemyShell:
			// Second stomp: Shell → Dead.
			enemyData.State = component.EnemyDead
			enemyData.CleanupTimer = cleanupDelay
			removeEnemyBody(space, enemyEntry)
		}
	}
}

// removeEnemyBody safely removes the enemy's physics body and shapes via a
// post-step callback (cannot modify the space during a collision callback).
func removeEnemyBody(space *cp.Space, enemyEntry *donburi.Entry) {
	enemyBd := component.Body.Get(enemyEntry)
	if enemyBd.Body == nil {
		return
	}
	space.AddPostStepCallback(func(s *cp.Space, _ any, _ any) {
		for _, sh := range enemyBd.Shapes {
			s.RemoveShape(sh)
		}
		s.RemoveBody(enemyBd.Body)
		enemyBd.Body = nil
		enemyBd.Shapes = nil
	}, enemyBd.Body, nil)
}

// stopPatrol zeroes patrol speed so the AI system no longer moves the enemy.
func stopPatrol(enemyEntry *donburi.Entry) {
	if enemyEntry.HasComponent(component.Patrol) {
		component.Patrol.Get(enemyEntry).Speed = 0
	}
}
