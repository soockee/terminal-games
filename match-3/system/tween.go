package system

import (
	"github.com/soockee/terminal-games/match-3/component"
	"github.com/soockee/terminal-games/match-3/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var tileQuery = donburi.NewQuery(filter.Contains(tags.Tile))

// StartTween initiates a position tween on an entity, reading its current PixelPos
// as the start and animating to (endX, endY) over the given duration.
func StartTween(entry *donburi.Entry, endX, endY, duration float64, ease component.EaseFunc) {
	pos := component.PixelPos.Get(entry)
	tw := component.Tween.Get(entry)
	tw.StartX, tw.StartY = pos.X, pos.Y
	tw.EndX, tw.EndY = endX, endY
	tw.Elapsed = 0
	tw.Duration = duration
	tw.Active = true
	tw.Ease = ease
}

// StartSwapTween initiates position tweens on both swap tiles so they animate
// towards each other's positions.
func StartSwapTween(grid *component.GridData, phase *component.PhaseData, ease component.EaseFunc, duration float64) {
	a := grid.Cells[phase.SwapA[0]][phase.SwapA[1]]
	b := grid.Cells[phase.SwapB[0]][phase.SwapB[1]]

	posA := component.PixelPos.Get(a)
	posB := component.PixelPos.Get(b)

	twA := component.Tween.Get(a)
	twA.StartX, twA.StartY = posA.X, posA.Y
	twA.EndX, twA.EndY = posB.X, posB.Y
	twA.Elapsed = 0
	twA.Duration = duration
	twA.Active = true
	twA.Ease = ease

	twB := component.Tween.Get(b)
	twB.StartX, twB.StartY = posB.X, posB.Y
	twB.EndX, twB.EndY = posA.X, posA.Y
	twB.Elapsed = 0
	twB.Duration = duration
	twB.Active = true
	twB.Ease = ease
}

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
