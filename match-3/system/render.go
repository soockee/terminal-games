package system

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/soockee/terminal-games/match-3/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

var renderTileQuery = donburi.NewQuery(filter.Contains(component.Tile))

// DrawEntities renders all tile gems at their PixelPos with Sprite.
// Gems are 16×16 sprites scaled 2x to fill 32×32 board cells.
func DrawEntities(e *ecs.ECS, screen *ebiten.Image) {
	renderTileQuery.Each(e.World, func(entry *donburi.Entry) {
		sprite, ok := getSpriteImage(entry)
		if !ok {
			return
		}
		pos := component.PixelPos.Get(entry)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(pos.X, pos.Y)
		screen.DrawImage(sprite, op)
	})

	// Draw selection highlight.
	boardEntry, ok := component.BoardGrid.First(e.World)
	if !ok {
		return
	}
	grid := component.BoardGrid.Get(boardEntry)
	phase := component.BoardPhase.Get(boardEntry)
	input := component.BoardInput.Get(boardEntry)
	display := component.BoardDisplay.Get(boardEntry)
	ts := float64(display.TileSize)

	// Draw keyboard cursor (thin border).
	if input.CursorCol >= 0 && input.CursorCol < grid.Cols && input.CursorRow >= 0 && input.CursorRow < grid.Rows {
		cx := display.OffsetX + float64(input.CursorCol*display.TileSize)
		cy := display.OffsetY + float64(input.CursorRow*display.TileSize)
		drawBorder(screen, cx, cy, ts, ts, 1, color.RGBA{200, 200, 200, 120})
	}

	// Draw selected tile (thicker highlighted border).
	if phase.SelectedCol >= 0 && phase.SelectedRow >= 0 {
		sx := display.OffsetX + float64(phase.SelectedCol*display.TileSize)
		sy := display.OffsetY + float64(phase.SelectedRow*display.TileSize)
		drawBorder(screen, sx, sy, ts, ts, 2, color.RGBA{255, 255, 100, 200})
	}
}

// DrawBackground renders the scrolling tiled background from the ECS entity.
func DrawBackground(e *ecs.ECS, screen *ebiten.Image) {
	entry, ok := component.ScrollingBG.First(e.World)
	if !ok {
		return
	}
	bg := component.ScrollingBG.Get(entry)
	if bg.Tile == nil {
		return
	}

	screen.Fill(color.RGBA{15, 10, 30, 255})

	tileW := bg.Tile.Bounds().Dx()
	tileH := bg.Tile.Bounds().Dy()
	screenW := screen.Bounds().Dx()
	screenH := screen.Bounds().Dy()

	cols := screenW/tileW + 2
	rows := screenH/tileH + 1

	for row := 0; row <= rows; row++ {
		for col := 0; col <= cols; col++ {
			op := &ebiten.DrawImageOptions{}
			x := float64(col*tileW) + bg.OffsetX
			y := float64(row * tileH)
			op.GeoM.Translate(x, y)
			op.ColorScale.ScaleAlpha(0.3)
			screen.DrawImage(bg.Tile, op)
		}
	}
}

// UpdateBackground advances the scrolling background offset.
func UpdateBackground(e *ecs.ECS) {
	entry, ok := component.ScrollingBG.First(e.World)
	if !ok {
		return
	}
	bg := component.ScrollingBG.Get(entry)
	if bg.Tile == nil {
		return
	}
	bg.OffsetX -= bg.Speed
	tileW := float64(bg.Tile.Bounds().Dx())
	if bg.OffsetX <= -tileW {
		bg.OffsetX += tileW
	}
}

// drawBorder draws a rectangular border of the given thickness.
func drawBorder(screen *ebiten.Image, x, y, w, h, thickness float64, clr color.RGBA) {
	// Top
	ebitenutil.DrawRect(screen, x, y, w, thickness, clr)
	// Bottom
	ebitenutil.DrawRect(screen, x, y+h-thickness, w, thickness, clr)
	// Left
	ebitenutil.DrawRect(screen, x, y, thickness, h, clr)
	// Right
	ebitenutil.DrawRect(screen, x+w-thickness, y, thickness, h, clr)
}

