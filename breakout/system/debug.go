package system

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/soockee/terminal-games/breakout/component"
	"github.com/yohamta/donburi/ecs"
)

func DrawDebug(ecs *ecs.ECS, screen *ebiten.Image) {
	settings, ok := GetSettings(ecs)
	if !ok {
		return
	}
	if !settings.Debug {
		return
	}
	space := component.Space.Get(component.Space.MustFirst(ecs.World))

	for y := 0; y < space.Height(); y++ {

		for x := 0; x < space.Width(); x++ {

			cell := space.Cell(x, y)

			cw := float32(space.CellWidth())
			ch := float32(space.CellHeight())
			cx := float32(cell.X) * cw
			cy := float32(cell.Y) * ch

			strokeWidth := float32(2.0)

			var drawColor color.Color
			drawColor = color.RGBA{20, 20, 20, 255}
			if cell.IsOccupied() {
				drawColor = color.RGBA{255, 255, 0, 255}
			}

			vector.StrokeLine(screen, cx, cy, cx+cw, cy, strokeWidth, drawColor, false)

			vector.StrokeLine(screen, cx+cw, cy, cx+cw, cy+ch, strokeWidth, drawColor, false)

			vector.StrokeLine(screen, cx+cw, cy+ch, cx, cy+ch, strokeWidth, drawColor, false)

			vector.StrokeLine(screen, cx, cy+ch, cx, cy, strokeWidth, drawColor, false)
		}
	}
}
