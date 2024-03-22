package event

import (
	"github.com/yohamta/donburi/features/events"
)

type Mouse struct{}

var MouseEvent = events.NewEventType[*Mouse]()
