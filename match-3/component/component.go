package component

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/yohamta/donburi"
)

// ---- Phase (board state machine) ----

type Phase int

const (
	PhaseIdle Phase = iota
	PhaseSelected
	PhaseSwapping
	PhaseReverting
	PhaseMatching
	PhaseCollapsing
	PhaseRefilling
)

// ---- EaseFunc (tween interpolation) ----

type EaseFunc int

const (
	EaseLinear EaseFunc = iota
	EaseOutQuad
)

// ---- Board components (split by responsibility) ----

// CellType values (matches LDtk IntGrid: 1=empty, 2=playable, 3=blocked).
const (
	CellEmpty    = 1
	CellPlayable = 2
	CellBlocked  = 3
)

// GridData holds the board's topology and tile references.
type GridData struct {
	Cols     int
	Rows     int
	CellType [][]int            // [col][row] IntGrid value
	Cells    [][]*donburi.Entry // [col][row] tile entity (nil if empty/blocked)
}

var BoardGrid = donburi.NewComponentType[GridData]()

// PhaseData tracks the match-3 state machine progress.
type PhaseData struct {
	Phase Phase

	SelectedCol int // -1 = none
	SelectedRow int

	SwapA [2]int // [col, row]
	SwapB [2]int

	ChainDepth     int
	ReshuffleTimer float64 // countdown (seconds) to show "Reshuffling" message
}

var BoardPhase = donburi.NewComponentType[PhaseData]()

// InputData holds UI-only input tracking.
type InputData struct {
	CursorCol     int // keyboard cursor position
	CursorRow     int
	AutoPlay      bool    // when true, automatically pick valid swaps
	AutoPlayDelay float64 // seconds to wait before next autoplay move
}

var BoardInput = donburi.NewComponentType[InputData]()

// DisplayData holds rendering layout and sprite data.
type DisplayData struct {
	OffsetX    float64 // board pixel origin on screen
	OffsetY    float64
	TileSize   int
	GemSprites []*ebiten.Image // pre-sliced per-color sprites
}

var BoardDisplay = donburi.NewComponentType[DisplayData]()

// RulesData holds per-level game rule parameters.
type RulesData struct {
	NumColors     int
	ScoreTarget   int
	TimeLimit     float64 // seconds, 0 = unlimited
	TimeRemaining float64
}

var BoardRules = donburi.NewComponentType[RulesData]()

// ---- ScrollingBG (singleton) ----

type ScrollingBGData struct {
	Tile    *ebiten.Image
	OffsetX float64
	Speed   float64 // pixels per frame
}

var ScrollingBG = donburi.NewComponentType[ScrollingBGData]()

// ---- GridPos (per-tile) ----

type GridPosData struct {
	Col, Row int
}

var GridPos = donburi.NewComponentType[GridPosData]()

// ---- GemType (per-tile) ----

type GemTypeData struct {
	Color int // 0–7
}

var GemType = donburi.NewComponentType[GemTypeData]()

// ---- PixelPos (per-tile, interpolated position for rendering) ----

type PixelPosData struct {
	X, Y float64
}

var PixelPos = donburi.NewComponentType[PixelPosData]()

// ---- Tween (per-tile animation) ----

type TweenData struct {
	StartX, StartY float64
	EndX, EndY     float64
	Elapsed        float64
	Duration       float64
	Active         bool
	Ease           EaseFunc
}

var Tween = donburi.NewComponentType[TweenData]()

// ---- Sprite rendering ----

type SpriteData struct {
	Image *ebiten.Image
}

var Sprite = donburi.NewComponentType[SpriteData]()

// ---- Score (singleton) ----

type ScoreData struct {
	Value  int
	Target int
}

var Score = donburi.NewComponentType[ScoreData]()

// ---- GameState (singleton) ----

type GameStateData struct {
	Dead      bool
	Started   bool
	Restart   bool
	Won       bool
	Paused    bool
	WinScreen bool // true when this level is the final "Win_screen" celebration
}

var GameState = donburi.NewComponentType[GameStateData]()

// ---- Camera (singleton) ----

type CameraData struct {
	X      float64
	ScaleX float64
	ScaleY float64
}

var Camera = donburi.NewComponentType[CameraData]()

// ---- Debug (singleton toggle) ----

type DebugData struct {
	Enabled bool
}

var Debug = donburi.NewComponentType[DebugData]()

// ---- Audio (singleton) ----

type AudioData struct {
	Ctx     *audio.Context
	BGMusic *audio.Player
	SFX     map[string]*audio.Player
}

var Audio = donburi.NewComponentType[AudioData]()

// ---- GemQuads (tileset mapping) ----

// GemQuads maps Color index (0–7) to a 16×16 rect in match_3_art.png.
// The tileset is 6 cols × 16 rows of 16×16 tiles.
// Selected for maximum visual contrast.
var GemQuads = [8]image.Rectangle{
	image.Rect(0, 0*16, 16, 1*16),   // 0: red (row 1)
	image.Rect(0, 4*16, 16, 5*16),   // 1: yellow (row 5)
	image.Rect(0, 7*16, 16, 8*16),   // 2: green (row 8)
	image.Rect(0, 8*16, 16, 9*16),   // 3: light blue (row 9)
	image.Rect(0, 9*16, 16, 10*16),  // 4: dark blue (row 10)
	image.Rect(0, 13*16, 16, 14*16), // 5: dark grey (row 14)
	image.Rect(0, 15*16, 16, 16*16), // 6: teal (row 16)
	image.Rect(0, 4*16, 16, 5*16),   // 7: reserved (row 5)
}
