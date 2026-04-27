package component

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
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

// ---- Fallback color (when no sprite is set) ----

type ColorData struct {
	Color color.RGBA
}

var Color = donburi.NewComponentType[ColorData]()

// ---- Resolv Space (singleton) ----

var Space = donburi.NewComponentType[resolv.Space]()

// ---- Debug (singleton toggle) ----

type DebugData struct {
	Enabled bool
}

var Debug = donburi.NewComponentType[DebugData]()
