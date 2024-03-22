package system

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/soockee/terminal-games/ldtk-snake/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func UpdateTextField(ecs *ecs.ECS) {

}

func DrawTextField(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.TextField.Each(ecs.World, func(e *donburi.Entry) {
		component.DrawScaledSprite(screen, component.Sprite.Get(e).Images[0], e)
		obj := component.Object.Get(e)

		// leftAlligned := obj.Center().Mult(resolv.NewVector(0.5, 1))
		textData := component.Text.Get(e)
		if textData.IsAnimated {
			// do animation
		} else {
			if textData.Identifier == "Score" {
				drawSingleLineCentered(screen, obj, textData.Text[0], 0.5, 1.1)
			}
			if textData.Identifier == "Time" {
				drawSingleLineCentered(screen, obj, textData.Text[0], 0.9, 1.1)
			}
			if textData.Identifier == "Leaderboard" {
				drawFieldTextCovered(screen, obj, true, textData.Text...)
			}
		}
	})
}

func drawSingleLineCentered(screen *ebiten.Image, obj *resolv.Object, txt string, offsetX float64, offsetY float64) {
	shiftX := obj.Center().X * offsetX
	shiftY := obj.Center().Y * offsetY

	// Draw the text
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(4, 4)
	op.GeoM.Translate(shiftX, shiftY)
	op.ColorScale.SetR(0)
	op.ColorScale.SetG(0)
	op.ColorScale.SetB(0)
	op.ColorScale.SetA(255)

	text.DrawWithOptions(screen, txt, assets.NormalFont, op)

	// Move to the next line
}

func drawFieldText(screen *ebiten.Image, x, y int, textLines ...string) {
	f := assets.NormalFont
	lineHeight := 10
	for _, txt := range textLines {

		// Draw the text
		text.Draw(screen, txt, f, x, y, color.RGBA{0, 0, 0, 255})

		// Move to the next line
		y += lineHeight

		y += 10
	}
}

func drawFieldTextCovered(screen *ebiten.Image, obj *resolv.Object, allign bool, textLines ...string) {

	// dynamically calculate fontsize based on width and height of textfield

	h := 20.0

	// x := int(obj.Position.X)
	dy := 0.0

	shiftX := obj.Center().X * 0.9
	shiftY := obj.Center().Y * 1.1

	if allign {
		shiftX = shiftX - (obj.Size.X * 0.15)
		shiftY = shiftY - (obj.Size.Y * 0.15)
	}

	for _, txt := range textLines {
		// Draw the text
		op := &ebiten.DrawImageOptions{}
		if !allign {
			op.GeoM.Scale(4, 4)
		} else {
			op.GeoM.Scale(1, 1)
		}
		op.GeoM.Translate(shiftX, shiftY+dy)
		op.ColorScale.SetR(0)
		op.ColorScale.SetG(0)
		op.ColorScale.SetB(0)
		op.ColorScale.SetA(255)

		text.DrawWithOptions(screen, txt, assets.NormalFont, op)

		// Move to the next line
		dy += h
	}
}
