package scene

import (
	"sync"

	"github.com/soockee/terminal-games/breakout/component"
	pkgevents "github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/layers"
	"github.com/soockee/terminal-games/breakout/system"
	"github.com/yohamta/donburi/ecs"
)

type GameOverScene struct {
	ecs  *ecs.ECS
	once *sync.Once
}

func NewGameOverScene(ecs *ecs.ECS) *GameOverScene {
	return &GameOverScene{
		ecs:  ecs,
		once: &sync.Once{},
	}
}

func (s *GameOverScene) configure() {
	s.ecs.AddSystem(system.UpdateButton)

	s.ecs.AddRenderer(layers.Default, system.DrawTextField)
	s.ecs.AddRenderer(layers.Default, system.DrawButton)

	// Subscribe events.
	pkgevents.UpdateSettingEvent.Subscribe(s.ecs.World, system.OnSettingsEvent)
	pkgevents.InteractionEvent.Subscribe(s.ecs.World, system.HandleButtonClick)
}

func (s *GameOverScene) GetId() string {
	return component.GameOverScene
}

func (s *GameOverScene) getEcs() *ecs.ECS {
	return s.ecs
}
func (s *GameOverScene) getOnce() *sync.Once {
	return s.once
}
