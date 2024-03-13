package event

import (
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi/features/events"
)

type Collect struct {
	Type component.CollectableType
}

var CollectEvent = events.NewEventType[*Collect]()
