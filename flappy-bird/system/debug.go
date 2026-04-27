package system

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/flappy-bird/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var debugColor = color.RGBA{0, 255, 0, 180}

// DrawDebug renders collision boxes and FPS when debug mode is enabled (F3).
func DrawDebug(e *ecs.ECS, screen *ebiten.Image) {
	entry, ok := component.Debug.First(e.World)
	if !ok {
		return
	}
	d := component.Debug.Get(entry)
	if !d.Enabled {
		return
	}

	proj := NewProjection(e.World)

	// Draw collision shape for every entity that has a Shape component.
	component.Shape.Each(e.World, func(entry *donburi.Entry) {
		s := component.Shape.Get(entry)
		bounds := s.Shape.Bounds()

		// Circle shapes get a circle outline; everything else gets a rect.
		if circle, ok := s.Shape.(*resolv.Circle); ok {
			cx := float32(proj.WorldToScreenX(circle.Position().X))
			cy := float32(proj.WorldToScreenY(circle.Position().Y))
			r := float32(float64(circle.Radius()) * proj.ScaleX)
			vector.StrokeCircle(screen, cx, cy, r, 1, debugColor, false)
		} else {
			x := float32(proj.WorldToScreenX(bounds.Min.X))
			y := float32(proj.WorldToScreenY(bounds.Min.Y))
			w := float32(proj.WorldToScreenW(bounds.Width()))
			h := float32(proj.WorldToScreenH(bounds.Height()))
			vector.StrokeRect(screen, x, y, w, h, 1, debugColor, false)
		}

		// Label with tag name
		label := entityLabel(entry)
		if label != "" {
			ebitenutil.DebugPrintAt(screen, label, int(proj.WorldToScreenX(bounds.Min.X)), int(proj.WorldToScreenY(bounds.Min.Y))-12)
		}
	})

	// FPS / TPS counter in top-left
	fps := fmt.Sprintf("FPS: %.0f  TPS: %.0f", ebiten.ActualFPS(), ebiten.ActualTPS())
	ebitenutil.DebugPrintAt(screen, fps, 4, 4)
}

func entityLabel(entry *donburi.Entry) string {
	switch {
	default:
		return ""
	}
}
