package system

import (
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/match-3/archetype"
	"github.com/soockee/terminal-games/match-3/component"
	"github.com/soockee/terminal-games/match-3/event"
	"github.com/yohamta/donburi/ecs"
)

// UpdateBoard drives the match-3 state machine: swap validation,
// match detection, gravity, and refill.
func UpdateBoard(e *ecs.ECS) {
	boardEntry, ok := component.Board.First(e.World)
	if !ok {
		return
	}
	board := component.Board.Get(boardEntry)

	// Tick down reshuffle notification timer.
	if board.ReshuffleTimer > 0 {
		board.ReshuffleTimer -= 1.0 / 60.0
		if board.ReshuffleTimer < 0 {
			board.ReshuffleTimer = 0
		}
	}

	switch board.Phase {
	case component.PhaseSwapping:
		if tweensComplete(board) {
			finishSwap(board, e)
		}

	case component.PhaseReverting:
		if tweensComplete(board) {
			board.ChainDepth = 0
			if !hasValidMoves(board) {
				reshuffleBoard(board)
			}
			board.Phase = component.PhaseIdle
		}

	case component.PhaseMatching:
		matches := findMatches(board)
		if len(matches) > 0 {
			board.ChainDepth++
			removeMatches(board, matches, e)
			if scoreEntry, ok := component.Score.First(e.World); ok {
				score := component.Score.Get(scoreEntry)
				score.Value += len(matches) * 10 * board.ChainDepth
			}
			if board.ChainDepth > 1 {
				chain := board.ChainDepth
				if chain > 8 {
					chain = 8
				}
				event.AudioEvent.Publish(e.World, event.AudioEventData{Name: fmt.Sprintf("chain_%d", chain)})
			} else {
				event.AudioEvent.Publish(e.World, event.AudioEventData{Name: "match"})
			}
			board.Phase = component.PhaseCollapsing
		} else {
			board.ChainDepth = 0
			if !hasValidMoves(board) {
				reshuffleBoard(board)
			}
			board.Phase = component.PhaseIdle
		}

	case component.PhaseCollapsing:
		if tweensComplete(board) {
			collapse(board)
			if tweensComplete(board) {
				board.Phase = component.PhaseRefilling
			}
		}

	case component.PhaseRefilling:
		if tweensComplete(board) {
			refill(board, e)
			board.Phase = component.PhaseMatching // Check for cascades.
		}
	}
}

