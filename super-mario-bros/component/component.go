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

// GamePhase enumerates the mutually exclusive states of the game.
// Transitions go through methods on GameStateData so the rules live in one place.
type GamePhase int

const (
	PhaseIdle    GamePhase = iota // awaiting first input
	PhasePlaying                  // active gameplay
	PhasePaused                   // explicitly paused
	PhaseDead                     // player died, awaiting restart input
	PhaseWon                      // player won, awaiting restart input
)

type GameStateData struct {
	Phase            GamePhase
	RestartRequested bool // set by input to signal a level reload
}

var GameState = donburi.NewComponentType[GameStateData]()

// IsActive reports whether gameplay systems should run this tick.
func (g *GameStateData) IsActive() bool {
	return g.Phase == PhasePlaying
}

// Start moves Idle → Playing on the first input.
func (g *GameStateData) Start() {
	if g.Phase == PhaseIdle {
		g.Phase = PhasePlaying
	}
}

// TogglePause flips between Playing and Paused; ignored from any other phase.
func (g *GameStateData) TogglePause() {
	switch g.Phase {
	case PhasePlaying:
		g.Phase = PhasePaused
	case PhasePaused:
		g.Phase = PhasePlaying
	}
}

// Resume moves Paused → Playing.
func (g *GameStateData) Resume() {
	if g.Phase == PhasePaused {
		g.Phase = PhasePlaying
	}
}

// Die moves Playing → Dead.
func (g *GameStateData) Die() {
	if g.Phase == PhasePlaying {
		g.Phase = PhaseDead
	}
}

// Win moves Playing → Won.
func (g *GameStateData) Win() {
	if g.Phase == PhasePlaying {
		g.Phase = PhaseWon
	}
}

// RequestRestart marks the game for reload; only valid from terminal phases.
func (g *GameStateData) RequestRestart() {
	if g.Phase == PhaseDead || g.Phase == PhaseWon {
		g.RestartRequested = true
	}
}

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
