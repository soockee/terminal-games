// Command levelgen generates new match-3 level .ldtkl files and updates the
// main match-3.ldtk project file to reference them.
//
// Usage:
//
//	go run ./cmd/levelgen -count 3 -startDifficulty 2
//
// This inserts levels before the Win_screen entry.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

// Grid constants: the LDtk project uses a 17×12 cell grid at 32px/cell.
// Viewport is 514×360 pixels → visible area is cols 0–15, rows 0–10.
const (
	gridCols = 17
	gridRows = 12

	// Playable area must stay within viewport (514×360 at 32px tiles).
	maxVisibleCol = 15 // 16*32 = 512 ≤ 514
	maxVisibleRow = 10 // 11*32 = 352 ≤ 360
)

// Board shape generators produce a 12×17 (rows×cols) grid of cell values.
// All shapes MUST keep playable cells within the visible viewport.
type shapeFunc func(difficulty int) []int

var shapes = []shapeFunc{
	shapeRect,
	shapeDiamond,
	shapePlus,
	shapeT,
	shapeU,
	shapeH,
	shapeSteps,
	shapeArrow,
	shapeRing,
	shapeWideRect,
}

func main() {
	count := flag.Int("count", 3, "number of levels to generate")
	startDiff := flag.Int("startDifficulty", 2, "starting difficulty (1=easy, 5=hard)")
	projectPath := flag.String("project", "assets/ldtk/match-3.ldtk", "path to the LDtk project file")
	flag.Parse()

	if *count < 1 {
		log.Fatal("count must be >= 1")
	}

	absProject, err := filepath.Abs(*projectPath)
	if err != nil {
		log.Fatalf("resolving project path: %v", err)
	}

	// Read and parse the project file.
	projectData, err := os.ReadFile(absProject)
	if err != nil {
		log.Fatalf("reading project: %v", err)
	}

	var project map[string]interface{}
	if err := json.Unmarshal(projectData, &project); err != nil {
		log.Fatalf("parsing project JSON: %v", err)
	}

	levels := project["levels"].([]interface{})
	nextUid := int(project["nextUid"].(float64))

	// Find Win_screen index (we insert before it).
	winIdx := -1
	for i, l := range levels {
		lm := l.(map[string]interface{})
		if lm["identifier"].(string) == "Win_screen" {
			winIdx = i
			break
		}
	}
	if winIdx < 0 {
		log.Fatal("Win_screen level not found in project")
	}

	// Determine the next level number from existing levels.
	maxLevelNum := -1
	for _, l := range levels {
		lm := l.(map[string]interface{})
		id := lm["identifier"].(string)
		if strings.HasPrefix(id, "Level_") {
			var n int
			if _, err := fmt.Sscanf(id, "Level_%d", &n); err == nil && n > maxLevelNum {
				maxLevelNum = n
			}
		}
	}

	// Determine worldX for new levels (each level is 544px apart).
	const levelSpacing = 544
	lastPlayableWorldX := 0
	for i := 0; i < winIdx; i++ {
		lm := levels[i].(map[string]interface{})
		wx := int(lm["worldX"].(float64))
		if wx > lastPlayableWorldX {
			lastPlayableWorldX = wx
		}
	}

	ldtklDir := filepath.Join(filepath.Dir(absProject), "match-3")

	// Generate levels.
	newLevels := make([]map[string]interface{}, 0, *count)
	newFiles := make([]string, 0, *count)

	for i := 0; i < *count; i++ {
		levelNum := maxLevelNum + 1 + i
		identifier := fmt.Sprintf("Level_%d", levelNum)
		uid := nextUid + i
		worldX := lastPlayableWorldX + levelSpacing*(i+1)

		// Difficulty ramps up (clamped 1–5).
		diff := *startDiff + i
		if diff > 5 {
			diff = 5
		}

		numColors, scoreTarget, timeLimit := difficultyParams(diff)

		// Pick a shape — distribute shapes so adjacent levels feel different.
		shapeIdx := (levelNum + i) % len(shapes)
		grid := shapes[shapeIdx](diff)

		// Validate: ensure grid is within viewport and has enough playable cells.
		grid = enforceViewport(grid)
		playable := countPlayable(grid)
		if playable < 20 {
			grid = shapeRect(diff)
			grid = enforceViewport(grid)
			playable = countPlayable(grid)
		}

		// Adjust score target based on actual playable area.
		scoreTarget = adjustScoreTarget(scoreTarget, playable, diff)

		iid := generateIID()
		layerIID := generateIID()

		// Build the .ldtkl content.
		ldtkl := buildLdtkl(identifier, iid, uid, worldX, numColors, scoreTarget, timeLimit, grid, layerIID)

		// Write the .ldtkl file.
		filename := identifier + ".ldtkl"
		filePath := filepath.Join(ldtklDir, filename)
		data, err := json.MarshalIndent(ldtkl, "", "\t")
		if err != nil {
			log.Fatalf("marshaling %s: %v", identifier, err)
		}
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			log.Fatalf("writing %s: %v", filePath, err)
		}
		fmt.Printf("Created %s\n", filePath)
		newFiles = append(newFiles, filename)

		// Write simplified export (data.json + Board.csv).
		writeSimplified(ldtklDir, identifier, iid, worldX, numColors, scoreTarget, timeLimit, grid)

		// Build the project-level entry.
		entry := buildProjectLevelEntry(identifier, iid, uid, worldX, numColors, scoreTarget, timeLimit, filename)
		newLevels = append(newLevels, entry)

		fmt.Printf("  Level %d: %s shape, %d playable, %d colors, target=%d, time=%ds\n",
			levelNum, shapeName(shapeIdx), playable, numColors, scoreTarget, timeLimit)
	}

	// Insert new levels before Win_screen in the project.
	updatedLevels := make([]interface{}, 0, len(levels)+len(newLevels))
	updatedLevels = append(updatedLevels, levels[:winIdx]...)
	for _, nl := range newLevels {
		updatedLevels = append(updatedLevels, nl)
	}
	// Update Win_screen worldX.
	winLevel := levels[winIdx].(map[string]interface{})
	winLevel["worldX"] = float64(lastPlayableWorldX + levelSpacing*(*count+1))
	updatedLevels = append(updatedLevels, levels[winIdx:]...)

	project["levels"] = updatedLevels
	project["nextUid"] = float64(nextUid + *count)

	// Write updated project file.
	out, err := json.MarshalIndent(project, "", "\t")
	if err != nil {
		log.Fatalf("marshaling project: %v", err)
	}
	if err := os.WriteFile(absProject, out, 0644); err != nil {
		log.Fatalf("writing project: %v", err)
	}
	fmt.Printf("\nUpdated %s (added %d levels)\n", absProject, *count)
}

