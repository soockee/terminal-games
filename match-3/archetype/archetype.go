package archetype

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/game-template/component"
	"github.com/yohamta/donburi"
)

// NewSpace creates a singleton entity holding the resolv.Space.
func NewSpace(w donburi.World, width, height, cellW, cellH int) {
	e := w.Entry(w.Create(component.Space))
	space := resolv.NewSpace(width, height, cellW, cellH)
	component.Space.Set(e, space)
}

// NewDebug creates the singleton debug toggle entity.
func NewDebug(w donburi.World) {
	e := w.Entry(w.Create(component.Debug))
	component.Debug.Set(e, &component.DebugData{Enabled: false})
}

// NewCamera creates the singleton camera entity with the given virtual screen scale.
func NewCamera(w donburi.World, scaleX, scaleY float64) {
	e := w.Entry(w.Create(component.Camera))
	component.Camera.Set(e, &component.CameraData{ScaleX: scaleX, ScaleY: scaleY})
}

// NewScore creates the singleton score entity with a target for win detection.
func NewScore(w donburi.World, target int) {
	e := w.Entry(w.Create(component.Score))
	component.Score.Set(e, &component.ScoreData{Target: target})
}

// NewGameOver creates the singleton game-over state entity.
func NewGameOver(w donburi.World) {
	e := w.Entry(w.Create(component.GameOver))
	component.GameOver.Set(e, &component.GameOverData{})
}

// NewAudio creates the singleton audio entity with pre-decoded players.
func NewAudio(w donburi.World, ctx *audio.Context, bgm *audio.Player, sfx map[string]*audio.Player) {
	e := w.Entry(w.Create(component.Audio))
	component.Audio.Set(e, &component.AudioData{Ctx: ctx, BGMusic: bgm, SFX: sfx})
}
