package scene

import (
	"sync"

	"github.com/soockee/terminal-games/breakout/component"
	pkgevents "github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/layers"
	"github.com/soockee/terminal-games/breakout/system"
	"github.com/yohamta/donburi/ecs"
)

type StartScene struct {
	ecs  *ecs.ECS
	once *sync.Once
}

func NewStartScene(ecs *ecs.ECS) *StartScene {
	return &StartScene{
		ecs:  ecs,
		once: &sync.Once{},
	}
}

func (s *StartScene) configure() {
	s.ecs.AddSystem(system.UpdateButton)

	s.ecs.AddRenderer(layers.Default, system.DrawButton)

	// Subscribe events.
	pkgevents.UpdateSettingEvent.Subscribe(s.ecs.World, system.OnSettingsEvent)
	pkgevents.InteractionEvent.Subscribe(s.ecs.World, system.HandleButtonClick)
}

func (s *StartScene) GetId() string {
	return component.StartScene
}
func (s *StartScene) getEcs() *ecs.ECS {
	return s.ecs
}
func (s *StartScene) getOnce() *sync.Once {
	return s.once
}