func shapeName(idx int) string {
	names := []string{"rect", "diamond", "plus", "T", "U", "H", "steps", "arrow", "ring", "wideRect"}
	if idx < len(names) {
		return names[idx]
	}
	return "unknown"
}

// difficultyParams returns game params based on difficulty 1–5.
// Score targets are base values adjusted later by board size.
func difficultyParams(diff int) (numColors, scoreTarget, timeLimit int) {
	switch diff {
	case 1:
		return 4, 400, 0
	case 2:
		return 5, 800, 0
	case 3:
		return 5, 1500, 0
	case 4:
		return 6, 2000, 90
	case 5:
		return 6, 3000, 60
	default:
		return 5, 1000, 0
	}
}

// adjustScoreTarget scales the target based on playable cell count.
// Ensures the target is fun: achievable in a few minutes of active play.
func adjustScoreTarget(baseTarget, playable, diff int) int {
	// Each match of 3 scores 30 pts base. Cascades multiply.
	// Scale target so it takes roughly 2–5 minutes of play.
	scale := float64(playable) / 64.0
	target := int(float64(baseTarget) * scale)

	minTarget := playable * 5
	maxTarget := playable * 35

	if target < minTarget {
		target = minTarget
	}
	if target > maxTarget {
		target = maxTarget
	}
	// Round to nearest 50.
	target = ((target + 25) / 50) * 50
	if target < 100 {
		target = 100
	}
	return target
}

func countPlayable(grid []int) int {
	n := 0
	for _, v := range grid {
		if v == 2 {
			n++
		}
	}
	return n
}

