package archetype

import (
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/yohamta/donburi"
)

var (
	Space = newArchetype(
		component.Space,
	)
)

func NewSpace(w donburi.World, width int, height int, cellWidth int, cellHeight int) *donburi.Entry {
	space := Space.SpawnInWorld(w)
	spaceData := resolv.NewSpace(width, height, cellWidth, cellHeight)
	component.Space.Set(space, spaceData)

	return space
}
