package event

import (
	input "github.com/quasilyte/ebitengine-input"
	"github.com/yohamta/donburi/features/events"
)

type Move struct {
	Action input.Action
}

var MoveEvent = events.NewEventType[*Move]()
