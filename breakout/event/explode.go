package event

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type Explode struct {
	Brick *donburi.Entry
}

var ExplodeEvent = events.NewEventType[*Explode]()
