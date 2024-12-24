package scene

import (
	"sync"

	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"
	pkgevents "github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/factory"
	"github.com/soockee/terminal-games/breakout/layers"
	"github.com/soockee/terminal-games/breakout/system"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	decs "github.com/yohamta/donburi/ecs"
)

type GameOverScene struct {
	ecs         *decs.ECS
	ldtkProject *assets.LDtkProject
	once        *sync.Once
}

func NewGameOverScene(ecs *decs.ECS, project *assets.LDtkProject) *GameOverScene {
	return &GameOverScene{
		ecs:         ecs,
		ldtkProject: project,
		once:        &sync.Once{},
	}
}

func (s *GameOverScene) configure() {
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

	CreateEntities(s, space)

	//gamedata := component.GameState.Get(component.GameState.MustFirst(s.ecs.World))

	component.Text.Each(s.ecs.World, func(e *donburi.Entry) {
		//textfield := component.Text.Get(e)
	})

	// Subscribe events.
	pkgevents.UpdateSettingEvent.Subscribe(s.ecs.World, system.OnSettingsEvent)
	pkgevents.InteractionEvent.Subscribe(s.ecs.World, system.HandleButtonClick)
}

func (s *GameOverScene) GetId() string {
	return component.GameOverScene
}
func (s *GameOverScene) getLdtkProject() *assets.LDtkProject {
	return s.ldtkProject
}
func (s *GameOverScene) getEcs() *ecs.ECS {
	return s.ecs
}
func (s *GameOverScene) getOnce() *sync.Once {
	return s.once
}
