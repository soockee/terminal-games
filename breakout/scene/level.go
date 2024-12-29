package scene

import (
	"sync"

	"github.com/soockee/terminal-games/breakout/archetype"
	"github.com/soockee/terminal-games/breakout/assets"
	pkgevents "github.com/soockee/terminal-games/breakout/event"

	"github.com/soockee/terminal-games/breakout/layers"
	"github.com/soockee/terminal-games/breakout/system"
	"github.com/yohamta/donburi/ecs"
)

type LevelScene struct {
	ecs         *ecs.ECS
	ldtkProject *assets.LDtkProject
	once        *sync.Once
	level       string
}

func NewLevelScene(ecs *ecs.ECS, project *assets.LDtkProject, level string) *LevelScene {
	return &LevelScene{
		ecs:         ecs,
		ldtkProject: project,
		once:        &sync.Once{},
		level:       level,
	}

}

func (s *LevelScene) configure() {
	s.ecs.AddSystem(system.ProcessEvents)
	s.ecs.AddSystem(system.UpdatePlayer)
	s.ecs.AddSystem(system.UpdateBall)
	s.ecs.AddSystem(system.UpdateBrick)
	s.ecs.AddSystem(system.UpdateExplosion)

	s.ecs.AddRenderer(layers.Default, system.DrawWall)
	s.ecs.AddRenderer(layers.Default, system.DrawPlayer)
	s.ecs.AddRenderer(layers.Default, system.DrawBall)
	s.ecs.AddRenderer(layers.Default, system.DrawBrick)
	s.ecs.AddRenderer(layers.Default, system.DrawExplosion)

	CreateEntities(s)
	// start gametime
	archetype.NewGameState(s.ecs)

	// Subscribe events.
	pkgevents.UpdateSettingEvent.Subscribe(s.ecs.World, system.OnSettingsEvent)
	pkgevents.MoveEvent.Subscribe(s.ecs.World, system.OnMoveEvent)
	pkgevents.CollideEvent.Subscribe(s.ecs.World, system.OnCollideEvent)
	pkgevents.ExplodeEvent.Subscribe(s.ecs.World, system.OnExplodeEvent)
	pkgevents.ReleaseEvent.Subscribe(s.ecs.World, system.OnReleaseEvent)

}

func (s *LevelScene) GetId() string {
	return s.level
}
func (s *LevelScene) getLdtkProject() *assets.LDtkProject {
	return s.ldtkProject
}
func (s *LevelScene) getEcs() *ecs.ECS {
	return s.ecs
}
func (s *LevelScene) getOnce() *sync.Once {
	return s.once
}
