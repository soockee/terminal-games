package system

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/donburi/ecs"

	"github.com/soockee/terminal-games/super-mario-bros/component"
)

// DrawDebug renders debug info and FPS when debug mode is enabled (F3).
func DrawDebug(e *ecs.ECS, screen *ebiten.Image) {
	entry, ok := component.Debug.First(e.World)
	if !ok {
		return
	}
	d := component.Debug.Get(entry)
	if !d.Enabled {
		return
	}

	// FPS / TPS counter in top-left
	fps := fmt.Sprintf("FPS: %.0f  TPS: %.0f", ebiten.ActualFPS(), ebiten.ActualTPS())
	ebitenutil.DebugPrintAt(screen, fps, 4, 4)

	// TODO: add game-specific debug overlays here (e.g. collision boxes, grid).
}
