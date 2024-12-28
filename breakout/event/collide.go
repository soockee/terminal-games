package event

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type Collide struct {
	CollideWith  *donburi.Entry
	Collider     *donburi.Entry
	Intersection resolv.IntersectionSet
}

var CollideEvent = events.NewEventType[*Collide]()
