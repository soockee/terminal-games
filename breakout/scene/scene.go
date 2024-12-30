package scene

import (
	"log/slog"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/breakout/archetype"
	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/event"
	"github.com/soockee/terminal-games/breakout/layers"
	"github.com/yohamta/donburi/ecs"
)

type Scene interface {
	configure()
	GetId() string
	getEcs() *ecs.ECS
	getOnce() *sync.Once
}

func Update(s Scene) error {
	s.getOnce().Do(s.configure)
	s.getEcs().Update()
	return nil
}

func Draw(s Scene, screen *ebiten.Image) {
	s.getEcs().Draw(screen)
}

func Layout(height, width int) (int, int) {
	return height, width
}

func CreateScene(sceneId string, ecs *ecs.ECS, project *assets.LDtkProject) Scene {
	level := project.Project.LevelByIdentifier(sceneId)
	entities := project.GetEntities(level.Identifier)

	archetype.NewSpace(ecs.World, level.Width, level.Height, level.Layers[layers.Default].CellWidth, level.Layers[layers.Default].CellHeight)
	// Create entities
	for _, entity := range entities {
		data := &event.CreateEntityData{
			Tags:       entity.Tags,
			Identifier: entity.Identifier,
			X:          float64(entity.Position[0]),
			Y:          float64(entity.Position[1]),
			W:          float64(entity.Width),
			H:          float64(entity.Height),
		}
		event.CreateEntityEvent.Publish(ecs.World, data)
	}

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
		return NewStartScene(ecs)
	case component.LevelClearScene:
		return NewLevelClearScene(ecs)
	case component.GameOverScene:
		return NewGameOverScene(ecs)
	case component.Level_0:
		fallthrough
	case component.Level_1:
		return NewLevelScene(ecs, sceneId)

	default:
		slog.Error("invalid sceneId for creation", slog.Any("sceneId", sceneId))
		panic(0)
	}
}
