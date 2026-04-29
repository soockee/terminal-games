package system

import (
	"github.com/soockee/terminal-games/match-3/component"
	"github.com/soockee/terminal-games/match-3/rules"
	"github.com/yohamta/donburi/ecs"
)

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

// boardColorGrid projects the board's ECS entities into a pure [][]int color grid.
func boardColorGrid(grid *component.GridData) [][]int {
	return rules.ColorGrid(grid.Cols, grid.Rows, func(col, row int) int {
		if grid.CellType[col][row] != component.CellPlayable || grid.Cells[col][row] == nil {
			return -1
		}
		return component.GemType.Get(grid.Cells[col][row]).Color
	})
}
