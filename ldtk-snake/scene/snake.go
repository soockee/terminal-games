package scene

import (
	"sync"

	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/soockee/terminal-games/ldtk-snake/component"
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
}

func NewSnakeScene(ecs *decs.ECS, project *assets.LDtkProject) *SnakeScene {
	return &SnakeScene{
		ecs:         ecs,
		ldtkProject: project,
		once:        &sync.Once{},
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

	cellWidth := s.ldtkProject.Project.Levels[s.getLevelId()].Width / s.ldtkProject.Project.Levels[s.getLevelId()].Layers[layers.Default].CellWidth
	CellHeight := s.ldtkProject.Project.Levels[s.getLevelId()].Height / s.ldtkProject.Project.Levels[s.getLevelId()].Layers[layers.Default].CellHeight
	space := factory.CreateSpace(
		s.ecs,
		s.ldtkProject.Project.Levels[s.getLevelId()].Width,
		s.ldtkProject.Project.Levels[s.getLevelId()].Height,
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

func (s *SnakeScene) GetId() component.SceneId {
	return component.SnakeScene
}
func (s *SnakeScene) getLevelId() int {
	return int(s.GetId())
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
