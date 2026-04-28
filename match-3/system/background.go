package system

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// ScrollingBG holds the state for a looping scrolling background.
type ScrollingBG struct {
	tile    *ebiten.Image
	offsetX float64
	speed   float64 // pixels per frame
}

// NewScrollingBG creates a procedurally generated symmetric tile pattern
// and returns a ScrollingBG that scrolls left at the given speed.
func NewScrollingBG(screenW, screenH int, speed float64) *ScrollingBG {
	const tileSize = 64
	tile := generateSymmetricTile(tileSize)
	return &ScrollingBG{
		tile:  tile,
		speed: speed,
	}
}

// Update advances the scroll offset.
func (bg *ScrollingBG) Update() {
	bg.offsetX -= bg.speed
	tileW := float64(bg.tile.Bounds().Dx())
	if bg.offsetX <= -tileW {
		bg.offsetX += tileW
	}
}

// Draw renders the scrolling tiled background.
func (bg *ScrollingBG) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{15, 10, 30, 255})

	tileW := bg.tile.Bounds().Dx()
	tileH := bg.tile.Bounds().Dy()
	screenW := screen.Bounds().Dx()
	screenH := screen.Bounds().Dy()

	// Number of tiles needed to cover screen + 1 extra for scroll.
	cols := screenW/tileW + 2
	rows := screenH/tileH + 1

	for row := 0; row <= rows; row++ {
		for col := 0; col <= cols; col++ {
			op := &ebiten.DrawImageOptions{}
			x := float64(col*tileW) + bg.offsetX
			y := float64(row * tileH)
			op.GeoM.Translate(x, y)
			op.ColorScale.ScaleAlpha(0.3)
			screen.DrawImage(bg.tile, op)
		}
	}
}

// generateSymmetricTile creates a 4-fold symmetric diamond/grid pattern.
func generateSymmetricTile(size int) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	half := size / 2
	for y := range size {
		for x := range size {
			// Mirror coordinates to get 4-fold symmetry.
			mx := x
			if mx >= half {
				mx = size - 1 - mx
			}
			my := y
			if my >= half {
				my = size - 1 - my
			}

			// Diamond pattern based on Manhattan distance from center.
			dist := math.Abs(float64(mx-half/2)) + math.Abs(float64(my-half/2))
			norm := dist / float64(half)

			// Create a subtle gradient pattern.
			var c color.RGBA
			switch {
			case norm < 0.3:
				c = color.RGBA{80, 50, 120, 255} // inner diamond - purple
			case norm < 0.6:
				c = color.RGBA{40, 30, 80, 255} // mid ring - dark purple
			default:
				c = color.RGBA{20, 15, 50, 255} // outer - very dark
			}

			// Add grid lines at edges.
			if x == 0 || y == 0 || x == size-1 || y == size-1 {
				c = color.RGBA{50, 35, 90, 255}
			}

			img.SetRGBA(x, y, c)
		}
	}

	return ebiten.NewImageFromImage(img)
}
