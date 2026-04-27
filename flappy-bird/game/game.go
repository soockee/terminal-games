package game

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/flappy-bird/archetype"
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

	// VirtualW and VirtualH define the fixed internal rendering resolution.
	// All world coordinates are projected onto this virtual screen.
	// If zero, defaults to the LDtk level dimensions.
	VirtualW int
	VirtualH int
}

// LevelConfig holds per-level overrides. Systems and Renderers here are
// appended after the game-wide ones.
type LevelConfig struct {
	archetype.SpawnConfig
	Systems   []SystemFunc
	Renderers []RendererFunc
}

// LoadedLevel is the result of Build: a fully populated ECS world plus
// the Ebiten images needed for rendering.
type LoadedLevel struct {
	ECS         *ecs.ECS
	BGImage     *ebiten.Image
	LayerImages []*ebiten.Image
	ScreenW     int     // virtual screen width (for Layout)
	ScreenH     int     // virtual screen height (for Layout)
	WorldW      int     // LDtk level width in world pixels
	WorldH      int     // LDtk level height in world pixels
	ScaleX      float64 // virtual / world horizontal scale
	ScaleY      float64 // virtual / world vertical scale
}

// AudioConfig defines the audio assets for a game. Paths are relative to the
// embedded FS (e.g. "audio/jump.wav").
type AudioConfig struct {
	BGMusicPath string            // path to looping background music (MP3)
	SFX         map[string]string // SFX name → file path (WAV)
}

// AudioState holds the long-lived audio objects that survive level reloads.
type AudioState struct {
	Ctx     *audio.Context
	BGMusic *audio.Player
	SFXData map[string][]byte // raw decoded PCM bytes per SFX, for fast re-creation
}

const sampleRate = 44100

// InitAudio creates the audio context, decodes BGM and SFX from the given FS.
// Call once at startup; the returned AudioState is reused across level reloads.
func InitAudio(afs fs.FS, cfg AudioConfig) (*AudioState, error) {
	ctx := audio.NewContext(sampleRate)
	state := &AudioState{Ctx: ctx, SFXData: make(map[string][]byte)}

	// Decode and start BGM (looping).
	if cfg.BGMusicPath != "" {
		f, err := afs.Open(cfg.BGMusicPath)
		if err != nil {
			return nil, fmt.Errorf("open bgm %s: %w", cfg.BGMusicPath, err)
		}
		defer f.Close()
		decoded, err := mp3.DecodeF32(f)
		if err != nil {
			return nil, fmt.Errorf("decode bgm: %w", err)
		}
		loop := audio.NewInfiniteLoop(decoded, decoded.Length())
		player, err := ctx.NewPlayerF32(loop)
		if err != nil {
			return nil, fmt.Errorf("bgm player: %w", err)
		}
		player.Play()
		state.BGMusic = player
	}

	// Decode SFX (WAV) into raw bytes for rewind+replay players.
	for name, path := range cfg.SFX {
		f, err := afs.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open sfx %s: %w", path, err)
		}
		decoded, err := wav.DecodeF32(f)
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("decode sfx %s: %w", path, err)
		}
		raw, err := io.ReadAll(decoded)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("read sfx %s: %w", path, err)
		}
		state.SFXData[name] = raw
	}

	return state, nil
}

// Build creates a fresh ECS world from an LDtk level, a game-wide config,
// and a per-level config. If audioState is non-nil, an Audio singleton is
// injected with pre-decoded SFX players.
func Build(level *ldtkgo.Level, world *ldtkgo.World, gc GameConfig, lc LevelConfig, audioState *AudioState) *LoadedLevel {
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

	// Virtual screen: default to LDtk dimensions if not configured.
	vw, vh := gc.VirtualW, gc.VirtualH
	if vw == 0 {
		vw = level.Width
	}
	if vh == 0 {
		vh = level.Height
	}
	scaleX := float64(vw) / float64(level.Width)
	scaleY := float64(vh) / float64(level.Height)

	archetype.NewCamera(e.World, scaleX, scaleY)

	// Spawn all entities (tag-based and definition-only) through the archetype layer.
	archetype.SpawnEntities(e.World, level, world, lc.SpawnConfig)

	// Inject audio singleton if audio is initialised.
	if audioState != nil {
		sfxPlayers := make(map[string]*audio.Player, len(audioState.SFXData))
		for name, raw := range audioState.SFXData {
			p, err := audioState.Ctx.NewPlayerF32(bytes.NewReader(raw))
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
