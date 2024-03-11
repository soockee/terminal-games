package main

import (
	"log"

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

	ebiten.SetWindowSize(ldtkProject.Project.Levels[int(component.StartScene)].Width, ldtkProject.Project.Levels[int(component.StartScene)].Height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("LDtk Snake")
	g := game.NewGame(ldtkProject)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
