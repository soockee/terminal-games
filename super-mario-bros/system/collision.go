package system

import (
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/soockee/terminal-games/super-mario-bros/physics"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
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

	space.OnFootGroundBegin(makeFootGroundBeginHandler(w))
	space.OnFootGroundSeparate(makeFootGroundSeparateHandler(w))
	space.OnPlayerEnemyContact(makePlayerEnemyContactHandler(w, space))
}

func makeFootGroundBeginHandler(w donburi.World) func(*physics.Body) {
	return func(playerPhysBody *physics.Body) {
		collisionPlayerQuery.Each(w, func(entry *donburi.Entry) {
			if component.Body.Get(entry).Body == playerPhysBody {
				component.Player.Get(entry).GroundContacts++
			}
		})
	}
}

func makeFootGroundSeparateHandler(w donburi.World) func(*physics.Body) {
	return func(playerPhysBody *physics.Body) {
		collisionPlayerQuery.Each(w, func(entry *donburi.Entry) {
			if component.Body.Get(entry).Body == playerPhysBody {
				pd := component.Player.Get(entry)
				pd.GroundContacts--
				if pd.GroundContacts < 0 {
					pd.GroundContacts = 0
				}
			}
		})
	}
}

func makePlayerEnemyContactHandler(w donburi.World, space *physics.Space) func(*physics.Body, *physics.Body, float64) bool {
	return func(playerPhysBody, enemyPhysBody *physics.Body, normalY float64) bool {
		var playerEntry *donburi.Entry
		collisionPlayerQuery.Each(w, func(e *donburi.Entry) {
			if component.Body.Get(e).Body == playerPhysBody {
				playerEntry = e
			}
		})
		var enemyEntry *donburi.Entry
		collisionEnemyQuery.Each(w, func(e *donburi.Entry) {
			if component.Body.Get(e).Body == enemyPhysBody {
				enemyEntry = e
			}
		})
		if playerEntry == nil || enemyEntry == nil {
			return true
		}

		enemyData := component.Enemy.Get(enemyEntry)
		kind := ClassifyContact(normalY, enemyData.State)
		switch kind {
		case ContactIgnore:
			return false
		case ContactStomp:
			result := ClassifyStomp(enemyData.Type, enemyData.State)
			applyStompResult(w, space, playerEntry, enemyEntry, result)
		case ContactSideHit:
			if gsEntry, ok := component.GameState.First(w); ok {
				component.GameState.Get(gsEntry).Die()
			}
		}
		return false
	}
}

// ---- Pure contact classification (no ECS, no physics world needed) ----

// ContactKind classifies a player/enemy collision outcome.
type ContactKind int

const (
	ContactIgnore  ContactKind = iota // enemy already dead, no effect
	ContactStomp                      // player lands on enemy from above
	ContactSideHit                    // side contact with a live enemy — player dies
)

// ClassifyContact returns the ContactKind for a player/enemy collision.
// normalY is the Y component of the collision normal; positive means the
// player is above the enemy (a stomp).
func ClassifyContact(normalY float64, state component.EnemyState) ContactKind {
	if state == component.EnemyDead {
		return ContactIgnore
	}
	if normalY > 0 || state == component.EnemyShell {
		// Shell side-hit is treated as a second stomp.
		return ContactStomp
	}
	return ContactSideHit
}

// StompResult describes what should happen to an enemy after a successful stomp.
type StompResult struct {
	NextState    component.EnemyState
	EnemyAnim    string // animation name to play ("squashed", "shell", "")
	RemoveBody   bool
	StopPatrol   bool
	ZeroVelocity bool    // zero horizontal velocity (Kooper Alive→Shell)
	CleanupTimer float64 // seconds before removal (0 = none)
	Score        int
}

const cleanupDelay = 0.5

// ClassifyStomp returns the StompResult for the given enemy type and state.
func ClassifyStomp(enemyType component.EnemyType, state component.EnemyState) StompResult {
	const points = 100
	switch enemyType {
	case component.EnemyGoomba:
		return StompResult{
			NextState:    component.EnemyDead,
			EnemyAnim:    "squashed",
			RemoveBody:   true,
			StopPatrol:   true,
			CleanupTimer: cleanupDelay,
			Score:        points,
		}
	case component.EnemyKooper:
		switch state {
		case component.EnemyAlive:
			return StompResult{
				NextState:    component.EnemyShell,
				EnemyAnim:    "shell",
				StopPatrol:   true,
				ZeroVelocity: true,
				Score:        points,
			}
		case component.EnemyShell:
			return StompResult{
				NextState:    component.EnemyDead,
				RemoveBody:   true,
				CleanupTimer: cleanupDelay,
				Score:        points,
			}
		}
	}
	return StompResult{}
}

// ---- Application layer (ECS + physics mutations) ----

func applyStompResult(w donburi.World, space *physics.Space, playerEntry, enemyEntry *donburi.Entry, result StompResult) {
	// Bounce the player upward (60% of normal jump force) and play stomp anim.
	pd := component.Player.Get(playerEntry)
	playerBody := component.Body.Get(playerEntry).Body
	vx, _ := playerBody.Velocity()
	playerBody.SetVelocity(vx, -pd.JumpForce*0.6)
	component.Animation.Get(playerEntry).Play("stomp")

	if result.Score > 0 {
		if scoreEntry, ok := component.Score.First(w); ok {
			component.Score.Get(scoreEntry).Value += result.Score
		}
	}

	enemyData := component.Enemy.Get(enemyEntry)
	enemyData.State = result.NextState
	if result.CleanupTimer > 0 {
		enemyData.CleanupTimer = result.CleanupTimer
	}
	if result.EnemyAnim != "" {
		component.Animation.Get(enemyEntry).Play(result.EnemyAnim)
	}
	if result.StopPatrol {
		stopPatrol(enemyEntry)
	}
	if result.ZeroVelocity {
		if enemyBody := component.Body.Get(enemyEntry).Body; enemyBody.IsAlive() {
			_, vy := enemyBody.Velocity()
			enemyBody.SetVelocity(0, vy)
		}
	}
	if result.RemoveBody {
		enemyBd := component.Body.Get(enemyEntry)
		space.RemoveBody(enemyBd.Body, func() {
			enemyBd.Body = nil
		})
	}
}

// stopPatrol zeroes patrol speed so the AI system no longer moves the enemy.
func stopPatrol(enemyEntry *donburi.Entry) {
	if enemyEntry.HasComponent(component.Patrol) {
		component.Patrol.Get(enemyEntry).Speed = 0
	}
}
