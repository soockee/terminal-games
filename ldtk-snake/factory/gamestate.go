package factory

import (
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateGamestate(ecs *ecs.ECS, scene *component.Scene) *donburi.Entry {
	gamestate := archetype.Gamestate.Spawn(ecs)
	if scene != nil {
		component.Gamestate.SetValue(gamestate, component.GamestateData{
			CurrentScene: *scene,
		})
	}
	return gamestate
}
