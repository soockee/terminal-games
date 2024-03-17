package event

import (
	"github.com/yohamta/donburi/features/events"
)

type CollideWithType int

const (
	CollideWall CollideWithType = iota
	CollideBody
)

type Collide struct {
	Type CollideWithType
}

var CollideEvent = events.NewEventType[*Collide]()
