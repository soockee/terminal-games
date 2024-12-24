package event

import (
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/yohamta/donburi/features/events"
)

type Collect struct {
	Type component.CollectableType
}

var CollectEvent = events.NewEventType[*Collect]()