// enforceViewport zeroes out any playable/blocked cells outside the visible area.
func enforceViewport(grid []int) []int {
	for r := 0; r < gridRows; r++ {
		for c := 0; c < gridCols; c++ {
			if c > maxVisibleCol || r > maxVisibleRow {
				grid[r*gridCols+c] = 0
			}
		}
	}
	return grid
}

// --- Shape generators ---
// All produce a flat [gridRows*gridCols] array (row-major).
// Playable cells are centered within the visible viewport.

func makeGrid() []int {
	return make([]int, gridCols*gridRows)
}

func fillRect(grid []int, startCol, startRow, endCol, endRow int) {
	for r := startRow; r <= endRow && r <= maxVisibleRow; r++ {
		for c := startCol; c <= endCol && c <= maxVisibleCol; c++ {
			if r >= 0 && c >= 0 {
				grid[r*gridCols+c] = 2
			}
		}
	}
}

// addStrategicBlockers adds blockers in interesting patterns (not random noise).
func addStrategicBlockers(grid []int, difficulty int) {
	if difficulty <= 1 {
		return
	}

	// Number of blockers scales with difficulty.
	numBlockers := (difficulty - 1) * 2
	playable := countPlayable(grid)
	if numBlockers > playable/6 {
		numBlockers = playable / 6
	}

	patterns := []func([]int, int){
		placeBlockerPair,
		placeBlockerLine,
		placeBlockerCorners,
	}

	for b := 0; b < numBlockers; b++ {
		patterns[b%len(patterns)](grid, b)
	}
}

func placeBlockerPair(grid []int, seed int) {
	minC, maxC, minR, maxR := playableBounds(grid)
	cx := (minC + maxC) / 2
	cy := (minR + maxR) / 2

	offsets := [][2]int{{1, 0}, {2, 0}, {0, 1}, {0, 2}, {1, 1}, {2, 1}}
	off := offsets[seed%len(offsets)]

	candidates := [][2]int{
		{cx + off[0], cy + off[1]},
		{cx - off[0], cy - off[1]},
	}

	for _, pos := range candidates {
		c, r := pos[0], pos[1]
		if canBlock(grid, c, r) {
			grid[r*gridCols+c] = 3
		}
	}
}

func placeBlockerLine(grid []int, seed int) {
	minC, maxC, minR, maxR := playableBounds(grid)
	if maxC-minC < 5 || maxR-minR < 5 {
		return
	}

	if seed%2 == 0 {
		r := minR + 2 + seed%(maxR-minR-3)
		if r > maxR-2 {
			r = (minR + maxR) / 2
		}
		startC := minC + 2 + seed%(maxC-minC-4)
		for c := startC; c < startC+2 && c <= maxC-2; c++ {
			if canBlock(grid, c, r) {
				grid[r*gridCols+c] = 3
			}
		}
	} else {
		c := minC + 2 + seed%(maxC-minC-3)
		if c > maxC-2 {
			c = (minC + maxC) / 2
		}
		startR := minR + 2 + seed%(maxR-minR-4)
		for r := startR; r < startR+2 && r <= maxR-2; r++ {
			if canBlock(grid, c, r) {
				grid[r*gridCols+c] = 3
			}
		}
	}
}

func placeBlockerCorners(grid []int, seed int) {
	minC, maxC, minR, maxR := playableBounds(grid)
	corners := [][2]int{
		{minC + 1, minR + 1},
		{maxC - 1, minR + 1},
		{minC + 1, maxR - 1},
		{maxC - 1, maxR - 1},
	}
	corner := corners[seed%len(corners)]
	if canBlock(grid, corner[0], corner[1]) {
		grid[corner[1]*gridCols+corner[0]] = 3
	}
}

func canBlock(grid []int, c, r int) bool {
	if c < 0 || c >= gridCols || r < 0 || r >= gridRows {
		return false
	}
	idx := r*gridCols + c
	if grid[idx] != 2 {
		return false
	}

	grid[idx] = 3
	colOK := longestRunInCol(grid, c) >= 3
	rowOK := longestRunInRow(grid, r) >= 3
	grid[idx] = 2

	return colOK && rowOK
}

