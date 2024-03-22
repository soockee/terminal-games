package scene

import (
	"log/slog"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/ldtkgo"
	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/factory"
	resolv "github.com/soockee/terminal-games/ldtk-snake/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type Scene interface {
	configure()
	GetId() string
	getLdtkProject() *assets.LDtkProject
	getEcs() *ecs.ECS
	getOnce() *sync.Once
}

var TagsMapping = map[string]func(*ecs.ECS, *assets.LDtkProject, *ldtkgo.Entity) *donburi.Entry{
	tags.Snake.Name():     factory.CreateSnake,
	tags.Wall.Name():      factory.CreateWall,
	tags.Button.Name():    factory.CreateButton,
	tags.TextField.Name(): factory.CreateTextField,
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
	project.Renderer.Render(project.Project.LevelByIdentifier(sceneId))
	switch sceneId {
	case component.StartScene:
		return NewStartScene(ecs, project)
	case component.LevelClearScene:
		return NewLevelClearScene(ecs, project)
	case component.GameOverScene:
		return NewGameOverScene(ecs, project)
	case component.Level_0:
		fallthrough
	case component.Level_1:
		fallthrough
	case component.Level_2:
		fallthrough
	case component.Level_3:
		fallthrough
	case component.Level_4:
		return NewSnakeScene(ecs, project, sceneId)

	default:
		slog.Error("invalid sceneId for creation", slog.Any("sceneId", sceneId))
		panic(0)
	}
}

func CreateEntities[T Scene](s T, space *donburi.Entry) {
	level := s.getLdtkProject().Project.LevelByIdentifier(s.GetId())
	entities := s.getLdtkProject().GetEntities(level.Identifier)

	for _, entity := range entities {
		for name, f := range TagsMapping {
			for _, ldtkTag := range entity.Tags {
				if name == ldtkTag {
					resolv.Add(space, f(s.getEcs(), s.getLdtkProject(), entity))
				}
			}
		}
	}
}

func DrawLevel[T Scene](s T, screen *ebiten.Image) {
	level := s.getLdtkProject().Project.LevelByIdentifier(s.GetId())

	if level.BGImage != nil {
		opt := &ebiten.DrawImageOptions{}
		bgImage := level.BGImage
		opt.GeoM.Translate(-bgImage.CropRect[0], -bgImage.CropRect[1])
		opt.GeoM.Scale(bgImage.ScaleX, bgImage.ScaleY)
		img := s.getLdtkProject().Renderer.BGImage
		screen.DrawImage(img, opt)
	}

	for _, layer := range s.getLdtkProject().Renderer.RenderedLayers {
		if s.getLdtkProject().ActiveLayers[layer.Layer.Identifier] {
			screen.DrawImage(layer.Image, &ebiten.DrawImageOptions{})
		}
	}
}
