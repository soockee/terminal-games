package component

import (
	"time"

	"github.com/yohamta/donburi"
)

type GameData struct {
	IsGameOver bool
	TotalScore int
	TotalTime  time.Duration
	Score      int
	Start      time.Time
	End        time.Time
}

var GameState = donburi.NewComponentType[GameData]()
