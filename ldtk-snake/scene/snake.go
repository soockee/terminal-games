package scene

import (
	"sync"

	"github.com/soockee/terminal-games/ldtk-snake/assets"
	pkgevents "github.com/soockee/terminal-games/ldtk-snake/event"
	"github.com/soockee/terminal-games/ldtk-snake/tags"

	"github.com/soockee/terminal-games/ldtk-snake/factory"
	"github.com/soockee/terminal-games/ldtk-snake/layers"
	"github.com/soockee/terminal-games/ldtk-snake/system"
	"github.com/yohamta/donburi/ecs"
	decs "github.com/yohamta/donburi/ecs"
)

type SnakeScene struct {
	ecs         *decs.ECS
	ldtkProject *assets.LDtkProject
	once        *sync.Once
	level       string
}

func NewSnakeScene(ecs *decs.ECS, project *assets.LDtkProject, level string) *SnakeScene {
	return &SnakeScene{
		ecs:         ecs,
		ldtkProject: project,
		once:        &sync.Once{},
		level:       level,
	}

}

func (s *SnakeScene) configure() {
	s.ecs.AddSystem(system.UpdateSnake)
	s.ecs.AddSystem(system.ProcessEvents)
	s.ecs.AddSystem(system.UpdateFood)
	s.ecs.AddSystem(system.UpdateObjects)

	s.ecs.AddRenderer(layers.Default, system.DrawSnake)
	s.ecs.AddRenderer(layers.Default, system.DrawFood)
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

	factory.CreateFood(s.ecs.World, s.ldtkProject, s.ldtkProject.Project.EntityDefinitionByIdentifier(tags.Food.Name()))

	// Subscribe events.
	pkgevents.UpdateSettingEvent.Subscribe(s.ecs.World, system.OnSettingsEvent)
	pkgevents.MoveEvent.Subscribe(s.ecs.World, system.OnMoveEvent)
	pkgevents.CollectEvent.Subscribe(s.ecs.World, system.OnPickupEvent)
	pkgevents.CollideEvent.Subscribe(s.ecs.World, system.OnCollideEvent)

}

func (s *SnakeScene) GetId() string {
	return s.level
}
func (s *SnakeScene) getLdtkProject() *assets.LDtkProject {
	return s.ldtkProject
}
func (s *SnakeScene) getEcs() *ecs.ECS {
	return s.ecs
}
func (s *SnakeScene) getOnce() *sync.Once {
	return s.once
}
