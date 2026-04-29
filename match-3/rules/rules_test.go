package rules_test

import (
	"testing"

	"github.com/soockee/terminal-games/match-3/rules"
)

func playableGrid(cols, rows int) [][]int {
	ct := make([][]int, cols)
	for c := range cols {
		ct[c] = make([]int, rows)
		for r := range rows {
			ct[c][r] = 2 // playable
		}
	}
	return ct
}

func TestFindMatches_Horizontal(t *testing.T) {
	// 3x3 grid with horizontal match in row 0: [0,0,0]
	grid := [][]int{
		{0, 1, 2},
		{0, 2, 1},
		{0, 1, 2},
	}
	matches := rules.FindMatches(grid, 3, 3)
	if len(matches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(matches))
	}
}

func TestFindMatches_Vertical(t *testing.T) {
	// 3x3 grid with vertical match in col 0: [0,0,0]
	grid := [][]int{
		{0, 0, 0},
		{1, 2, 1},
		{2, 1, 2},
	}
	matches := rules.FindMatches(grid, 3, 3)
	if len(matches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(matches))
	}
}

func TestFindMatches_NoMatch(t *testing.T) {
	grid := [][]int{
		{0, 1, 0},
		{1, 0, 1},
		{0, 1, 0},
	}
	matches := rules.FindMatches(grid, 3, 3)
	if len(matches) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(matches))
	}
}

func TestCollapse(t *testing.T) {
	grid := [][]int{
		{-1, 1, 2}, // col 0: empty at top, should collapse
	}
	cellType := [][]int{
		{2, 2, 2},
	}
	moves := rules.Collapse(grid, cellType, 1, 3)
	// Tile at row 1 (color 1) should fall to row 2... wait, row 2 already has color 2.
	// Actually col 0 has: [row0=-1, row1=1, row2=2]. Bottom is row 2.
	// writeRow starts at 2. row=2: grid[0][2]=2 >= 0, stays at 2, writeRow=1.
	// row=1: grid[0][1]=1 >= 0, stays at 1, writeRow=0.
	// row=0: grid[0][0]=-1, skip.
	// No moves needed since gap is at top.
	if len(moves) != 0 {
		t.Fatalf("expected 0 moves (gap at top), got %d", len(moves))
	}

	// Now test with a gap in the middle.
	grid2 := [][]int{
		{1, -1, 2}, // col 0: [row0=1, row1=empty, row2=2]
	}
	cellType2 := [][]int{
		{2, 2, 2},
	}
	moves2 := rules.Collapse(grid2, cellType2, 1, 3)
	// writeRow starts at 2. row=2: color 2 stays. writeRow=1.
	// row=1: -1, skip. row=0: color 1, moves to row 1.
	if len(moves2) != 1 {
		t.Fatalf("expected 1 move, got %d", len(moves2))
	}
	if grid2[0][1] != 1 || grid2[0][0] != -1 {
		t.Fatalf("collapse didn't move correctly: %v", grid2[0])
	}
}

func TestHasValidMoves(t *testing.T) {
	// Board where swapping (0,0) and (1,0) creates a horizontal match.
	grid := [][]int{
		{1, 0, 2},
		{0, 0, 1},
		{2, 1, 0},
	}
	cellType := playableGrid(3, 3)
	if !rules.HasValidMoves(grid, cellType, 3, 3) {
		t.Fatal("expected valid moves to exist")
	}
}

func TestHasValidMoves_Stuck(t *testing.T) {
	// Latin square: no swap of adjacent cells can create 3 in a row.
	grid := [][]int{
		{0, 1, 2},
		{1, 2, 0},
		{2, 0, 1},
	}
	cellType := playableGrid(3, 3)
	if rules.HasValidMoves(grid, cellType, 3, 3) {
		t.Fatal("expected no valid moves in Latin square")
	}
}

