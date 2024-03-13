package scene

import (
	"sync"

	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	pkgevents "github.com/soockee/terminal-games/ldtk-snake/event"
	"github.com/soockee/terminal-games/ldtk-snake/factory"
	"github.com/soockee/terminal-games/ldtk-snake/layers"
	"github.com/soockee/terminal-games/ldtk-snake/system"
	"github.com/yohamta/donburi/ecs"
	decs "github.com/yohamta/donburi/ecs"
)

type StartScene struct {
	ecs         *decs.ECS
	ldtkProject *assets.LDtkProject
	once        *sync.Once
}

func NewStartScene(ecs *decs.ECS, project *assets.LDtkProject) *StartScene {
	return &StartScene{
		ecs:         ecs,
		ldtkProject: project,
		once:        &sync.Once{},
	}
}

func (s *StartScene) configure() {
	s.ecs.AddSystem(system.UpdateObjects)
	s.ecs.AddSystem(system.ProcessEvents)
	s.ecs.AddSystem(system.UpdateButton)

	s.ecs.AddRenderer(layers.Default, system.DrawDebug)
	s.ecs.AddRenderer(layers.Default, system.DrawHelp)
	s.ecs.AddRenderer(layers.Default, system.DrawButton)

	cellWidth := s.ldtkProject.Project.Levels[s.getLevelId()].Width / s.ldtkProject.Project.Levels[s.getLevelId()].Layers[s.getLevelId()].CellWidth
	CellHeight := s.ldtkProject.Project.Levels[s.getLevelId()].Height / s.ldtkProject.Project.Levels[s.getLevelId()].Layers[s.getLevelId()].CellHeight
	space := factory.CreateSpace(
		s.ecs,
		s.ldtkProject.Project.Levels[s.getLevelId()].Width,
		s.ldtkProject.Project.Levels[s.getLevelId()].Height,
		cellWidth,
		CellHeight,
	)

	// s.createEntities(s.ecs, space)
	CreateEntities(s, space)

	// Subscribe events.
	pkgevents.UpdateSettingEvent.Subscribe(s.ecs.World, system.OnSettingsEvent)
	pkgevents.InteractionEvent.Subscribe(s.ecs.World, system.HandleButtonClick)
}

func (s *StartScene) GetId() component.SceneId {
	return component.StartScene
}
func (s *StartScene) getLevelId() int {
	return int(s.GetId())
}
func (s *StartScene) getLdtkProject() *assets.LDtkProject {
	return s.ldtkProject
}
func (s *StartScene) getEcs() *ecs.ECS {
	return s.ecs
}
func (s *StartScene) Once() *sync.Once {
	return s.once
}
