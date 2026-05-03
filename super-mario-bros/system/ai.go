package system

import (
	"math"

	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var patrolQuery = donburi.NewQuery(
	filter.Contains(component.Patrol, component.Position, component.Body),
)

// UpdatePatrol sets the horizontal velocity on patrolling entities' physics
// bodies. Physics handles gravity and ground collision for the Y axis.
func UpdatePatrol(e *ecs.ECS) {
	patrolQuery.Each(e.World, func(entry *donburi.Entry) {
		patrol := component.Patrol.Get(entry)
		bd := component.Body.Get(entry)
		if !bd.Body.IsAlive() {
			return
		}
		pos := component.Position.Get(entry)

		if len(patrol.Waypoints) < 2 {
			_, vy := bd.Body.Velocity()
			bd.Body.SetVelocity(0, vy)
			return
		}

		target := patrol.Waypoints[patrol.Current]
		w, _ := bd.Body.Size()
		entityCenterX := pos.X + w/2
		targetCenterX := target.X + w/2
		dx := targetCenterX - entityCenterX

		if math.Abs(dx) <= 1.0 {
			_, vy := bd.Body.Velocity()
			bd.Body.SetVelocity(0, vy)
			advance(patrol)
		} else {
			dir := 1.0
			if dx < 0 {
				dir = -1.0
			}
			_, vy := bd.Body.Velocity()
			bd.Body.SetVelocity(dir*patrol.Speed, vy)

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
				p.Current = 1
				if p.Current >= len(p.Waypoints) {
					p.Current = 0
				}
				p.Forward = true
			}
		}
	}
}
