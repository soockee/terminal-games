package scene

import (
	"log/slog"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/breakout/archetype"
	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/layers"
	"github.com/yohamta/donburi/ecs"
)

type Scene interface {
	configure()
	GetId() string
	getLdtkProject() *assets.LDtkProject
	getEcs() *ecs.ECS
	getOnce() *sync.Once
}

func Update(s Scene) error {
	s.getOnce().Do(s.configure)
	s.getEcs().Update()
	return nil
}

func Draw(s Scene, screen *ebiten.Image) {
	DrawLevel(s, screen)
	s.getEcs().Draw(screen)
}

func Layout(s Scene) (int, int) {
	return s.getLdtkProject().Project.LevelByIdentifier(s.GetId()).Width, s.getLdtkProject().Project.LevelByIdentifier(s.GetId()).Height
}

func CreateScene(sceneId string, ecs *ecs.ECS, project *assets.LDtkProject) Scene {
	level := project.Project.LevelByIdentifier(sceneId)

	cellWidth := level.Width / level.Layers[layers.Default].CellWidth
	CellHeight := level.Height / level.Layers[layers.Default].CellHeight
	archetype.NewSpace(
		ecs.World,
		level.Width,
		level.Height,
		cellWidth,
		CellHeight,
	)

	switch sceneId {
	case component.StartScene:
		return NewStartScene(ecs, project)
	case component.LevelClearScene:
		return NewLevelClearScene(ecs, project)
	case component.GameOverScene:
		return NewGameOverScene(ecs, project)
	case component.Level_0:
		return NewLevelScene(ecs, project, sceneId)

	default:
		slog.Error("invalid sceneId for creation", slog.Any("sceneId", sceneId))
		panic(0)
	}
}

func CreateEntities[T Scene](s T) {
	level := s.getLdtkProject().Project.LevelByIdentifier(s.GetId())
	entities := s.getLdtkProject().GetEntities(level.Identifier)

	for _, entity := range entities {
		for _, ldtkTag := range entity.Tags {
			if archetype.TagsMapping[ldtkTag] == nil {
				slog.Error("tag not found: noop", slog.Any("tag", ldtkTag))
				continue
			}
			archetype.TagsMapping[ldtkTag](s.getEcs().World, s.getLdtkProject(), entity)
		}
	}
}

func DrawLevel[T Scene](s T, screen *ebiten.Image) {
	level := s.getLdtkProject().Project.LevelByIdentifier(s.GetId())
	opt := assets.NewDefaultDrawOptions()
	s.getLdtkProject().Renderer.Render(level, screen, opt)
}
