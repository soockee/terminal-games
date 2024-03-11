package system

import (
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/events"
)

func ProcessEvents(ecs *ecs.ECS) {
	events.ProcessAllEvents(ecs.World)
}
