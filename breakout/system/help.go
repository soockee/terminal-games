package system

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/soockee/terminal-games/breakout/assets"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/yohamta/donburi/ecs"
)

func DrawHelp(ecs *ecs.ECS, screen *ebiten.Image) {
	settings, ok := GetSettings(ecs)
	if !ok {
		return
	}
	if settings.ShowHelpText {
		// top left croner
		// space := component.Space.Get(component.Space.MustFirst(ecs.World))
		// cell := space.Cell(4, 4)
		// cw := space.CellWidth()
		// ch := space.CellHeight()
		// cx := cell.X * cw
		// cy := cell.Y * ch

		// center
		space := component.Space.Get(component.Space.MustFirst(ecs.World))
		shift := float64(space.HeightInCells()) * 0.75
		cell := space.Cell(4, int(math.Ceil(shift)))
		cw := space.CellWidth()
		ch := space.CellHeight()
		cx := float32(cell.X * cw)
		cy := float32(cell.Y * ch)

		drawHelpText(screen, cx, cy,
			"~ Snake ~",
			"Move Snake: W,A,S,D or Arrow Keys",
			"Halt Snake:Space",
			"",
			"F1: Toggle Debug View",
			"F2: Show / Hide help text",
			fmt.Sprintf("%d FPS (frames per second)", int(ebiten.ActualFPS())),
			fmt.Sprintf("%d TPS (ticks per second)", int(ebiten.ActualTPS())),
		)
	}
}

func drawHelpText(screen *ebiten.Image, x, y float32, textLines ...string) {
	lineSpacingInPixels := 10.0
	f := assets.NormalFont

	for _, txt := range textLines {
		// Measure text width
		textWidth, textHeight := text.Measure(txt, f, lineSpacingInPixels)

		// Draw filled rectangle around the text
		vector.DrawFilledRect(screen, x, y, float32(textWidth), float32(textHeight), color.RGBA{0, 0, 0, 180}, false)

		colorScale := ebiten.ColorScale{}
		colorScale.Scale(0.5, 0.5, 0.6, 1.0)

		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(x+1), float64(y+1))
		op.ColorScale = colorScale

		// Draw the text
		text.Draw(screen, txt, f, op)

		op.ColorScale.Reset()
		op.GeoM.Reset()
		op.GeoM.Translate(float64(x), float64(y))
		colorScale.Scale(0.5, 0.5, 0.6, 1.0)

		text.Draw(screen, txt, f, op)

		// Move to the next line
		y += float32(lineSpacingInPixels)

		y += 10
	}
}