// tweensComplete checks if all active tweens are done.
func tweensComplete(board *component.BoardData) bool {
	for c := range board.Cols {
		for r := range board.Rows {
			entry := board.Cells[c][r]
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
func finishSwap(board *component.BoardData, e *ecs.ECS) {
	a := board.SwapA
	b := board.SwapB

	// Actually swap the entries in the grid.
	board.Cells[a[0]][a[1]], board.Cells[b[0]][b[1]] = board.Cells[b[0]][b[1]], board.Cells[a[0]][a[1]]

	// Update GridPos components.
	if entryA := board.Cells[a[0]][a[1]]; entryA != nil {
		gp := component.GridPos.Get(entryA)
		gp.Col, gp.Row = a[0], a[1]
	}
	if entryB := board.Cells[b[0]][b[1]]; entryB != nil {
		gp := component.GridPos.Get(entryB)
		gp.Col, gp.Row = b[0], b[1]
	}

	// Check if swap creates matches.
	matches := findMatches(board)
	if len(matches) > 0 {
		board.ChainDepth = 1
		removeMatches(board, matches, e)
		if scoreEntry, ok := component.Score.First(e.World); ok {
			score := component.Score.Get(scoreEntry)
			score.Value += len(matches) * 10 * board.ChainDepth
		}
		event.AudioEvent.Publish(e.World, event.AudioEventData{Name: "match"})
		board.Phase = component.PhaseCollapsing
	} else {
		// Invalid swap: reverse it.
		board.Cells[a[0]][a[1]], board.Cells[b[0]][b[1]] = board.Cells[b[0]][b[1]], board.Cells[a[0]][a[1]]
		if entryA := board.Cells[a[0]][a[1]]; entryA != nil {
			gp := component.GridPos.Get(entryA)
			gp.Col, gp.Row = a[0], a[1]
		}
		if entryB := board.Cells[b[0]][b[1]]; entryB != nil {
			gp := component.GridPos.Get(entryB)
			gp.Col, gp.Row = b[0], b[1]
		}
		// Tween back to original positions.
		startSwapTweens(board, component.EaseOutQuad, 0.12)
		event.AudioEvent.Publish(e.World, event.AudioEventData{Name: "invalid_swap"})
		board.Phase = component.PhaseReverting
	}
}

type gridPos struct{ c, r int }

// findMatches returns all matched tile positions (runs of 3+).
func findMatches(board *component.BoardData) []gridPos {
	matched := map[gridPos]bool{}

	// Horizontal runs.
	for r := range board.Rows {
		for c := 0; c < board.Cols-2; c++ {
			if board.Cells[c][r] == nil {
				continue
			}
			color := component.GemType.Get(board.Cells[c][r]).Color
			run := 1
			for nc := c + 1; nc < board.Cols && board.Cells[nc][r] != nil; nc++ {
				if component.GemType.Get(board.Cells[nc][r]).Color == color {
					run++
				} else {
					break
				}
			}
			if run >= 3 {
				for i := 0; i < run; i++ {
					matched[gridPos{c + i, r}] = true
				}
			}
		}
	}

	// Vertical runs.
	for c := range board.Cols {
		for r := 0; r < board.Rows-2; r++ {
			if board.Cells[c][r] == nil {
				continue
			}
			color := component.GemType.Get(board.Cells[c][r]).Color
			run := 1
			for nr := r + 1; nr < board.Rows && board.Cells[c][nr] != nil; nr++ {
				if component.GemType.Get(board.Cells[c][nr]).Color == color {
					run++
				} else {
					break
				}
			}
			if run >= 3 {
				for i := 0; i < run; i++ {
					matched[gridPos{c, r + i}] = true
				}
			}
		}
	}

	result := make([]gridPos, 0, len(matched))
	for pos := range matched {
		result = append(result, pos)
	}
	return result
}

// removeMatches destroys matched tile entities and clears cells.
func removeMatches(board *component.BoardData, matches []gridPos, e *ecs.ECS) {
	for _, pos := range matches {
		entry := board.Cells[pos.c][pos.r]
		if entry == nil {
			continue
		}
		e.World.Remove(entry.Entity())
		board.Cells[pos.c][pos.r] = nil
	}
}

// collapse drops tiles down to fill gaps.
func collapse(board *component.BoardData) {
	for c := range board.Cols {
		// Walk from bottom up, pull tiles down into empty playable slots.
		writeRow := board.Rows - 1
		for r := board.Rows - 1; r >= 0; r-- {
			if board.CellType[c][r] != component.CellPlayable {
				writeRow = r - 1
				continue
			}
			if board.Cells[c][r] != nil {
				if r != writeRow {
					entry := board.Cells[c][r]
					board.Cells[c][writeRow] = entry
					board.Cells[c][r] = nil

					gp := component.GridPos.Get(entry)
					gp.Col, gp.Row = c, writeRow

					px := board.OffsetX + float64(c*board.TileSize)
					py := board.OffsetY + float64(writeRow*board.TileSize)

					tw := component.Tween.Get(entry)
					pos := component.PixelPos.Get(entry)
					tw.StartX, tw.StartY = pos.X, pos.Y
					tw.EndX, tw.EndY = px, py
					tw.Elapsed = 0
					tw.Duration = 0.12
					tw.Active = true
					tw.Ease = component.EaseOutQuad
				}
				writeRow--
			}
		}
	}
}

// refill spawns new tiles for empty playable cells above the board and tweens them down.
func refill(board *component.BoardData, e *ecs.ECS) {
	for c := range board.Cols {
		// Count how many empty playable cells exist in this column (for staggering spawn height).
		spawnOffset := 0
		for r := range board.Rows {
			if board.CellType[c][r] != component.CellPlayable {
				continue
			}
			if board.Cells[c][r] != nil {
				continue
			}

			spawnOffset++
			color := rand.Intn(board.NumColors)
			px := board.OffsetX + float64(c*board.TileSize)
			py := board.OffsetY + float64(r*board.TileSize)

			// Spawn position: above board, staggered by spawn order.
			spawnY := board.OffsetY - float64(spawnOffset*board.TileSize)

			var sprite *ebiten.Image
			if color < len(board.GemSprites) {
				sprite = board.GemSprites[color]
			}

			entry := archetype.NewTile(e.World, c, r, color, sprite, px, spawnY)
			board.Cells[c][r] = entry

			// Tween from spawn position to final position.
			tw := component.Tween.Get(entry)
			tw.StartX, tw.StartY = px, spawnY
			tw.EndX, tw.EndY = px, py
			tw.Elapsed = 0
			tw.Duration = 0.08 * float64(spawnOffset)
			tw.Active = true
			tw.Ease = component.EaseOutQuad
		}
	}
}

// hasValidMoves checks if any adjacent swap would produce a match.
func hasValidMoves(board *component.BoardData) bool {
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
				// Simulate swap.
				colorA := component.GemType.Get(board.Cells[c][r]).Color
				colorB := component.GemType.Get(board.Cells[nc][nr]).Color
				if colorA == colorB {
					continue
				}
				// Temporarily swap colors and check.
				component.GemType.Get(board.Cells[c][r]).Color = colorB
				component.GemType.Get(board.Cells[nc][nr]).Color = colorA
				matches := findMatches(board)
				// Swap back.
				component.GemType.Get(board.Cells[c][r]).Color = colorA
				component.GemType.Get(board.Cells[nc][nr]).Color = colorB
				if len(matches) > 0 {
					return true
				}
			}
		}
	}
	return false
}

// reshuffleBoard randomizes gem colors until the board is match-free with valid moves.
func reshuffleBoard(board *component.BoardData) {
	board.ReshuffleTimer = 2.0 // show "Reshuffling" message for 2 seconds

	for attempts := 0; attempts < 100; attempts++ {
		// Randomize all gem colors.
		for c := range board.Cols {
			for r := range board.Rows {
				if board.CellType[c][r] != component.CellPlayable || board.Cells[c][r] == nil {
					continue
				}
				color := rand.Intn(board.NumColors)
				gem := component.GemType.Get(board.Cells[c][r])
				gem.Color = color
				if color < len(board.GemSprites) {
					component.Sprite.Get(board.Cells[c][r]).Image = board.GemSprites[color]
				}
			}
		}

		// Check: no existing matches and at least one valid move.
		if len(findMatches(board)) == 0 && hasValidMoves(board) {
			return
		}
	}
}
