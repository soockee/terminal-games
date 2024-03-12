package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/game"
)

func main() {
	ldtkProject, err := assets.NewLDtkProject("ldtk/simple.ldtk")
	if err != nil {
		panic(err)
	}

	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(h))

	ebiten.SetWindowSize(ldtkProject.Project.Levels[int(component.StartScene)].Width, ldtkProject.Project.Levels[int(component.StartScene)].Height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("LDtk Snake")
	g := game.NewGame(ldtkProject)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
