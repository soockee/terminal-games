package main

import (
	"github/soockee/terminal-games/tictacgoe/internal/tictacgoe"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	game, err := tictacgoe.NewGame()
	if err != nil {
		log.Err(err)
	}
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("TicTacGoe")
	if err := ebiten.RunGame(game); err != nil {
		log.Err(err)
	}
}
