package systems

import (
	"image"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/ldtkgo"
	"github.com/soockee/terminal-games/ldtk-snake/helper"
)

const (
	ActionMoveUp input.Action = iota
	ActionMoveDown
	ActionMoveRight
	ActionMoveLeft
	ActionClick
)

type Snake struct {
	input     *input.Handler
	snakehead *ldtkgo.Entity
	direction input.Action // todo: not used yet
	tile      *ebiten.Image
}

type SnakeBody struct {
	next *ldtkgo.Entity // todo: not used yet
}

func NewSnake(head *ldtkgo.Entity, renderer *helper.EbitenRenderer) *Snake {
	snakeheadKeymap := input.Keymap{
		ActionMoveUp:    {input.KeyGamepadUp, input.KeyUp, input.KeyW},
		ActionMoveDown:  {input.KeyGamepadDown, input.KeyDown, input.KeyS},
		ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
		ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
		ActionClick:     {input.KeyTouchTap, input.KeyMouseLeft},
	}

	tileset := renderer.Tilesets[head.TileRect.Tileset.Path]
	tileRect := head.TileRect
	tile := tileset.SubImage(image.Rect(tileRect.X, tileRect.Y, tileRect.X+tileRect.W, tileRect.Y+tileRect.H)).(*ebiten.Image)

	snakehead := &Snake{
		input:     InputSytem.NewHandler(0, snakeheadKeymap),
		snakehead: head,
		tile:      tile,
	}

	return snakehead
}

func (s *Snake) move(action input.Action) {
	switch action {
	case ActionMoveLeft:
		s.snakehead.Position[0]--
	case ActionMoveRight:
		s.snakehead.Position[0]++
	case ActionMoveUp:
		s.snakehead.Position[1]--
	case ActionMoveDown:
		s.snakehead.Position[1]++
	default:
		slog.Warn("invalid move key")
	}
}

func (s *Snake) Update() {
	if s.input.ActionIsPressed(ActionMoveUp) {
		slog.Debug("Up")
		s.move(ActionMoveUp)
	}
	if s.input.ActionIsPressed(ActionMoveDown) {
		slog.Debug("Down")
		s.move(ActionMoveDown)
	}
	if s.input.ActionIsPressed(ActionMoveLeft) {
		slog.Debug("Left")
		s.move(ActionMoveLeft)
	}
	if s.input.ActionIsPressed(ActionMoveRight) {
		slog.Debug("Right")
		s.move(ActionMoveRight)
	}
}

func (s *Snake) Draw(screen *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(s.snakehead.Position[0]), float64(s.snakehead.Position[1]))
	screen.DrawImage(s.tile, opt)
}
