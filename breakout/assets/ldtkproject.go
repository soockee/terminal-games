package assets

import (
	"embed"
	"fmt"
	"image"
	"io/fs"
	"path/filepath"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/ldtkgo"
	"github.com/yohamta/ganim8/v2"
)

var (
	//go:embed ldtk/*
	assetsFS embed.FS
)

type LDtkProject struct {
	Project      *ldtkgo.Project
	Basepath     string
	Renderer     *Renderer
	ActiveLayers map[string]bool
}

func NewLDtkProject(path string) (*LDtkProject, error) {
	assetsFS, err := fs.Sub(assetsFS, "ldtk")
	if err != nil {
		return nil, err
	}

	ldtkProject, err := ldtkgo.Open(path, assetsFS)

	if err != nil {
		return nil, err
	}

	renderer, err := New(assetsFS, ldtkProject)
	if err != nil {
		return nil, err
	}

	return &LDtkProject{
		Project:  ldtkProject,
		Basepath: filepath.Dir(path),
		Renderer: renderer,
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

func (ldtk LDtkProject) GetEntitiesByTag(tag string) []*ldtkgo.EntityDefinition {
	entityDefinitions := []*ldtkgo.EntityDefinition{}
	for _, definition := range ldtk.Project.EntityDefinitions {
		if slices.Contains(definition.Tags, tag) {
			entityDefinitions = append(entityDefinitions, definition)
		}
	}
	return entityDefinitions
}

func (ldtk LDtkProject) GetSpritesByTag(tag string) map[string]*ebiten.Image {
	entityDefinitions := ldtk.GetEntitiesByTag(tag)
	sprites := map[string]*ebiten.Image{}
	for _, e := range entityDefinitions {
		s := ldtk.GetSprite(e.TileRect)
		sprites[e.Identifier] = s
	}
	return sprites
}

func (ldtk LDtkProject) GetSpriteByIdentifier(identifier string) *ebiten.Image {
	entityDefinition := ldtk.Project.EntityDefinitionByIdentifier(identifier)
	return ldtk.GetSprite(entityDefinition.TileRect)
}

func (ldtk LDtkProject) GetSpriteByDefinition(entityDefinition *ldtkgo.EntityDefinition) *ebiten.Image {
	return ldtk.GetSprite(entityDefinition.TileRect)
}

func (ldtk LDtkProject) GetSpriteByEntityInstance(entity *ldtkgo.Entity) *ebiten.Image {
	return ldtk.GetSprite(entity.TileRect)
}

func (ldtk LDtkProject) GetSprite(tileRect *ldtkgo.TileRect) *ebiten.Image {
	tileset := ldtk.Renderer.Tilesets[tileRect.Tileset.Path]

	//tileset, err := ldtk.Renderer.Loader.LoadImage(tileRect.Tileset.Path)
	t := tileRect
	subImageRect := image.Rect(t.X, t.Y, t.X+t.W, t.Y+t.H)
	sprite := tileset.SubImage(subImageRect).(*ebiten.Image)
	return sprite
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
	return ldtk.GetAnimatedSprite(entityDefinition.TileRect, entityDefinition.Width, entityDefinition.Height), nil
}

func (ldtk LDtkProject) GetAnimatedSprite(tileRect *ldtkgo.TileRect, frameW, frameH int) *ganim8.Animation {
	tileset := ldtk.Renderer.Tilesets[tileRect.Tileset.Path]

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
	return animation
}
