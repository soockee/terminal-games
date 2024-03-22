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
		component.Velocity,
	)
	SnakeBody = newArchetype(
		tags.SnakeBody,

		component.Object,
		component.Sprite,
		component.SnakeBody,
		component.Velocity,
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
		tags.Collectable,

		component.Object,
		component.Sprite,
		component.Collectable,
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
	TextField = newArchetype(
		tags.TextField,

		component.Object,
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

	Mouse = newArchetype(
		tags.Mouse,
		tags.Animated,
		tags.Collectable,

		component.Mouse,
		component.Object,
		component.Animation,
		component.Collectable,
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
