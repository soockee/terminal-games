//go:build autoplay

package system

import (
	"github.com/soockee/terminal-games/match-3/component"
	"github.com/soockee/terminal-games/match-3/rules"
)

// tryAutoPlay evaluates all valid swaps, scores them by gems matched (simulating
// cascades), and picks the best move. Returns true if a move was made.
func tryAutoPlay(grid *component.GridData, phase *component.PhaseData, input *component.InputData) bool {
	type candidate struct {
		c1, r1, c2, r2 int
		score          int
	}

	colorGrid := boardColorGrid(grid)
	var best candidate
	dirs := [2][2]int{{1, 0}, {0, 1}}

	for c := range grid.Cols {
		for r := range grid.Rows {
			if grid.CellType[c][r] != component.CellPlayable || grid.Cells[c][r] == nil {
				continue
			}
			for _, d := range dirs {
				nc, nr := c+d[0], r+d[1]
				if nc >= grid.Cols || nr >= grid.Rows {
					continue
				}
				if grid.CellType[nc][nr] != component.CellPlayable || grid.Cells[nc][nr] == nil {
					continue
				}
				score := rules.SimulateSwap(colorGrid, grid.CellType, grid.Cols, grid.Rows, c, r, nc, nr)
				if score > best.score {
					best = candidate{c, r, nc, nr, score}
				}
			}
		}
	}

	if best.score == 0 {
		return false
	}

	phase.SwapA = [2]int{best.c1, best.r1}
	phase.SwapB = [2]int{best.c2, best.r2}
	phase.SelectedCol = -1
	phase.SelectedRow = -1
	StartSwapTween(grid, phase, component.EaseOutQuad, 0.15)
	phase.Phase = component.PhaseSwapping
	return true
}

const autoPlayEnabled = true
