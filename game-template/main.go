package main

import (
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/game-template/assets"
	"github.com/soockee/terminal-games/game-template/component"
	"github.com/soockee/terminal-games/game-template/event"
	"github.com/soockee/terminal-games/game-template/game"
	"github.com/soockee/terminal-games/game-template/system"
	"github.com/yohamta/donburi"
)

// Game holds the loaded level state, LDtk data, and audio.
type Game struct {
	world *ldtkgo.World

	gameConfig   game.GameConfig
	defaultLevel game.LevelConfig
	levelConfigs map[string]game.LevelConfig
	audioState   *game.AudioState

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

	// TODO: replace with your LDtk project file name.
	world, err := ldtkgo.LoadWorld("game.ldtk", ldtkFS)
	if err != nil {
		log.Fatalf("loading world: %v", err)
	}

	// Virtual screen: the logical resolution the game renders at.
	// The LDtk level dimensions define the world; VirtualW/H can differ.
	const virtualW, virtualH = 0, 0 // 0 = use LDtk level size

	gameConfig := game.GameConfig{
		VirtualW: virtualW,
		VirtualH: virtualH,
		Systems: []game.SystemFunc{
			system.UpdateInput,
			system.UpdateMovement,
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
	}

	defaultLevel := game.LevelConfig{}

	// Per-level overrides keyed by LDtk level identifier.
	levelConfigs := map[string]game.LevelConfig{}

	// Audio — configure paths relative to the embedded FS.
	// audioState := game.InitAudio(assets.FS, game.AudioConfig{
	//     BGMusicPath: "audio/bgm.mp3",
	//     SFX: map[string]string{
	//         "jump": "audio/jump.wav",
	//     },
	// })
	var audioState *game.AudioState

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

	// Window = 2× virtual screen for crisp display.
	ebiten.SetWindowSize(g.loaded.ScreenW*2, g.loaded.ScreenH*2)
	ebiten.SetWindowTitle("Game Template — ECS + LDtk")

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

	g.loaded = game.Build(level, donburi.NewWorld(), g.gameConfig, lc, g.audioState)
	system.SubscribeAudioEvents(g.loaded.ECS.World)

	// Start BGM on first load.
	if g.audioState != nil && g.audioState.BGMusic != nil && !g.audioState.BGMusic.IsPlaying() {
		g.audioState.BGMusic.Play()
	}
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

	// Pause toggle (P key).
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if entry, ok := component.GameState.First(g.loaded.ECS.World); ok {
			gs := component.GameState.Get(entry)
			if gs.Started && !gs.Won && !gs.Dead {
				gs.Paused = !gs.Paused
				event.AudioEvent.Publish(g.loaded.ECS.World, event.AudioEventData{Name: "pause"})
			}
		}
	}

	// Skip ECS update while paused.
	if entry, ok := component.GameState.First(g.loaded.ECS.World); ok {
		gs := component.GameState.Get(entry)
		if gs.Paused {
			system.ProcessEvents(g.loaded.ECS)
			return nil
		}
	}

	g.loaded.ECS.Update()

	// Check restart.
	if entry, ok := component.GameState.First(g.loaded.ECS.World); ok {
		gs := component.GameState.Get(entry)
		if gs.Restart {
			g.loadLevel(0)
			return nil
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background image scaled to virtual screen.
	if g.loaded.BGImage != nil {
		op := &ebiten.DrawImageOptions{}
		bw := float64(g.loaded.BGImage.Bounds().Dx())
		bh := float64(g.loaded.BGImage.Bounds().Dy())
		op.GeoM.Scale(float64(g.loaded.ScreenW)/bw, float64(g.loaded.ScreenH)/bh)
		screen.DrawImage(g.loaded.BGImage, op)
	}

	// Layer images scaled by virtual-to-world ratio.
	for _, img := range g.loaded.LayerImages {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(g.loaded.ScaleX, g.loaded.ScaleY)
		screen.DrawImage(img, op)
	}

	// ECS renderers (entities, score, HUD, debug).
	g.loaded.ECS.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.loaded.ScreenW, g.loaded.ScreenH
}
