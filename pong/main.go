package main

import (
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/pong/assets"
	"github.com/soockee/terminal-games/pong/game"
	"github.com/soockee/terminal-games/pong/system"
)

// Game holds the loaded level state and LDtk data.
type Game struct {
	world *ldtkgo.World

	gameConfig   game.GameConfig
	defaultLevel game.LevelConfig
	levelConfigs map[string]game.LevelConfig

	loaded *game.LoadedLevel
	level  *ldtkgo.Level

	// Level switching keys 1-9.
	levelKeys []ebiten.Key
}

func main() {
	ldtkFS, err := fs.Sub(assets.FS, "ldtk")
	if err != nil {
		log.Fatalf("sub fs: %v", err)
	}
	world, err := ldtkgo.LoadWorld("pong.ldtk", ldtkFS)
	if err != nil {
		log.Fatalf("loading world: %v", err)
	}

	gameConfig := game.GameConfig{
		Systems: []game.SystemFunc{
			system.UpdateInput,
			system.UpdateMovement,
			system.UpdateCollision,
			system.UpdateScore,
		},
		Renderers: []game.RendererFunc{
			system.DrawEntities,
			system.DrawScore,
			system.DrawDebug,
		},
	}

	defaultLevel := game.LevelConfig{
		BallSpeed:   4.0,
		PaddleSpeed: 5.0,
	}

	// Per-level overrides keyed by LDtk level identifier.
	// Levels not listed here use defaultLevel as-is.
	levelConfigs := map[string]game.LevelConfig{}

	g := &Game{
		world:        world,
		gameConfig:   gameConfig,
		defaultLevel: defaultLevel,
		levelConfigs: levelConfigs,
		levelKeys: []ebiten.Key{
			ebiten.Key1, ebiten.Key2, ebiten.Key3,
			ebiten.Key4, ebiten.Key5, ebiten.Key6,
			ebiten.Key7, ebiten.Key8, ebiten.Key9,
		},
	}

	g.loadLevel(0)

	ebiten.SetWindowSize(g.loaded.ScreenW, g.loaded.ScreenH)
	ebiten.SetWindowTitle("Pong — ECS + ldtk-super-simple-loader")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// loadLevel builds a fresh ECS world from the given level index.
func (g *Game) loadLevel(index int) {
	if index < 0 || index >= len(g.world.Levels) {
		return
	}

	level := g.world.Levels[index]
	g.level = level

	lc := g.defaultLevel
	if override, ok := g.levelConfigs[level.Identifier]; ok {
		lc = game.MergeLevelConfig(g.defaultLevel, override)
	}

	g.loaded = game.Build(level, g.gameConfig, lc)
}

func (g *Game) Update() error {
	// Level switching: keys 1-9
	for i, key := range g.levelKeys {
		if i >= len(g.world.Levels) {
			break
		}
		if ebiten.IsKeyPressed(key) && g.world.Levels[i] != g.level {
			g.loadLevel(i)
			return nil
		}
	}

	g.loaded.ECS.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	if g.loaded.BGImage != nil {
		op := &ebiten.DrawImageOptions{}
		bw := float64(g.loaded.BGImage.Bounds().Dx())
		bh := float64(g.loaded.BGImage.Bounds().Dy())
		op.GeoM.Scale(float64(g.loaded.ScreenW)/bw, float64(g.loaded.ScreenH)/bh)
		screen.DrawImage(g.loaded.BGImage, op)
	}

	// Layer images
	for _, img := range g.loaded.LayerImages {
		screen.DrawImage(img, nil)
	}

	// ECS renderers (entities, score)
	g.loaded.ECS.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.loaded.ScreenW, g.loaded.ScreenH
}
