package archetype

import (
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/layers"
	"github.com/soockee/terminal-games/breakout/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var (
	// Game-specific archetypes
	Player = newArchetype(
		tags.Player,

		component.Player,
		component.Collidable,
		component.Sprite,
		component.Velocity,
	)

	Wall = newArchetype(
		tags.Wall,
		component.Collidable,
	)

	Ball = newArchetype(
		tags.Ball,

		component.Ball,
		component.Sprite,
		component.Velocity,
	)

	// Generic archetypes
	Space = newArchetype(
		component.Space,
	)
	Settings = newArchetype(
		component.Settings,
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

	// UI archetypes
	Button = newArchetype(
		tags.Button,

		component.Sprite,
		component.Button,
	)
	TextField = newArchetype(
		tags.TextField,

		component.Sprite,
		component.Text,
	)
)

type Archetype struct {
	components []donburi.IComponentType
}

func newArchetype(cs ...donburi.IComponentType) *Archetype {
	return &Archetype{
		components: cs,
	}
}

func (a *Archetype) Spawn(ecs *ecs.ECS, cs ...donburi.IComponentType) *donburi.Entry {
	e := ecs.World.Entry(ecs.Create(
		layers.Default,
		append(a.components, cs...)...,
	))
	return e
}

func (a *Archetype) SpawnInWorld(world donburi.World, cs ...donburi.IComponentType) *donburi.Entry {
	e := world.Entry(world.Create(
		append(a.components, cs...)...,
	))
	return e
}
