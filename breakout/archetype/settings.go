package archetype

import (
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var (
	Settings = newArchetype(
		component.Settings,
	)
)

func NewSettings(ecs *ecs.ECS) *donburi.Entry {
	settings := Settings.Spawn(ecs)
	component.Settings.SetValue(settings, component.SettingsData{
		ShowHelpText: false,
	})

	return settings
}