// DrawScore renders the current score and timer.
func DrawScore(e *ecs.ECS, screen *ebiten.Image) {
	entry, ok := component.Score.First(e.World)
	if !ok {
		return
	}
	score := component.Score.Get(entry)

	face := FontFace(14)
	scoreText := fmt.Sprintf("Score: %d / %d", score.Value, score.Target)
	drawText(screen, scoreText, face, 10, 18, color.White)

	// Display timer if level has a time limit.
	boardEntry, boardOK := component.BoardRules.First(e.World)
	if !boardOK {
		return
	}
	lvl := component.BoardRules.Get(boardEntry)
	if lvl.TimeLimit > 0 {
		secs := int(lvl.TimeRemaining)
		timer := fmt.Sprintf("Time: %d:%02d", secs/60, secs%60)
		drawText(screen, timer, face, 10, 36, color.White)
	}

	// Display chain multiplier when active.
	phase := component.BoardPhase.Get(boardEntry)
	if phase.ChainDepth > 1 {
		chain := fmt.Sprintf("Chain x%d!", phase.ChainDepth)
		drawText(screen, chain, face, 10, 54, color.RGBA{255, 200, 50, 255})
	}
}

// DrawHUD renders game state overlays.
func DrawHUD(e *ecs.ECS, screen *ebiten.Image) {
	entry, ok := component.GameState.First(e.World)
	if !ok {
		return
	}
	gs := component.GameState.Get(entry)
	bounds := screen.Bounds()
	cx := float64(bounds.Dx()) / 2
	cy := float64(bounds.Dy()) / 2

	// Win screen: show congratulatory message.
	if gs.WinScreen {
		bigFace := FontFace(24)
		smallFace := FontFace(14)
		drawTextCentered(screen, "Congratulations!", bigFace, cx, cy-30, color.RGBA{255, 220, 50, 255})
		drawTextCentered(screen, "You beat all levels!", smallFace, cx, cy+5, color.White)
		drawTextCentered(screen, "Press Space to play again", smallFace, cx, cy+30, color.RGBA{200, 200, 200, 255})
		return
	}

	smallFace := FontFace(12)
	medFace := FontFace(16)

	if !gs.Started {
		drawTextCentered(screen, "Tap or press Space to start", smallFace, cx, cy, color.White)
		return
	}

	// Show reshuffle notification.
	if boardEntry, boardOK := component.BoardPhase.First(e.World); boardOK {
		phase := component.BoardPhase.Get(boardEntry)
		if phase.ReshuffleTimer > 0 {
			drawTextCentered(screen, "No moves! Reshuffling...", smallFace, cx, cy-30, color.RGBA{255, 150, 50, 255})
		}
	}

	if gs.Won {
		drawTextCentered(screen, "LEVEL COMPLETE!", medFace, cx, cy-20, color.RGBA{100, 255, 100, 255})
		drawTextCentered(screen, "Tap or press Space for next level", smallFace, cx, cy+10, color.White)
		return
	}
	if gs.Dead {
		drawTextCentered(screen, "GAME OVER", medFace, cx, cy-20, color.RGBA{255, 80, 80, 255})
		drawTextCentered(screen, "Tap or press Space to retry", smallFace, cx, cy+10, color.White)
		return
	}
	if gs.Paused {
		drawTextCentered(screen, "PAUSED", medFace, cx, cy-10, color.White)
		drawTextCentered(screen, "Press P to resume", smallFace, cx, cy+15, color.RGBA{200, 200, 200, 255})
		return
	}
}

// drawText draws text at (x, y) with the given face and color.
func drawText(screen *ebiten.Image, str string, face *text.GoTextFace, x, y float64, clr color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, str, face, op)
}

// drawTextCentered draws text centered at (cx, cy).
func drawTextCentered(screen *ebiten.Image, str string, face *text.GoTextFace, cx, cy float64, clr color.Color) {
	w, h := text.Measure(str, face, 0)
	drawText(screen, str, face, cx-w/2, cy-h/2, clr)
}

func getSpriteImage(entry *donburi.Entry) (*ebiten.Image, bool) {
	if !entry.HasComponent(component.Sprite) {
		return nil, false
	}
	s := component.Sprite.Get(entry)
	if s.Image == nil {
		return nil, false
	}
	return s.Image, true
}