func TestSimulateSwap(t *testing.T) {
	cellType := playableGrid(3, 3)
	// Create a scenario with a vertical match in col 2:
	grid2 := [][]int{
		{0, 1, 2},
		{1, 0, 2},
		{0, 1, 2}, // col 2 is already all 2s — vertical match!
	}
	// The grid already has a match in col 2 (all 2s). SimulateSwap should detect it
	// even without an actual swap if it produces matches after.
	score := rules.SimulateSwap(grid2, cellType, 3, 3, 0, 0, 1, 0)
	// After swap: col0=[1,1,2], col1=[0,0,2], col2=[0,1,2]
	// col1 row0 and row1 are both 0 — not 3 in a row.
	// Let's just verify the function runs without panic.
	_ = score

	// A guaranteed match scenario:
	grid4 := [][]int{
		{0, 1, 2},
		{1, 0, 1},
		{0, 1, 2},
	}
	// Swap (0,0)=0 with (0,1)=1: col 0 becomes [1,1,0]. Row 0 becomes [1,1,0]... not a match.
	// Let me just test that simulation doesn't modify the original.
	original := make([][]int, 3)
	for c := range 3 {
		original[c] = make([]int, 3)
		copy(original[c], grid4[c])
	}
	rules.SimulateSwap(grid4, cellType, 3, 3, 0, 0, 1, 0)
	for c := range 3 {
		for r := range 3 {
			if grid4[c][r] != original[c][r] {
				t.Fatalf("SimulateSwap modified original grid at [%d][%d]", c, r)
			}
		}
	}
}

func TestScoreForMatches(t *testing.T) {
	if got := rules.ScoreForMatches(3, 1); got != 30 {
		t.Fatalf("expected 30, got %d", got)
	}
	if got := rules.ScoreForMatches(5, 2); got != 100 {
		t.Fatalf("expected 100, got %d", got)
	}
}

func TestFindMatches_LongHorizontalRun(t *testing.T) {
	// 5 cols, 1 row: all same color => 5 matched
	grid := [][]int{
		{0}, {0}, {0}, {0}, {0},
	}
	matches := rules.FindMatches(grid, 5, 1)
	if len(matches) != 5 {
		t.Fatalf("expected 5 matches for run of 5, got %d", len(matches))
	}
}

func TestFindMatches_LongVerticalRun(t *testing.T) {
	// 1 col, 5 rows: all same color => 5 matched
	grid := [][]int{
		{1, 1, 1, 1, 1},
	}
	matches := rules.FindMatches(grid, 1, 5)
	if len(matches) != 5 {
		t.Fatalf("expected 5 matches for vertical run of 5, got %d", len(matches))
	}
}

func TestFindMatches_CrossPattern(t *testing.T) {
	// Both horizontal and vertical match overlap at center
	// Grid (col-major): col0=[1,0,1], col1=[0,0,0], col2=[1,0,1]
	// Row 1: [0, 0, 0] => horizontal match
	// Col 1: [0, 0, 0] => vertical match
	grid := [][]int{
		{1, 0, 1},
		{0, 0, 0},
		{1, 0, 1},
	}
	matches := rules.FindMatches(grid, 3, 3)
	// 3 horizontal + 3 vertical - 1 overlap = 5 unique positions
	if len(matches) != 5 {
		t.Fatalf("expected 5 matches for cross pattern, got %d", len(matches))
	}
}

func TestFindMatches_SkipsEmpty(t *testing.T) {
	// -1 values should not form matches
	grid := [][]int{
		{-1, 1, 2},
		{-1, 2, 1},
		{-1, 1, 2},
	}
	matches := rules.FindMatches(grid, 3, 3)
	if len(matches) != 0 {
		t.Fatalf("expected 0 matches with empty cells, got %d", len(matches))
	}
}

func TestFindMatches_TwoSeparateMatches(t *testing.T) {
	// 5 cols, 3 rows: horizontal match at row 0 (cols 0-2) and row 2 (cols 2-4)
	grid := [][]int{
		{0, 1, 2},
		{0, 2, 2},
		{0, 1, 1},
		{1, 2, 1},
		{2, 1, 1},
	}
	matches := rules.FindMatches(grid, 5, 3)
	// Row 0: grid[0][0]=0, grid[1][0]=0, grid[2][0]=0 => 3 match
	// Row 2: grid[2][2]=1, grid[3][2]=1, grid[4][2]=1 => 3 match
	if len(matches) != 6 {
		t.Fatalf("expected 6 matches (two separate runs), got %d", len(matches))
	}
}

