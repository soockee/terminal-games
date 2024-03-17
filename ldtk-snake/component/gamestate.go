package component

import (
	"time"

	"github.com/yohamta/donburi"
)

type GameData struct {
	IsGameOver bool
	Score      int
	Start      time.Time
	End        time.Time
}

var GameState = donburi.NewComponentType[GameData]()
