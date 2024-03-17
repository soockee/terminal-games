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
	GetId() component.SceneId
	getLdtkProject() *assets.LDtkProject
	getLevelId() int
	getEcs() *ecs.ECS
	configure()
	Once() *sync.Once
}

var TagsMapping = map[string]func(*ecs.ECS, *assets.LDtkProject, *ldtkgo.Entity) *donburi.Entry{
	tags.Snake.Name():  factory.CreateSnake,
	tags.Wall.Name():   factory.CreateWall,
	tags.Button.Name(): factory.CreateButton,
}

func Update(s Scene) error {
	s.Once().Do(s.configure)
	s.getEcs().Update()
	return nil
}

func Draw(s Scene, screen *ebiten.Image) {
	DrawLevel(s, screen)
	s.getEcs().Draw(screen)
}

func Layout(s Scene) (int, int) {
	return s.getLdtkProject().Project.Levels[s.getLevelId()].Width, s.getLdtkProject().Project.Levels[s.getLevelId()].Height
}

func CreateScene(sceneId component.SceneId, ecs *ecs.ECS, project *assets.LDtkProject) Scene {
	project.Renderer.Render(project.Project.Levels[sceneId])
	switch sceneId {
	case component.SnakeScene:
		return NewSnakeScene(ecs, project)
	case component.StartScene:
		return NewStartScene(ecs, project)
	default:
		slog.Error("invalid sceneId for creation", slog.Any("sceneId", sceneId))
		panic(0)
	}
}

func CreateEntities[T Scene](s T, space *donburi.Entry) {
	entities := s.getLdtkProject().GetEntities(s.getLevelId())

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
	level := s.getLdtkProject().Project.Levels[s.getLevelId()]

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