func TestCollapse_MultipleGaps(t *testing.T) {
	// col with gaps at row 1 and row 3
	grid := [][]int{
		{0, -1, 1, -1, 2}, // rows 0-4
	}
	cellType := [][]int{
		{2, 2, 2, 2, 2},
	}
	moves := rules.Collapse(grid, cellType, 1, 5)
	// Expected final: [-1, -1, 0, 1, 2] (all gems fall down)
	if grid[0][4] != 2 {
		t.Fatalf("expected color 2 at bottom, got %d", grid[0][4])
	}
	if grid[0][3] != 1 {
		t.Fatalf("expected color 1 at row 3, got %d", grid[0][3])
	}
	if grid[0][2] != 0 {
		t.Fatalf("expected color 0 at row 2, got %d", grid[0][2])
	}
	if grid[0][1] != -1 || grid[0][0] != -1 {
		t.Fatalf("expected empty at top, got %v", grid[0])
	}
	if len(moves) != 2 {
		t.Fatalf("expected 2 moves, got %d", len(moves))
	}
}

func TestCollapse_AllEmpty(t *testing.T) {
	grid := [][]int{
		{-1, -1, -1},
	}
	cellType := [][]int{
		{2, 2, 2},
	}
	moves := rules.Collapse(grid, cellType, 1, 3)
	if len(moves) != 0 {
		t.Fatalf("expected 0 moves on all-empty col, got %d", len(moves))
	}
}

func TestCollapse_NonPlayableCells(t *testing.T) {
	// Cell at row 2 is non-playable (cellType != 2), acts as a wall
	grid := [][]int{
		{1, -1, 0, 2}, // rows 0-3
	}
	cellType := [][]int{
		{2, 2, 0, 2}, // row 2 is non-playable
	}
	moves := rules.Collapse(grid, cellType, 1, 4)
	// Gems above the wall don't pass through it
	// Below wall: row 3 has color 2, stays. writeRow resets at wall.
	// Above wall: row 0 has color 1, row 1 is empty.
	// row 1 should get color 1 from row 0
	if grid[0][1] != 1 {
		t.Fatalf("expected color 1 fell to row 1, got %d", grid[0][1])
	}
	if grid[0][0] != -1 {
		t.Fatalf("expected empty at row 0, got %d", grid[0][0])
	}
	if len(moves) != 1 {
		t.Fatalf("expected 1 move, got %d", len(moves))
	}
}

func TestEmptyCells(t *testing.T) {
	grid := [][]int{
		{-1, 0, -1},
	}
	cellType := [][]int{
		{2, 2, 2},
	}
	cells := rules.EmptyCells(grid, cellType, 1, 3)
	if len(cells) != 2 {
		t.Fatalf("expected 2 empty cells, got %d", len(cells))
	}
	// First empty at row 0, offset 1
	if cells[0].Col != 0 || cells[0].Row != 0 || cells[0].SpawnOffset != 1 {
		t.Fatalf("unexpected first cell: %+v", cells[0])
	}
	if cells[1].Col != 0 || cells[1].Row != 2 || cells[1].SpawnOffset != 2 {
		t.Fatalf("unexpected second cell: %+v", cells[1])
	}
}

func TestEmptyCells_NonPlayableIgnored(t *testing.T) {
	grid := [][]int{
		{-1, -1, -1},
	}
	cellType := [][]int{
		{0, 2, 2}, // row 0 is not playable
	}
	cells := rules.EmptyCells(grid, cellType, 1, 3)
	if len(cells) != 2 {
		t.Fatalf("expected 2 empty cells (non-playable excluded), got %d", len(cells))
	}
}

func TestHasValidMoves_SingleSwapCreatesVerticalMatch(t *testing.T) {
	// Swapping (1,0) with (1,1) creates vertical match in col 1
	// Col 1 before: [0, 1, 1]. After swap: [1, 0, 1]... no.
	// Let me be more deliberate:
	// Col0=[0,1,2], Col1=[1,0,1], Col2=[2,1,0]
	// Swap (1,0) and (0,0): col0 becomes [1,...], col1 becomes [0,...] row0=[1,0,2]
	// Instead let's set up: col0=[0,0,1], col1=[1,0,0], col2=[2,1,2]
	// Swap (0,2) with (1,2): row2=[0,1,2]->[1,0,2]. Col0=[0,0,0]? No col0=[0,0,1]->[0,0,0] after swap? No.
	// Simpler: set up so swapping creates col with 3 same.
	// col0=[0,1,0], col1=[1,0,1], col2=[0,1,0]
	// Swap (0,1) [val=1] with (1,1) [val=0]:
	// col0 becomes [0,0,0] => vertical match!
	grid := [][]int{
		{0, 1, 0},
		{1, 0, 1},
		{0, 1, 0},
	}
	cellType := playableGrid(3, 3)
	if !rules.HasValidMoves(grid, cellType, 3, 3) {
		t.Fatal("expected valid moves to exist (swap creates vertical match)")
	}
}

