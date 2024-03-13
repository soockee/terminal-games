package system

import (
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/event"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// move temporarily uses a speed of type int whiel figuring out the collision
func OnSettingsEvent(w donburi.World, e *event.UpdateSetting) {
	entity, _ := component.Settings.First(w)
	settings := component.Settings.Get(entity)

	switch e.Action {
	case component.ActionDebug:
		settings.Debug = !settings.Debug
	case component.ActionHelp:
		settings.ShowHelpText = !settings.ShowHelpText
	}

}

func GetSettings(ecs *ecs.ECS) (*component.SettingsData, bool) {
	ent, ok := component.Settings.First(ecs.World)
	return component.Settings.Get(ent), ok
}
