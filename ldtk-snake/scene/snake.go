package scene

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/config"
	pkgevents "github.com/soockee/terminal-games/ldtk-snake/event"
	"github.com/soockee/terminal-games/ldtk-snake/factory"
	"github.com/soockee/terminal-games/ldtk-snake/layers"
	dresolv "github.com/soockee/terminal-games/ldtk-snake/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/system"
	"github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/yohamta/donburi"
	decs "github.com/yohamta/donburi/ecs"
)

type SnakeScene struct {
	ecs           *decs.ECS
	onceConfigure sync.Once
}

func NewSnakeScene(ecs *decs.ECS) *SnakeScene {
	return &SnakeScene{
		ecs: ecs,
	}
}
func (s *SnakeScene) GetId() component.Scene {
	return component.SnakeScene
}

func (s *SnakeScene) Update() error {
	s.onceConfigure.Do(s.configure)
	s.ecs.Update()
	return nil
}

func (s *SnakeScene) Draw(screen *ebiten.Image) {
	s.drawLevel(screen)
	s.ecs.Draw(screen)
}

func (s *SnakeScene) Layout(w, h int) (int, int) {
	return config.C.LDtkProject.Levels[config.C.CurrentLevel].Width, config.C.LDtkProject.Levels[config.C.CurrentLevel].Height
}

func (s *SnakeScene) drawLevel(screen *ebiten.Image) {

	level := config.C.LDtkProject.Levels[config.C.CurrentLevel]

	if config.C.BGImage != nil {
		opt := &ebiten.DrawImageOptions{}
		bgImage := level.BGImage
		opt.GeoM.Translate(-bgImage.CropRect[0], -bgImage.CropRect[1])
		opt.GeoM.Scale(bgImage.ScaleX, bgImage.ScaleY)
		screen.DrawImage(config.C.BGImage, opt)
	}

	for _, layer := range config.C.EbitenRenderer.RenderedLayers {
		if config.C.ActiveLayers[layer.Layer.Identifier] {
			screen.DrawImage(layer.Image, &ebiten.DrawImageOptions{})
		}
	}
}

func (s *SnakeScene) configure() {
	config.C.CurrentLevel = config.LevelMapping[config.SnakeLevel1]
	config.RenderLevel()

	s.ecs.AddSystem(system.UpdateSnake)
	s.ecs.AddSystem(system.ProcessEvents)
	s.ecs.AddSystem(system.UpdateFood)
	s.ecs.AddSystem(system.UpdateObjects)

	s.ecs.AddRenderer(layers.Default, system.DrawSnake)
	s.ecs.AddRenderer(layers.Default, system.DrawFood)
	s.ecs.AddRenderer(layers.Default, system.DrawWall)
	s.ecs.AddRenderer(layers.Default, system.DrawDebug)

	factory.CreateSettings(s.ecs)

	space := factory.CreateSpace(s.ecs)

	entities := config.C.GetEntities()

	Tags := map[string]func(*decs.ECS, string) *donburi.Entry{
		tags.Snake.Name(): factory.CreateSnake,
		tags.Wall.Name():  factory.CreateWall,
	}
	for _, entity := range entities {
		for name, f := range Tags {
			for _, ldtkTag := range entity.Tags {
				if name == ldtkTag {
					dresolv.Add(space, f(s.ecs, entity.IID))
				}
			}
		}
	}

	// Subscribe events.
	pkgevents.UpdateSettingEvent.Subscribe(s.ecs.World, system.HandleSettingsEvent)
	pkgevents.MoveEvent.Subscribe(s.ecs.World, system.HandleMoveEvent)
}