func longestRunInCol(grid []int, c int) int {
	maxRun, run := 0, 0
	for r := 0; r < gridRows; r++ {
		if grid[r*gridCols+c] == 2 {
			run++
		} else {
			if run > maxRun {
				maxRun = run
			}
			run = 0
		}
	}
	if run > maxRun {
		maxRun = run
	}
	return maxRun
}

func longestRunInRow(grid []int, r int) int {
	maxRun, run := 0, 0
	for c := 0; c < gridCols; c++ {
		if grid[r*gridCols+c] == 2 {
			run++
		} else {
			if run > maxRun {
				maxRun = run
			}
			run = 0
		}
	}
	if run > maxRun {
		maxRun = run
	}
	return maxRun
}

func playableBounds(grid []int) (minC, maxC, minR, maxR int) {
	minC, minR = gridCols, gridRows
	maxC, maxR = 0, 0
	for r := 0; r < gridRows; r++ {
		for c := 0; c < gridCols; c++ {
			if grid[r*gridCols+c] == 2 {
				if c < minC {
					minC = c
				}
				if c > maxC {
					maxC = c
				}
				if r < minR {
					minR = r
				}
				if r > maxR {
					maxR = r
				}
			}
		}
	}
	return
}

// --- Shape implementations ---

// shapeRect: classic centered rectangle. Good cascades, easy to read.
func shapeRect(difficulty int) []int {
	grid := makeGrid()
	var w, h int
	switch difficulty {
	case 1:
		w, h = 6, 6
	case 2:
		w, h = 7, 7
	case 3:
		w, h = 8, 8
	case 4:
		w, h = 8, 9
	default:
		w, h = 9, 9
	}
	startC := (maxVisibleCol + 1 - w) / 2
	startR := (maxVisibleRow + 1 - h) / 2
	fillRect(grid, startC, startR, startC+w-1, startR+h-1)
	addStrategicBlockers(grid, difficulty)
	return grid
}

// shapeDiamond: diamond centered in viewport. Interesting cascade patterns.
func shapeDiamond(difficulty int) []int {
	grid := makeGrid()
	var radius int
	switch {
	case difficulty <= 2:
		radius = 3
	case difficulty <= 4:
		radius = 4
	default:
		radius = 5
	}

	cx := (maxVisibleCol + 1) / 2
	cy := (maxVisibleRow + 1) / 2

	for r := 0; r <= maxVisibleRow; r++ {
		for c := 0; c <= maxVisibleCol; c++ {
			if abs(c-cx)+abs(r-cy) <= radius {
				grid[r*gridCols+c] = 2
			}
		}
	}
	addStrategicBlockers(grid, difficulty)
	return grid
}

// shapePlus: plus/cross shape. Matches cross at the intersection.
func shapePlus(difficulty int) []int {
	grid := makeGrid()
	cx := (maxVisibleCol + 1) / 2
	cy := (maxVisibleRow + 1) / 2

	var armLen, armWidth int
	switch {
	case difficulty <= 2:
		armLen, armWidth = 3, 2
	case difficulty <= 4:
		armLen, armWidth = 4, 2
	default:
		armLen, armWidth = 4, 3
	}

	halfW := armWidth / 2

	// Vertical arm.
	fillRect(grid, cx-halfW, cy-armLen, cx+halfW, cy+armLen)
	// Horizontal arm.
	fillRect(grid, cx-armLen, cy-halfW, cx+armLen, cy+halfW)

	addStrategicBlockers(grid, difficulty)
	return grid
}

// shapeT: T-shape. Wide top bar with narrow stem — cascades funnel down.
func shapeT(difficulty int) []int {
	grid := makeGrid()
	cx := (maxVisibleCol + 1) / 2

	var topWidth, topHeight, stemWidth, stemHeight int
	switch {
	case difficulty <= 2:
		topWidth, topHeight = 8, 3
		stemWidth, stemHeight = 3, 4
	case difficulty <= 4:
		topWidth, topHeight = 9, 3
		stemWidth, stemHeight = 3, 5
	default:
		topWidth, topHeight = 10, 3
		stemWidth, stemHeight = 4, 5
	}

	topLeft := cx - topWidth/2
	fillRect(grid, topLeft, 1, topLeft+topWidth-1, topHeight)

	stemLeft := cx - stemWidth/2
	fillRect(grid, stemLeft, topHeight+1, stemLeft+stemWidth-1, topHeight+stemHeight)

	addStrategicBlockers(grid, difficulty)
	return grid
}

