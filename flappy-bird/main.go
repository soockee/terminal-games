package main

import (
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/flappy-bird/archetype"
	"github.com/soockee/terminal-games/flappy-bird/assets"
	"github.com/soockee/terminal-games/flappy-bird/component"
	"github.com/soockee/terminal-games/flappy-bird/game"
	"github.com/soockee/terminal-games/flappy-bird/system"
)

// Game holds the loaded level state and LDtk data.
type Game struct {
	world *ldtkgo.World

	gameConfig   game.GameConfig
	defaultLevel game.LevelConfig
	levelConfigs map[string]game.LevelConfig

	audioState *game.AudioState

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
	world, err := ldtkgo.LoadWorld("flappy-bird.ldtk", ldtkFS)
	if err != nil {
		log.Fatalf("loading world: %v", err)
	}

	gameConfig := game.GameConfig{
		Systems: []game.SystemFunc{
			system.UpdateInput,
			system.UpdateMovement,
			system.UpdateCollision,
			system.UpdateCamera,
			system.UpdateScore,
			system.ProcessEvents,
		},
		Renderers: []game.RendererFunc{
			system.DrawEntities,
			system.DrawScore,
			system.DrawHUD,
			system.DrawDebug,
		},
		VirtualW: 514,
		VirtualH: 360,
	}

	defaultLevel := game.LevelConfig{
		SpawnConfig: archetype.SpawnConfig{
			BirdVelX:             2,
			BirdJump:             8,
			BirdGravity:          0.4,
			PipeMinGapVertical:   130,
			PipeMinGapHorizontal: 200,
			PipeInterval:         200,
			PipeCount:            25,
		},
	}

	// Per-level overrides keyed by LDtk level identifier.
	// Levels not listed here use defaultLevel as-is.
	levelConfigs := map[string]game.LevelConfig{}

	// Audio: decode BGM + SFX once, reuse across level reloads.
	audioCfg := game.AudioConfig{
		BGMusicPath: "audio/marios_way.mp3",
		SFX: map[string]string{
			component.SFXJump:      "audio/jump.wav",
			component.SFXHurt:      "audio/hurt.wav",
			component.SFXExplosion: "audio/explosion.wav",
		},
	}
	audioState, err := game.InitAudio(ldtkFS, audioCfg)
	if err != nil {
		log.Fatalf("init audio: %v", err)
	}

	g := &Game{
		world:        world,
		gameConfig:   gameConfig,
		defaultLevel: defaultLevel,
		levelConfigs: levelConfigs,
		audioState:   audioState,
		levelKeys: []ebiten.Key{
			ebiten.Key1, ebiten.Key2, ebiten.Key3,
			ebiten.Key4, ebiten.Key5, ebiten.Key6,
			ebiten.Key7, ebiten.Key8, ebiten.Key9,
		},
	}

	g.loadLevel(0)

	ebiten.SetWindowSize(g.loaded.ScreenW*2, g.loaded.ScreenH*2)
	ebiten.SetWindowTitle("flappy-bird — ECS + ldtk-super-simple-loader")

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
		lc = override
	}

	g.loaded = game.Build(level, g.world, g.gameConfig, lc, g.audioState)
	system.SubscribeAudioEvents(g.loaded.ECS.World)
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

	// Restart or win-switch on signal from input/score systems.
	if goEntry, ok := component.GameState.First(g.loaded.ECS.World); ok {
		go_ := component.GameState.Get(goEntry)
		if go_.Restart {
			// From win screen or death screen → back to gameplay.
			g.loadLevel(0)
		} else if go_.Won && !go_.Dead && g.level.Identifier != "Level_1" {
			// Just won on gameplay level → switch to win screen.
			g.loadLevel(1)
			// Mark the win level so HUD shows congratulations.
			if goEntry2, ok2 := component.GameState.First(g.loaded.ECS.World); ok2 {
				go2 := component.GameState.Get(goEntry2)
				go2.Won = true
				go2.Started = true
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	sx := g.loaded.ScaleX
	sy := g.loaded.ScaleY

	// Background
	if g.loaded.BGImage != nil {
		op := &ebiten.DrawImageOptions{}
		bw := float64(g.loaded.BGImage.Bounds().Dx())
		bh := float64(g.loaded.BGImage.Bounds().Dy())
		op.GeoM.Scale(float64(g.loaded.ScreenW)/bw, float64(g.loaded.ScreenH)/bh)
		screen.DrawImage(g.loaded.BGImage, op)
	}

	// Layer images (authored at LDtk resolution, scale to virtual)
	for _, img := range g.loaded.LayerImages {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(sx, sy)
		screen.DrawImage(img, op)
	}

	// ECS renderers (entities, score)
	g.loaded.ECS.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.loaded.ScreenW, g.loaded.ScreenH
}