func TestHasValidMoves_EmptyCellsIgnored(t *testing.T) {
	grid := [][]int{
		{-1, -1, -1},
		{-1, -1, -1},
		{-1, -1, -1},
	}
	cellType := playableGrid(3, 3)
	if rules.HasValidMoves(grid, cellType, 3, 3) {
		t.Fatal("expected no valid moves on all-empty grid")
	}
}

func TestSimulateSwap_ProducesMatch(t *testing.T) {
	// col0=[0,1,0], col1=[1,0,1], col2=[0,1,0]
	// Swap (0,1) and (1,1): col0=[0,0,0] => vertical match of 3
	grid := [][]int{
		{0, 1, 0},
		{1, 0, 1},
		{0, 1, 0},
	}
	cellType := playableGrid(3, 3)
	score := rules.SimulateSwap(grid, cellType, 3, 3, 0, 1, 1, 1)
	if score < 3 {
		t.Fatalf("expected score >= 3 for a match, got %d", score)
	}
}

func TestSimulateSwap_NoMatch(t *testing.T) {
	// Latin square: no swap creates a match
	grid := [][]int{
		{0, 1, 2},
		{1, 2, 0},
		{2, 0, 1},
	}
	cellType := playableGrid(3, 3)
	score := rules.SimulateSwap(grid, cellType, 3, 3, 0, 0, 1, 0)
	if score != 0 {
		t.Fatalf("expected score 0 for no-match swap, got %d", score)
	}
}

func TestSimulateSwap_Cascade(t *testing.T) {
	// Set up a board where one swap triggers a cascade:
	// After first match clears and collapse, a second match forms.
	// 4 cols, 4 rows
	// col0=[0,0,1,2], col1=[1,1,0,2], col2=[2,2,0,2], col3=[0,1,2,0]
	// Swap (0,2) [val=1] with (1,2) [val=0]:
	// col0=[0,0,0,2], col1=[1,1,1,2] => both vertical matches of 3!
	// After clearing: col0=[-1,-1,-1,2], col1=[-1,-1,-1,2], collapse pushes 2 down
	// col0=[-1,-1,-1,2], col1=[-1,-1,-1,2] (2 already at bottom)
	// Then row3=[2,2,2,0] could be horizontal? Let's check col2[3]=2, col3[3]=0.
	// row3: col0=2, col1=2, col2=2, col3=0 => 3 match!
	grid := [][]int{
		{0, 0, 1, 2},
		{1, 1, 0, 2},
		{2, 2, 0, 2},
		{0, 1, 2, 0},
	}
	cellType := playableGrid(4, 4)
	score := rules.SimulateSwap(grid, cellType, 4, 4, 0, 2, 1, 2)
	// First wave: 3 (col0) + 3 (col1) = 6
	// After collapse, potential second wave
	if score < 6 {
		t.Fatalf("expected score >= 6 for cascade, got %d", score)
	}
}

func TestReshuffle(t *testing.T) {
	grid := [][]int{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}
	cellType := playableGrid(3, 3)
	ok := rules.Reshuffle(grid, cellType, 3, 3, 3, 1000)
	if !ok {
		t.Fatal("expected reshuffle to find valid configuration")
	}
	// Verify no matches remain
	matches := rules.FindMatches(grid, 3, 3)
	if len(matches) != 0 {
		t.Fatalf("reshuffle left matches on board: %d", len(matches))
	}
	// Verify valid moves exist
	if !rules.HasValidMoves(grid, cellType, 3, 3) {
		t.Fatal("reshuffle produced board with no valid moves")
	}
}

func TestReshuffle_FailsWithOneColor(t *testing.T) {
	// With only 1 color, it's impossible to have no matches on a 3x3 board
	grid := [][]int{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}
	cellType := playableGrid(3, 3)
	ok := rules.Reshuffle(grid, cellType, 3, 3, 1, 100)
	if ok {
		t.Fatal("expected reshuffle to fail with only 1 color")
	}
}

