package archetypes

import (
	"github.com/soockee/terminal-games/ldtk-snake/components"
	"github.com/soockee/terminal-games/ldtk-snake/layers"
	"github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var (
	Snake = newArchetype(
		tags.Snake,
		components.Snake,
		components.Object,
		components.Sprite,
		components.Control,
	)
	Space = newArchetype(
		components.Space,
	)
	Wall = newArchetype(
		tags.Wall,
		components.Object,
		components.Sprite,
	)
	Settings = newArchetype(
		components.Settings,
		components.Control,
	)
)

type archetype struct {
	components []donburi.IComponentType
}

func newArchetype(cs ...donburi.IComponentType) *archetype {
	return &archetype{
		components: cs,
	}
}

func (a *archetype) Spawn(ecs *ecs.ECS, cs ...donburi.IComponentType) *donburi.Entry {
	e := ecs.World.Entry(ecs.Create(
		layers.Default,
		append(a.components, cs...)...,
	))
	return e
}
