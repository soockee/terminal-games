package tictacgoe

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/rs/zerolog/log"
	"golang.org/x/image/font"
)

type Win struct{}

func NewWin() *Win {
	b := &Win{}
	return b
}

const (
	WinBoxHeight    = 150
	WinBoxWidth     = 200
	WinBoxTopMargin = 25
)

func (win *Win) Update(input *Input) error {
	return nil
}

func (win *Win) Draw(winImage *ebiten.Image, winnerImage *ebiten.Image, winner int, font font.Face) {
	winImage.Fill(frameColor)
	var msg string
	switch winner {
	case 0:
		msg = "Cross Won"
	case 1:
		msg = "Circle Won"
	case -1:
		msg = "Even"
	default:
		log.Fatal().Msg("undefined draw state")
	}
	text.Draw(winImage, msg, font, 10, 24, color.Black)
}

func (win *Win) Size() (int, int) {
	return WinBoxWidth, WinBoxHeight
}
