package event

import "github.com/yohamta/donburi/features/events"

type CreateEntityData struct {
	// Entity *ldtkgo.Entity
	Tags       []string
	Identifier string
	X          float64
	Y          float64
	W          float64
	H          float64
}

var CreateEntityEvent = events.NewEventType[*CreateEntityData]()
