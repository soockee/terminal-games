package event

import (
	"github.com/yohamta/donburi/features/events"
)

type SceneStateData struct {
	CurrentScene string
	NextScene    string
}

var SceneStateEvent = events.NewEventType[*SceneStateData]()
