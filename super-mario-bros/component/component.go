package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/yohamta/donburi"
)

// ---- Position (world-space coordinates) ----

type PositionData struct {
	X, Y float64
}

var Position = donburi.NewComponentType[PositionData]()

// ---- Sprite rendering ----

type SpriteData struct {
	Image *ebiten.Image
}

var Sprite = donburi.NewComponentType[SpriteData]()

// ---- Score (singleton) ----

type ScoreData struct {
	Value  int
	Target int // number of objectives to win (0 = no win condition)
}

var Score = donburi.NewComponentType[ScoreData]()

// ---- GameState (singleton) ----

type GameStateData struct {
	Dead    bool
	Started bool // false until first input
	Restart bool // set by input to signal a level reload
	Won     bool // true when win condition is met
	Paused  bool // true while the game is paused
}

var GameState = donburi.NewComponentType[GameStateData]()

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
