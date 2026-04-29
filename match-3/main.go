package main

import (
	"image/png"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/match-3/archetype"
	"github.com/soockee/terminal-games/match-3/assets"
	"github.com/soockee/terminal-games/match-3/component"
	"github.com/soockee/terminal-games/match-3/event"
	"github.com/soockee/terminal-games/match-3/game"
	"github.com/soockee/terminal-games/match-3/system"
	"github.com/yohamta/donburi"
)

// Game holds the loaded level state, LDtk data, and audio.
type Game struct {
	world *ldtkgo.World

	gameConfig   game.GameConfig
	defaultLevel game.LevelConfig
	levelConfigs map[string]game.LevelConfig
	audioState   *game.AudioState

	tileSheet *ebiten.Image

	loaded     *game.LoadedLevel
	level      *ldtkgo.Level
	levelIndex int

	// Level switching keys 1-9.
	levelKeys []ebiten.Key
}

func main() {
	ldtkFS, err := fs.Sub(assets.FS, "ldtk")
	if err != nil {
		log.Fatalf("sub fs: %v", err)
	}

	world, err := ldtkgo.LoadWorld("match-3.ldtk", ldtkFS)
	if err != nil {
		log.Fatalf("loading world: %v", err)
	}

	// Load the gem tileset.
	tileSheet := loadTileSheet()

	// Virtual screen: the logical resolution the game renders at.
	const virtualW, virtualH = 0, 0 // 0 = use LDtk level size

	gameConfig := game.GameConfig{
		VirtualW: virtualW,
		VirtualH: virtualH,
		Systems: []game.SystemFunc{
			system.UpdateBackground,
			system.UpdateInput,
			system.UpdateBoard,
			system.UpdateTween,
			system.UpdateScore,
			system.ProcessEvents,
		},
		Renderers: []game.RendererFunc{
			system.DrawBackground,
			system.DrawEntities,
			system.DrawScore,
			system.DrawHUD,
			system.DrawDebug,
		},
	}

	defaultLevel := game.LevelConfig{}
	levelConfigs := map[string]game.LevelConfig{}

	audioState := game.InitAudio(assets.FS, game.AudioConfig{
		BGMusicPath: "ldtk/audio/bgm.ogg",
		SFX: map[string]string{
			"swap":         "ldtk/audio/swap.mp3",
			"match":        "ldtk/audio/match.wav",
			"invalid_swap": "ldtk/audio/invalid_swap.wav",
			"select":       "ldtk/audio/select.mp3",
			"pause":        "ldtk/audio/pause.wav",
			"chain_1":      "ldtk/audio/chain_1.wav",
			"chain_2":      "ldtk/audio/chain_2.wav",
			"chain_3":      "ldtk/audio/chain_3.wav",
			"chain_4":      "ldtk/audio/chain_4.wav",
			"chain_5":      "ldtk/audio/chain_5.wav",
			"chain_6":      "ldtk/audio/chain_6.wav",
			"chain_7":      "ldtk/audio/chain_7.wav",
			"chain_8":      "ldtk/audio/chain_8.wav",
		},
	})

	g := &Game{
		world:        world,
		gameConfig:   gameConfig,
		defaultLevel: defaultLevel,
		levelConfigs: levelConfigs,
		audioState:   audioState,
		tileSheet:    tileSheet,
		levelKeys: []ebiten.Key{
			ebiten.Key1, ebiten.Key2, ebiten.Key3,
			ebiten.Key4, ebiten.Key5, ebiten.Key6,
			ebiten.Key7, ebiten.Key8, ebiten.Key9,
		},
	}

	g.loadLevel(0)

	ebiten.SetWindowSize(g.loaded.ScreenW*2, g.loaded.ScreenH*2)
	ebiten.SetWindowTitle("Match-3")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func loadTileSheet() *ebiten.Image {
	f, err := assets.FS.Open("ldtk/tilesets/match_3_art.png")
	if err != nil {
		log.Fatalf("open tileset: %v", err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		log.Fatalf("decode tileset: %v", err)
	}
	return ebiten.NewImageFromImage(img)
}

// loadLevel builds a fresh ECS world from the given level index.
func (g *Game) loadLevel(index int) {
	if index < 0 || index >= len(g.world.Levels) {
		return
	}

	g.levelIndex = index
	level := g.world.Levels[index]
	g.level = level

	lc := g.defaultLevel
	if override, ok := g.levelConfigs[level.Identifier]; ok {
		lc = override
	}

	g.loaded = game.Build(level, donburi.NewWorld(), g.gameConfig, lc, g.audioState, g.tileSheet)
	system.SubscribeAudioEvents(g.loaded.ECS.World)

	// Create scrolling background entity.
	tile := system.GenerateSymmetricTile(64)
	archetype.NewScrollingBG(g.loaded.ECS.World, tile, 0.3)

	// Mark the win screen level so the HUD can show the congratulatory message.
	if level.Identifier == "Win_screen" {
		if entry, ok := component.GameState.First(g.loaded.ECS.World); ok {
			gs := component.GameState.Get(entry)
			gs.WinScreen = true
			gs.Started = true // no gameplay needed, auto-start
		}
	}

	if g.audioState != nil && g.audioState.BGMusic != nil && !g.audioState.BGMusic.IsPlaying() {
		g.audioState.BGMusic.Play()
	}
}

func (g *Game) Update() error {
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

	if entry, ok := component.GameState.First(g.loaded.ECS.World); ok {
		gs := component.GameState.Get(entry)
		if gs.Restart {
			if gs.WinScreen {
				// Restart from beginning after win screen.
				g.loadLevel(0)
			} else if gs.Won {
				// Advance to next level on win.
				next := g.levelIndex + 1
				if next >= len(g.world.Levels) {
					next = 0 // wrap to first level
				}
				g.loadLevel(next)
			} else {
				// Retry current level on game over.
				g.loadLevel(g.levelIndex)
			}
			return nil
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Use LDtk background image if the level has one (e.g., win screen),
	// otherwise the ECS DrawBackground renderer handles the scrolling BG.
	if g.loaded.BGImage != nil {
		op := &ebiten.DrawImageOptions{}
		bw := float64(g.loaded.BGImage.Bounds().Dx())
		bh := float64(g.loaded.BGImage.Bounds().Dy())
		op.GeoM.Scale(float64(g.loaded.ScreenW)/bw, float64(g.loaded.ScreenH)/bh)
		screen.DrawImage(g.loaded.BGImage, op)
	}

	g.loaded.ECS.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.loaded.ScreenW, g.loaded.ScreenH
}
