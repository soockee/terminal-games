package factory

import (
	"time"

	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateGameState(ecs *ecs.ECS) *donburi.Entry {
	gamestate := archetype.GameState.Spawn(ecs)
	start := time.Now()
	component.GameState.SetValue(gamestate, component.GameData{
		IsGameOver: false,
		Score:      0,
		Start:      start,
		End:        start,
	})
	return gamestate
}

func FinalizeGameState(ecs *ecs.ECS, gamedata *component.GameData) *donburi.Entry {
	gamestate := archetype.GameState.Spawn(ecs)
	component.GameState.SetValue(gamestate, component.GameData{
		IsGameOver: false,
		Score:      gamedata.Score,
		Start:      gamedata.Start,
		End:        time.Now(),
	})
	return gamestate
}

func ResetGameState(ecs *ecs.ECS) *donburi.Entry {
	gameState, _ := component.GameState.First(ecs.World)
	ecs.World.Remove(gameState.Entity())
	gamestate := archetype.GameState.Spawn(ecs)
	start := time.Now()
	component.GameState.SetValue(gamestate, component.GameData{
		IsGameOver: false,
		Score:      0,
		Start:      start,
		End:        start,
	})
	return gamestate
}
