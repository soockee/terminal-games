package factory

import (
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateSettings(ecs *ecs.ECS) *donburi.Entry {
	settings := archetype.Settings.Spawn(ecs)
	component.Settings.SetValue(settings, component.SettingsData{
		ShowHelpText: false,
	})

	return settings
}
