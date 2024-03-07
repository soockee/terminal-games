package tictacgoe

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
)

// Board represents the game board.
type Board struct {
	size  int
	tiles map[int]*Tile
}

// NewBoard generates a new Board with giving a size.
func NewBoard(size int) (*Board, error) {
	b := &Board{
		size:  size,
		tiles: map[int]*Tile{},
	}
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			b.tiles[j*size+i] = NewTile(TileStateEmpty, i, j)
		}
	}
	return b, nil
}

func init() {
	log.Debug().Msg("Init Board Object")
}

func (b *Board) Draw(boardImage *ebiten.Image) {
	boardImage.Fill(frameColor)
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			op := &ebiten.DrawImageOptions{}
			x := i*tileSize + (i+1)*tileMargin
			y := j*tileSize + (j+1)*tileMargin
			op.GeoM.Translate(float64(x), float64(y))
			op.ColorM.ScaleWithColor(tileColor)
			boardImage.DrawImage(tileImage, op)
		}
	}
	for _, tile := range b.tiles {
		tile.Draw(boardImage)
	}
}

// Size returns the board size.
func (b *Board) Size() (int, int) {
	x := b.size*tileSize + (b.size+1)*tileMargin
	y := x
	return x, y
}
