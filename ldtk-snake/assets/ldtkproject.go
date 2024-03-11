package assets

import (
	"embed"
	"image"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/ldtkgo"
)

var (
	//go:embed *
	assetsFS embed.FS
)

type LDtkProject struct {
	Project      *ldtkgo.Project
	Basepath     string
	Renderer     *EbitenRenderer
	ActiveLayers map[string]bool
}

func NewLDtkProject(path string) (*LDtkProject, error) {
	ldtkProject, err := ldtkgo.Open(path, assetsFS)
	loader := NewEmbedLoader(filepath.Dir(path))
	if err != nil {
		return nil, err
	}
	return &LDtkProject{
		Project:  ldtkProject,
		Basepath: filepath.Dir(path),
		Renderer: NewEbitenRenderer(loader),
		ActiveLayers: map[string]bool{
			"Entities":   false,
			"Background": true,
		},
	}, nil

}

func (ldtk LDtkProject) GetEntities(level int) []*ldtkgo.Entity {
	entities := []*ldtkgo.Entity{}
	for _, layer := range ldtk.Project.Levels[level].Layers {
		entities = append(entities, layer.Entities...)
	}
	return entities
}

// GetEntityByName returns the first found entity by name
func (ldtk LDtkProject) GetEntityByName(name string, level int) *ldtkgo.Entity {
	for _, layer := range ldtk.Project.Levels[level].Layers {
		for _, entity := range layer.Entities {
			if entity.Identifier == name {
				return entity
			}
		}
	}
	return nil
}

func (ldtk LDtkProject) GetEntityByIID(iid string, level int) *ldtkgo.Entity {
	for _, layer := range ldtk.Project.Levels[level].Layers {
		for _, entity := range layer.Entities {
			if entity.IID == iid {
				return entity
			}
		}
	}
	return nil
}

func (ldtk LDtkProject) loadRequiredTileset(entity *ldtkgo.Entity, level int) (*ebiten.Image, error) {
	var tileset *ebiten.Image
	var err error
	for _, layer := range ldtk.Project.Levels[level].Layers {
		if layer.Identifier == "Entities" {
			for _, e := range layer.Entities {
				if entity == e {
					tileset, err = ldtk.Renderer.Loader.LoadImage(e.TileRect.Tileset.Path)
				}
			}
			break
		}
	}
	return tileset, err
}

func (ldtk LDtkProject) GetSprite(entity *ldtkgo.Entity) (*ebiten.Image, error) {
	tileset, err := ldtk.Renderer.Loader.LoadImage(entity.TileRect.Tileset.Path)
	tileRect := entity.TileRect
	subImageRect := image.Rect(tileRect.X, tileRect.Y, tileRect.X+tileRect.W, tileRect.Y+tileRect.H)
	sprite := tileset.SubImage(subImageRect).(*ebiten.Image)
	return sprite, err
}

func (ldtk *LDtkProject) RenderLevel(currentLevel int) {
	if currentLevel >= len(ldtk.Project.Levels) {
		currentLevel = 0
	}

	if currentLevel < 0 {
		currentLevel = len(ldtk.Project.Levels) - 1
	}

	level := ldtk.Project.Levels[currentLevel]

	ldtk.Renderer.Render(level)
}
