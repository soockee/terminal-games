package system

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/soockee/terminal-games/pong/component"
	"github.com/soockee/terminal-games/pong/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// DrawEntities renders all paddles, balls, and walls.
func DrawEntities(e *ecs.ECS, screen *ebiten.Image) {
	// Draw paddles
	tags.Paddle.Each(e.World, func(entry *donburi.Entry) {
		s := component.Shape.Get(entry)
		bounds := s.Shape.Bounds()

		if sprite, ok := getSpriteImage(entry); ok {
			op := &ebiten.DrawImageOptions{}
			sw := float64(sprite.Bounds().Dx())
			sh := float64(sprite.Bounds().Dy())
			op.GeoM.Scale(bounds.Width()/sw, bounds.Height()/sh)
			op.GeoM.Translate(bounds.Min.X, bounds.Min.Y)
			screen.DrawImage(sprite, op)
		} else {
			c := component.Color.Get(entry).Color
			vector.DrawFilledRect(screen, float32(bounds.Min.X), float32(bounds.Min.Y),
				float32(bounds.Width()), float32(bounds.Height()), c, false)
		}
	})

	// Draw ball (circle-aware)
	if ballEntry, ok := tags.Ball.First(e.World); ok {
		s := component.Shape.Get(ballEntry)
		bounds := s.Shape.Bounds()
		cx := float32((bounds.Min.X + bounds.Max.X) / 2)
		cy := float32((bounds.Min.Y + bounds.Max.Y) / 2)
		r := float32(bounds.Width() / 2)

		if sprite, ok := getSpriteImage(ballEntry); ok {
			op := &ebiten.DrawImageOptions{}
			sw := float64(sprite.Bounds().Dx())
			sh := float64(sprite.Bounds().Dy())
			op.GeoM.Scale(bounds.Width()/sw, bounds.Height()/sh)
			op.GeoM.Translate(bounds.Min.X, bounds.Min.Y)
			screen.DrawImage(sprite, op)
		} else {
			c := component.Color.Get(ballEntry).Color
			vector.DrawFilledCircle(screen, cx, cy, r, c, false)
		}
	}

	// Draw walls (debug: thin white outlines; remove or keep as you like)
	tags.Wall.Each(e.World, func(entry *donburi.Entry) {
		s := component.Shape.Get(entry)
		bounds := s.Shape.Bounds()
		vector.StrokeRect(screen, float32(bounds.Min.X), float32(bounds.Min.Y),
			float32(bounds.Width()), float32(bounds.Height()),
			1, color.RGBA{60, 60, 60, 255}, false)
	})
}

// DrawScore renders the score at the top center.
func DrawScore(e *ecs.ECS, screen *ebiten.Image) {
	var leftScore, rightScore int
	tags.Paddle.Each(e.World, func(entry *donburi.Entry) {
		p := component.Paddle.Get(entry)
		switch p.Side {
		case component.SideLeft:
			leftScore = p.Score
		case component.SideRight:
			rightScore = p.Score
		}
	})

	space := component.Space.Get(component.Space.MustFirst(e.World))
	text := fmt.Sprintf("%d - %d", leftScore, rightScore)
	ebitenutil.DebugPrintAt(screen, text, space.Width()/2-20, 10)
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
