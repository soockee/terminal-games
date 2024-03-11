package config

import (
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/soockee/ldtkgo"
	"github.com/soockee/terminal-games/ldtk-snake/util"
)

type Config struct {
	LDtkProject    *ldtkgo.Project
	CurrentLevel   int
	EbitenRenderer *util.EbitenRenderer
	BGImage        *ebiten.Image
	ActiveLayers   map[string]bool
	AssetBasePath  string
}

const (
	assetBasePath   string = "assets/ldtk"
	projectFileName string = "simple.ldtk"
)

var C *Config

type LevelName int

const (
	StartScreen LevelName = iota
	SnakeLevel1
)

var levelMapping = map[LevelName]int{
	SnakeLevel1: 0,
	StartScreen: 1,
}

var activeLayerMapping = map[string]bool{
	"Entities":   false,
	"Background": true,
}

func init() {

	dir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	var ldtkProject *ldtkgo.Project
	path := assetBasePath + "/" + projectFileName
	ldtkProject, err = ldtkgo.Open(path, os.DirFS(dir))

	if err != nil {
		panic(err)
	}

	C = &Config{
		LDtkProject:    ldtkProject,
		CurrentLevel:   levelMapping[StartScreen],
		EbitenRenderer: util.NewEbitenRenderer(util.NewDiskLoader(assetBasePath)),
		ActiveLayers:   activeLayerMapping,
		AssetBasePath:  assetBasePath,
	}
	RenderLevel()
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

func (c *Config) loadRequiredTileset(entity *ldtkgo.Entity) *ebiten.Image {
	// requiredTilessets := []string{}
	for _, layer := range c.LDtkProject.Levels[c.CurrentLevel].Layers {
		if layer.Identifier == "Entities" {
			for _, e := range layer.Entities {
				if entity == e {
					return c.EbitenRenderer.Loader.LoadTileset(e.TileRect.Tileset.Path)
				}
			}
			break
		}
	}
	return nil
}

func (c *Config) GetSprite(entity *ldtkgo.Entity) *ebiten.Image {
	tileset := c.EbitenRenderer.Loader.LoadTileset(entity.TileRect.Tileset.Path)
	tileRect := entity.TileRect
	subImageRect := image.Rect(tileRect.X, tileRect.Y, tileRect.X+tileRect.W, tileRect.Y+tileRect.H)
	sprite := tileset.SubImage(subImageRect).(*ebiten.Image)
	return sprite
}

func RelativeCrop(source *ebiten.Image, r image.Rectangle) *ebiten.Image {
	rx, ry := source.Bounds().Min.X+r.Min.X, source.Bounds().Min.Y+r.Min.Y
	return source.SubImage(image.Rect(rx, ry, rx+r.Max.X, ry+r.Max.Y)).(*ebiten.Image)
}

func RenderLevel() {
	if C.CurrentLevel >= len(C.LDtkProject.Levels) {
		C.CurrentLevel = 0
	}

	if C.CurrentLevel < 0 {
		C.CurrentLevel = len(C.LDtkProject.Levels) - 1
	}

	level := C.LDtkProject.Levels[C.CurrentLevel]

	if level.BGImage != nil {
		C.BGImage, _, _ = ebitenutil.NewImageFromFile(assetBasePath + "/" + level.BGImage.Path)
	} else {
		C.BGImage = nil
	}

	C.EbitenRenderer.Render(level)
}
