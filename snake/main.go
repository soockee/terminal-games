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

func NewGame() *Game {
	game := &Game{}

	game.board = NewBoard()

	return game
}

func (g *Game) Update() error {
	g.board.ticks++

	g.board.UpdateActors()

	if g.board.hitWall || g.board.hitBody {
		slog.Info("Game Lost")
		os.Exit(0)
	}
	return nil
}

func (g *Game) DrawGrid(screen *ebiten.Image) {
	screen.Fill(color.RGBA{247, 246, 187, 128})
	g.DrawInputOverlay(screen)
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

func (g *Game) DrawInputOverlay(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, float32(fieldHeight), OverlayTypeMapping[UpField], false)
	vector.DrawFilledRect(screen, 0, float32(fieldHeight), float32(fieldWidth), float32(fieldHeight), OverlayTypeMapping[LeftField], false)
	vector.DrawFilledRect(screen, float32(fieldWidth), float32(fieldHeight), float32(fieldWidth), float32(fieldHeight), OverlayTypeMapping[RightField], false)
	vector.DrawFilledRect(screen, 0, float32(fieldHeight*2), ScreenWidth, float32(fieldHeight), OverlayTypeMapping[DownField], false)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawGrid(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenHeight, ScreenHeight
}
