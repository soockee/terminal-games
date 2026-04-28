package game

import (
	"bytes"
	"io"
	"io/fs"
	"log"
	"math/rand"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/match-3/archetype"
	"github.com/soockee/terminal-games/match-3/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// SystemFunc is a function registered as an ECS update system.
type SystemFunc = ecs.System

// RendererFunc is a function registered as an ECS draw renderer.
type RendererFunc func(e *ecs.ECS, screen *ebiten.Image)

// GameConfig holds the systems and renderers that every level shares,
// plus the virtual screen dimensions.
type GameConfig struct {
	Systems   []SystemFunc
	Renderers []RendererFunc
	VirtualW  int // logical render width  (0 → use LDtk level width)
	VirtualH  int // logical render height (0 → use LDtk level height)
}

// LevelConfig holds per-level overrides. Systems and Renderers here are
// appended after the game-wide ones.
type LevelConfig struct {
	archetype.SpawnConfig
	Systems   []SystemFunc
	Renderers []RendererFunc
}

// AudioConfig defines the asset paths for background music and SFX.
type AudioConfig struct {
	BGMusicPath string            // path inside embed FS (mp3)
	SFX         map[string]string // name → path inside embed FS (wav)
}

// AudioState holds the decoded audio context and players.
// The Game struct owns this across level reloads.
type AudioState struct {
	Ctx     *audio.Context
	BGMusic *audio.Player
	SFXData map[string][]byte // raw decoded PCM per SFX (for re-creating players)
}

// InitAudio decodes audio assets from the given FS. Returns nil if cfg is empty.
func InitAudio(fsys fs.FS, cfg AudioConfig) *AudioState {
	if cfg.BGMusicPath == "" && len(cfg.SFX) == 0 {
		return nil
	}
	const sampleRate = 44100
	ctx := audio.NewContext(sampleRate)
	state := &AudioState{Ctx: ctx, SFXData: make(map[string][]byte)}

	// Background music (ogg/mp3, infinite loop).
	if cfg.BGMusicPath != "" {
		f, err := fsys.Open(cfg.BGMusicPath)
		if err != nil {
			log.Printf("audio: open bgm: %v", err)
		} else {
			defer f.Close()
			var decoded io.ReadSeeker
			var length int64
			ext := strings.ToLower(filepath.Ext(cfg.BGMusicPath))
			switch ext {
			case ".ogg":
				d, err2 := vorbis.DecodeWithSampleRate(sampleRate, f)
				if err2 != nil {
					log.Printf("audio: decode bgm ogg: %v", err2)
				} else {
					decoded, length = d, d.Length()
				}
			case ".mp3":
				d, err2 := mp3.DecodeWithSampleRate(sampleRate, f)
				if err2 != nil {
					log.Printf("audio: decode bgm mp3: %v", err2)
				} else {
					decoded, length = d, d.Length()
				}
			default:
				log.Printf("audio: unsupported bgm format: %s", ext)
			}
			if decoded != nil {
				loop := audio.NewInfiniteLoop(decoded, length)
				p, err := ctx.NewPlayer(loop)
				if err == nil {
					state.BGMusic = p
				}
			}
		}
	}

	// SFX (wav/mp3/ogg).
	for name, path := range cfg.SFX {
		f, err := fsys.Open(path)
		if err != nil {
			log.Printf("audio: open sfx %s: %v", name, err)
			continue
		}
		var raw []byte
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".wav":
			decoded, err2 := wav.DecodeWithSampleRate(sampleRate, f)
			if err2 != nil {
				log.Printf("audio: decode sfx %s (wav): %v", name, err2)
				f.Close()
				continue
			}
			raw, _ = io.ReadAll(decoded)
		case ".mp3":
			decoded, err2 := mp3.DecodeWithSampleRate(sampleRate, f)
			if err2 != nil {
				log.Printf("audio: decode sfx %s (mp3): %v", name, err2)
				f.Close()
				continue
			}
			raw, _ = io.ReadAll(decoded)
		case ".ogg":
			decoded, err2 := vorbis.DecodeWithSampleRate(sampleRate, f)
			if err2 != nil {
				log.Printf("audio: decode sfx %s (ogg): %v", name, err2)
				f.Close()
				continue
			}
			raw, _ = io.ReadAll(decoded)
		default:
			log.Printf("audio: unsupported sfx format for %s: %s", name, ext)
			f.Close()
			continue
		}
		f.Close()
		state.SFXData[name] = raw
	}

	return state
}

