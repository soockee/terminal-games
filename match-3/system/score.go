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

	boardEntry, boardOK := component.BoardRules.First(e.World)
	if !boardOK {
		return
	}
	lvl := component.BoardRules.Get(boardEntry)

	// Timer countdown.
	if lvl.TimeLimit > 0 {
		lvl.TimeRemaining -= 1.0 / 60.0
		if lvl.TimeRemaining <= 0 {
			lvl.TimeRemaining = 0
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
