package system

import (
	"github.com/soockee/terminal-games/match-3/component"
	"github.com/soockee/terminal-games/match-3/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var tileQuery = donburi.NewQuery(filter.Contains(tags.Tile))

// UpdateTween advances all active tweens by 1/60s and updates PixelPos.
func UpdateTween(e *ecs.ECS) {
	const dt = 1.0 / 60.0

	tileQuery.Each(e.World, func(entry *donburi.Entry) {
		tw := component.Tween.Get(entry)
		if !tw.Active {
			return
		}

		tw.Elapsed += dt
		pos := component.PixelPos.Get(entry)

		if tw.Elapsed >= tw.Duration {
			// Snap directly to end position to avoid floating-point drift
			// from the interpolation formula (a + (b-a)*1.0 != b in IEEE 754).
			tw.Elapsed = tw.Duration
			tw.Active = false
			pos.X = tw.EndX
			pos.Y = tw.EndY
			return
		}

		t := tw.Elapsed / tw.Duration
		t = applyEase(tw.Ease, t)

		pos.X = tw.StartX + (tw.EndX-tw.StartX)*t
		pos.Y = tw.StartY + (tw.EndY-tw.StartY)*t
	})
}

func applyEase(ease component.EaseFunc, t float64) float64 {
	switch ease {
	case component.EaseOutQuad:
		return 1 - (1-t)*(1-t)
	default: // EaseLinear
		return t
	}
}
