package system

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/match-3/archetype"
	"github.com/soockee/terminal-games/match-3/component"
	"github.com/soockee/terminal-games/match-3/rules"
	"github.com/yohamta/donburi/ecs"
)

// collapse drops tiles down to fill gaps. Returns true if any tiles moved.
func collapse(grid *component.GridData, display *component.DisplayData) bool {
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
	return len(moves) > 0
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
