package system

import (
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var playerMovementQuery = donburi.NewQuery(
	filter.Contains(component.Player, component.Body, component.Animation),
)

// UpdateMovement handles player horizontal velocity and jump impulse.
//
// Horizontal: sets DesiredVelX on PlayerData; the player body's velocity update
// func injects it at the start of Space.Step(), so the solver can correct for
// walls afterwards.
// Vertical:   Chipmunk gravity handles falling; jump overrides Y velocity.
// Ground detection comes from the foot sensor contact counter.
func UpdateMovement(e *ecs.ECS) {
	gsEntry, ok := component.GameState.First(e.World)
	if !ok {
		return
	}
	gs := component.GameState.Get(gsEntry)
	if !gs.IsActive() {
		return
	}

	playerMovementQuery.Each(e.World, func(entry *donburi.Entry) {
		pd := component.Player.Get(entry)
		bd := component.Body.Get(entry)
		anim := component.Animation.Get(entry)
		body := bd.Body

		grounded := pd.GroundContacts > 0

		// Store desired horizontal velocity for the velocity update callback.
		pd.DesiredVelX = pd.MoveDir * pd.MoveSpeed

		// Jump: override Y velocity when grounded.
		if pd.JumpInput && grounded {
			vx, _ := body.Velocity()
			body.SetVelocity(vx, -pd.JumpForce)
		}

		// Animation state.
		_, velY := body.Velocity()
		if pd.MoveDir != 0 {
			anim.FlipH = pd.MoveDir < 0
		}
		if !grounded {
			if velY < 0 {
				anim.Play("jump")
			} else {
				anim.Play("fall")
			}
		} else if pd.MoveDir != 0 {
			anim.Play("walk")
		} else if anim.FlipH || anim.Current == "walk" || anim.Current == "idle_dir" {
			anim.Play("idle_dir")
		} else {
			anim.Play("idle")
		}
	})
}
