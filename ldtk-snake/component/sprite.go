package component

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type SpriteData struct {
	Image *ebiten.Image
}

var Sprite = donburi.NewComponentType[SpriteData]()

func ScaleSpriteToMatchBox(o *resolv.Object, sprite *SpriteData, op *ebiten.DrawImageOptions) *ebiten.DrawImageOptions {
	scaleX := o.Size.X / float64(sprite.Image.Bounds().Dx())
	scaleY := o.Size.Y / float64(sprite.Image.Bounds().Dy())
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(o.Position.X), float64(o.Position.Y))
	return op
}

func DrawRotatedSprite(screen *ebiten.Image, e *donburi.Entry, angle float64) {
	o := Object.Get(e)
	sprite := Sprite.Get(e)
	op := &ebiten.DrawImageOptions{}
	halfW := float64(sprite.Image.Bounds().Dx() / 2)
	halfH := float64(sprite.Image.Bounds().Dy() / 2)

	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(angle * math.Pi / 180.0)
	op.GeoM.Translate(halfW, halfH)

	op = ScaleSpriteToMatchBox(o, sprite, op)

	screen.DrawImage(sprite.Image, op)
}

func DrawScaledSprite(screen *ebiten.Image, e *donburi.Entry) {
	o := Object.Get(e)
	sprite := Sprite.Get(e)
	op := &ebiten.DrawImageOptions{}
	op = ScaleSpriteToMatchBox(o, sprite, op)
	screen.DrawImage(sprite.Image, op)
}

func DrawRepeatedSprite(screen *ebiten.Image, e *donburi.Entry) {
	o := Object.Get(e)
	sprite := Sprite.Get(e)
	xTimes := o.Size.X / float64(o.Space.CellWidth)
	yTimes := o.Size.Y / float64(o.Space.CellHeight)
	for i := 0; i < int(xTimes); i++ {
		dx := float64(o.Space.CellWidth * i)
		for j := 0; j < int(yTimes); j++ {
			dy := float64(o.Space.CellWidth * j)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(o.Position.X+dx, o.Position.Y+dy)
			screen.DrawImage(sprite.Image, op)
		}
	}
}
