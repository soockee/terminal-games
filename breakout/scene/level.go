package scene

import (
	"sync"

	"github.com/soockee/terminal-games/breakout/assets"
	pkgevents "github.com/soockee/terminal-games/breakout/event"

	"github.com/soockee/terminal-games/breakout/factory"
	"github.com/soockee/terminal-games/breakout/layers"
	"github.com/soockee/terminal-games/breakout/system"
	"github.com/yohamta/donburi/ecs"
	decs "github.com/yohamta/donburi/ecs"
)

type LevelScene struct {
	ecs         *decs.ECS
	ldtkProject *assets.LDtkProject
	once        *sync.Once
	level       string
}

func NewLevelScene(ecs *decs.ECS, project *assets.LDtkProject, level string) *LevelScene {
	return &LevelScene{
		ecs:         ecs,
		ldtkProject: project,
		once:        &sync.Once{},
		level:       level,
	}

}

func (s *LevelScene) configure() {
	s.ecs.AddSystem(system.ProcessEvents)
	s.ecs.AddSystem(system.UpdateObjects)

	s.ecs.AddRenderer(layers.Default, system.DrawWall)

	level := s.ldtkProject.Project.LevelByIdentifier(s.GetId())

	cellWidth := level.Width / level.Layers[layers.Default].CellWidth
	CellHeight := level.Height / level.Layers[layers.Default].CellHeight
	space := factory.CreateSpace(
		s.ecs,
		level.Width,
		level.Height,
		cellWidth,
		CellHeight,
	)

	CreateEntities(s, space)
	// start gametime
	factory.CreateGameState(s.ecs)

	// Subscribe events.
	pkgevents.UpdateSettingEvent.Subscribe(s.ecs.World, system.OnSettingsEvent)
	pkgevents.CollectEvent.Subscribe(s.ecs.World, system.OnPickupEvent)
	pkgevents.CollideEvent.Subscribe(s.ecs.World, system.OnCollideEvent)

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
