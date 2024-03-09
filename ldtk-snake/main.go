package main

import (
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/ldtkgo"
	"github.com/soockee/terminal-games/ldtk-snake/scene"
)

type Game struct {
	bounds      image.Rectangle
	scene       scene.Scene
	ldtkProject *ldtkgo.Project
}

func NewGame() *Game {
	g := &Game{
		bounds: image.Rectangle{},
	}

	dir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	g.ldtkProject, err = ldtkgo.Open("assets/ldtk/simple.ldtk", os.DirFS(dir))

	if err != nil {
		panic(err)
	}

	g.scene = scene.NewScene(g.ldtkProject)

	return g
}

func (g *Game) Update() error {
	g.scene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.scene.Draw(screen)
}

func (g *Game) Layout(width, height int) (int, int) {
	g.bounds = image.Rect(0, 0, width, height)
	return width, height
}

func main() {
	g := NewGame()
	ebiten.SetWindowSize(g.ldtkProject.WorldGridWidth, g.ldtkProject.WorldGridHeight)

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	ebiten.SetWindowTitle("LDtk Snake")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
