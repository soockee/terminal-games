package factory

import (
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/archetype"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/config"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateSpace(ecs *ecs.ECS) *donburi.Entry {
	space := archetype.Space.Spawn(ecs)

	cfg := config.C
	var spaceData *resolv.Space
	if cfg.LDtkProject.Levels[0].Layers[0] != nil {
		cellWidth := cfg.LDtkProject.WorldGridWidth / cfg.LDtkProject.Levels[0].Layers[0].CellWidth
		CellHeight := cfg.LDtkProject.WorldGridHeight / cfg.LDtkProject.Levels[0].Layers[0].CellHeight
		spaceData = resolv.NewSpace(cfg.LDtkProject.WorldGridWidth, cfg.LDtkProject.WorldGridWidth, cellWidth, CellHeight)
	}
	component.Space.Set(space, spaceData)

	return space
}
