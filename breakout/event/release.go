package event

import (
	input "github.com/quasilyte/ebitengine-input"
	"github.com/yohamta/donburi/features/events"
)

type Release struct {
	Action input.Action
}

var ReleaseEvent = events.NewEventType[*Release]()
