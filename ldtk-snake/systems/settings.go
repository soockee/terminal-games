package systems

import (
	"log/slog"

	"github.com/soockee/terminal-games/ldtk-snake/components"
	"github.com/yohamta/donburi/ecs"
)

func UpdateSettings(ecs *ecs.ECS) {
	ent, ok := components.Settings.First(ecs.World)
	if !ok {
		slog.Warn("settings object not found")
		return
	}
	settings := components.Settings.Get(ent)
	control := components.Control.Get(ent)

	if control.InputHandler.ActionIsJustPressed(components.ActionDebug) {
		settings.Debug = !settings.Debug
	}
	if control.InputHandler.ActionIsJustPressed(components.ActionDebug) {
		settings.ShowHelpText = !settings.ShowHelpText
	}
}

func GetSettings(ecs *ecs.ECS) (*components.SettingsData, bool) {
	ent, ok := components.Settings.First(ecs.World)
	return components.Settings.Get(ent), ok
}
