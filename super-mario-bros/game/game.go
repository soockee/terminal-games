package game

import (
	"bytes"
	"image/color"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/super-mario-bros/archetype"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/soockee/terminal-games/super-mario-bros/event"
	"github.com/soockee/terminal-games/super-mario-bros/system"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// SystemFunc is a function registered as an ECS update system.
type SystemFunc = ecs.System

// RendererFunc is a function registered as an ECS draw renderer.
type RendererFunc func(e *ecs.ECS, screen *ebiten.Image)

// GameConfig holds the systems and renderers that every level shares,
// plus the virtual screen dimensions.
type GameConfig struct {
	Systems   []SystemFunc
	Renderers []RendererFunc
	VirtualW  int // logical render width  (0 → use LDtk level width)
	VirtualH  int // logical render height (0 → use LDtk level height)
}

// LevelConfig holds per-level overrides. Systems and Renderers here are
// appended after the game-wide ones.
type LevelConfig struct {
	Systems   []SystemFunc
	Renderers []RendererFunc
}

// AudioConfig defines the asset paths for background music and SFX.
type AudioConfig struct {
	BGMusicPath string            // path inside embed FS (mp3/ogg)
	SFX         map[string]string // name → path inside embed FS (wav/mp3/ogg)
}

// AudioState holds the decoded audio context and players.
// The Game struct owns this across level reloads.
type AudioState struct {
	Ctx     *audio.Context
	BGMusic *audio.Player
	SFXData map[string][]byte // raw decoded PCM per SFX (for re-creating players)
}

// Game holds the loaded level state, LDtk data, and audio.
// It implements ebiten.Game.
type Game struct {
	world *ldtkgo.World

	gameConfig   GameConfig
	defaultLevel LevelConfig
	levelConfigs map[string]LevelConfig
	audioState   *AudioState

	loaded *loadedLevel
	level  *ldtkgo.Level

	// Level switching keys 1-9.
	levelKeys []ebiten.Key
}

// New creates a Game from an LDtk world, a game-wide config, a default
// level config, per-level overrides, and optional audio state.
// It loads the first level before returning.
func New(world *ldtkgo.World, gc GameConfig, defaultLevel LevelConfig, levelConfigs map[string]LevelConfig, audioState *AudioState) *Game {
	g := &Game{
		world:        world,
		gameConfig:   gc,
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
	return g
}

// ScreenSize returns the virtual screen dimensions (for window sizing).
func (g *Game) ScreenSize() (w, h int) {
	return g.loaded.screenW, g.loaded.screenH
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

	g.loaded = build(level, donburi.NewWorld(), g.gameConfig, lc, g.audioState)
	system.SubscribeAudioEvents(g.loaded.e.World)

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
		if entry, ok := component.GameState.First(g.loaded.e.World); ok {
			gs := component.GameState.Get(entry)
			before := gs.Phase
			gs.TogglePause()
			if gs.Phase != before {
				event.AudioEvent.Publish(g.loaded.e.World, event.AudioEventData{Name: "pause"})
			}
		}
	}

	// Skip ECS update while paused.
	if entry, ok := component.GameState.First(g.loaded.e.World); ok {
		gs := component.GameState.Get(entry)
		if gs.Phase == component.PhasePaused {
			system.ProcessEvents(g.loaded.e)
			return nil
		}
	}

	g.loaded.e.Update()

	// Check restart.
	if entry, ok := component.GameState.First(g.loaded.e.World); ok {
		gs := component.GameState.Get(entry)
		if gs.RestartRequested {
			g.loadLevel(0)
			return nil
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Fill with the level background color.
	if g.loaded.bgColor != nil {
		screen.Fill(g.loaded.bgColor)
	}

	// Background image scaled to virtual screen.
	if g.loaded.bgImage != nil {
		op := &ebiten.DrawImageOptions{}
		bw := float64(g.loaded.bgImage.Bounds().Dx())
		bh := float64(g.loaded.bgImage.Bounds().Dy())
		op.GeoM.Scale(float64(g.loaded.screenW)/bw, float64(g.loaded.screenH)/bh)
		screen.DrawImage(g.loaded.bgImage, op)
	}

	// Layer images scaled by virtual-to-world ratio.
	for _, img := range g.loaded.layerImages {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(g.loaded.scaleX, g.loaded.scaleY)
		screen.DrawImage(img, op)
	}

	// ECS renderers (entities, score, HUD, debug).
	g.loaded.e.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.loaded.screenW, g.loaded.screenH
}

// InitAudio decodes audio assets from the given FS. Returns nil if cfg is empty.
func InitAudio(fsys fs.FS, cfg AudioConfig) *AudioState {
	if cfg.BGMusicPath == "" && len(cfg.SFX) == 0 {
		return nil
	}
	const sampleRate = 44100
	ctx := audio.NewContext(sampleRate)
	state := &AudioState{Ctx: ctx, SFXData: make(map[string][]byte)}

	// Background music (ogg/mp3, infinite loop).
	if cfg.BGMusicPath != "" {
		f, err := fsys.Open(cfg.BGMusicPath)
		if err != nil {
			log.Printf("audio: open bgm: %v", err)
		} else {
			defer f.Close()
			var decoded io.ReadSeeker
			var length int64
			ext := strings.ToLower(filepath.Ext(cfg.BGMusicPath))
			switch ext {
			case ".ogg":
				d, err2 := vorbis.DecodeWithSampleRate(sampleRate, f)
				if err2 != nil {
					log.Printf("audio: decode bgm ogg: %v", err2)
				} else {
					decoded, length = d, d.Length()
				}
			case ".mp3":
				d, err2 := mp3.DecodeWithSampleRate(sampleRate, f)
				if err2 != nil {
					log.Printf("audio: decode bgm mp3: %v", err2)
				} else {
					decoded, length = d, d.Length()
				}
			default:
				log.Printf("audio: unsupported bgm format: %s", ext)
			}
			if decoded != nil {
				loop := audio.NewInfiniteLoop(decoded, length)
				p, err := ctx.NewPlayer(loop)
				if err == nil {
					state.BGMusic = p
				}
			}
		}
	}

	// SFX (wav/mp3/ogg).
	for name, path := range cfg.SFX {
		f, err := fsys.Open(path)
		if err != nil {
			log.Printf("audio: open sfx %s: %v", name, err)
			continue
		}
		var raw []byte
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".wav":
			decoded, err2 := wav.DecodeWithSampleRate(sampleRate, f)
			if err2 != nil {
				log.Printf("audio: decode sfx %s (wav): %v", name, err2)
				f.Close()
				continue
			}
			raw, _ = io.ReadAll(decoded)
		case ".mp3":
			decoded, err2 := mp3.DecodeWithSampleRate(sampleRate, f)
			if err2 != nil {
				log.Printf("audio: decode sfx %s (mp3): %v", name, err2)
				f.Close()
				continue
			}
			raw, _ = io.ReadAll(decoded)
		case ".ogg":
			decoded, err2 := vorbis.DecodeWithSampleRate(sampleRate, f)
			if err2 != nil {
				log.Printf("audio: decode sfx %s (ogg): %v", name, err2)
				f.Close()
				continue
			}
			raw, _ = io.ReadAll(decoded)
		default:
			log.Printf("audio: unsupported sfx format for %s: %s", name, ext)
			f.Close()
			continue
		}
		f.Close()
		state.SFXData[name] = raw
	}

	return state
}

// loadedLevel is the result of build: a fully populated ECS world plus
// the Ebiten images needed for rendering.
type loadedLevel struct {
	e           *ecs.ECS
	bgColor     color.Color
	bgImage     *ebiten.Image
	layerImages []*ebiten.Image
	screenW     int
	screenH     int
	scaleX      float64 // virtualW / worldW
	scaleY      float64 // virtualH / worldH
}

// build creates a fresh ECS world from an LDtk level, a game-wide config,
// a per-level config, and optional audio state.
func build(level *ldtkgo.Level, w donburi.World, gc GameConfig, lc LevelConfig, audioState *AudioState) *loadedLevel {
	e := ecs.NewECS(w)

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

	// Compute virtual screen dimensions and scale.
	vw, vh := gc.VirtualW, gc.VirtualH
	if vw == 0 {
		vw = level.Width
	}
	if vh == 0 {
		vh = level.Height
	}
	scaleX := float64(vw) / float64(level.Width)
	scaleY := float64(vh) / float64(level.Height)

	// Singletons.
	archetype.NewDebug(e.World)
	archetype.NewCamera(e.World, scaleX, scaleY)
	archetype.NewScore(e.World, 0)
	archetype.NewGameState(e.World)

	// Physics space — Mario-style downward gravity.
	archetype.NewPhysicsSpace(e.World, 980)

	// Spawn ECS entities from LDtk level data.
	spawnEntities(e.World, level)

	// Register physics collision handlers.
	system.SetupCollisionHandlers(e.World)

	// Audio.
	if audioState != nil {
		sfxPlayers := make(map[string]*audio.Player)
		for name, raw := range audioState.SFXData {
			p, err := audioState.Ctx.NewPlayer(bytes.NewReader(raw))
			if err == nil {
				sfxPlayers[name] = p
			}
		}
		archetype.NewAudio(e.World, audioState.Ctx, audioState.BGMusic, sfxPlayers)
	}

	// Prepare images.
	loaded := &loadedLevel{
		e:       e,
		bgColor: level.BGColor,
		screenW: vw,
		screenH: vh,
		scaleX:  scaleX,
		scaleY:  scaleY,
	}

	if level.BGImage != nil {
		loaded.bgImage = ebiten.NewImageFromImage(level.BGImage)
	}
	for _, l := range level.LoadedLayers {
		if l.Type == ldtkgo.LayerTiles {
			loaded.layerImages = append(loaded.layerImages, ebiten.NewImageFromImage(l.Image))
		}
	}

	return loaded
}
