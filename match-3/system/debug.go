package system

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/soockee/terminal-games/match-3/component"
	"github.com/yohamta/donburi/ecs"
)

// DrawDebug renders FPS and board info when debug mode is enabled (F3).
func DrawDebug(e *ecs.ECS, screen *ebiten.Image) {
	entry, ok := component.Debug.First(e.World)
	if !ok {
		return
	}
	d := component.Debug.Get(entry)
	if !d.Enabled {
		return
	}

	fps := fmt.Sprintf("FPS: %.0f  TPS: %.0f", ebiten.ActualFPS(), ebiten.ActualTPS())
	ebitenutil.DebugPrintAt(screen, fps, 4, 4)

	if boardEntry, ok := component.Board.First(e.World); ok {
		board := component.Board.Get(boardEntry)
		info := fmt.Sprintf("Phase: %d  Sel: (%d,%d)", board.Phase, board.SelectedCol, board.SelectedRow)
		ebitenutil.DebugPrintAt(screen, info, 4, 20)
	}
}