// LoadedLevel is the result of Build: a fully populated ECS world plus
// the Ebiten images needed for rendering.
type LoadedLevel struct {
	ECS         *ecs.ECS
	BGImage     *ebiten.Image
	LayerImages []*ebiten.Image
	ScreenW     int
	ScreenH     int
	WorldW      int     // LDtk level pixel width
	WorldH      int     // LDtk level pixel height
	ScaleX      float64 // virtualW / worldW
	ScaleY      float64 // virtualH / worldH
}

// Build creates a fresh ECS world from an LDtk level, a game-wide config,
// a per-level config, and optional audio state.
func Build(level *ldtkgo.Level, w donburi.World, gc GameConfig, lc LevelConfig, audioState *AudioState, tileSheet *ebiten.Image) *LoadedLevel {
	e := ecs.NewECS(w)

	// Register game-wide systems, then level-specific systems.
	for _, s := range gc.Systems {
		e.AddSystem(s)
	}
	for _, s := range lc.Systems {
		e.AddSystem(s)
	}

	// Register game-wide renderers, then level-specific renderers.
	for _, r := range gc.Renderers {
		e.AddRenderer(ecs.LayerDefault, r)
	}
	for _, r := range lc.Renderers {
		e.AddRenderer(ecs.LayerDefault, r)
	}

	// Compute virtual screen dimensions and scale.
	vw, vh := gc.VirtualW, gc.VirtualH
	if vw == 0 {
		vw = level.Width
	}
	if vh == 0 {
		vh = level.Height
	}
	scaleX := float64(vw) / float64(level.Width)
	scaleY := float64(vh) / float64(level.Height)

	// Singletons.
	archetype.NewDebug(e.World)
	archetype.NewCamera(e.World, scaleX, scaleY)
	archetype.NewGameState(e.World)

	// Read level custom fields.
	numColors := intFieldOr(level.CustomFields, "NumColors", 6)
	scoreTarget := intFieldOr(level.CustomFields, "ScoreTarget", 1000)
	timeLimit := intFieldOr(level.CustomFields, "TimeLimit", 0)

	archetype.NewScore(e.World, scoreTarget)

	// Build the board from the IntGrid layer named "Board".
	boardData := buildBoard(level, numColors, scoreTarget, timeLimit, tileSheet, e.World)
	archetype.NewBoard(e.World, boardData)

	// Audio.
	if audioState != nil {
		sfxPlayers := make(map[string]*audio.Player)
		for name, raw := range audioState.SFXData {
			p, err := audioState.Ctx.NewPlayer(bytes.NewReader(raw))
			if err == nil {
				sfxPlayers[name] = p
			}
		}
		archetype.NewAudio(e.World, audioState.Ctx, audioState.BGMusic, sfxPlayers)
	}

	// Prepare images.
	loaded := &LoadedLevel{
		ECS:     e,
		ScreenW: vw,
		ScreenH: vh,
		WorldW:  level.Width,
		WorldH:  level.Height,
		ScaleX:  scaleX,
		ScaleY:  scaleY,
	}

	if level.BGImage != nil {
		loaded.BGImage = ebiten.NewImageFromImage(level.BGImage)
	}
	for _, l := range level.LoadedLayers {
		loaded.LayerImages = append(loaded.LayerImages, ebiten.NewImageFromImage(l.Image))
	}

	return loaded
}

