package archetype

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/soockee/terminal-games/match-3/component"
	"github.com/soockee/terminal-games/match-3/tags"
	"github.com/yohamta/donburi"
)

// NewBoard creates the singleton Board entity.
func NewBoard(w donburi.World, data *component.BoardData) *donburi.Entry {
	e := w.Entry(w.Create(component.Board))
	component.Board.Set(e, data)
	return e
}

// NewTile creates a gem tile entity with all required components.
func NewTile(w donburi.World, col, row int, color int, img *ebiten.Image, px, py float64) *donburi.Entry {
	e := w.Entry(w.Create(
		tags.Tile,
		component.GridPos,
		component.GemType,
		component.PixelPos,
		component.Tween,
		component.Sprite,
	))
	component.GridPos.Set(e, &component.GridPosData{Col: col, Row: row})
	component.GemType.Set(e, &component.GemTypeData{Color: color})
	component.PixelPos.Set(e, &component.PixelPosData{X: px, Y: py})
	component.Tween.Set(e, &component.TweenData{})
	component.Sprite.Set(e, &component.SpriteData{Image: img})
	return e
}

// NewDebug creates the singleton debug toggle entity.
func NewDebug(w donburi.World) {
	e := w.Entry(w.Create(component.Debug))
	component.Debug.Set(e, &component.DebugData{Enabled: false})
}

// NewCamera creates the singleton camera entity.
func NewCamera(w donburi.World, scaleX, scaleY float64) {
	e := w.Entry(w.Create(component.Camera))
	component.Camera.Set(e, &component.CameraData{ScaleX: scaleX, ScaleY: scaleY})
}

// NewScore creates the singleton score entity.
func NewScore(w donburi.World, target int) {
	e := w.Entry(w.Create(component.Score))
	component.Score.Set(e, &component.ScoreData{Target: target})
}

// NewGameState creates the singleton game state entity.
func NewGameState(w donburi.World) {
	e := w.Entry(w.Create(component.GameState))
	component.GameState.Set(e, &component.GameStateData{})
}

// NewAudio creates the singleton audio entity.
func NewAudio(w donburi.World, ctx *audio.Context, bgm *audio.Player, sfx map[string]*audio.Player) {
	e := w.Entry(w.Create(component.Audio))
	component.Audio.Set(e, &component.AudioData{Ctx: ctx, BGMusic: bgm, SFX: sfx})
}
