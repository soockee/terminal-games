package event

import "github.com/yohamta/donburi/features/events"

// AudioEvent carries the name of an SFX to play (e.g. "jump", "hit").
type AudioEventData struct {
	Name string
}

var AudioEvent = events.NewEventType[AudioEventData]()
