package archetype

import (
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/layers"
	"github.com/soockee/terminal-games/ldtk-snake/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var (
	Snake = newArchetype(
		tags.Snake,
		component.Snake,
		component.Object,
		component.Sprite,
	)
	Space = newArchetype(
		component.Space,
	)
	Wall = newArchetype(
		tags.Wall,
		component.Object,
		component.Sprite,
	)
	Food = newArchetype(
		tags.Food,
		component.Object,
		component.Sprite,
	)
	Settings = newArchetype(
		component.Settings,
	)
	Button = newArchetype(
		tags.Button,
		component.Object,
		component.Sprite,
		component.Button,
	)
	Controls = newArchetype(
		component.Control,
	)
	GameState = newArchetype(
		component.Gamestate,
	)
	SceneState = newArchetype(
		component.SceneState,
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