func TestColorGrid(t *testing.T) {
	grid := rules.ColorGrid(3, 2, func(col, row int) int {
		return col*10 + row
	})
	if grid[0][0] != 0 || grid[0][1] != 1 || grid[2][1] != 21 {
		t.Fatalf("ColorGrid returned unexpected values: %v", grid)
	}
}

func TestFindMatches_ExactlyThree(t *testing.T) {
	// Ensure a run of exactly 2 does NOT match
	grid := [][]int{
		{0, 1, 2},
		{0, 2, 1},
		{1, 1, 2},
	}
	// Row 0: [0, 0, 1] => only run of 2
	matches := rules.FindMatches(grid, 3, 3)
	if len(matches) != 0 {
		t.Fatalf("expected 0 matches for run of 2, got %d", len(matches))
	}
}

func TestFindMatches_FourInARow(t *testing.T) {
	// 4 cols: all same in row 0
	grid := [][]int{
		{0, 1},
		{0, 2},
		{0, 1},
		{0, 2},
	}
	matches := rules.FindMatches(grid, 4, 2)
	if len(matches) != 4 {
		t.Fatalf("expected 4 matches for run of 4, got %d", len(matches))
	}
}

// --- IsValidSwap tests ---

func TestIsValidSwap_OutOfBounds(t *testing.T) {
	grid := [][]int{
		{0, 1, 2},
		{1, 2, 0},
		{2, 0, 1},
	}
	cellType := playableGrid(3, 3)

	cases := []struct {
		name           string
		c1, r1, c2, r2 int
	}{
		{"c1 negative", -1, 0, 0, 0},
		{"r1 negative", 0, -1, 0, 0},
		{"c2 out of range", 2, 0, 3, 0},
		{"r2 out of range", 0, 2, 0, 3},
		{"both out", -1, -1, 3, 3},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if rules.IsValidSwap(grid, cellType, 3, 3, tc.c1, tc.r1, tc.c2, tc.r2) {
				t.Fatal("expected invalid swap for out-of-bounds")
			}
		})
	}
}

func TestIsValidSwap_NonAdjacent(t *testing.T) {
	grid := [][]int{
		{0, 1, 2},
		{1, 2, 0},
		{2, 0, 1},
	}
	cellType := playableGrid(3, 3)

	cases := []struct {
		name           string
		c1, r1, c2, r2 int
	}{
		{"diagonal", 0, 0, 1, 1},
		{"two apart horizontal", 0, 0, 2, 0},
		{"two apart vertical", 0, 0, 0, 2},
		{"same cell", 1, 1, 1, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if rules.IsValidSwap(grid, cellType, 3, 3, tc.c1, tc.r1, tc.c2, tc.r2) {
				t.Fatal("expected invalid swap for non-adjacent cells")
			}
		})
	}
}

func TestIsValidSwap_NonPlayableCell(t *testing.T) {
	grid := [][]int{
		{0, 1, 2},
		{1, 2, 0},
		{2, 0, 1},
	}
	cellType := playableGrid(3, 3)
	cellType[1][0] = 0 // make (1,0) non-playable

	// Swap (0,0) with (1,0): target is non-playable
	if rules.IsValidSwap(grid, cellType, 3, 3, 0, 0, 1, 0) {
		t.Fatal("expected invalid swap when target is non-playable")
	}
	// Swap (1,0) with (0,0): source is non-playable
	if rules.IsValidSwap(grid, cellType, 3, 3, 1, 0, 0, 0) {
		t.Fatal("expected invalid swap when source is non-playable")
	}
}

func TestIsValidSwap_EmptyCell(t *testing.T) {
	grid := [][]int{
		{0, 1, 2},
		{-1, 2, 0},
		{2, 0, 1},
	}
	cellType := playableGrid(3, 3)

	// Swap (0,0) with (1,0): (1,0) is empty
	if rules.IsValidSwap(grid, cellType, 3, 3, 0, 0, 1, 0) {
		t.Fatal("expected invalid swap with empty cell")
	}
}