// buildBoard reads the "Board" IntGrid layer and creates tile entities.
func buildBoard(level *ldtkgo.Level, numColors, scoreTarget, timeLimit int, tileSheet *ebiten.Image, w donburi.World) *component.BoardData {
	ig := level.IntGrids["Board"]
	if ig == nil {
		log.Println("warning: no 'Board' IntGrid layer found, creating empty board")
		return &component.BoardData{SelectedCol: -1, SelectedRow: -1, TileSize: 32}
	}

	cols := ig.Width
	rows := ig.Height
	tileSize := 32
	if ig.Def != nil {
		tileSize = ig.Def.GridSize
	}

	// Convert IntGrid [row][col] to our [col][row] layout.
	cellType := make([][]int, cols)
	cells := make([][]*donburi.Entry, cols)
	for c := range cols {
		cellType[c] = make([]int, rows)
		cells[c] = make([]*donburi.Entry, rows)
	}

	for r := range rows {
		for c := range cols {
			v := ig.Grid[r][c]
			if v == 0 {
				cellType[c][r] = component.CellEmpty
			} else {
				cellType[c][r] = v
			}
		}
	}

	// Pre-slice gem sprites from tileset.
	gemSprites := make([]*ebiten.Image, numColors)
	for i := range numColors {
		rect := component.GemQuads[i]
		gemSprites[i] = tileSheet.SubImage(rect).(*ebiten.Image)
	}

	// Create tile entities for playable cells (match-free initial fill).
	for c := range cols {
		for r := range rows {
			if cellType[c][r] != component.CellPlayable {
				continue
			}
			color := randomMatchFreeColor(cellType, cells, c, r, numColors, w)
			px := float64(c * tileSize)
			py := float64(r * tileSize)
			entry := archetype.NewTile(w, c, r, color, gemSprites[color], px, py)
			cells[c][r] = entry
		}
	}

	return &component.BoardData{
		Cols:          cols,
		Rows:          rows,
		CellType:      cellType,
		Cells:         cells,
		Phase:         component.PhaseIdle,
		SelectedCol:   -1,
		SelectedRow:   -1,
		NumColors:     numColors,
		ScoreTarget:   scoreTarget,
		TimeLimit:     float64(timeLimit),
		TimeRemaining: float64(timeLimit),
		OffsetX:       0,
		OffsetY:       0,
		TileSize:      tileSize,
		GemSprites:    gemSprites,
	}
}

// randomMatchFreeColor picks a random color that doesn't create a match at (c, r).
func randomMatchFreeColor(cellType [][]int, cells [][]*donburi.Entry, c, r, numColors int, w donburi.World) int {
	for attempts := 0; attempts < 100; attempts++ {
		color := rand.Intn(numColors)
		if !wouldMatch(cells, c, r, color, w) {
			return color
		}
	}
	return rand.Intn(numColors)
}

// wouldMatch checks if placing a gem of the given color at (c, r) creates a match.
func wouldMatch(cells [][]*donburi.Entry, c, r, color int, w donburi.World) bool {
	// Check horizontal: 2 to the left.
	if c >= 2 && cells[c-1][r] != nil && cells[c-2][r] != nil {
		c1 := component.GemType.Get(cells[c-1][r]).Color
		c2 := component.GemType.Get(cells[c-2][r]).Color
		if c1 == color && c2 == color {
			return true
		}
	}
	// Check vertical: 2 above.
	if r >= 2 && cells[c][r-1] != nil && cells[c][r-2] != nil {
		c1 := component.GemType.Get(cells[c][r-1]).Color
		c2 := component.GemType.Get(cells[c][r-2]).Color
		if c1 == color && c2 == color {
			return true
		}
	}
	return false
}

// intFieldOr reads an int custom field with a fallback default.
func intFieldOr(fields map[string]interface{}, key string, def int) int {
	if fields == nil {
		return def
	}
	v, ok := fields[key]
	if !ok {
		return def
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	default:
		return def
	}
}

// gemColorAt returns the color of the gem at (c, r) or -1 if empty.
func gemColorAt(cells [][]*donburi.Entry, c, r int) int {
	if cells[c][r] == nil {
		return -1
	}
	return component.GemType.Get(cells[c][r]).Color
}
