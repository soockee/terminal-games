package system

import (
	"github.com/soockee/terminal-games/match-3/component"
	"github.com/yohamta/donburi/ecs"
)

// UpdateScore handles timer countdown and win/loss detection.
func UpdateScore(e *ecs.ECS) {
	gsEntry, gsOK := component.GameState.First(e.World)
	if !gsOK {
		return
	}
	gs := component.GameState.Get(gsEntry)
	if !gs.Started || gs.Won || gs.Dead {
		return
	}

	boardEntry, boardOK := component.Board.First(e.World)
	if !boardOK {
		return
	}
	board := component.Board.Get(boardEntry)

	// Timer countdown.
	if board.TimeLimit > 0 {
		board.TimeRemaining -= 1.0 / 60.0
		if board.TimeRemaining <= 0 {
			board.TimeRemaining = 0
			gs.Dead = true
			return
		}
	}

	// Win condition: score >= target.
	scoreEntry, scoreOK := component.Score.First(e.World)
	if !scoreOK {
		return
	}
	score := component.Score.Get(scoreEntry)
	if score.Value >= score.Target {
		gs.Won = true
	}
}
