package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/ldtkgo"
	"github.com/soockee/terminal-games/ldtk-snake/helper"
	"github.com/soockee/terminal-games/ldtk-snake/systems"
)

type Scene struct {
	EbitenRenderer *helper.EbitenRenderer
	BGImage        *ebiten.Image
	CurrentLevel   int
	ActiveLayers   []bool
	ldtkProject    *ldtkgo.Project
	systems        []systems.System
}

func NewScene(ldtkProject *ldtkgo.Project, ebitenRenderer *helper.EbitenRenderer) Scene {

	s := Scene{
		ActiveLayers:   []bool{true, true, true, true},
		ldtkProject:    ldtkProject,
		EbitenRenderer: ebitenRenderer,
	}

	s.RenderLevel()

	s.systems = append(s.systems, systems.NewSnake(helper.GetEntityByName("Snakehead", s.CurrentLevel, s.ldtkProject), s.EbitenRenderer))

	return s
}
func (s *Scene) RenderLevel() {

	if s.CurrentLevel >= len(s.ldtkProject.Levels) {
		s.CurrentLevel = 0
	}

	if s.CurrentLevel < 0 {
		s.CurrentLevel = len(s.ldtkProject.Levels) - 1
	}

	level := s.ldtkProject.Levels[s.CurrentLevel]

	if level.BGImage != nil {
		s.BGImage, _, _ = ebitenutil.NewImageFromFile(level.BGImage.Path)
	} else {
		s.BGImage = nil
	}

	s.EbitenRenderer.Render(level)
}

func (g *Scene) Update() error {
	for _, system := range g.systems {
		system.Update()
	}
	return nil

}

func (g *Scene) Draw(screen *ebiten.Image) {

	level := g.ldtkProject.Levels[g.CurrentLevel]

	screen.Fill(level.BGColor)

	if g.BGImage != nil {
		opt := &ebiten.DrawImageOptions{}
		bgImage := level.BGImage
		opt.GeoM.Translate(-bgImage.CropRect[0], -bgImage.CropRect[1])
		opt.GeoM.Scale(bgImage.ScaleX, bgImage.ScaleY)
		screen.DrawImage(g.BGImage, opt)
	}

	for i, layer := range g.EbitenRenderer.RenderedLayers {
		if g.ActiveLayers[i] {
			screen.DrawImage(layer.Image, &ebiten.DrawImageOptions{})
		}
	}

	for _, system := range g.systems {
		system.Draw(screen)
	}
}

func (g *Scene) Layout(w, h int) (int, int) {
	return g.ldtkProject.WorldGridWidth, g.ldtkProject.WorldGridWidth
}
