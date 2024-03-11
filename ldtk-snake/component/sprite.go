package component

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type SpriteData struct {
	Image *ebiten.Image
}

var Sprite = donburi.NewComponentType[SpriteData]()

func ScaleSpriteToMatchBox(o *resolv.Object, sprite *SpriteData) *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	scaleX := o.Size.X / float64(sprite.Image.Bounds().Dx())
	scaleY := o.Size.Y / float64(sprite.Image.Bounds().Dy())
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(o.Position.X), float64(o.Position.Y))
	return op
}

func DrawSprite(screen *ebiten.Image, e *donburi.Entry) {
	o := Object.Get(e)
	sprite := Sprite.Get(e)
	op := ScaleSpriteToMatchBox(o, sprite)
	screen.DrawImage(sprite.Image, op)
}
