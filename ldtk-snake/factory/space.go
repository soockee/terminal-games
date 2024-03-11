package factory

import (
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateSpace(ecs *ecs.ECS, width int, height int, cellWidth int, cellHeight int) *donburi.Entry {
	space := archetype.Space.Spawn(ecs)
	spaceData := resolv.NewSpace(width, height, cellWidth, cellHeight)
	component.Space.Set(space, spaceData)

	return space
}
