package system

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/soockee/terminal-games/super-mario-bros/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/ganim8/v2"
)

// animationQuery matches every entity that has both an Animation and a Position.
var animationQuery = donburi.NewQuery(
	filter.Contains(component.Animation, component.Position),
)

// UpdateAnimation advances the current animation for every animated entity.
// Must be called in the Update phase (not Draw) so frame timing is correct.
func UpdateAnimation(e *ecs.ECS) {
	animationQuery.Each(e.World, func(entry *donburi.Entry) {
		anim := component.Animation.Get(entry)
		cur := anim.CurrentAnimation()
		if cur == nil {
			return
		}
		cur.Update()
	})
}

// DrawAnimation renders the current animation for every animated entity at
// its Position, using the camera projection for world-to-screen mapping.
func DrawAnimation(e *ecs.ECS, screen *ebiten.Image) {
	proj := NewProjection(e.World)

	animationQuery.Each(e.World, func(entry *donburi.Entry) {
		anim := component.Animation.Get(entry)
		cur := anim.CurrentAnimation()
		if cur == nil {
			return
		}

		pos := component.Position.Get(entry)
		sx := proj.WorldToScreenX(pos.X)
		sy := proj.WorldToScreenY(pos.Y)

		scX := proj.ScaleX
		ox := 0.0
		if anim.FlipH {
			scX = -proj.ScaleX
			ox = 1.0 // normalized origin: right edge of frame
		}
		drawOpts := ganim8.DrawOpts(sx, sy, 0, scX, proj.ScaleY, ox, 0)
		cur.Draw(screen, drawOpts)
	})
}
