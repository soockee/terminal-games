package factory

import (
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateGameState(ecs *ecs.ECS) *donburi.Entry {
	gamestate := archetype.GameState.Spawn(ecs)
	component.Gamestate.SetValue(gamestate, component.GameStateData{})
	return gamestate
}