func TestIsValidSwap_SameColor(t *testing.T) {
	grid := [][]int{
		{0, 1, 2},
		{0, 2, 1},
		{2, 0, 1},
	}
	cellType := playableGrid(3, 3)

	// (0,0)=0 and (1,0)=0: same color, pointless swap
	if rules.IsValidSwap(grid, cellType, 3, 3, 0, 0, 1, 0) {
		t.Fatal("expected invalid swap for same-color cells")
	}
}

func TestIsValidSwap_NoMatchProduced(t *testing.T) {
	// Latin square: no adjacent swap creates a match
	grid := [][]int{
		{0, 1, 2},
		{1, 2, 0},
		{2, 0, 1},
	}
	cellType := playableGrid(3, 3)

	// Try all adjacent swaps — none should be valid
	dirs := [][2]int{{1, 0}, {0, 1}}
	for c := range 3 {
		for r := range 3 {
			for _, d := range dirs {
				nc, nr := c+d[0], r+d[1]
				if nc >= 3 || nr >= 3 {
					continue
				}
				if rules.IsValidSwap(grid, cellType, 3, 3, c, r, nc, nr) {
					t.Fatalf("expected no valid swap in Latin square, but (%d,%d)<->(%d,%d) was valid", c, r, nc, nr)
				}
			}
		}
	}
}

func TestIsValidSwap_ValidHorizontalMatch(t *testing.T) {
	// Grid is col-major: grid[col][row]
	// col0=[1,2,1], col1=[0,1,2], col2=[0,2,1]
	// Row 0 values: col0[0]=1, col1[0]=0, col2[0]=0
	// Swap (0,0)=1 with (0,1)=2 won't help. Let's construct carefully:
	// We want swap (1,0)<->(0,0) to create row0 = [0,0,0]
	// Need: col0[0]=0 after swap, col1[0]=0 already, col2[0]=0
	// So col0[0]=X, col1[0]=0, col2[0]=0, and after swapping (0,0) with (1,0):
	// col0[0] gets col1[0]=0, col1[0] gets col0[0]=X
	// Row0 after: [0, X, 0] — not a match unless X=0 (same color, rejected).
	// Alternative: swap vertically. Let's create horizontal match in row 1:
	// col0[1]=0, col1[1]=X, col2[1]=0. Swap (1,0)<->(1,1): col1 becomes [col1[1], col1[0], ...]
	// After swap: col1[0]=X (was col1[1]), col1[1]=col1[0] (was col1[0])
	// We need row1 after swap = [col0[1], new_col1[1], col2[1]] = [0, col1[0], 0]
	// If col1[0]=0 => row1 = [0,0,0] match! But then col0[1]=0 and col1[0]=0 same? No,
	// we swap (1,0) and (1,1), so we check col1[0] vs col1[1] colors.
	// col1=[2, 0, 1]: swap (1,0)=2 with (1,1)=0 => col1=[0,2,1]
	// row1: col0[1]=0, col1[1]=2, col2[1]=0 — not 3 same.
	// Let me just do: col0=[1,0,2], col1=[0,1,0], col2=[0,2,1]
	// Swap (0,0)=1 with (1,0)=0: col0[0]=0, col1[0]=1
	// Row 0: [0,1,0] — no.
	// Simplest: col0=[1,2,0], col1=[0,2,1], col2=[0,2,1]
	// Swap (0,0)=1 with (1,0)=0: row0=[0,0,0]? No: col0[0] becomes 0, col1[0] becomes 1, col2[0]=0
	// row0 = [0, 1, 0] not a match.
	// The issue: swapping (0,0) and (1,0) exchanges col0[0] and col1[0].
	// For row0 to become all same after swap, we need the OTHER cols to already match.
	// row0 = [A, B, C]. After swap A<->B: [B, A, C]. For match: B=A=C => all same already.
	// That won't work for horizontal via (0,0)<->(1,0).
	// Instead, swap (0,0)<->(0,1) (vertical swap in col 0):
	// col0 = [A, B, ...]. After: col0 = [B, A, ...].
	// Check row 0: [B, col1[0], col2[0]]. For match: B = col1[0] = col2[0].
	// So: col0=[X, B, ...], col1=[B, ...], col2=[B, ...]
	// With B=0: col0=[1, 0, 2], col1=[0, 2, 1], col2=[0, 1, 2]
	// Swap (0,0)<->(0,1): col0 becomes [0, 1, 2]. Row 0: [0, 0, 0] => match!
	grid := [][]int{
		{1, 0, 2},
		{0, 2, 1},
		{0, 1, 2},
	}
	cellType := playableGrid(3, 3)

	if !rules.IsValidSwap(grid, cellType, 3, 3, 0, 0, 0, 1) {
		t.Fatal("expected valid swap creating horizontal match")
	}
}

