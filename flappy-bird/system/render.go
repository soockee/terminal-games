package system

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/soockee/terminal-games/flappy-bird/component"
	"github.com/soockee/terminal-games/flappy-bird/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// DrawEntities renders the bird, pipes, and ground.
func DrawEntities(e *ecs.ECS, screen *ebiten.Image) {
	proj := NewProjection(e.World)

	// Draw pipes.
	tags.Pipe.Each(e.World, func(entry *donburi.Entry) {
		drawEntity(screen, entry, proj)
	})

	// Draw bird.
	if entry, ok := tags.Bird.First(e.World); ok {
		drawEntity(screen, entry, proj)
	}

	// Draw ground.
	if gEntry, ok := tags.Ground.First(e.World); ok {
		g := component.Ground.Get(gEntry)
		drawGround(screen, g.Tile, g.Y, proj)
	}
}

func drawEntity(screen *ebiten.Image, entry *donburi.Entry, proj Projection) {
	s := component.Shape.Get(entry)
	bounds := s.Shape.Bounds()
	sx := proj.WorldToScreenX(bounds.Min.X)
	sy := proj.WorldToScreenY(bounds.Min.Y)
	sw := proj.WorldToScreenW(bounds.Width())
	sh := proj.WorldToScreenH(bounds.Height())

	// Pipes use native aspect ratio, anchored at the gap edge.
	if entry.HasComponent(component.Pipe) {
		if sprite, ok := getSpriteImage(entry); ok {
			drawPipe(screen, sprite, sx, sy, sw, sh, component.Pipe.Get(entry).FlipY)
		}
		return
	}

	if sprite, ok := getSpriteImage(entry); ok {
		op := &ebiten.DrawImageOptions{}
		iw := float64(sprite.Bounds().Dx())
		ih := float64(sprite.Bounds().Dy())
		op.GeoM.Scale(sw/iw, sh/ih)
		op.GeoM.Translate(sx, sy)
		screen.DrawImage(sprite, op)
	} else {
		c := component.Color.Get(entry).Color
		vector.DrawFilledRect(screen, float32(sx), float32(sy),
			float32(sw), float32(sh), c, false)
	}
}

// drawPipe draws a pipe sprite stretched to exactly fill the collision box.
// Top pipes (flipY) are drawn upside-down so the opening faces the gap.
// The sprite is taller than the collision box — excess extends off-screen,
// avoiding any visual gaps without distorting the texture.
func drawPipe(screen *ebiten.Image, sprite *ebiten.Image, sx, sy, screenW, screenH float64, flipY bool) {
	iw := float64(sprite.Bounds().Dx())
	scale := screenW / iw // uniform scale to match collision width

	op := &ebiten.DrawImageOptions{}
	if flipY {
		// Opening at gap top edge (sy + screenH). Sprite extends upward off-screen.
		op.GeoM.Scale(scale, -scale)
		op.GeoM.Translate(sx, sy+screenH)
	} else {
		// Opening at gap bottom edge (sy). Sprite extends downward off-screen.
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(sx, sy)
	}
	screen.DrawImage(sprite, op)
}

// DrawScore renders the current score at the top center.
func DrawScore(e *ecs.ECS, screen *ebiten.Image) {
	sw := screen.Bounds().Dx()
	score := 0
	if entry, ok := component.Score.First(e.World); ok {
		score = component.Score.Get(entry).Value
	}
	text := fmt.Sprintf("%d", score)
	ebitenutil.DebugPrintAt(screen, text, sw/2, 10)
}

// pauseButtonBounds returns the pause button rectangle [x, y, w, h] in virtual screen coordinates.
func pauseButtonBounds(screenW int) (x, y, w, h int) {
	return screenW - 44, 4, 40, 40
}

