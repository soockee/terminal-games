package factory

import (
	input "github.com/quasilyte/ebitengine-input"
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var inputMap = input.Keymap{
	component.ActionMoveBoost:    {input.KeyGamepadA, input.KeySpace, input.KeyMouseRight},
	component.ActionMovePosition: {input.KeyMouseLeft, input.KeyTouchDrag},

	component.ActionClick: {input.KeyTouchTap, input.KeyMouseLeft},
	component.ActionDebug: {input.KeyF1},
	component.ActionHelp:  {input.KeyF2},
}

func CreateControl(ecs *ecs.ECS) *donburi.Entry {
	control := archetype.Controls.Spawn(ecs)
	component.Control.SetValue(control, component.ControlData{
		InputHandler: component.InputSytem.NewHandler(0, inputMap),
	})
	return control
}
