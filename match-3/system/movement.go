package system

import (
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/match-3/archetype"
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
	if !tweensComplete(grid) {
		return
	}
	finishSwap(grid, phase, display, e)
}

// phaseReverting waits for revert tweens, then returns to idle.
func phaseReverting(grid *component.GridData, phase *component.PhaseData, display *component.DisplayData) {
	if !tweensComplete(grid) {
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
func phaseCollapsing(grid *component.GridData, phase *component.PhaseData, display *component.DisplayData) {
	if !tweensComplete(grid) {
		return
	}
	collapse(grid, display)
	if tweensComplete(grid) {
		phase.Phase = component.PhaseRefilling
	}
}

// phaseRefilling spawns new tiles and checks for cascades.
func phaseRefilling(grid *component.GridData, phase *component.PhaseData, display *component.DisplayData, lvl *component.RulesData, e *ecs.ECS) {
	if !tweensComplete(grid) {
		return
	}
	refill(grid, display, lvl, e)
	phase.Phase = component.PhaseMatching // Check for cascades.
}

// tweensComplete checks if all active tweens are done.
func tweensComplete(grid *component.GridData) bool {
	for c := range grid.Cols {
		for r := range grid.Rows {
			entry := grid.Cells[c][r]
			if entry == nil {
				continue
			}
			tw := component.Tween.Get(entry)
			if tw.Active {
				return false
			}
		}
	}
	return true
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

type gridPos struct{ c, r int }

// findMatches returns all matched tile positions (runs of 3+).
func findMatches(grid *component.GridData) []gridPos {
	colorGrid := boardColorGrid(grid)
	matches := rules.FindMatches(colorGrid, grid.Cols, grid.Rows)
	result := make([]gridPos, len(matches))
	for i, m := range matches {
		result[i] = gridPos{m.Col, m.Row}
	}
	return result
}

// removeMatches destroys matched tile entities and clears cells.
func removeMatches(grid *component.GridData, matches []gridPos, e *ecs.ECS) {
	for _, pos := range matches {
		entry := grid.Cells[pos.c][pos.r]
		if entry == nil {
			continue
		}
		e.World.Remove(entry.Entity())
		grid.Cells[pos.c][pos.r] = nil
	}
}

// collapse drops tiles down to fill gaps.
func collapse(grid *component.GridData, display *component.DisplayData) {
	colorGrid := boardColorGrid(grid)
	moves := rules.Collapse(colorGrid, grid.CellType, grid.Cols, grid.Rows)

	for _, m := range moves {
		entry := grid.Cells[m.Col][m.FromRow]
		grid.Cells[m.Col][m.ToRow] = entry
		grid.Cells[m.Col][m.FromRow] = nil

		gp := component.GridPos.Get(entry)
		gp.Col, gp.Row = m.Col, m.ToRow

		endX := display.OffsetX + float64(m.Col*display.TileSize)
		endY := display.OffsetY + float64(m.ToRow*display.TileSize)
		StartTween(entry, endX, endY, 0.12, component.EaseOutQuad)
	}
}

// refill spawns new tiles for empty playable cells above the board and tweens them down.
func refill(grid *component.GridData, display *component.DisplayData, lvl *component.RulesData, e *ecs.ECS) {
	colorGrid := boardColorGrid(grid)
	empties := rules.EmptyCells(colorGrid, grid.CellType, grid.Cols, grid.Rows)

	for _, ec := range empties {
		color := rand.Intn(lvl.NumColors)
		px := display.OffsetX + float64(ec.Col*display.TileSize)
		py := display.OffsetY + float64(ec.Row*display.TileSize)

		// Spawn position: above board, staggered by spawn order.
		spawnY := display.OffsetY - float64(ec.SpawnOffset*display.TileSize)

		var sprite *ebiten.Image
		if color < len(display.GemSprites) {
			sprite = display.GemSprites[color]
		}

		entry := archetype.NewTile(e.World, ec.Col, ec.Row, color, sprite, px, spawnY)
		grid.Cells[ec.Col][ec.Row] = entry

		// Tween from spawn position to final position.
		duration := 0.08 * float64(ec.SpawnOffset)
		StartTween(entry, px, py, duration, component.EaseOutQuad)
	}
}

// hasValidMoves checks if any adjacent swap would produce a match.
func hasValidMoves(grid *component.GridData) bool {
	colorGrid := boardColorGrid(grid)
	return rules.HasValidMoves(colorGrid, grid.CellType, grid.Cols, grid.Rows)
}

// reshuffleBoard randomizes gem colors until the board is match-free with valid moves.
func reshuffleBoard(grid *component.GridData, phase *component.PhaseData, display *component.DisplayData) {
	phase.ReshuffleTimer = 2.0 // show "Reshuffling" message for 2 seconds

	colorGrid := boardColorGrid(grid)
	if rules.Reshuffle(colorGrid, grid.CellType, grid.Cols, grid.Rows, len(display.GemSprites), 100) {
		// Apply reshuffled colors back to entities.
		for c := range grid.Cols {
			for r := range grid.Rows {
				if grid.CellType[c][r] != component.CellPlayable || grid.Cells[c][r] == nil {
					continue
				}
				color := colorGrid[c][r]
				gem := component.GemType.Get(grid.Cells[c][r])
				gem.Color = color
				if color < len(display.GemSprites) {
					component.Sprite.Get(grid.Cells[c][r]).Image = display.GemSprites[color]
				}
			}
		}
	}
}

// boardColorGrid projects the board's ECS entities into a pure [][]int color grid.
func boardColorGrid(grid *component.GridData) [][]int {
	return rules.ColorGrid(grid.Cols, grid.Rows, func(col, row int) int {
		if grid.CellType[col][row] != component.CellPlayable || grid.Cells[col][row] == nil {
			return -1
		}
		return component.GemType.Get(grid.Cells[col][row]).Color
	})
}
