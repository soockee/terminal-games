package system

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// GenerateSymmetricTile creates a procedurally generated symmetric tile pattern
// for use as a scrolling background.
func GenerateSymmetricTile(tileSize int) *ebiten.Image {
	return generateSymmetricTile(tileSize)
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
