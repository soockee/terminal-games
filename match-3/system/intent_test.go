package system

import (
	"testing"

	"github.com/soockee/terminal-games/match-3/component"
	"github.com/yohamta/donburi"
)

func TestResolveIntent_OutOfBounds(t *testing.T) {
	w := donburi.NewWorld()
	entry := w.Entry(w.Create(component.GridPos, component.GemType, component.PixelPos, component.Tween, component.Sprite))

	cellType := [][]int{{2, 2, 2}}
	cells := [][]*donburi.Entry{{entry, entry, entry}}
	phase := &component.PhaseData{Phase: component.PhaseIdle}

	cases := []struct {
		name     string
		col, row int
	}{
		{"negative col", -1, 0},
		{"negative row", 0, -1},
		{"col too large", 1, 0},
		{"row too large", 0, 3},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			intent := ResolveIntent(tc.col, tc.row, phase, cellType, cells)
			if intent != nil {
				t.Fatalf("expected nil for out-of-bounds (%d,%d), got %v", tc.col, tc.row, intent)
			}
		})
	}
}

func TestResolveIntent_NonPlayableCell(t *testing.T) {
	w := donburi.NewWorld()
	entry := w.Entry(w.Create(component.GridPos, component.GemType, component.PixelPos, component.Tween, component.Sprite))

	cellType := [][]int{{1, 2, 2}} // row 0 is not playable
	cells := [][]*donburi.Entry{{entry, entry, entry}}
	phase := &component.PhaseData{Phase: component.PhaseIdle}

	intent := ResolveIntent(0, 0, phase, cellType, cells)
	if intent != nil {
		t.Fatal("expected nil for non-playable cell")
	}
}

func TestResolveIntent_NilEntry(t *testing.T) {
	cellType := [][]int{{2, 2, 2}}
	cells := [][]*donburi.Entry{{nil, nil, nil}} // no tile entities
	phase := &component.PhaseData{Phase: component.PhaseIdle}

	intent := ResolveIntent(0, 0, phase, cellType, cells)
	if intent != nil {
		t.Fatal("expected nil for nil entry")
	}
}

func TestResolveIntent_IdlePhaseSelectsTile(t *testing.T) {
	w := donburi.NewWorld()
	entry := w.Entry(w.Create(component.GridPos, component.GemType, component.PixelPos, component.Tween, component.Sprite))

	cellType := [][]int{{2, 2}, {2, 2}}
	cells := [][]*donburi.Entry{{entry, entry}, {entry, entry}}
	phase := &component.PhaseData{Phase: component.PhaseIdle}

	intent := ResolveIntent(1, 0, phase, cellType, cells)
	sel, ok := intent.(SelectTile)
	if !ok {
		t.Fatalf("expected SelectTile, got %T", intent)
	}
	if sel.Col != 1 || sel.Row != 0 {
		t.Fatalf("expected SelectTile{1,0}, got %+v", sel)
	}
}

func TestResolveIntent_SelectedPhaseDeselects(t *testing.T) {
	w := donburi.NewWorld()
	entry := w.Entry(w.Create(component.GridPos, component.GemType, component.PixelPos, component.Tween, component.Sprite))

	cellType := [][]int{{2, 2}, {2, 2}}
	cells := [][]*donburi.Entry{{entry, entry}, {entry, entry}}
	phase := &component.PhaseData{
		Phase:       component.PhaseSelected,
		SelectedCol: 1,
		SelectedRow: 0,
	}

	// Tap same tile → deselect
	intent := ResolveIntent(1, 0, phase, cellType, cells)
	if _, ok := intent.(Deselect); !ok {
		t.Fatalf("expected Deselect, got %T", intent)
	}
}

func TestResolveIntent_SelectedPhaseAdjacentSwap(t *testing.T) {
	w := donburi.NewWorld()
	entry := w.Entry(w.Create(component.GridPos, component.GemType, component.PixelPos, component.Tween, component.Sprite))

	cellType := [][]int{{2, 2}, {2, 2}}
	cells := [][]*donburi.Entry{{entry, entry}, {entry, entry}}
	phase := &component.PhaseData{
		Phase:       component.PhaseSelected,
		SelectedCol: 0,
		SelectedRow: 0,
	}

	// Tap adjacent tile → initiate swap
	intent := ResolveIntent(1, 0, phase, cellType, cells)
	swap, ok := intent.(InitiateSwap)
	if !ok {
		t.Fatalf("expected InitiateSwap, got %T", intent)
	}
	if swap.FromCol != 0 || swap.FromRow != 0 || swap.ToCol != 1 || swap.ToRow != 0 {
		t.Fatalf("unexpected swap coords: %+v", swap)
	}
}

func TestResolveIntent_SelectedPhaseDistantChangesSelection(t *testing.T) {
	w := donburi.NewWorld()
	entry := w.Entry(w.Create(component.GridPos, component.GemType, component.PixelPos, component.Tween, component.Sprite))

	cellType := [][]int{{2, 2, 2}, {2, 2, 2}, {2, 2, 2}}
	cells := [][]*donburi.Entry{{entry, entry, entry}, {entry, entry, entry}, {entry, entry, entry}}
	phase := &component.PhaseData{
		Phase:       component.PhaseSelected,
		SelectedCol: 0,
		SelectedRow: 0,
	}

	// Tap distant tile → change selection
	intent := ResolveIntent(2, 2, phase, cellType, cells)
	cs, ok := intent.(ChangeSelection)
	if !ok {
		t.Fatalf("expected ChangeSelection, got %T", intent)
	}
	if cs.Col != 2 || cs.Row != 2 {
		t.Fatalf("unexpected change selection: %+v", cs)
	}
}

func TestResolveIntent_SwappingPhaseReturnsNil(t *testing.T) {
	w := donburi.NewWorld()
	entry := w.Entry(w.Create(component.GridPos, component.GemType, component.PixelPos, component.Tween, component.Sprite))

	cellType := [][]int{{2, 2}}
	cells := [][]*donburi.Entry{{entry, entry}}
	phase := &component.PhaseData{Phase: component.PhaseSwapping}

	// During animation phases, classifyTap returns nil
	intent := ResolveIntent(0, 0, phase, cellType, cells)
	if intent != nil {
		t.Fatalf("expected nil during PhaseSwapping, got %T", intent)
	}
}