// shapeU: U-shape. Two columns connected at bottom — gravity pulls to base.
func shapeU(difficulty int) []int {
	grid := makeGrid()
	cx := (maxVisibleCol + 1) / 2

	var outerWidth, height, wallWidth, baseHeight int
	switch {
	case difficulty <= 2:
		outerWidth, height = 8, 7
		wallWidth, baseHeight = 2, 2
	case difficulty <= 4:
		outerWidth, height = 9, 8
		wallWidth, baseHeight = 2, 2
	default:
		outerWidth, height = 10, 9
		wallWidth, baseHeight = 2, 3
	}

	left := cx - outerWidth/2
	startR := (maxVisibleRow + 1 - height) / 2

	fillRect(grid, left, startR, left+wallWidth-1, startR+height-1)
	fillRect(grid, left+outerWidth-wallWidth, startR, left+outerWidth-1, startR+height-1)
	fillRect(grid, left, startR+height-baseHeight, left+outerWidth-1, startR+height-1)

	addStrategicBlockers(grid, difficulty)
	return grid
}

// shapeH: H-shape. Two columns connected by middle crossbar.
func shapeH(difficulty int) []int {
	grid := makeGrid()
	cx := (maxVisibleCol + 1) / 2
	cy := (maxVisibleRow + 1) / 2

	var outerWidth, height, wallWidth, barHeight int
	switch {
	case difficulty <= 2:
		outerWidth, height = 8, 7
		wallWidth, barHeight = 2, 2
	case difficulty <= 4:
		outerWidth, height = 9, 8
		wallWidth, barHeight = 2, 2
	default:
		outerWidth, height = 10, 9
		wallWidth, barHeight = 3, 3
	}

	left := cx - outerWidth/2
	startR := cy - height/2

	fillRect(grid, left, startR, left+wallWidth-1, startR+height-1)
	fillRect(grid, left+outerWidth-wallWidth, startR, left+outerWidth-1, startR+height-1)
	fillRect(grid, left, cy-barHeight/2, left+outerWidth-1, cy+barHeight/2)

	addStrategicBlockers(grid, difficulty)
	return grid
}

// shapeSteps: Staircase pattern. Asymmetric cascades.
func shapeSteps(difficulty int) []int {
	grid := makeGrid()

	var stepW, stepH, numSteps int
	switch {
	case difficulty <= 2:
		stepW, stepH, numSteps = 3, 3, 3
	case difficulty <= 4:
		stepW, stepH, numSteps = 3, 2, 4
	default:
		stepW, stepH, numSteps = 3, 2, 5
	}

	totalW := stepW * numSteps
	totalH := stepH * numSteps
	startC := (maxVisibleCol + 1 - totalW) / 2
	startR := (maxVisibleRow + 1 - totalH) / 2

	for s := 0; s < numSteps; s++ {
		sc := startC + s*stepW
		sr := startR + s*stepH
		// Each step overlaps with next vertically for column connectivity.
		endR := sr + stepH + stepH - 1
		if endR > startR+totalH-1 {
			endR = startR + totalH - 1
		}
		fillRect(grid, sc, sr, sc+stepW-1, endR)
	}

	addStrategicBlockers(grid, difficulty)
	return grid
}

// shapeArrow: Arrow pointing down. Tapers to a point — funnel effect.
func shapeArrow(difficulty int) []int {
	grid := makeGrid()
	cx := (maxVisibleCol + 1) / 2
	cy := (maxVisibleRow + 1) / 2

	var baseWidth, height int
	switch {
	case difficulty <= 2:
		baseWidth, height = 6, 7
	case difficulty <= 4:
		baseWidth, height = 7, 8
	default:
		baseWidth, height = 8, 9
	}

	startR := cy - height/2

	// Top rectangle (shaft).
	shaftH := height / 2
	fillRect(grid, cx-baseWidth/4, startR, cx+baseWidth/4, startR+shaftH-1)

	// Arrow head (triangle pointing down).
	for r := startR + shaftH; r < startR+height; r++ {
		dist := r - (startR + shaftH)
		halfW := baseWidth/2 - dist
		if halfW < 0 {
			halfW = 0
		}
		fillRect(grid, cx-halfW, r, cx+halfW, r)
	}

	addStrategicBlockers(grid, difficulty)
	return grid
}

