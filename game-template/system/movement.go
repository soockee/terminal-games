package system

import (
	"github.com/soockee/terminal-games/game-template/component"
	"github.com/yohamta/donburi/ecs"
)

// UpdateMovement applies velocity and physics to game entities.
// Gates on Started && !Dead && !Won so nothing moves before the game begins
// or after it ends.
func UpdateMovement(e *ecs.ECS) {
	entry, ok := component.GameOver.First(e.World)
	if !ok {
		return
	}
	go_ := component.GameOver.Get(entry)
	if !go_.Started || go_.Dead || go_.Won {
		return
	}

	// TODO: add your movement/physics logic here.
}
