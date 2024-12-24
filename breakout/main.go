package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/game"
)

func main() {
	ldtkProject, err := assets.NewLDtkProject("simple.ldtk")
	if err != nil {
		panic(err)
	}

	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(h))

	ebiten.SetWindowSize(ldtkProject.Project.LevelByIdentifier(component.StartScene).Width, ldtkProject.Project.LevelByIdentifier(component.StartScene).Height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Breakout")
	g := game.NewGame(ldtkProject)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
