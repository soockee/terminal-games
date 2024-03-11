package event

import (
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi/features/events"
)

type Gamestate struct {
	CurrentScene component.Scene
}

var GamestateEvent = events.NewEventType[*Gamestate]()
