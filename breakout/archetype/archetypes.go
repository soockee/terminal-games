package archetype

import (
	"github.com/solarlune/ldtkgo"
	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/layers"
	"github.com/soockee/terminal-games/breakout/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var (
	TagsMapping = map[string]func(donburi.World, *assets.LDtkProject, *ldtkgo.Entity) *donburi.Entry{
		tags.Button.Name():     NewButton,
		tags.Collidable.Name(): NewWall,
		tags.Player.Name():     NewPlayer,
		tags.Ball.Name():       NewBall,
		tags.Wall.Name():       NewWall,
		tags.Brick.Name():      NewBrick,
	}
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
