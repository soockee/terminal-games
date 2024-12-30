package system

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"

	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/soockee/terminal-games/breakout/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateTextField(ecs *ecs.ECS) {

}

func DrawTextField(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.TextField.Each(ecs.World, func(e *donburi.Entry) {
		sprite := component.Sprite.Get(e)
		t := component.Text.Get(e)

		component.DrawScaledSprite(screen, sprite.Images[0], t.Shape)
		DrawText(screen, t.Shape,
			"Finish! This is a very long string that should be wrapped around the text box. And this is the last line of the text box.",
			"It also contains a linebreak",
		)
	})
}

func DrawText(screen *ebiten.Image, shape resolv.IShape, textLines ...string) {
	lineSpacingInPixels := 10.0
	f := assets.NormalFont

	Y := shape.Bounds().Min.Y
	X := shape.Bounds().Min.X
	for _, txt := range textLines {
		// Measure text width
		textWidth, textHeight := text.Measure(txt, f, lineSpacingInPixels)

		// Draw filled rectangle around the text
		vector.DrawFilledRect(screen, float32(X), float32(Y), float32(textWidth), float32(textHeight), color.RGBA{0, 0, 0, 180}, false)

		colorScale := ebiten.ColorScale{}
		op := &text.DrawOptions{}

		// Draw the text
		op.ColorScale.Reset()
		op.GeoM.Reset()
		op.GeoM.Translate(X, Y)
		colorScale.Scale(0.5, 0.5, 0.6, 1.0)

		text.Draw(screen, txt, f, op)

		// Move to the next line
		Y += lineSpacingInPixels

		Y += 10
	}
}
