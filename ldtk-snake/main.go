package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/ldtk-snake/config"
	"github.com/soockee/terminal-games/ldtk-snake/scene"
)

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type Game struct {
	bounds image.Rectangle
	scene  Scene
}

func NewGame() *Game {
	g := &Game{
		bounds: image.Rectangle{},
		scene:  &scene.StartScene{},
	}

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
	ebiten.SetWindowSize(config.C.LDtkProject.Levels[config.C.CurrentLevel].Width, config.C.LDtkProject.Levels[config.C.CurrentLevel].Height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("LDtk Snake")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
