package tictacgoe

import "github.com/hajimehoshi/ebiten/v2"

// turnstate 0 = cross , 1 = circle
type Turn struct {
	turnState int
}

func NewTurn(turnState int) *Turn {
	return &Turn{
		turnState: turnState,
	}
}

var (
	turnImage = ebiten.NewImage(TurnBoxWidth-TurnBoxPadding, TurnBoxHeight-TurnBoxPadding)
)

const (
	TurnBoxHeight  = 30
	TurnBoxWidth   = 200
	TurnBoxPadding = 4
	TurnBoxMargin  = 35
)

// Size returns the TurnBox size.
func (turnState *Turn) Size() (int, int) {
	x := TurnBoxWidth + TurnBoxPadding
	y := TurnBoxHeight + TurnBoxPadding
	return x, y
}

// Draw draws the current tile to the given boardImage.
func (turnState *Turn) Draw(turnBoxImage *ebiten.Image) {
	turnBoxImage.Fill(frameColor)
	turnImage.Fill(turnBoxColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(TurnBoxPadding), float64(TurnBoxPadding))
	op.ColorM.ScaleWithColor(turnBoxColor)
	turnBoxImage.DrawImage(turnImage, op)

	switch turnState.turnState {
	case 0:
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(0), float64(TurnBoxHeight/2))
		// linear scaling
		sx, sy := calcDownScaleLinear(float64(crossImage.Bounds().Max.Y), float64(TurnBoxHeight))
		op.GeoM.Scale(sx, sy)
		turnBoxImage.DrawImage(crossImage, op)
	case 1:
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(0), float64(TurnBoxHeight/2))
		sx, sy := calcDownScaleLinear(float64(crossImage.Bounds().Max.Y), float64(TurnBoxHeight))
		op.GeoM.Scale(sx, sy)
		turnBoxImage.DrawImage(circleImage, op)
	}
}

// calcLinearScale calculates a linear scale
func calcDownScaleLinear(x float64, px float64) (float64, float64) {
	sx := px / x
	sy := sx
	return sx, sy
}
