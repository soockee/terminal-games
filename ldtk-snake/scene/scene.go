package scene

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/ldtkgo"
)

type Scene struct {
	EbitenRenderer *EbitenRenderer
	BGImage        *ebiten.Image
	CurrentLevel   int
	ActiveLayers   []bool
	ldtkProject    *ldtkgo.Project
}

func NewScene(ldtkProject *ldtkgo.Project) Scene {

	s := Scene{
		ActiveLayers: []bool{true, true, true, true},
		ldtkProject:  ldtkProject,
	}

	s.EbitenRenderer = NewEbitenRenderer(NewDiskLoader("assets/ldtk"))

	s.RenderLevel()

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

	return nil

}

func (g *Scene) Draw(screen *ebiten.Image) {

	level := g.ldtkProject.Levels[g.CurrentLevel]

	screen.Fill(level.BGColor) // We want to use the BG Color when possible

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

	// We'll additionally render the entities onscreen.
	for _, layer := range level.Layers {
		// In truth, we don't have to check to see if it's an entity layer before looping through,
		// because only Entity layers have entities in the Entities slice.
		for _, entity := range layer.Entities {

			if entity.TileRect != nil {

				tileset := g.EbitenRenderer.Tilesets[entity.TileRect.Tileset.Path]
				tileRect := entity.TileRect
				tile := tileset.SubImage(image.Rect(tileRect.X, tileRect.Y, tileRect.X+tileRect.W, tileRect.Y+tileRect.H)).(*ebiten.Image)

				opt := &ebiten.DrawImageOptions{}
				opt.GeoM.Translate(float64(entity.Position[0]), float64(entity.Position[1]))

				screen.DrawImage(tile, opt)

			}
		}
	}
}

func (g *Scene) Layout(w, h int) (int, int) { return 256, 256 }
