// Package rules implements pure match-3 board logic on a [][]int color grid.
// A cell value of -1 means empty/unplayable. Values >= 0 are gem colors.
package rules

import "math/rand"

// GridPos identifies a cell on the board.
type GridPos struct{ Col, Row int }

// FindMatches returns all positions involved in runs of 3+ same-colored gems.
func FindMatches(grid [][]int, cols, rows int) []GridPos {
	matched := map[GridPos]bool{}

	// Horizontal runs.
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
					matched[GridPos{c + i, r}] = true
				}
			}
		}
	}

	// Vertical runs.
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
					matched[GridPos{c, r + i}] = true
				}
			}
		}
	}

	result := make([]GridPos, 0, len(matched))
	for pos := range matched {
		result = append(result, pos)
	}
	return result
}

// Collapse drops gems down to fill gaps. cellType marks playable cells (2 = playable).
// Returns a list of moves: {col, fromRow, toRow} for each tile that fell.
type FallMove struct {
	Col     int
	FromRow int
	ToRow   int
}

func Collapse(grid [][]int, cellType [][]int, cols, rows int) []FallMove {
	var moves []FallMove
	for c := range cols {
		writeRow := rows - 1
		for r := rows - 1; r >= 0; r-- {
			if cellType[c][r] != 2 { // not playable
				writeRow = r - 1
				continue
			}
			if grid[c][r] >= 0 {
				if r != writeRow {
					grid[c][writeRow] = grid[c][r]
					grid[c][r] = -1
					moves = append(moves, FallMove{Col: c, FromRow: r, ToRow: writeRow})
				}
				writeRow--
			}
		}
	}
	return moves
}

// EmptyCells returns all playable cells that are empty (for refill).
type EmptyCell struct {
	Col         int
	Row         int
	SpawnOffset int // 1-based offset for stagger
}

func EmptyCells(grid [][]int, cellType [][]int, cols, rows int) []EmptyCell {
	var cells []EmptyCell
	for c := range cols {
		spawnOffset := 0
		for r := range rows {
			if cellType[c][r] != 2 {
				continue
			}
			if grid[c][r] >= 0 {
				continue
			}
			spawnOffset++
			cells = append(cells, EmptyCell{Col: c, Row: r, SpawnOffset: spawnOffset})
		}
	}
	return cells
}

// HasValidMoves checks if any adjacent swap would produce a match.
func HasValidMoves(grid [][]int, cellType [][]int, cols, rows int) bool {
	dirs := [2][2]int{{1, 0}, {0, 1}}
	for c := range cols {
		for r := range rows {
			if cellType[c][r] != 2 || grid[c][r] < 0 {
				continue
			}
			for _, d := range dirs {
				nc, nr := c+d[0], r+d[1]
				if nc >= cols || nr >= rows {
					continue
				}
				if cellType[nc][nr] != 2 || grid[nc][nr] < 0 {
					continue
				}
				colorA, colorB := grid[c][r], grid[nc][nr]
				if colorA == colorB {
					continue
				}
				// Temporarily swap and check.
				grid[c][r], grid[nc][nr] = colorB, colorA
				matches := FindMatches(grid, cols, rows)
				grid[c][r], grid[nc][nr] = colorA, colorB
				if len(matches) > 0 {
					return true
				}
			}
		}
	}
	return false
}

// Reshuffle randomizes gem colors until the board is match-free with valid moves.
// Returns true if a valid configuration was found within maxAttempts.
func Reshuffle(grid [][]int, cellType [][]int, cols, rows, numColors, maxAttempts int) bool {
	for attempts := range maxAttempts {
		_ = attempts
		for c := range cols {
			for r := range rows {
				if cellType[c][r] != 2 || grid[c][r] < 0 {
					continue
				}
				grid[c][r] = rand.Intn(numColors)
			}
		}
		if len(FindMatches(grid, cols, rows)) == 0 && HasValidMoves(grid, cellType, cols, rows) {
			return true
		}
	}
	return false
}

// SimulateSwap scores a swap by counting total gems cleared including cascades.
func SimulateSwap(grid [][]int, cellType [][]int, cols, rows, c1, r1, c2, r2 int) int {
	// Make a copy to simulate on.
	sim := make([][]int, cols)
	for c := range cols {
		sim[c] = make([]int, rows)
		copy(sim[c], grid[c])
	}

	sim[c1][r1], sim[c2][r2] = sim[c2][r2], sim[c1][r1]

	total := 0
	for {
		matched := FindMatches(sim, cols, rows)
		if len(matched) == 0 {
			break
		}
		total += len(matched)
		for _, pos := range matched {
			sim[pos.Col][pos.Row] = -1
		}
		Collapse(sim, cellType, cols, rows)
	}
	return total
}

// IsValidSwap checks whether a swap between two cells is a legal move.
// A swap is invalid if: out of bounds, non-adjacent, either cell is
// non-playable or empty, or the swap doesn't produce any match (including cascades).
func IsValidSwap(grid [][]int, cellType [][]int, cols, rows, c1, r1, c2, r2 int) bool {
	// Bounds check.
	if c1 < 0 || c1 >= cols || r1 < 0 || r1 >= rows {
		return false
	}
	if c2 < 0 || c2 >= cols || r2 < 0 || r2 >= rows {
		return false
	}

	// Adjacency check (Manhattan distance must be 1).
	dc := c2 - c1
	if dc < 0 {
		dc = -dc
	}
	dr := r2 - r1
	if dr < 0 {
		dr = -dr
	}
	if dc+dr != 1 {
		return false
	}

	// Playability check.
	if cellType[c1][r1] != 2 || cellType[c2][r2] != 2 {
		return false
	}

	// Empty cell check.
	if grid[c1][r1] < 0 || grid[c2][r2] < 0 {
		return false
	}

	// Same color — no point in swapping.
	if grid[c1][r1] == grid[c2][r2] {
		return false
	}

	// Produces a match (including cascades).
	return SimulateSwap(grid, cellType, cols, rows, c1, r1, c2, r2) > 0
}

// ScoreForMatches calculates the score earned for a set of matched gems.
func ScoreForMatches(matchCount, chainDepth int) int {
	return matchCount * 10 * chainDepth
}

// ColorGrid extracts a [][]int color grid from a board's cell information.
// cellColor should return the color at (col, row), or -1 if empty/unplayable.
func ColorGrid(cols, rows int, cellColor func(col, row int) int) [][]int {
	grid := make([][]int, cols)
	for c := range cols {
		grid[c] = make([]int, rows)
		for r := range rows {
			grid[c][r] = cellColor(c, r)
		}
	}
	return grid
}
