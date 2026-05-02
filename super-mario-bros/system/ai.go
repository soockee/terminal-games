package system

import (
	"math"

	"github.com/jakecoffman/cp/v2"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

// patrolQuery matches entities with Patrol, Position, and Body.
var patrolQuery = donburi.NewQuery(
	filter.Contains(component.Patrol, component.Position, component.Body),
)

// UpdatePatrol sets the horizontal velocity on patrolling entities' physics
// bodies. Physics handles gravity and ground collision for the Y axis.
// Waypoint arrival is checked by X distance only.
func UpdatePatrol(e *ecs.ECS) {
	patrolQuery.Each(e.World, func(entry *donburi.Entry) {
		patrol := component.Patrol.Get(entry)
		bd := component.Body.Get(entry)
		if bd.Body == nil {
			return
		}
		pos := component.Position.Get(entry)

		if len(patrol.Waypoints) < 2 {
			bd.Body.SetVelocityVector(cp.Vector{X: 0, Y: bd.Body.Velocity().Y})
			return
		}

		target := patrol.Waypoints[patrol.Current]
		// Compare entity center-X with waypoint center-X.
		// Waypoints store top-left positions, so their center is at X + W/2.
		entityCenterX := pos.X + bd.W/2
		targetCenterX := target.X + bd.W/2
		dx := targetCenterX - entityCenterX

		if math.Abs(dx) <= 1.0 {
			// Arrived at waypoint X — stop horizontal, advance.
			bd.Body.SetVelocityVector(cp.Vector{X: 0, Y: bd.Body.Velocity().Y})
			advance(patrol)
		} else {
			// Move horizontally toward target. Preserve vertical velocity (gravity).
			dir := 1.0
			if dx < 0 {
				dir = -1.0
			}
			bd.Body.SetVelocityVector(cp.Vector{
				X: dir * patrol.Speed,
				Y: bd.Body.Velocity().Y,
			})

			// Flip sprite to face movement direction.
			// Default sprite faces left; flip when moving right.
			if entry.HasComponent(component.Animation) {
				anim := component.Animation.Get(entry)
				anim.FlipH = dx > 0
			}
		}
	})
}

// advance moves to the next waypoint, handling ping-pong and loop modes.
func advance(p *component.PatrolData) {
	if p.Forward {
		p.Current++
		if p.Current >= len(p.Waypoints) {
			if p.Loop {
				p.Current = 0
			} else {
				// Ping-pong: reverse direction.
				p.Current = len(p.Waypoints) - 2
				if p.Current < 0 {
					p.Current = 0
				}
				p.Forward = false
			}
		}
	} else {
		p.Current--
		if p.Current < 0 {
			if p.Loop {
				p.Current = len(p.Waypoints) - 1
			} else {
				// Ping-pong: reverse direction.
				p.Current = 1
				if p.Current >= len(p.Waypoints) {
					p.Current = 0
				}
				p.Forward = true
			}
		}
	}
}
