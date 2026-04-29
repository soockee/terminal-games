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

	boardEntry, ok := component.BoardGrid.First(e.World)
	if !ok {
		return
	}
	grid := component.BoardGrid.Get(boardEntry)
	phase := component.BoardPhase.Get(boardEntry)
	input := component.BoardInput.Get(boardEntry)
	display := component.BoardDisplay.Get(boardEntry)

	// Toggle autoplay on F4 (only functional with -tags=autoplay build).
	if inpututil.IsKeyJustPressed(ebiten.KeyF4) && autoPlayEnabled {
		input.AutoPlay = !input.AutoPlay
	}

	// Only accept input during Idle or Selected phases.
	if phase.Phase != component.PhaseIdle && phase.Phase != component.PhaseSelected {
		return
	}

	// Autoplay: find and execute a valid swap automatically.
	if input.AutoPlay && phase.Phase == component.PhaseIdle {
		// Stop autoplay if score target is reached.
		if scoreEntry, ok := component.Score.First(e.World); ok {
			s := component.Score.Get(scoreEntry)
			if s.Target > 0 && s.Value >= s.Target {
				input.AutoPlay = false
				return
			}
		}
		// Cooldown between moves so autoplay has the same pacing as a player.
		if input.AutoPlayDelay > 0 {
			input.AutoPlayDelay -= 1.0 / float64(ebiten.TPS())
			return
		}
		if tryAutoPlay(grid, phase, input) {
			input.AutoPlayDelay = 0.4 // seconds before next move
			return
		}
		input.AutoPlay = false
		return
	}

	// Map raw input to intent.
	intent := mapInput(grid, phase, input, display, e)
	if intent == nil {
		return
	}

	// Execute intent.
	executeIntent(intent, grid, phase, input, e)
}

// mapInput converts raw hardware input into a typed Intent.
// Returns nil if no actionable input was detected.
func mapInput(grid *component.GridData, phase *component.PhaseData, input *component.InputData, display *component.DisplayData, e *ecs.ECS) Intent {
	// Check for game state (start / restart).
	if gsEntry, gsOK := component.GameState.First(e.World); gsOK {
		gs := component.GameState.Get(gsEntry)
		if !gs.Started {
			if actionPressed() {
				return StartGame{}
			}
			return nil
		}
		if gs.Won || gs.Dead || gs.WinScreen {
			if actionPressed() {
				return RestartGame{}
			}
			return nil
		}
	}

	// Keyboard navigation.
	if intent := mapKeyboardInput(grid, phase, input); intent != nil {
		return intent
	}

	// Mouse/touch tap → grid coordinates.
	col, row, tapped := inputToGrid(display)
	if !tapped {
		return nil
	}

	// Validate: must be a playable cell with a tile.
	if col < 0 || col >= grid.Cols || row < 0 || row >= grid.Rows {
		return nil
	}
	if grid.CellType[col][row] != component.CellPlayable || grid.Cells[col][row] == nil {
		return nil
	}

	return classifyTap(phase, col, row)
}

// classifyTap determines the intent for a tap at (col, row) given current board phase.
func classifyTap(phase *component.PhaseData, col, row int) Intent {
	switch phase.Phase {
	case component.PhaseIdle:
		return SelectTile{Col: col, Row: row}

	case component.PhaseSelected:
		if col == phase.SelectedCol && row == phase.SelectedRow {
			return Deselect{}
		}
		dc := col - phase.SelectedCol
		dr := row - phase.SelectedRow
		if abs(dc)+abs(dr) != 1 {
			return ChangeSelection{Col: col, Row: row}
		}
		return InitiateSwap{
			FromCol: phase.SelectedCol, FromRow: phase.SelectedRow,
			ToCol: col, ToRow: row,
		}
	}
	return nil
}

// executeIntent applies an intent to the board state.
func executeIntent(intent Intent, grid *component.GridData, phase *component.PhaseData, input *component.InputData, e *ecs.ECS) {
	switch i := intent.(type) {
	case StartGame:
		if gsEntry, ok := component.GameState.First(e.World); ok {
			gs := component.GameState.Get(gsEntry)
			gs.Started = true
		}

	case RestartGame:
		if gsEntry, ok := component.GameState.First(e.World); ok {
			gs := component.GameState.Get(gsEntry)
			gs.Restart = true
		}

	case SelectTile:
		phase.SelectedCol = i.Col
		phase.SelectedRow = i.Row
		phase.Phase = component.PhaseSelected
		event.AudioEvent.Publish(e.World, event.AudioEventData{Name: "select"})

	case Deselect:
		phase.SelectedCol = -1
		phase.SelectedRow = -1
		phase.Phase = component.PhaseIdle

	case ChangeSelection:
		phase.SelectedCol = i.Col
		phase.SelectedRow = i.Row

	case InitiateSwap:
		phase.SwapA = [2]int{i.FromCol, i.FromRow}
		phase.SwapB = [2]int{i.ToCol, i.ToRow}
		phase.SelectedCol = -1
		phase.SelectedRow = -1
		StartSwapTween(grid, phase, component.EaseOutQuad, 0.15)
		phase.Phase = component.PhaseSwapping
		event.AudioEvent.Publish(e.World, event.AudioEventData{Name: "swap"})
	}
}

// actionPressed returns true if any action button was just pressed (Space, click, or touch).
func actionPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

// inputToGrid converts mouse/touch position to grid col/row.
func inputToGrid(display *component.DisplayData) (col, row int, ok bool) {
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

	col = int(float64(px)-display.OffsetX) / display.TileSize
	row = int(float64(py)-display.OffsetY) / display.TileSize
	return col, row, true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// mapKeyboardInput processes arrow key navigation and Enter/Space selection.
// Returns an Intent if a keyboard action produced one, nil otherwise.
func mapKeyboardInput(grid *component.GridData, phase *component.PhaseData, input *component.InputData) Intent {
	moved := false

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		input.CursorCol--
		moved = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		input.CursorCol++
		moved = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		input.CursorRow--
		moved = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		input.CursorRow++
		moved = true
	}

	// Clamp cursor to board bounds.
	if input.CursorCol < 0 {
		input.CursorCol = 0
	}
	if input.CursorCol >= grid.Cols {
		input.CursorCol = grid.Cols - 1
	}
	if input.CursorRow < 0 {
		input.CursorRow = 0
	}
	if input.CursorRow >= grid.Rows {
		input.CursorRow = grid.Rows - 1
	}

	// Enter/Space acts as a tap on the cursor position.
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		col, row := input.CursorCol, input.CursorRow
		if col < 0 || col >= grid.Cols || row < 0 || row >= grid.Rows {
			return nil
		}
		if grid.CellType[col][row] != component.CellPlayable || grid.Cells[col][row] == nil {
			return nil
		}
		return classifyTap(phase, col, row)
	}

	// Arrow keys moved the cursor but didn't produce an intent that needs execution.
	if moved {
		return nil
	}
	return nil
}
