package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/soockee/terminal-games/match-3/component"
	"github.com/soockee/terminal-games/match-3/event"
	"github.com/yohamta/donburi/ecs"
)

// UpdateInput handles debug toggle and tile selection/swap input.
func UpdateInput(e *ecs.ECS) {
	// Toggle debug on F3.
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		if entry, ok := component.Debug.First(e.World); ok {
			d := component.Debug.Get(entry)
			d.Enabled = !d.Enabled
		}
	}

	boardEntry, ok := component.Board.First(e.World)
	if !ok {
		return
	}
	board := component.Board.Get(boardEntry)

	// Toggle autoplay on F4 (only functional with -tags=autoplay build).
	if inpututil.IsKeyJustPressed(ebiten.KeyF4) && autoPlayEnabled {
		board.AutoPlay = !board.AutoPlay
	}

	// Only accept input during Idle or Selected phases.
	if board.Phase != component.PhaseIdle && board.Phase != component.PhaseSelected {
		return
	}

	// Check for game state (start / restart).
	if gsEntry, gsOK := component.GameState.First(e.World); gsOK {
		gs := component.GameState.Get(gsEntry)
		if !gs.Started {
			actionPressed := inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
				inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
				len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
			if actionPressed {
				gs.Started = true
			}
			return
		}
		if gs.Won || gs.Dead || gs.WinScreen {
			actionPressed := inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
				inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
				len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
			if actionPressed {
				gs.Restart = true
			}
			return
		}
	}

	// Detect tap/click position → grid coordinates.
	col, row, tapped := inputToGrid(board)

	// Autoplay: find and execute a valid swap automatically.
	if board.AutoPlay && board.Phase == component.PhaseIdle {
		if tryAutoPlay(board) {
			return
		}
		// No valid move found — dead end. Stop autoplay.
		board.AutoPlay = false
		return
	}

	// Keyboard navigation: arrow keys move cursor, Enter/Space selects.
	if handleKeyboardInput(board, e) {
		return
	}

	if !tapped {
		return
	}

	// Validate: must be a playable cell with a tile.
	if col < 0 || col >= board.Cols || row < 0 || row >= board.Rows {
		return
	}
	if board.CellType[col][row] != component.CellPlayable || board.Cells[col][row] == nil {
		return
	}

	switch board.Phase {
	case component.PhaseIdle:
		board.SelectedCol = col
		board.SelectedRow = row
		board.Phase = component.PhaseSelected
		event.AudioEvent.Publish(e.World, event.AudioEventData{Name: "select"})

	case component.PhaseSelected:
		if col == board.SelectedCol && row == board.SelectedRow {
			// Deselect.
			board.SelectedCol = -1
			board.SelectedRow = -1
			board.Phase = component.PhaseIdle
			return
		}

		// Check adjacency (Manhattan distance == 1).
		dc := col - board.SelectedCol
		dr := row - board.SelectedRow
		if abs(dc)+abs(dr) != 1 {
			// Non-adjacent: change selection.
			board.SelectedCol = col
			board.SelectedRow = row
			return
		}

		// Initiate swap.
		board.SwapA = [2]int{board.SelectedCol, board.SelectedRow}
		board.SwapB = [2]int{col, row}
		board.SelectedCol = -1
		board.SelectedRow = -1
		startSwapTweens(board, component.EaseOutQuad, 0.15)
		board.Phase = component.PhaseSwapping
		event.AudioEvent.Publish(e.World, event.AudioEventData{Name: "swap"})
	}
}

// inputToGrid converts mouse/touch position to grid col/row.
func inputToGrid(board *component.BoardData) (col, row int, ok bool) {
	var px, py int

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		px, py = ebiten.CursorPosition()
		ok = true
	} else {
		touchIDs := inpututil.AppendJustPressedTouchIDs(nil)
		if len(touchIDs) > 0 {
			px, py = ebiten.TouchPosition(touchIDs[0])
			ok = true
		}
	}

	if !ok {
		return 0, 0, false
	}

	col = int(float64(px)-board.OffsetX) / board.TileSize
	row = int(float64(py)-board.OffsetY) / board.TileSize
	return col, row, true
}

// startSwapTweens initiates position tweens on both swap tiles.
func startSwapTweens(board *component.BoardData, ease component.EaseFunc, duration float64) {
	a := board.Cells[board.SwapA[0]][board.SwapA[1]]
	b := board.Cells[board.SwapB[0]][board.SwapB[1]]

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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// handleKeyboardInput processes arrow key navigation and Enter/Space selection.
// Returns true if a keyboard action was consumed.
func handleKeyboardInput(board *component.BoardData, e *ecs.ECS) bool {
	moved := false

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		board.CursorCol--
		moved = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		board.CursorCol++
		moved = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		board.CursorRow--
		moved = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		board.CursorRow++
		moved = true
	}

	// Clamp cursor to board bounds.
	if board.CursorCol < 0 {
		board.CursorCol = 0
	}
	if board.CursorCol >= board.Cols {
		board.CursorCol = board.Cols - 1
	}
	if board.CursorRow < 0 {
		board.CursorRow = 0
	}
	if board.CursorRow >= board.Rows {
		board.CursorRow = board.Rows - 1
	}

	// Enter/Space acts as a tap on the cursor position.
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		col, row := board.CursorCol, board.CursorRow
		if col < 0 || col >= board.Cols || row < 0 || row >= board.Rows {
			return moved
		}
		if board.CellType[col][row] != component.CellPlayable || board.Cells[col][row] == nil {
			return moved
		}

		switch board.Phase {
		case component.PhaseIdle:
			board.SelectedCol = col
			board.SelectedRow = row
			board.Phase = component.PhaseSelected
			event.AudioEvent.Publish(e.World, event.AudioEventData{Name: "select"})

		case component.PhaseSelected:
			if col == board.SelectedCol && row == board.SelectedRow {
				board.SelectedCol = -1
				board.SelectedRow = -1
				board.Phase = component.PhaseIdle
				return true
			}

			dc := col - board.SelectedCol
			dr := row - board.SelectedRow
			if abs(dc)+abs(dr) != 1 {
				board.SelectedCol = col
				board.SelectedRow = row
				return true
			}

			board.SwapA = [2]int{board.SelectedCol, board.SelectedRow}
			board.SwapB = [2]int{col, row}
			board.SelectedCol = -1
			board.SelectedRow = -1
			startSwapTweens(board, component.EaseOutQuad, 0.15)
			board.Phase = component.PhaseSwapping
			event.AudioEvent.Publish(e.World, event.AudioEventData{Name: "swap"})
		}
		return true
	}

	return moved
}