// shapeRing: Hollow rectangle frame — matches happen along the edges.
func shapeRing(difficulty int) []int {
	grid := makeGrid()
	cx := (maxVisibleCol + 1) / 2
	cy := (maxVisibleRow + 1) / 2

	var outerW, outerH, thickness int
	switch {
	case difficulty <= 2:
		outerW, outerH, thickness = 7, 7, 2
	case difficulty <= 4:
		outerW, outerH, thickness = 8, 8, 2
	default:
		outerW, outerH, thickness = 9, 9, 2
	}

	left := cx - outerW/2
	top := cy - outerH/2

	fillRect(grid, left, top, left+outerW-1, top+outerH-1)

	// Clear inner rect.
	innerLeft := left + thickness
	innerTop := top + thickness
	innerRight := left + outerW - 1 - thickness
	innerBottom := top + outerH - 1 - thickness
	if innerLeft <= innerRight && innerTop <= innerBottom {
		for r := innerTop; r <= innerBottom; r++ {
			for c := innerLeft; c <= innerRight; c++ {
				grid[r*gridCols+c] = 0
			}
		}
	}

	addStrategicBlockers(grid, difficulty)
	return grid
}

// shapeWideRect: Wide but short rectangle — emphasizes horizontal matches.
func shapeWideRect(difficulty int) []int {
	grid := makeGrid()
	var w, h int
	switch {
	case difficulty <= 2:
		w, h = 10, 5
	case difficulty <= 4:
		w, h = 11, 6
	default:
		w, h = 12, 7
	}
	startC := (maxVisibleCol + 1 - w) / 2
	startR := (maxVisibleRow + 1 - h) / 2
	fillRect(grid, startC, startR, startC+w-1, startR+h-1)
	addStrategicBlockers(grid, difficulty)
	return grid
}

// --- Utility ---

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// --- LDtk JSON builders ---

func buildLdtkl(identifier, iid string, uid, worldX, numColors, scoreTarget, timeLimit int, grid []int, layerIID string) map[string]interface{} {
	return map[string]interface{}{
		"__header__": map[string]interface{}{
			"fileType":   "LDtk Project JSON",
			"app":        "LDtk",
			"doc":        "https://ldtk.io/json",
			"schema":     "https://ldtk.io/files/JSON_SCHEMA.json",
			"appAuthor":  "Sebastien 'deepnight' Benard",
			"appVersion": "1.5.3",
			"url":        "https://ldtk.io",
		},
		"identifier":        identifier,
		"iid":               iid,
		"uid":               uid,
		"worldX":            worldX,
		"worldY":            0,
		"worldDepth":        0,
		"pxWid":             514,
		"pxHei":             360,
		"__bgColor":         "#696A79",
		"bgColor":           nil,
		"useAutoIdentifier": true,
		"bgRelPath":         nil,
		"bgPos":             nil,
		"bgPivotX":          0.5,
		"bgPivotY":          0.5,
		"__smartColor":      "#ADADB5",
		"__bgPos":           nil,
		"externalRelPath":   nil,
		"fieldInstances":    buildFieldInstances(numColors, scoreTarget, timeLimit),
		"layerInstances": []interface{}{
			buildBoardLayer(layerIID, uid, grid),
		},
		"__neighbours": []interface{}{},
	}
}

