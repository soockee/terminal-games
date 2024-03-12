package system

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/soockee/terminal-games/ldtk-snake/assets"
	"github.com/soockee/terminal-games/ldtk-snake/component"
	"github.com/yohamta/donburi/ecs"
	"golang.org/x/image/font"
)

const lineHeight = 10 // Adjust this as needed

func DrawHelp(ecs *ecs.ECS, screen *ebiten.Image) {
	settings, ok := GetSettings(ecs)
	if !ok {
		return
	}
	if settings.ShowHelpText {
		spaceEntry, ok := component.Space.First(ecs.World)
		if !ok {
			return
		}
		space := component.Space.Get(spaceEntry)
		cell := space.Cell(4, 4)
		cw := space.CellWidth
		ch := space.CellHeight
		cx := cell.X * cw
		cy := cell.Y * ch

		drawText(screen, cx, cy,
			"~ Snake ~",
			"Move Snake: W,A,S,D or Arrow Keys",
			"",
			"F1: Toggle Debug View",
			"F2: Show / Hide help text",
			fmt.Sprintf("%d FPS (frames per second)", int(ebiten.ActualFPS())),
			fmt.Sprintf("%d TPS (ticks per second)", int(ebiten.ActualTPS())),
		)
	}
}

func drawText(screen *ebiten.Image, x, y int, textLines ...string) {
	f := assets.NormalFont
	for _, txt := range textLines {
		// Measure text width
		textWidth := font.MeasureString(f, txt)

		// Calculate rectangle dimensions
		rectWidth := float32(textWidth.Round())

		// Draw filled rectangle around the text
		vector.DrawFilledRect(screen, float32(x), float32(y-lineHeight-5), rectWidth, lineHeight*2, color.RGBA{0, 0, 0, 180}, false)

		// Draw the text
		text.Draw(screen, txt, f, x+1, y+1, color.RGBA{0, 0, 150, 255})
		text.Draw(screen, txt, f, x, y, color.RGBA{100, 150, 255, 255})

		// Move to the next line
		y += lineHeight

		y += 10
	}
}
