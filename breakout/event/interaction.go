package event

import (
	input "github.com/quasilyte/ebitengine-input"
	"github.com/yohamta/donburi/features/events"
)

type Interaction struct {
	Action   input.Action
	Position input.Vec
}

var InteractionEvent = events.NewEventType[*Interaction]()
