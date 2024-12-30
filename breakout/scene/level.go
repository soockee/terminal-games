package scene

import (
	"sync"

	"github.com/soockee/terminal-games/breakout/archetype"
	pkgevents "github.com/soockee/terminal-games/breakout/event"

	"github.com/soockee/terminal-games/breakout/layers"
	"github.com/soockee/terminal-games/breakout/system"
	"github.com/yohamta/donburi/ecs"
)

type LevelScene struct {
	ecs   *ecs.ECS
	once  *sync.Once
	level string
}

func NewLevelScene(ecs *ecs.ECS, level string) *LevelScene {
	return &LevelScene{
		ecs:   ecs,
		once:  &sync.Once{},
		level: level,
	}

}

func (s *LevelScene) configure() {
	s.ecs.AddSystem(system.UpdatePlayer)
	s.ecs.AddSystem(system.UpdateBall)
	s.ecs.AddSystem(system.UpdateBrick)

	s.ecs.AddRenderer(layers.Default, system.DrawWall)
	s.ecs.AddRenderer(layers.Default, system.DrawPlayer)
	s.ecs.AddRenderer(layers.Default, system.DrawBall)
	s.ecs.AddRenderer(layers.Default, system.DrawBrick)

	// start gametime
	archetype.NewGameState(s.ecs)

	// Subscribe events.
	pkgevents.UpdateSettingEvent.Subscribe(s.ecs.World, system.OnSettingsEvent)
	pkgevents.MoveEvent.Subscribe(s.ecs.World, system.OnMoveEvent)
	pkgevents.CollideEvent.Subscribe(s.ecs.World, system.OnCollideEvent)
	pkgevents.ReleaseEvent.Subscribe(s.ecs.World, system.OnReleaseEvent)

}

func (s *LevelScene) GetId() string {
	return s.level
}
func (s *LevelScene) getEcs() *ecs.ECS {
	return s.ecs
}
func (s *LevelScene) getOnce() *sync.Once {
	return s.once
}
