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
	"github.com/yohamta/donburi/ecs"
	decs "github.com/yohamta/donburi/ecs"
)

type StartScene struct {
	ecs  *decs.ECS
	once sync.Once
}

func NewStartScene(ecs *decs.ECS) *StartScene {
	return &StartScene{
		ecs: ecs,
	}
}

func (s *StartScene) GetId() component.Scene {
	return component.StartScreen
}

func (s *StartScene) Update() error {
	s.once.Do(s.configure)
	s.ecs.Update()
	return nil
}

func (s *StartScene) Draw(screen *ebiten.Image) {
	s.drawLevel(screen)
	s.ecs.Draw(screen)
}

func (s *StartScene) drawLevel(screen *ebiten.Image) {

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

func (s *StartScene) Layout(w, h int) (int, int) {
	return config.C.LDtkProject.Levels[config.C.CurrentLevel].Width, config.C.LDtkProject.Levels[config.C.CurrentLevel].Height
}

func (s *StartScene) configure() {
	s.ecs.AddSystem(system.UpdateObjects)
	s.ecs.AddSystem(system.ProcessEvents)
	s.ecs.AddSystem(system.UpdateButton)

	s.ecs.AddRenderer(layers.Default, system.DrawDebug)
	s.ecs.AddRenderer(layers.Default, system.DrawButton)

	factory.CreateSettings(s.ecs)
	space := factory.CreateSpace(s.ecs)

	createEntities(s.ecs, space)

	// Subscribe events.
	pkgevents.UpdateSettingEvent.Subscribe(s.ecs.World, system.HandleSettingsEvent)
	pkgevents.InteractionEvent.Subscribe(s.ecs.World, system.HandleButtonClick)
}

func createEntities(ecs *ecs.ECS, space *donburi.Entry) {
	entities := config.C.GetEntities()

	Tags := map[string]func(*decs.ECS, string) *donburi.Entry{
		tags.Button.Name(): factory.CreateButton,
	}
	for _, entity := range entities {
		for name, f := range Tags {
			for _, ldtkTag := range entity.Tags {
				if name == ldtkTag {
					dresolv.Add(space, f(ecs, entity.IID))
				}
			}
		}
	}
}