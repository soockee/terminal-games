package system

import (
	"github.com/soockee/terminal-games/game-template/component"
	"github.com/yohamta/donburi/ecs"
)

// UpdateMovement applies velocity and physics to game entities.
// Gates on Started && !Dead && !Won && !Paused so nothing moves before
// the game begins, after it ends, or while paused.
func UpdateMovement(e *ecs.ECS) {
	entry, ok := component.GameState.First(e.World)
	if !ok {
		return
	}
	gs := component.GameState.Get(entry)
	if !gs.Started || gs.Dead || gs.Won || gs.Paused {
		return
	}

	// TODO: add your movement/physics logic here.
}
