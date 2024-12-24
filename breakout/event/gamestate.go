package event

import (
	"time"

	"github.com/yohamta/donburi/features/events"
)

type GameStateData struct {
	IsGameOver bool
	Time       time.Duration
	Score      int
}

var GameStateEvent = events.NewEventType[*GameStateData]()
