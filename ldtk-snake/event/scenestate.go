package event

import (
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi/features/events"
)

type SceneStateData struct {
	CurrentScene component.SceneId
}

var SceneStateEvent = events.NewEventType[*SceneStateData]()
