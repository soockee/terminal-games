package factory

import (
	input "github.com/quasilyte/ebitengine-input"
	"github.com/soockee/terminal-games/ldtk-snake/archetypes"
	"github.com/soockee/terminal-games/ldtk-snake/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateSettings(ecs *ecs.ECS) *donburi.Entry {
	settings := archetypes.Settings.Spawn(ecs)
	components.Settings.SetValue(settings, components.SettingsData{
		ShowHelpText: true,
	})

	inputHandler := components.InputSytem.NewHandler(components.SettingsHandler, input.Keymap{
		components.ActionDebug: {input.KeyF1},
		components.ActionHelp:  {input.KeyF1},
	})
	components.Control.SetValue(settings, components.ControlData{
		InputHandler: inputHandler,
	})
	return settings
}
