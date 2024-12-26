package factory

import (
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/breakout/archetype"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/yohamta/donburi"
)

func CreateSpace(w donburi.World, width int, height int, cellWidth int, cellHeight int) *donburi.Entry {
	space := archetype.Space.SpawnInWorld(w)
	spaceData := resolv.NewSpace(width, height, cellWidth, cellHeight)
	component.Space.Set(space, spaceData)

	return space
}
