package util

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func RelativeCrop(source *ebiten.Image, r image.Rectangle) *ebiten.Image {
	rx, ry := source.Bounds().Min.X+r.Min.X, source.Bounds().Min.Y+r.Min.Y
	return source.SubImage(image.Rect(rx, ry, rx+r.Max.X, ry+r.Max.Y)).(*ebiten.Image)
}
