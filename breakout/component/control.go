package component

import (
	"log/slog"

	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

const (
	ActionMoveBoost input.Action = iota
	ActionMoveLeft
	ActionMoveRight
	ActionMoveMobile
	ActionClick
	ActionDebug
	ActionHelp
)

type ControlData struct {
	InputHandler *input.Handler
	LastPosition *resolv.Vector
}

var Control = donburi.NewComponentType[ControlData]()

var InputSytem *input.System = &input.System{}

func init() {
	slog.Info("initialize inputsystem")
	InputSytem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
}
