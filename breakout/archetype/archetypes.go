package archetype

import (
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/layers"
	"github.com/soockee/terminal-games/breakout/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var (
	Space = newArchetype(
		component.Space,
	)
	Settings = newArchetype(
		component.Settings,
	)
	Button = newArchetype(
		tags.Button,

		component.ConvexPolygon,
		component.Sprite,
		component.Button,
	)
	TextField = newArchetype(
		tags.TextField,

		component.ConvexPolygon,
		component.Sprite,
		component.Text,
	)
	Controls = newArchetype(
		component.Control,
	)
	SceneState = newArchetype(
		component.SceneState,
	)
	GameState = newArchetype(
		component.GameState,
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

func (a *archetype) SpawnInWorld(world donburi.World, cs ...donburi.IComponentType) *donburi.Entry {
	e := world.Entry(world.Create(
		append(a.components, cs...)...,
	))
	return e
}
