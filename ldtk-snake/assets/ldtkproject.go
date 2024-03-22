package assets

import (
	"embed"
	"fmt"
	"image"
	"path/filepath"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/ldtkgo"
	"github.com/yohamta/ganim8/v2"
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

func (ldtk LDtkProject) GetEntities(level string) []*ldtkgo.Entity {
	entities := []*ldtkgo.Entity{}
	for _, layer := range ldtk.Project.LevelByIdentifier(level).Layers {
		entities = append(entities, layer.Entities...)
	}
	return entities
}

// GetEntityByName returns the first found entity by name
func (ldtk LDtkProject) GetEntityByName(name string, level string) *ldtkgo.Entity {
	for _, layer := range ldtk.Project.LevelByIdentifier(level).Layers {
		for _, entity := range layer.Entities {
			if entity.Identifier == name {
				return entity
			}
		}
	}
	return nil
}

func (ldtk LDtkProject) loadRequiredTileset(entity *ldtkgo.Entity, level string) (*ebiten.Image, error) {
	var tileset *ebiten.Image
	var err error
	for _, layer := range ldtk.Project.LevelByIdentifier(level).Layers {
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

func (ldtk LDtkProject) GetEntitiesByTag(tag string) []*ldtkgo.EntityDefinition {
	entityDefinitions := []*ldtkgo.EntityDefinition{}
	for _, definition := range ldtk.Project.EntityDefinitions {
		if slices.Contains(definition.Tags, tag) {
			entityDefinitions = append(entityDefinitions, definition)
		}
	}
	return entityDefinitions
}

func (ldtk LDtkProject) GetSpritesByTag(tag string) (map[string]*ebiten.Image, error) {
	entityDefinitions := ldtk.GetEntitiesByTag(tag)
	sprites := map[string]*ebiten.Image{}
	for _, e := range entityDefinitions {
		s, err := ldtk.GetSprite(e.TileRect)
		if err != nil {
			return nil, err
		}
		sprites[e.Identifier] = s
	}
	return sprites, nil
}

func (ldtk LDtkProject) GetSpriteByIdentifier(identifier string) (*ebiten.Image, error) {
	entityDefinition := ldtk.Project.EntityDefinitionByIdentifier(identifier)
	return ldtk.GetSprite(entityDefinition.TileRect)
}

func (ldtk LDtkProject) GetSpriteByDefinition(entityDefinition *ldtkgo.EntityDefinition) (*ebiten.Image, error) {
	return ldtk.GetSprite(entityDefinition.TileRect)
}

func (ldtk LDtkProject) GetSpriteByEntityInstance(entity *ldtkgo.Entity) (*ebiten.Image, error) {
	return ldtk.GetSprite(entity.TileRect)
}

func (ldtk LDtkProject) GetSprite(tileRect *ldtkgo.TileRect) (*ebiten.Image, error) {
	tileset, err := ldtk.Renderer.Loader.LoadImage(tileRect.Tileset.Path)
	t := tileRect
	subImageRect := image.Rect(t.X, t.Y, t.X+t.W, t.Y+t.H)
	sprite := tileset.SubImage(subImageRect).(*ebiten.Image)
	return sprite, err
}

func (ldtk LDtkProject) IsAnimated(identifier string) bool {
	entityDefinition := ldtk.Project.EntityDefinitionByIdentifier(identifier)
	return slices.Contains(entityDefinition.Tags, "Animated")
}
func (ldtk LDtkProject) GetAnimatedSpriteByIdentifier(identifier string) (*ganim8.Animation, error) {
	if !ldtk.IsAnimated(identifier) {
		return nil, fmt.Errorf("entity is not animated")
	}
	return ldtk.GetAnimatedSpriteByDefinition(ldtk.Project.EntityDefinitionByIdentifier(identifier))
}

func (ldtk LDtkProject) GetAnimatedSpriteByDefinition(entityDefinition *ldtkgo.EntityDefinition) (*ganim8.Animation, error) {
	if !ldtk.IsAnimated(entityDefinition.Identifier) {
		return nil, fmt.Errorf("entity is not animated")
	}
	return ldtk.GetAnimatedSprite(entityDefinition.TileRect, entityDefinition.Width, entityDefinition.Height)
}

func (ldtk LDtkProject) GetAnimatedSprite(tileRect *ldtkgo.TileRect, frameW, frameH int) (*ganim8.Animation, error) {
	tileset, err := ldtk.Renderer.Loader.LoadImage(tileRect.Tileset.Path)

	t := tileRect
	subImageRect := image.Rect(t.X, t.Y, t.X+t.W, t.Y+t.H)
	sprite := tileset.SubImage(subImageRect).(*ebiten.Image)
	grid := ganim8.NewGrid(frameW, frameH, t.W, t.H, t.X, t.Y, tileRect.Tileset.Spacing)

	frameCount := sprite.Bounds().Dx() / frameW
	animationTime := time.Millisecond * 100
	frameRowSelection := fmt.Sprintf("1-%d", frameCount)
	// use only column 1
	frames := grid.Frames(frameRowSelection, 1)
	animation := ganim8.New(sprite, frames, animationTime)
	return animation, err
}

func (ldtk *LDtkProject) RenderLevel(currentlevel string) {
	level := ldtk.Project.LevelByIdentifier(currentlevel)

	ldtk.Renderer.Render(level)
}
