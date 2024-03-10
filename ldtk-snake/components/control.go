package components

import (
	"log/slog"

	input "github.com/quasilyte/ebitengine-input"
	"github.com/yohamta/donburi"
)

const (
	ActionMoveUp input.Action = iota
	ActionMoveDown
	ActionMoveRight
	ActionMoveLeft
	ActionClick
	ActionDebug
	ActionHelp
)
const (
	SnakeHandler uint8 = iota
	SettingsHandler
)

type ControlData struct {
	InputHandler *input.Handler
}

var Control = donburi.NewComponentType[ControlData]()

var InputSytem *input.System = &input.System{}

func init() {
	slog.Info("initialize inputsystem")
	InputSytem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
}
