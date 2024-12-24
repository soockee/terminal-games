package scene

import (
	"sync"

	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"
	pkgevents "github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/factory"
	"github.com/soockee/terminal-games/breakout/layers"
	"github.com/soockee/terminal-games/breakout/system"
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

	// s.createEntities(s.ecs, space)
	CreateEntities(s, space)

	// Subscribe events.
	pkgevents.UpdateSettingEvent.Subscribe(s.ecs.World, system.OnSettingsEvent)
	pkgevents.InteractionEvent.Subscribe(s.ecs.World, system.HandleButtonClick)
}

func (s *StartScene) GetId() string {
	return component.StartScene
}
func (s *StartScene) getLdtkProject() *assets.LDtkProject {
	return s.ldtkProject
}
func (s *StartScene) getEcs() *ecs.ECS {
	return s.ecs
}
func (s *StartScene) getOnce() *sync.Once {
	return s.once
}
