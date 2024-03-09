package config

import (
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/soockee/ldtkgo"
	"github.com/soockee/terminal-games/ldtk-snake/helper"
)

type Config struct {
	LDtkProject    *ldtkgo.Project
	CurrentLevel   int
	EbitenRenderer *helper.EbitenRenderer
	BGImage        *ebiten.Image
	ActiveLayers   []bool
}

var C *Config

func init() {

	dir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	var ldtkProject *ldtkgo.Project
	ldtkProject, err = ldtkgo.Open("assets/ldtk/simple.ldtk", os.DirFS(dir))

	if err != nil {
		panic(err)
	}

	C = &Config{
		LDtkProject:    ldtkProject,
		CurrentLevel:   0,
		EbitenRenderer: helper.NewEbitenRenderer(helper.NewDiskLoader("assets/ldtk")),
		ActiveLayers:   []bool{true, true, true, true},
	}
	renderLevel()
}

func (c *Config) GetEntities() []*ldtkgo.Entity {
	entities := []*ldtkgo.Entity{}
	for _, layer := range c.LDtkProject.Levels[c.CurrentLevel].Layers {
		entities = append(entities, layer.Entities...)
	}
	return entities
}

// GetEntityByName returns the first found entity by name
func (c *Config) GetEntityByName(name string, level int) *ldtkgo.Entity {
	for _, layer := range c.LDtkProject.Levels[level].Layers {
		for _, entity := range layer.Entities {
			if entity.Identifier == name {
				return entity
			}
		}
	}
	return nil
}

func (c *Config) GetEntityByIID(iid string, level int) *ldtkgo.Entity {
	for _, layer := range c.LDtkProject.Levels[level].Layers {
		for _, entity := range layer.Entities {
			if entity.IID == iid {
				return entity
			}
		}
	}
	return nil
}

func (c *Config) GetSprite(entity *ldtkgo.Entity) *ebiten.Image {
	tileset := c.EbitenRenderer.Tilesets[entity.TileRect.Tileset.Path]
	tileRect := entity.TileRect
	sprite := tileset.SubImage(image.Rect(tileRect.X, tileRect.Y, tileRect.X+tileRect.W, tileRect.Y+tileRect.H)).(*ebiten.Image)
	return sprite
}

func renderLevel() {
	if C.CurrentLevel >= len(C.LDtkProject.Levels) {
		C.CurrentLevel = 0
	}

	if C.CurrentLevel < 0 {
		C.CurrentLevel = len(C.LDtkProject.Levels) - 1
	}

	level := C.LDtkProject.Levels[C.CurrentLevel]

	if level.BGImage != nil {
		C.BGImage, _, _ = ebitenutil.NewImageFromFile(level.BGImage.Path)
	} else {
		C.BGImage = nil
	}

	C.EbitenRenderer.Render(level)
}
