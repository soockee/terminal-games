package main

import (
	"io/fs"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	ldtkgo "github.com/soockee/ldtk-super-simple-loader"
	"github.com/soockee/terminal-games/super-mario-bros/assets"
	"github.com/soockee/terminal-games/super-mario-bros/game"
	"github.com/soockee/terminal-games/super-mario-bros/system"
)

func main() {
	ldtkFS, err := fs.Sub(assets.FS, "ldtk")
	if err != nil {
		log.Fatalf("sub fs: %v", err)
	}

	world, err := ldtkgo.LoadWorld("super-mario-bros.ldtk", ldtkFS)
	if err != nil {
		log.Fatalf("loading world: %v", err)
	}

	g := game.New(world, game.GameConfig{
		VirtualW: 256,
		VirtualH: 144,
		Systems: []game.SystemFunc{
			system.UpdateInput,
			system.UpdateMovement,
			system.UpdatePatrol,
			system.UpdatePhysics,
			system.UpdateCollection,
			system.UpdateEnemyCleanup,
			system.UpdateAnimation,
			system.ProcessEvents,
		},
		Renderers: []game.RendererFunc{
			system.DrawEntities,
			system.DrawAnimation,
			system.DrawScore,
			system.DrawHUD,
			system.DrawDebug,
		},
	}, game.LevelConfig{}, nil, nil)

	w, h := g.ScreenSize()
	ebiten.SetWindowSize(w*3, h*3)
	ebiten.SetWindowTitle("Super Mario Bros")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
