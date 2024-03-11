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

type SnakeScene struct {
	ecs         *decs.ECS
	ldtkProject *assets.LDtkProject
	once        *sync.Once
}

func NewSnakeScene(ecs *decs.ECS, project *assets.LDtkProject) *SnakeScene {
	return &SnakeScene{
		ecs:         ecs,
		ldtkProject: project,
		once: &sync.Once{},
	}

}

func (s *SnakeScene) configure() {
	// config.C.CurrentLevel = config.LevelMapping[config.SnakeLevel1]
	// config.RenderLevel()

	s.ecs.AddSystem(system.UpdateSnake)
	s.ecs.AddSystem(system.ProcessEvents)
	s.ecs.AddSystem(system.UpdateFood)
	s.ecs.AddSystem(system.UpdateObjects)

	s.ecs.AddRenderer(layers.Default, system.DrawSnake)
	s.ecs.AddRenderer(layers.Default, system.DrawFood)
	s.ecs.AddRenderer(layers.Default, system.DrawWall)
	s.ecs.AddRenderer(layers.Default, system.DrawDebug)

	factory.CreateSettings(s.ecs)

	cellWidth := s.ldtkProject.Project.Levels[s.getLevelId()].Width / s.ldtkProject.Project.Levels[s.getLevelId()].Layers[s.getLevelId()].CellWidth
	CellHeight := s.ldtkProject.Project.Levels[s.getLevelId()].Height / s.ldtkProject.Project.Levels[s.getLevelId()].Layers[s.getLevelId()].CellHeight
	space := factory.CreateSpace(
		s.ecs,
		s.ldtkProject.Project.Levels[s.getLevelId()].Width,
		s.ldtkProject.Project.Levels[s.getLevelId()].Height,
		cellWidth,
		CellHeight,
	)

	CreateEntities(s, space)

	// entities := s.ldtkProject.GetEntities(s.getLevelId())
	// Tags := map[string]func(*decs.ECS, *ebiten.Image, *ldtkgo.Entity) *donburi.Entry{
	// 	tags.Snake.Name(): factory.CreateSnake,
	// 	tags.Wall.Name():  factory.CreateWall,
	// }
	// for _, entity := range entities {
	// 	for name, f := range Tags {
	// 		for _, ldtkTag := range entity.Tags {
	// 			if name == ldtkTag {
	// 				sprite, err := s.ldtkProject.GetSprite(entity)
	// 				if err != nil {
	// 					slog.Error("could not find sprite for entity")
	// 				}
	// 				dresolv.Add(space, f(s.ecs, sprite, entity))
	// 			}
	// 		}
	// 	}
	// }

	// Subscribe events.
	pkgevents.UpdateSettingEvent.Subscribe(s.ecs.World, system.HandleSettingsEvent)
	pkgevents.MoveEvent.Subscribe(s.ecs.World, system.HandleMoveEvent)
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
func (s *SnakeScene) Once() *sync.Once {
	return s.once
}
