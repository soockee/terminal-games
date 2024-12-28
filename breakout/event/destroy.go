package event

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type Destroy struct {
	CollideWith *donburi.Entry
}

var DestroyEvent = events.NewEventType[*Destroy]()
