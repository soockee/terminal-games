package game

import (
	"bytes"
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
	archetype.NewDebug(e.World)
	archetype.NewCamera(e.World, scaleX, scaleY)
	archetype.NewScore(e.World, 0)
	archetype.NewGameState(e.World)

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
