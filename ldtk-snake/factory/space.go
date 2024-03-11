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
	if cfg.LDtkProject.Levels[cfg.CurrentLevel].Layers[cfg.CurrentLevel] != nil {
		cellWidth := cfg.LDtkProject.Levels[cfg.CurrentLevel].Width / cfg.LDtkProject.Levels[cfg.CurrentLevel].Layers[cfg.CurrentLevel].CellWidth
		CellHeight := cfg.LDtkProject.Levels[cfg.CurrentLevel].Height / cfg.LDtkProject.Levels[cfg.CurrentLevel].Layers[cfg.CurrentLevel].CellHeight
		spaceData = resolv.NewSpace(cfg.LDtkProject.Levels[cfg.CurrentLevel].Width, cfg.LDtkProject.Levels[cfg.CurrentLevel].Height, cellWidth, CellHeight)
	}
	component.Space.Set(space, spaceData)

	return space
}
