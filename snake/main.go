package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(DebugFlag)

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Snake in Go")
	log.Debug().Msg("a debug message")
	game := &Game{
		board: NewBoard(),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Err(err).Msg("Failed to run game")
	}
}

type Game struct {
	board *Board
}

func (g *Game) Update() error {
	g.board.ticks++

	for dx, row := range g.board.cells {
		for dy := range row {
			if dx == 0 || dx == CellsDX-1 || dy == 0 || dy == CellsDY-1 {
				continue
			}
			g.board.cells[dx][dy] = EmptyCell
		}
	}
	g.board.cells[g.board.snake.snakeHeadDX][g.board.snake.snakeHeadDY] = SnakeHead

	g.board.UpdateActors()

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

			vector.DrawFilledRect(screen, rectX, rectY, rectWidth, rectHeight, CellMapping[cell], false)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawGrid(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenHeight, ScreenHeight
}
