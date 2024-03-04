package main

import (
	"image/color"
	"os"

	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Snake in Go")
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(h))

	slog.Debug("Creating Game...")
	game := &Game{
		board: NewBoard(),
	}

	if err := ebiten.RunGame(game); err != nil {
		slog.Any("err", err)
	}
}

type Game struct {
	board *Board
}

func (g *Game) Update() error {
	g.board.ticks++

	g.board.UpdateActors()

	if g.board.hitWall {
		slog.Info("Game Lost")
		os.Exit(0)
	}
	return nil
}

func (g *Game) DrawGrid(screen *ebiten.Image) {
	screen.Fill(color.White)
	for dx, row := range g.board.cells {
		for dy, cell := range row {
			rectX := float32(dx * GridSize)
			rectY := float32(dy * GridSize)
			rectWidth := float32(GridSize - Offset)
			rectHeight := float32(GridSize - Offset)

			vector.DrawFilledRect(screen, rectX, rectY, rectWidth, rectHeight, CellTypeMapping[cell.cellType], false)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawGrid(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenHeight, ScreenHeight
}
