package game

import (
	"bytes"
	"io"
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/game-template/archetype"
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
	archetype.SpawnConfig
	Systems   []SystemFunc
	Renderers []RendererFunc
}

// AudioConfig defines the asset paths for background music and SFX.
type AudioConfig struct {
	BGMusicPath string            // path inside embed FS (mp3)
	SFX         map[string]string // name → path inside embed FS (wav)
}

// AudioState holds the decoded audio context and players.
// The Game struct owns this across level reloads.
type AudioState struct {
	Ctx     *audio.Context
	BGMusic *audio.Player
	SFXData map[string][]byte // raw decoded PCM per SFX (for re-creating players)
}

// InitAudio decodes audio assets from the given FS. Returns nil if cfg is empty.
func InitAudio(fsys fs.FS, cfg AudioConfig) *AudioState {
	if cfg.BGMusicPath == "" && len(cfg.SFX) == 0 {
		return nil
	}
	const sampleRate = 44100
	ctx := audio.NewContext(sampleRate)
	state := &AudioState{Ctx: ctx, SFXData: make(map[string][]byte)}

	// Background music (mp3, infinite loop).
	if cfg.BGMusicPath != "" {
		f, err := fsys.Open(cfg.BGMusicPath)
		if err != nil {
			log.Printf("audio: open bgm: %v", err)
		} else {
			defer f.Close()
			decoded, err := mp3.DecodeWithSampleRate(sampleRate, f)
			if err != nil {
				log.Printf("audio: decode bgm: %v", err)
			} else {
				loop := audio.NewInfiniteLoop(decoded, decoded.Length())
				p, err := ctx.NewPlayer(loop)
				if err == nil {
					state.BGMusic = p
				}
			}
		}
	}

	// SFX (wav).
	for name, path := range cfg.SFX {
		f, err := fsys.Open(path)
		if err != nil {
			log.Printf("audio: open sfx %s: %v", name, err)
			continue
		}
		decoded, err := wav.DecodeWithSampleRate(sampleRate, f)
		f.Close()
		if err != nil {
			log.Printf("audio: decode sfx %s: %v", name, err)
			continue
		}
		raw, _ := io.ReadAll(decoded)
		state.SFXData[name] = raw
	}

	return state
}

// LoadedLevel is the result of Build: a fully populated ECS world plus
// the Ebiten images needed for rendering.
type LoadedLevel struct {
	ECS         *ecs.ECS
	BGImage     *ebiten.Image
	LayerImages []*ebiten.Image
	ScreenW     int
	ScreenH     int
	WorldW      int     // LDtk level pixel width
	WorldH      int     // LDtk level pixel height
	ScaleX      float64 // virtualW / worldW
	ScaleY      float64 // virtualH / worldH
}

// Build creates a fresh ECS world from an LDtk level, a game-wide config,
// a per-level config, and optional audio state.
func Build(level *ldtkgo.Level, w donburi.World, gc GameConfig, lc LevelConfig, audioState *AudioState) *LoadedLevel {
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
	archetype.NewSpace(e.World, level.Width, level.Height, 8, 8)
	archetype.NewDebug(e.World)
	archetype.NewCamera(e.World, scaleX, scaleY)
	archetype.NewScore(e.World, 0)
	archetype.NewGameOver(e.World)

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
	loaded := &LoadedLevel{
		ECS:     e,
		ScreenW: vw,
		ScreenH: vh,
		WorldW:  level.Width,
		WorldH:  level.Height,
		ScaleX:  scaleX,
		ScaleY:  scaleY,
	}

	if level.BGImage != nil {
		loaded.BGImage = ebiten.NewImageFromImage(level.BGImage)
	}
	for _, l := range level.LoadedLayers {
		loaded.LayerImages = append(loaded.LayerImages, ebiten.NewImageFromImage(l.Image))
	}

	return loaded
}
