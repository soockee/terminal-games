package tictacgoe

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
)

type TileState int

const (
	TileStateEmpty  = -1
	TileStateCross  = 0
	TileStateCircle = 1
)

// TileData represents a tile information like a value and a position.
type TileData struct {
	TileState TileState
	x         int
	y         int
}

// Tile represents a tile information including TileData and animation states.
type Tile struct {
	current TileData
}

// NewTile creates a new Tile object.
func NewTile(value TileState, x, y int) *Tile {
	return &Tile{
		current: TileData{
			TileState: value,
			x:         x,
			y:         y,
		},
	}
}

const (
	tileSize   = 300
	tileMargin = 4
)

func init() {
	tileImage.Fill(color.White)
}

// Draw draws the current tile to the given boardImage.
func (t *Tile) Draw(boardImage *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	x := t.current.x*tileSize + (t.current.x+1)*tileMargin
	y := t.current.y*tileSize + (t.current.y+1)*tileMargin
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorM.ScaleWithColor(tileColor)
	switch tileState := t.current.TileState; tileState {
	case TileStateEmpty:
		boardImage.DrawImage(tileImage, op)
	case TileStateCross:
		boardImage.DrawImage(crossImage, op)
	case TileStateCircle:
		boardImage.DrawImage(circleImage, op)
	default:
		log.Fatal().Msg("invalid tile state")
	}
}
