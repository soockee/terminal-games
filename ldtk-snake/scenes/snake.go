package scenes

import (
	"log/slog"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/ldtk-snake/config"
	"github.com/soockee/terminal-games/ldtk-snake/factory"
	"github.com/soockee/terminal-games/ldtk-snake/layers"
	dresolv "github.com/soockee/terminal-games/ldtk-snake/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/systems"
	"github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	decs "github.com/yohamta/donburi/ecs"
)

type SnakeScene struct {
	CurrentLevel int
	ecs          *decs.ECS
	once         sync.Once
}

func (s *SnakeScene) Update() {
	s.once.Do(s.configure)
	s.ecs.Update()
}

func (s *SnakeScene) DrawLevel(screen *ebiten.Image) {

	level := config.C.LDtkProject.Levels[s.CurrentLevel]

	screen.Fill(level.BGColor)

	if config.C.BGImage != nil {
		opt := &ebiten.DrawImageOptions{}
		bgImage := level.BGImage
		opt.GeoM.Translate(-bgImage.CropRect[0], -bgImage.CropRect[1])
		opt.GeoM.Scale(bgImage.ScaleX, bgImage.ScaleY)
		screen.DrawImage(config.C.BGImage, opt)
	}

	for i, layer := range config.C.EbitenRenderer.RenderedLayers {
		if config.C.ActiveLayers[i] {
			screen.DrawImage(layer.Image, &ebiten.DrawImageOptions{})
		}
	}
}

func (s *SnakeScene) Draw(screen *ebiten.Image) {
	// screen.Fill(color.RGBA{20, 20, 40, 255})
	// s.DrawLevel(screen)
	s.ecs.Draw(screen)
}

func (s *SnakeScene) Layout(w, h int) (int, int) {
	return config.C.LDtkProject.WorldGridWidth, config.C.LDtkProject.WorldGridWidth
}

func (s *SnakeScene) configure() {
	ecs := ecs.NewECS(donburi.NewWorld())

	ecs.AddSystem(systems.UpdateSnake)
	ecs.AddSystem(systems.UpdateObjects)
	ecs.AddSystem(systems.UpdateSettings)

	ecs.AddRenderer(layers.Default, systems.DrawSnake)
	ecs.AddRenderer(layers.Default, systems.DrawWall)
	ecs.AddRenderer(layers.Default, systems.DrawDebug)

	s.ecs = ecs

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
					slog.Info("Create resolv obj", slog.Any("entity", entity))
					dresolv.Add(space, f(ecs, entity.IID))
				}
			}
		}
	}
}
