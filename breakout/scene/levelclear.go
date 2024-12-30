package scene

import (
	"sync"

	"github.com/soockee/terminal-games/breakout/component"
	pkgevents "github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/layers"
	"github.com/soockee/terminal-games/breakout/system"
	"github.com/yohamta/donburi/ecs"
)

type LevelClearScene struct {
	ecs  *ecs.ECS
	once *sync.Once
}

func NewLevelClearScene(ecs *ecs.ECS) *LevelClearScene {
	return &LevelClearScene{
		ecs:  ecs,
		once: &sync.Once{},
	}
}

func (s *LevelClearScene) configure() {
	s.ecs.AddSystem(system.UpdateButton)

	s.ecs.AddRenderer(layers.Default, system.DrawButton)

	// Subscribe events.
	pkgevents.UpdateSettingEvent.Subscribe(s.ecs.World, system.OnSettingsEvent)
	pkgevents.InteractionEvent.Subscribe(s.ecs.World, system.HandleButtonClick)
}

func (s *LevelClearScene) GetId() string {
	return component.LevelClearScene
}

func (s *LevelClearScene) getEcs() *ecs.ECS {
	return s.ecs
}
func (s *LevelClearScene) getOnce() *sync.Once {
	return s.once
}
