package factory

import (
	"github.com/soockee/terminal-games/breakout/archetype"
	"github.com/soockee/terminal-games/breakout/component"
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
