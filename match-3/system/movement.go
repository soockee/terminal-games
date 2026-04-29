package system

import (
	"fmt"

	"github.com/soockee/terminal-games/match-3/component"
	"github.com/soockee/terminal-games/match-3/event"
	"github.com/soockee/terminal-games/match-3/rules"
	"github.com/yohamta/donburi/ecs"
)

// UpdateBoard drives the match-3 state machine: swap validation,
// match detection, gravity, and refill.
func UpdateBoard(e *ecs.ECS) {
	boardEntry, ok := component.BoardGrid.First(e.World)
	if !ok {
		return
	}
	grid := component.BoardGrid.Get(boardEntry)
	phase := component.BoardPhase.Get(boardEntry)
	display := component.BoardDisplay.Get(boardEntry)
	lvl := component.BoardRules.Get(boardEntry)

	// Tick down reshuffle notification timer.
	if phase.ReshuffleTimer > 0 {
		phase.ReshuffleTimer -= 1.0 / 60.0
		if phase.ReshuffleTimer < 0 {
			phase.ReshuffleTimer = 0
		}
	}

	switch phase.Phase {
	case component.PhaseSwapping:
		phaseSwapping(grid, phase, display, e)
	case component.PhaseReverting:
		phaseReverting(grid, phase, display)
	case component.PhaseMatching:
		phaseMatching(grid, phase, display, e)
	case component.PhaseCollapsing:
		phaseCollapsing(grid, phase, display)
	case component.PhaseRefilling:
		phaseRefilling(grid, phase, display, lvl, e)
	}
}

// phaseSwapping waits for swap tweens to complete, then validates the swap.
func phaseSwapping(grid *component.GridData, phase *component.PhaseData, display *component.DisplayData, e *ecs.ECS) {
	if !tweensComplete(phase) {
		return
	}
	finishSwap(grid, phase, display, e)
}

// phaseReverting waits for revert tweens, then returns to idle.
func phaseReverting(grid *component.GridData, phase *component.PhaseData, display *component.DisplayData) {
	if !tweensComplete(phase) {
		return
	}
	phase.ChainDepth = 0
	if !hasValidMoves(grid) {
		reshuffleBoard(grid, phase, display)
	}
	phase.Phase = component.PhaseIdle
}

// phaseMatching detects cascading matches, removes them, and transitions to collapsing.
func phaseMatching(grid *component.GridData, phase *component.PhaseData, display *component.DisplayData, e *ecs.ECS) {
	matches := findMatches(grid)
	if len(matches) > 0 {
		phase.ChainDepth++
		removeMatches(grid, matches, e)
		if scoreEntry, ok := component.Score.First(e.World); ok {
			score := component.Score.Get(scoreEntry)
			score.Value += rules.ScoreForMatches(len(matches), phase.ChainDepth)
		}
		if phase.ChainDepth > 1 {
			chain := phase.ChainDepth
			if chain > 8 {
				chain = 8
			}
			event.AudioEvent.Publish(e.World, event.AudioEventData{Name: fmt.Sprintf("chain_%d", chain)})
		} else {
			event.AudioEvent.Publish(e.World, event.AudioEventData{Name: "match"})
		}
		phase.Phase = component.PhaseCollapsing
	} else {
		phase.ChainDepth = 0
		if !hasValidMoves(grid) {
			reshuffleBoard(grid, phase, display)
		}
		phase.Phase = component.PhaseIdle
	}
}

// phaseCollapsing applies gravity after matches are removed.
// Stays in PhaseCollapsing until collapse produces no new tweens.
func phaseCollapsing(grid *component.GridData, phase *component.PhaseData, display *component.DisplayData) {
	if !tweensComplete(phase) {
		return
	}
	moved := collapse(grid, display)
	if !moved {
		phase.Phase = component.PhaseRefilling
	}
}

// phaseRefilling spawns new tiles and checks for cascades.
func phaseRefilling(grid *component.GridData, phase *component.PhaseData, display *component.DisplayData, lvl *component.RulesData, e *ecs.ECS) {
	if !tweensComplete(phase) {
		return
	}
	refill(grid, display, lvl, e)
	phase.Phase = component.PhaseMatching // Check for cascades.
}

// tweensComplete checks if all active tweens are done via the counter
// maintained by UpdateTween.
func tweensComplete(phase *component.PhaseData) bool {
	return phase.ActiveTweens == 0
}

// finishSwap completes a swap, checks for matches. If no match, reverses swap.
func finishSwap(grid *component.GridData, phase *component.PhaseData, display *component.DisplayData, e *ecs.ECS) {
	a := phase.SwapA
	b := phase.SwapB

	// Actually swap the entries in the grid.
	grid.Cells[a[0]][a[1]], grid.Cells[b[0]][b[1]] = grid.Cells[b[0]][b[1]], grid.Cells[a[0]][a[1]]

	// Update GridPos components.
	if entryA := grid.Cells[a[0]][a[1]]; entryA != nil {
		gp := component.GridPos.Get(entryA)
		gp.Col, gp.Row = a[0], a[1]
	}
	if entryB := grid.Cells[b[0]][b[1]]; entryB != nil {
		gp := component.GridPos.Get(entryB)
		gp.Col, gp.Row = b[0], b[1]
	}

	// Check if swap creates matches.
	matches := findMatches(grid)
	if len(matches) > 0 {
		phase.ChainDepth = 1
		removeMatches(grid, matches, e)
		if scoreEntry, ok := component.Score.First(e.World); ok {
			score := component.Score.Get(scoreEntry)
			score.Value += rules.ScoreForMatches(len(matches), phase.ChainDepth)
		}
		event.AudioEvent.Publish(e.World, event.AudioEventData{Name: "match"})
		phase.Phase = component.PhaseCollapsing
	} else {
		// Invalid swap: reverse it.
		grid.Cells[a[0]][a[1]], grid.Cells[b[0]][b[1]] = grid.Cells[b[0]][b[1]], grid.Cells[a[0]][a[1]]
		if entryA := grid.Cells[a[0]][a[1]]; entryA != nil {
			gp := component.GridPos.Get(entryA)
			gp.Col, gp.Row = a[0], a[1]
		}
		if entryB := grid.Cells[b[0]][b[1]]; entryB != nil {
			gp := component.GridPos.Get(entryB)
			gp.Col, gp.Row = b[0], b[1]
		}
		// Tween back to original positions.
		StartSwapTween(grid, phase, component.EaseOutQuad, 0.12)
		event.AudioEvent.Publish(e.World, event.AudioEventData{Name: "invalid_swap"})
		phase.Phase = component.PhaseReverting
	}
}
