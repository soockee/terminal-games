package component

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type SpriteData struct {
	Images map[int]*ebiten.Image
}

var Sprite = donburi.NewComponentType[SpriteData]()

func ScaleSpriteToMatchBox(shape resolv.IShape, dx, dy int, op *ebiten.DrawImageOptions) *ebiten.DrawImageOptions {
	scaleX := shape.Bounds().Width() / float64(dx)
	scaleY := shape.Bounds().Height() / float64(dy)
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(shape.Bounds().Min.X), float64(shape.Bounds().Min.Y))
	return op
}

func DrawRotatedSprite(screen *ebiten.Image, sprite *ebiten.Image, shape resolv.IShape, angle float64) {
	op := &ebiten.DrawImageOptions{}
	halfW := float64(sprite.Bounds().Dx() / 2)
	halfH := float64(sprite.Bounds().Dy() / 2)

	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(angle * math.Pi / 180.0)
	op.GeoM.Translate(halfW, halfH)

	op = ScaleSpriteToMatchBox(shape, sprite.Bounds().Dx(), sprite.Bounds().Dy(), op)

	screen.DrawImage(sprite, op)
}

func DrawRotatedSpriteWithScale(screen *ebiten.Image, sprite *ebiten.Image, shape resolv.IShape, angle float64, scaleFactor float64) {
	op := &ebiten.DrawImageOptions{}
	halfW := float64(sprite.Bounds().Dx() / 2)
	halfH := float64(sprite.Bounds().Dy() / 2)

	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(angle * math.Pi / 180.0)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Scale(scaleFactor, scaleFactor)
	op.GeoM.Translate(float64(shape.Position().X), float64(shape.Position().Y))

	screen.DrawImage(sprite, op)
}

func DrawScaledSprite(screen *ebiten.Image, sprite *ebiten.Image, shape resolv.IShape) {
	op := &ebiten.DrawImageOptions{}
	op = ScaleSpriteToMatchBox(shape, sprite.Bounds().Dx(), sprite.Bounds().Dy(), op)
	screen.DrawImage(sprite, op)
}

func DrawRepeatedSprite(screen *ebiten.Image, sprite *ebiten.Image, shape resolv.IShape) {
	xTimes := shape.Bounds().Width() / float64(sprite.Bounds().Dx())
	yTimes := shape.Bounds().Height() / float64(sprite.Bounds().Dy())

	// round to nearest whole number, otherwise cause rendering issues
	xTimes = math.Round(xTimes)
	yTimes = math.Round(yTimes)

	for i := 0; i < int(xTimes); i++ {
		// dx := float64(sprite.Bounds().Dx() * i)
		dx := float64(sprite.Bounds().Dx() * i)
		for j := 0; j < int(yTimes); j++ {
			dy := float64(sprite.Bounds().Dx() * j)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(shape.Bounds().Min.X+dx, shape.Bounds().Min.Y+dy)
			screen.DrawImage(sprite, op)
		}
	}
}

func DrawPlaceholder(screen *ebiten.Image, sprite *ebiten.Image, shape resolv.IShape, angle float64, c color.Color, fill bool) {
	op := &ebiten.DrawImageOptions{}
	halfW := float64(shape.Bounds().Width() / 2)
	halfH := float64(shape.Bounds().Height() / 2)

	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(angle * math.Pi / 180.0)
	op.GeoM.Translate(halfW, halfH)
	op.GeoM.Translate(shape.Bounds().Min.X, shape.Bounds().Min.Y)
	switch s := shape.(type) {
	case *resolv.Circle:
		if fill {
			vector.DrawFilledCircle(screen, float32(s.Bounds().Min.X), float32(s.Bounds().Min.Y), float32(s.Radius()), c, true)
		} else {
			vector.DrawFilledCircle(screen, float32(s.Bounds().Min.X), float32(s.Bounds().Min.Y), float32(s.Radius()), c, true)
		}
	case *resolv.ConvexPolygon:
		if fill {
			vector.DrawFilledRect(screen, float32(s.Bounds().Min.X), float32(s.Bounds().Min.Y), float32(s.Bounds().Width()), float32(s.Bounds().Height()), c, true)
		} else {
			vector.DrawFilledRect(screen, float32(s.Bounds().Min.X), float32(s.Bounds().Min.Y), float32(s.Bounds().Width()), float32(s.Bounds().Height()), c, false)
		}
	}
}