// drawPauseButton draws a pixel-art style pause/resume button.
func drawPauseButton(screen *ebiten.Image, bx, by, bw, bh int, paused bool) {
	// Background
	vector.DrawFilledRect(screen, float32(bx), float32(by), float32(bw), float32(bh), color.RGBA{20, 20, 30, 220}, false)
	// 1px pixel border
	border := color.RGBA{180, 220, 180, 255}
	vector.DrawFilledRect(screen, float32(bx), float32(by), float32(bw), 1, border, false)
	vector.DrawFilledRect(screen, float32(bx), float32(by+bh-1), float32(bw), 1, border, false)
	vector.DrawFilledRect(screen, float32(bx), float32(by), 1, float32(bh), border, false)
	vector.DrawFilledRect(screen, float32(bx+bw-1), float32(by), 1, float32(bh), border, false)

	cx := float32(bx + bw/2)
	cy := float32(by + bh/2)
	icon := color.RGBA{180, 220, 180, 255}

	if paused {
		// Pixel play triangle (pointing right): 3 rows of decreasing width
		vector.DrawFilledRect(screen, cx-4, cy-6, 3, 12, icon, false)
		vector.DrawFilledRect(screen, cx-1, cy-4, 3, 8, icon, false)
		vector.DrawFilledRect(screen, cx+2, cy-2, 3, 4, icon, false)
	} else {
		// Two pixel bars (pause symbol)
		vector.DrawFilledRect(screen, cx-6, cy-6, 4, 12, icon, false)
		vector.DrawFilledRect(screen, cx+2, cy-6, 4, 12, icon, false)
	}
}

// DrawHUD shows contextual instructions (start / game over / win + restart).
func DrawHUD(e *ecs.ECS, screen *ebiten.Image) {
	goEntry, ok := component.GameState.First(e.World)
	if !ok {
		return
	}
	go_ := component.GameState.Get(goEntry)

	sw := screen.Bounds().Dx()
	sh := screen.Bounds().Dy()

	// Draw pause button when game is active (started, not dead/won).
	if go_.Started && !go_.Dead && !go_.Won {
		bx, by, bw, bh := pauseButtonBounds(sw)
		drawPauseButton(screen, bx, by, bw, bh, go_.Paused)
	}

	if go_.Paused {
		ebitenutil.DebugPrintAt(screen, "PAUSED", sw/2-18, sh/2-8)
		ebitenutil.DebugPrintAt(screen, "Press Escape or P / tap II to resume", sw/2-108, sh/2+8)
	} else if go_.Won {
		ebitenutil.DebugPrintAt(screen, "YOU ABSOLUTE MADLAD!", sw/2-62, sh/2-32)
		ebitenutil.DebugPrintAt(screen, "Congratulations on your insanity.", sw/2-100, sh/2-16)
		ebitenutil.DebugPrintAt(screen, "All pipes cleared. Devil Teemo approves.", sw/2-120, sh/2)
		ebitenutil.DebugPrintAt(screen, "Tap or press Space to play again", sw/2-96, sh/2+24)
	} else if go_.Dead {
		ebitenutil.DebugPrintAt(screen, "GAME OVER", sw/2-30, sh/2-16)
		ebitenutil.DebugPrintAt(screen, "Tap or press Space to restart", sw/2-90, sh/2)
	} else if !go_.Started {
		ebitenutil.DebugPrintAt(screen, "Tap or press Space to start", sw/2-84, sh/2)
	}
}

func drawGround(screen *ebiten.Image, tile *ebiten.Image, groundY float64, proj Projection) {
	screenW := float64(screen.Bounds().Dx())
	tileW := float64(tile.Bounds().Dx()) * proj.ScaleX
	camX := proj.CamX * proj.ScaleX
	startX := -math.Mod(camX, tileW)
	sy := proj.WorldToScreenY(groundY)
	op := &ebiten.DrawImageOptions{}
	for x := startX; x < screenW+tileW; x += tileW {
		op.GeoM.Reset()
		op.GeoM.Scale(proj.ScaleX, proj.ScaleY)
		op.GeoM.Translate(x, sy)
		screen.DrawImage(tile, op)
	}
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
