package system

import (
	input "github.com/quasilyte/ebitengine-input"
)

const (
	ActionMoveUp input.Action = iota
	ActionMoveDown
	ActionMoveRight
	ActionMoveLeft
	ActionClick
)

type Snakehead struct {
	input *input.Handler
}

func NewSnake() *Snakehead {
	snakeheadKeymap := input.Keymap{
		ActionMoveUp:    {input.KeyGamepadUp, input.KeyUp, input.KeyW},
		ActionMoveDown:  {input.KeyGamepadDown, input.KeyDown, input.KeyS},
		ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
		ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
		ActionClick:     {input.KeyTouchTap, input.KeyMouseLeft},
	}

	snakehead := &Snakehead{
		input: InputSytem.NewHandler(0, snakeheadKeymap),
	}

	return snakehead
}
