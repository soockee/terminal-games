package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/pong/archetype"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// SystemFunc is a function registered as an ECS update system.
type SystemFunc = ecs.System

// RendererFunc is a function registered as an ECS draw renderer.
type RendererFunc func(e *ecs.ECS, screen *ebiten.Image)

// GameConfig holds the systems and renderers that every level shares.
type GameConfig struct {
	Systems   []SystemFunc
	Renderers []RendererFunc
}

// LevelConfig holds per-level overrides. Systems and Renderers here are
// appended after the game-wide ones.
type LevelConfig struct {
	Systems   []SystemFunc
	Renderers []RendererFunc
}

// LoadedLevel is the result of Build: a fully populated ECS world plus
// the Ebiten images needed for rendering.
type LoadedLevel struct {
	ECS         *ecs.ECS
	BGImage     *ebiten.Image
	LayerImages []*ebiten.Image
	ScreenW     int
	ScreenH     int
}

// Build creates a fresh ECS world from an LDtk level, a game-wide config,
// and a per-level config.
func Build(level *ldtkgo.Level, gc GameConfig, lc LevelConfig) *LoadedLevel {
	e := ecs.NewECS(donburi.NewWorld())

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

	// Singletons.
	archetype.NewSpace(e.World, level.Width, level.Height, 8, 8)
	archetype.NewDebug(e.World)

	// Prepare images.
	loaded := &LoadedLevel{
		ECS:     e,
		ScreenW: level.Width,
		ScreenH: level.Height,
	}

	if level.BGImage != nil {
		loaded.BGImage = ebiten.NewImageFromImage(level.BGImage)
	}
	for _, l := range level.LoadedLayers {
		loaded.LayerImages = append(loaded.LayerImages, ebiten.NewImageFromImage(l.Image))
	}

	return loaded
}

// MergeLevelConfig returns a copy of defaultCfg with any non-zero fields
// from override applied on top. Systems and Renderers are appended.
func MergeLevelConfig(defaultCfg, override LevelConfig) LevelConfig {
	merged := LevelConfig{}

	merged.Systems = append(merged.Systems, defaultCfg.Systems...)
	merged.Systems = append(merged.Systems, override.Systems...)
	merged.Renderers = append(merged.Renderers, defaultCfg.Renderers...)
	merged.Renderers = append(merged.Renderers, override.Renderers...)
	return merged
}
