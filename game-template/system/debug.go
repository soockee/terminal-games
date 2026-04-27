package system

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
	"github.com/soockee/terminal-games/game-template/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var debugColor = color.RGBA{0, 255, 0, 180}

// DrawDebug renders collision boxes and FPS when debug mode is enabled (F3).
// Uses the camera projection so overlays match the rendered positions.
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

		sx := float32(proj.WorldToScreenX(bounds.Min.X))
		sy := float32(proj.WorldToScreenY(bounds.Min.Y))
		sw := float32(proj.WorldToScreenW(bounds.Width()))
		sh := float32(proj.WorldToScreenH(bounds.Height()))

		// Circle shapes get a circle outline; everything else gets a rect.
		if circle, ok := s.Shape.(*resolv.Circle); ok {
			cr := float32(proj.WorldToScreenW(circle.Radius()))
			vector.StrokeCircle(screen, sx+sw/2, sy+sh/2, cr, 1, debugColor, false)
		} else {
			vector.StrokeRect(screen, sx, sy, sw, sh, 1, debugColor, false)
		}

		// Label with tag name
		label := entityLabel(entry)
		if label != "" {
			ebitenutil.DebugPrintAt(screen, label, int(sx), int(sy)-12)
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
