package archetype

import (
	"time"

	"github.com/soockee/terminal-games/breakout/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var (
	GameState = newArchetype(
		component.GameState,
	)
)

func NewGameState(ecs *ecs.ECS) *donburi.Entry {
	gamestate := GameState.Spawn(ecs)
	start := time.Now()
	component.GameState.SetValue(gamestate, component.GameData{
		IsGameOver: false,
		TotalScore: 0,
		TotalTime:  time.Duration(0),
		Score:      0,
		Start:      start,
		End:        start,
	})
	return gamestate
}

func ContinueLevelGameState(ecs *ecs.ECS, gamedata *component.GameData) *donburi.Entry {
	gamestate := GameState.Spawn(ecs)
	component.GameState.SetValue(gamestate, component.GameData{
		IsGameOver: false,
		// todo: calc total score
		TotalScore: gamedata.TotalScore,
		TotalTime:  gamedata.TotalTime,
		Score:      0,
		Start:      time.Now(),
		End:        time.Now(),
	})
	return gamestate
}

func AccumulateGameState(ecs *ecs.ECS, gamedata *component.GameData) *donburi.Entry {
	gamestate := GameState.Spawn(ecs)
	component.GameState.SetValue(gamestate, component.GameData{
		IsGameOver: false,
		TotalScore: gamedata.TotalScore + gamedata.Score,
		TotalTime:  gamedata.TotalTime + time.Since(gamedata.Start),
		Score:      gamedata.Score,
		Start:      gamedata.Start,
		End:        time.Now(),
	})
	return gamestate
}

func ResetGameState(ecs *ecs.ECS) *donburi.Entry {
	gameState := component.GameState.MustFirst(ecs.World)
	ecs.World.Remove(gameState.Entity())
	gamestate := GameState.Spawn(ecs)
	start := time.Now()
	component.GameState.SetValue(gamestate, component.GameData{
		IsGameOver: false,
		Score:      0,
		Start:      start,
		End:        start,
	})
	return gamestate
}