func buildBoardLayer(iid string, levelID int, grid []int) map[string]interface{} {
	gridIface := make([]interface{}, len(grid))
	for i, v := range grid {
		gridIface[i] = v
	}

	return map[string]interface{}{
		"__identifier":       "Board",
		"__type":             "IntGrid",
		"__cWid":             gridCols,
		"__cHei":             gridRows,
		"__gridSize":         32,
		"__opacity":          1,
		"__pxTotalOffsetX":   0,
		"__pxTotalOffsetY":   0,
		"__tilesetDefUid":    nil,
		"__tilesetRelPath":   nil,
		"iid":                iid,
		"levelId":            levelID,
		"layerDefUid":        4,
		"pxOffsetX":          0,
		"pxOffsetY":          0,
		"visible":            true,
		"optionalRules":      []interface{}{},
		"intGridCsv":         gridIface,
		"autoLayerTiles":     []interface{}{},
		"seed":               rand.Intn(9999999),
		"overrideTilesetUid": nil,
		"gridTiles":          []interface{}{},
		"entityInstances":    []interface{}{},
	}
}

func buildFieldInstances(numColors, scoreTarget, timeLimit int) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"__identifier":     "NumColors",
			"__type":           "Int",
			"__value":          numColors,
			"__tile":           nil,
			"defUid":           1,
			"realEditorValues": []interface{}{map[string]interface{}{"id": "V_Int", "params": []interface{}{numColors}}},
		},
		map[string]interface{}{
			"__identifier":     "ScoreTarget",
			"__type":           "Int",
			"__value":          scoreTarget,
			"__tile":           nil,
			"defUid":           2,
			"realEditorValues": []interface{}{map[string]interface{}{"id": "V_Int", "params": []interface{}{scoreTarget}}},
		},
		map[string]interface{}{
			"__identifier":     "TimeLimit",
			"__type":           "Int",
			"__value":          timeLimit,
			"__tile":           nil,
			"defUid":           3,
			"realEditorValues": []interface{}{map[string]interface{}{"id": "V_Int", "params": []interface{}{timeLimit}}},
		},
	}
}

func buildProjectLevelEntry(identifier, iid string, uid, worldX, numColors, scoreTarget, timeLimit int, filename string) map[string]interface{} {
	return map[string]interface{}{
		"identifier":        identifier,
		"iid":               iid,
		"uid":               uid,
		"worldX":            worldX,
		"worldY":            0,
		"worldDepth":        0,
		"pxWid":             514,
		"pxHei":             360,
		"__bgColor":         "#696A79",
		"bgColor":           nil,
		"useAutoIdentifier": true,
		"bgRelPath":         nil,
		"bgPos":             nil,
		"bgPivotX":          0.5,
		"bgPivotY":          0.5,
		"__smartColor":      "#ADADB5",
		"__bgPos":           nil,
		"externalRelPath":   "match-3/" + filename,
		"fieldInstances":    buildFieldInstances(numColors, scoreTarget, timeLimit),
		"layerInstances":    nil,
		"__neighbours":      []interface{}{},
	}
}

func generateIID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x10
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func writeSimplified(ldtklDir, identifier, iid string, worldX, numColors, scoreTarget, timeLimit int, grid []int) {
	dir := filepath.Join(ldtklDir, "simplified", identifier)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("creating simplified dir %s: %v", dir, err)
	}

	dataJSON := map[string]interface{}{
		"identifier":      identifier,
		"uniqueIdentifer": iid,
		"x":               worldX,
		"y":               0,
		"width":           514,
		"height":          360,
		"bgColor":         "#696A79",
		"neighbourLevels": []interface{}{},
		"customFields": map[string]interface{}{
			"NumColors":   numColors,
			"ScoreTarget": scoreTarget,
			"TimeLimit":   timeLimit,
		},
		"layers":   []interface{}{},
		"entities": map[string]interface{}{},
	}
	dataBytes, err := json.MarshalIndent(dataJSON, "", "\t")
	if err != nil {
		log.Fatalf("marshaling data.json for %s: %v", identifier, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "data.json"), dataBytes, 0644); err != nil {
		log.Fatalf("writing data.json for %s: %v", identifier, err)
	}

	var sb strings.Builder
	for r := 0; r < gridRows; r++ {
		for c := 0; c < gridCols; c++ {
			if c > 0 {
				sb.WriteByte(',')
			}
			fmt.Fprintf(&sb, "%d", grid[r*gridCols+c])
		}
		sb.WriteString(",\n")
	}
	if err := os.WriteFile(filepath.Join(dir, "Board.csv"), []byte(sb.String()), 0644); err != nil {
		log.Fatalf("writing Board.csv for %s: %v", identifier, err)
	}
}
