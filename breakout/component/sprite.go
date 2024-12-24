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

func ScaleSpriteToMatchBox(o *resolv.ConvexPolygon, dx, dy int, op *ebiten.DrawImageOptions) *ebiten.DrawImageOptions {
	scaleX := o.Bounds().Width() / float64(dx)
	scaleY := o.Bounds().Height() / float64(dy)
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(o.Position().X), float64(o.Position().Y))
	return op
}

func DrawRotatedSprite(screen *ebiten.Image, sprite *ebiten.Image, e *donburi.Entry, angle float64) {
	o := ConvexPolygon.Get(e)
	op := &ebiten.DrawImageOptions{}
	halfW := float64(sprite.Bounds().Dx() / 2)
	halfH := float64(sprite.Bounds().Dy() / 2)

	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(angle * math.Pi / 180.0)
	op.GeoM.Translate(halfW, halfH)

	op = ScaleSpriteToMatchBox(o, sprite.Bounds().Dx(), sprite.Bounds().Dy(), op)

	screen.DrawImage(sprite, op)
}

func DrawRotatedSpriteWithScale(screen *ebiten.Image, sprite *ebiten.Image, e *donburi.Entry, angle float64, scaleFactor float64) {
	o := ConvexPolygon.Get(e)
	op := &ebiten.DrawImageOptions{}
	halfW := float64(sprite.Bounds().Dx() / 2)
	halfH := float64(sprite.Bounds().Dy() / 2)

	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(angle * math.Pi / 180.0)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Scale(scaleFactor, scaleFactor)
	op.GeoM.Translate(float64(o.Position().X), float64(o.Position().Y))

	screen.DrawImage(sprite, op)
}

func DrawScaledSprite(screen *ebiten.Image, sprite *ebiten.Image, e *donburi.Entry) {
	o := ConvexPolygon.Get(e)
	op := &ebiten.DrawImageOptions{}
	op = ScaleSpriteToMatchBox(o, sprite.Bounds().Dx(), sprite.Bounds().Dy(), op)
	screen.DrawImage(sprite, op)
}

func DrawRepeatedSprite(screen *ebiten.Image, sprite *ebiten.Image, e *donburi.Entry) {
	o := ConvexPolygon.Get(e)
	xTimes := o.Bounds().Width() / float64(sprite.Bounds().Dx())
	yTimes := o.Bounds().Height() / float64(sprite.Bounds().Dy())
	// scaleX := float64(o.Space.CellWidth) / float64(sprite.Image.Bounds().Dx())
	// scaleY := float64(o.Space.CellHeight) / float64(sprite.Image.Bounds().Dy())
	for i := 0; i < int(xTimes); i++ {
		dx := float64(sprite.Bounds().Dx() * i)
		for j := 0; j < int(yTimes); j++ {
			dy := float64(sprite.Bounds().Dx() * j)
			op := &ebiten.DrawImageOptions{}
			// op.GeoM.Scale(scaleX, scaleY)
			op.GeoM.Translate(o.Position().X+dx, o.Position().Y+dy)
			screen.DrawImage(sprite, op)
		}
	}
}

func DrawPlaceholder(screen *ebiten.Image, o *resolv.ConvexPolygon, angle float64) {
	// op.GeoM.Scale(scaleX, scaleY)
	op := &ebiten.DrawImageOptions{}
	halfW := float64(o.Bounds().Width() / 2)
	halfH := float64(o.Bounds().Height() / 2)

	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(angle * math.Pi / 180.0)
	op.GeoM.Translate(halfW, halfH)
	op.GeoM.Translate(o.Position().X, o.Position().Y)
	rectImage := ebiten.NewImage(int(o.Bounds().Width()), int(o.Bounds().Height()))
	rectImage.Fill(color.White) // Change color as needed
	rect := rectImage.Bounds()
	vector.StrokeRect(screen, float32(o.Position().X), float32(o.Position().Y), float32(rect.Dx()), float32(rect.Dy()), 2, color.RGBA{255, 255, 255, 0}, false)

	// screen.DrawImage(rectImage, op)
}
