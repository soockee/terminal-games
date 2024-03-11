package event

import (
	"github.com/yohamta/donburi/features/events"
)

type Gamestate struct {
}

var GameStateEvent = events.NewEventType[*Gamestate]()