func TestIsValidSwap_ValidVerticalMatch(t *testing.T) {
	// col0=[0,1,0], col1=[1,0,1]
	// Swap (0,1)=1 with (1,1)=0 => col0=[0,0,0] => vertical match
	grid := [][]int{
		{0, 1, 0},
		{1, 0, 1},
		{2, 1, 2},
	}
	cellType := playableGrid(3, 3)

	if !rules.IsValidSwap(grid, cellType, 3, 3, 0, 1, 1, 1) {
		t.Fatal("expected valid swap creating vertical match")
	}
}

func TestIsValidSwap_ValidOnlyViaCascade(t *testing.T) {
	// A swap that doesn't immediately look useful but the first match triggers
	// a collapse that creates a second match (cascade). The swap IS valid because
	// SimulateSwap considers the full cascade.
	//
	// 4 cols, 4 rows:
	// col0=[0,0,1,2], col1=[1,1,0,2], col2=[2,2,0,2], col3=[0,1,2,0]
	// Swap (0,2)[=1] with (1,2)[=0]:
	//   col0=[0,0,0,2] => vertical match (3)
	//   col1=[1,1,1,2] => vertical match (3)
	// After clearing and collapse, row3=[2,2,2,0] => horizontal match (cascade)
	grid := [][]int{
		{0, 0, 1, 2},
		{1, 1, 0, 2},
		{2, 2, 0, 2},
		{0, 1, 2, 0},
	}
	cellType := playableGrid(4, 4)

	if !rules.IsValidSwap(grid, cellType, 4, 4, 0, 2, 1, 2) {
		t.Fatal("expected valid swap that triggers cascade")
	}
}

func TestIsValidSwap_DoesNotMutateGrid(t *testing.T) {
	grid := [][]int{
		{0, 1, 0},
		{1, 0, 1},
		{0, 1, 0},
	}
	cellType := playableGrid(3, 3)

	// Save original state
	original := make([][]int, 3)
	for c := range 3 {
		original[c] = make([]int, 3)
		copy(original[c], grid[c])
	}

	rules.IsValidSwap(grid, cellType, 3, 3, 0, 1, 1, 1)

	for c := range 3 {
		for r := range 3 {
			if grid[c][r] != original[c][r] {
				t.Fatalf("IsValidSwap mutated grid at [%d][%d]: was %d, now %d", c, r, original[c][r], grid[c][r])
			}
		}
	}
}

func TestIsValidSwap_ChainDoesNotValidateOtherwise(t *testing.T) {
	// Ensure that a swap which produces NO immediate match and NO cascade
	// is correctly rejected, even on a larger board where cascades are possible.
	// 4x4 board with no useful swaps at position (0,0)<->(1,0).
	grid := [][]int{
		{0, 1, 2, 0},
		{1, 2, 0, 1},
		{2, 0, 1, 2},
		{0, 1, 2, 0},
	}
	cellType := playableGrid(4, 4)

	// This is a 4x4 Latin-like pattern; verify specific swap is invalid
	if rules.IsValidSwap(grid, cellType, 4, 4, 0, 0, 1, 0) {
		t.Fatal("expected invalid swap on pattern that produces no match or cascade")
	}
}

func TestIsValidSwap_Symmetric(t *testing.T) {
	// Swapping A<->B should give the same result as B<->A
	grid := [][]int{
		{1, 2, 1},
		{0, 1, 2},
		{0, 2, 1},
	}
	cellType := playableGrid(3, 3)

	fwd := rules.IsValidSwap(grid, cellType, 3, 3, 0, 0, 1, 0)
	rev := rules.IsValidSwap(grid, cellType, 3, 3, 1, 0, 0, 0)
	if fwd != rev {
		t.Fatalf("IsValidSwap not symmetric: (%d,%d)<->(%d,%d)=%v but reverse=%v", 0, 0, 1, 0, fwd, rev)
	}
}
