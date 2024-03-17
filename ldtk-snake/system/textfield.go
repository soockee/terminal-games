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
		component.DrawScaledSprite(screen, e)
		obj := component.Object.Get(e)

		// leftAlligned := obj.Center().Mult(resolv.NewVector(0.5, 1))
		textData := component.Text.Get(e)
		if textData.IsAnimated {
			// do animation
		} else {
			drawFieldTextCovered(screen, obj, textData.Text...)
		}
	})
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

func drawFieldTextCovered	(screen *ebiten.Image, obj *resolv.Object, textLines ...string) {

	// dynamically calculate fontsize based on width and height of textfield

	h := obj.Size.Y / float64(len(textLines))
	f := text.FaceWithLineHeight(assets.NormalFont, obj.Size.Y)

	// x := int(obj.Position.X)
	dy := 0

	leftAlligned := obj.Center().X - (obj.Size.X * 0.25)
	topAlligned := obj.Center().Y - (obj.Size.Y * 0.25)

	for _, txt := range textLines {

		// Draw the text
		text.Draw(screen, txt, f, int(leftAlligned), int(topAlligned)+dy, color.RGBA{0, 0, 0, 255})

		// Move to the next line
		dy += int(h)
	}
}
