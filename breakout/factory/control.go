package factory

import (
	input "github.com/quasilyte/ebitengine-input"
	"github.com/soockee/terminal-games/breakout/archetype"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var inputMap = KeyboardMovementGameplay()

func CreateControl(ecs *ecs.ECS) *donburi.Entry {
	control := archetype.Controls.Spawn(ecs)
	component.Control.SetValue(control, component.ControlData{
		InputHandler: component.InputSytem.NewHandler(0, inputMap),
	})
	return control
}

func KeyboardMovementGameplay() input.Keymap {
	return input.Keymap{
		component.ActionReleaseBall: {input.KeyGamepadA, input.KeySpace},
		component.ActionMoveLeft:    {input.KeyA, input.KeyLeft, input.KeyGamepadLStickLeft},
		component.ActionMoveRight:   {input.KeyD, input.KeyRight, input.KeyGamepadLStickRight},

		component.ActionClick: {input.KeyTouchTap, input.KeyMouseLeft},
		component.ActionDebug: {input.KeyF1},
		component.ActionHelp:  {input.KeyF2},
	}
}
