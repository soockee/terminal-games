package event

import (
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi/features/events"
)

type Move struct {
	Action   input.Action
	Position resolv.Vector
	Boost    bool
}

var MoveEvent = events.NewEventType[*Move]()
