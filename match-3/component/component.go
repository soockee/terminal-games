package component

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

// ---- Shape (resolv collision shape, owned per-entity) ----

type ShapeData struct {
	Shape resolv.IShape
}

var Shape = donburi.NewComponentType[ShapeData]()

// ---- Sprite rendering ----

type SpriteData struct {
	Image *ebiten.Image
}

var Sprite = donburi.NewComponentType[SpriteData]()

// ---- Spawn position (for reset after scoring) ----

type SpawnPosData struct {
	X, Y float64
}

var SpawnPos = donburi.NewComponentType[SpawnPosData]()

// ---- Fallback color (when no sprite is set) ----

type ColorData struct {
	Color color.RGBA
}

var Color = donburi.NewComponentType[ColorData]()

// ---- Resolv Space (singleton) ----

var Space = donburi.NewComponentType[resolv.Space]()

// ---- Score (singleton) ----

type ScoreData struct {
	Value  int
	Target int // number of objectives to win (0 = no win condition)
}

var Score = donburi.NewComponentType[ScoreData]()

// ---- GameOver (singleton) ----

type GameOverData struct {
	Dead    bool
	Started bool // false until first input
	Restart bool // set by input to signal a level reload
	Won     bool // true when win condition is met
}

var GameOver = donburi.NewComponentType[GameOverData]()

// ---- Camera (singleton) ----

type CameraData struct {
	X      float64 // world-space X origin of the viewport
	ScaleX float64 // virtual screen / world scale (horizontal)
	ScaleY float64 // virtual screen / world scale (vertical)
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
