package tictacgoe

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type ResetButton struct {
	input Input
}

func NewResetButton(input Input) *ResetButton {
	b := &ResetButton{
		input: input,
	}
	return b
}

const (
	ResetButtonHeight       = 75
	ResetButtonBoxWidth     = 500
	ResetButtonBottomMargin = ResetButtonHeight + 3
)

func (resetButton *ResetButton) Update(game *Game) error {
	if game.input.spaceState == spaceStateSettled {
		for _, tile := range game.board.tiles {
			tile.current.TileState = -1
		}
		for i := 0; i < len(game.gameState.boards); i++ {
			game.gameState.boards[i] = 0
		}
		game.gameState.history = NewStack[int]()
		game.gameOver = false
		game.gameState.moveCounter = 0
	}
	return nil
}

func (resetButton *ResetButton) Draw(resetButtonImage *ebiten.Image, font font.Face) {
	msg := "Press Space To Reset"
	width := GetMessageDrawLength(msg, font)
	centerText := width / 2
	centerImageX := resetButtonImage.Bounds().Max.X / 2
	centerImageY := resetButtonImage.Bounds().Max.Y/2 + GetMessageMaxHeight(msg, font)
	text.Draw(resetButtonImage, msg, font, centerImageX-centerText, centerImageY, color.Black)
}

func (resetButton *ResetButton) Size() (int, int) {
	return ResetButtonBoxWidth, ResetButtonHeight
}
