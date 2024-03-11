package factory

import (
	"time"

	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/engine"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

func CreateCamera(ecs *ecs.ECS, startPosition math.Vec2) *donburi.Entry {
	camera := archetype.Camera.Spawn(ecs)
	cameraCamera := component.Camera.Get(camera)
	cameraCamera.MoveTimer = engine.NewTimer(time.Second * 3)
	transform.Transform.Get(camera).LocalPosition = startPosition
	return camera
}
