//go:build autoplay

package system

import (
	"github.com/soockee/terminal-games/match-3/component"
)

// tryAutoPlay evaluates all valid swaps, scores them by gems matched (simulating
// cascades), and picks the best move. Returns true if a move was made.
func tryAutoPlay(board *component.BoardData) bool {
	type candidate struct {
		c1, r1, c2, r2 int
		score          int
	}

	var best candidate
	dirs := [][2]int{{1, 0}, {0, 1}}

	for c := range board.Cols {
		for r := range board.Rows {
			if board.CellType[c][r] != component.CellPlayable || board.Cells[c][r] == nil {
				continue
			}
			for _, d := range dirs {
				nc, nr := c+d[0], r+d[1]
				if nc >= board.Cols || nr >= board.Rows {
					continue
				}
				if board.CellType[nc][nr] != component.CellPlayable || board.Cells[nc][nr] == nil {
					continue
				}
				score := simulateSwap(board, c, r, nc, nr)
				if score > best.score {
					best = candidate{c, r, nc, nr, score}
				}
			}
		}
	}

	if best.score == 0 {
		return false
	}

	board.SwapA = [2]int{best.c1, best.r1}
	board.SwapB = [2]int{best.c2, best.r2}
	board.SelectedCol = -1
	board.SelectedRow = -1
	startSwapTweens(board, component.EaseOutQuad, 0.15)
	board.Phase = component.PhaseSwapping
	return true
}

// simulateSwap scores a swap by counting total gems cleared including cascades.
func simulateSwap(board *component.BoardData, c1, r1, c2, r2 int) int {
	cols, rows := board.Cols, board.Rows
	grid := make([][]int, cols)
	for c := range cols {
		grid[c] = make([]int, rows)
		for r := range rows {
			if board.CellType[c][r] != component.CellPlayable || board.Cells[c][r] == nil {
				grid[c][r] = -1
			} else {
				grid[c][r] = component.GemType.Get(board.Cells[c][r]).Color
			}
		}
	}

	grid[c1][r1], grid[c2][r2] = grid[c2][r2], grid[c1][r1]

	total := 0
	for {
		matched := simFindMatches(grid, cols, rows)
		if len(matched) == 0 {
			break
		}
		total += len(matched)
		for _, pos := range matched {
			grid[pos[0]][pos[1]] = -1
		}
		simCollapse(grid, cols, rows, board)
	}

	return total
}

func simFindMatches(grid [][]int, cols, rows int) [][2]int {
	matched := map[[2]int]bool{}

	for r := range rows {
		for c := 0; c < cols-2; c++ {
			if grid[c][r] < 0 {
				continue
			}
			color := grid[c][r]
			run := 1
			for nc := c + 1; nc < cols && grid[nc][r] == color; nc++ {
				run++
			}
			if run >= 3 {
				for i := range run {
					matched[[2]int{c + i, r}] = true
				}
			}
		}
	}

	for c := range cols {
		for r := 0; r < rows-2; r++ {
			if grid[c][r] < 0 {
				continue
			}
			color := grid[c][r]
			run := 1
			for nr := r + 1; nr < rows && grid[c][nr] == color; nr++ {
				run++
			}
			if run >= 3 {
				for i := range run {
					matched[[2]int{c, r + i}] = true
				}
			}
		}
	}

	result := make([][2]int, 0, len(matched))
	for pos := range matched {
		result = append(result, pos)
	}
	return result
}

func simCollapse(grid [][]int, cols, rows int, board *component.BoardData) {
	for c := range cols {
		writeRow := rows - 1
		for r := rows - 1; r >= 0; r-- {
			if board.CellType[c][r] != component.CellPlayable {
				writeRow = r - 1
				continue
			}
			if grid[c][r] >= 0 {
				if r != writeRow {
					grid[c][writeRow] = grid[c][r]
					grid[c][r] = -1
				}
				writeRow--
			}
		}
	}
}

const autoPlayEnabled = true
