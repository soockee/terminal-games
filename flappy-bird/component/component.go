package component

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

// ---- Position & physics ----

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

type PlayerData struct {
	VelX    float64 // v: constant horizontal velocity
	VelY    float64 // vertical velocity, modified by gravity and jump
	Jump    float64 // j: upward impulse magnitude applied on jump
	Gravity float64 // g: downward acceleration applied each tick
}

var Player = donburi.NewComponentType[PlayerData]()

// ---- Fallback color (when no sprite is set) ----

type ColorData struct {
	Color color.RGBA
}

var Color = donburi.NewComponentType[ColorData]()

// ---- Resolv Space (singleton) ----

var Space = donburi.NewComponentType[resolv.Space]()

// ---- Ground (singleton) ----

type GroundData struct {
	Tile *ebiten.Image
	Y    float64
}

var Ground = donburi.NewComponentType[GroundData]()

// ---- Pipe ----

type PipeData struct {
	Passed bool // true once the bird has passed this pipe pair
	FlipY  bool // true for top pipes (sprite drawn upside-down)
}

var Pipe = donburi.NewComponentType[PipeData]()

// ---- Score (singleton) ----

type ScoreData struct {
	Value  int
	Target int // number of pipe pairs to win (0 = no win condition)
}

var Score = donburi.NewComponentType[ScoreData]()

// ---- GameOver (singleton) ----

type GameOverData struct {
	Dead    bool
	Started bool // false until first jump input
	Restart bool // set by input to signal a level reload
	Won     bool // true when all pipes are passed
	Paused  bool // true while the game is paused
}

var GameState = donburi.NewComponentType[GameOverData]()

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

// SFX name constants for compile-time safety.
const (
	SFXJump      = "jump"
	SFXHurt      = "hurt"
	SFXExplosion = "explosion"
)

type AudioData struct {
	Ctx     *audio.Context
	BGMusic *audio.Player
	SFX     map[string]*audio.Player
}

var Audio = donburi.NewComponentType[AudioData]()
