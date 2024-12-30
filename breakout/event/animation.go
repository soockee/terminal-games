package event

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type AnimationEnd struct {
	Animation *donburi.Entry
}

var AnimationEndEvent = events.NewEventType[*AnimationEnd]()
