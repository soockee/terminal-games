package event

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type Collide struct {
	Type donburi.IComponentType
}

var CollideEvent = events.NewEventType[*Collide]()
