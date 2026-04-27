package system

import (
	"github.com/soockee/terminal-games/flappy-bird/component"
	"github.com/soockee/terminal-games/flappy-bird/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// UpdateCamera tracks the bird horizontally, keeping it at screenW/4 from the left.
func UpdateCamera(e *ecs.ECS) {
	birdEntry, ok := tags.Bird.First(e.World)
	if !ok {
		return
	}
	camEntry, ok := component.Camera.First(e.World)
	if !ok {
		return
	}

	space := component.Space.Get(component.Space.MustFirst(e.World))
	screenW := float64(space.Width())

	birdX := component.Shape.Get(birdEntry).Shape.Bounds().Min.X
	cam := component.Camera.Get(camEntry)
	cam.X = birdX - screenW/4
	if cam.X < 0 {
		cam.X = 0
	}
}

// Projection converts world coordinates to virtual screen coordinates.
type Projection struct {
	CamX   float64
	ScaleX float64
	ScaleY float64
}

// NewProjection reads the current camera state from the world.
func NewProjection(w donburi.World) Projection {
	if entry, ok := component.Camera.First(w); ok {
		cam := component.Camera.Get(entry)
		sx, sy := cam.ScaleX, cam.ScaleY
		if sx == 0 {
			sx = 1
		}
		if sy == 0 {
			sy = 1
		}
		return Projection{CamX: cam.X, ScaleX: sx, ScaleY: sy}
	}
	return Projection{ScaleX: 1, ScaleY: 1}
}

// WorldToScreenX converts a world-space X to virtual-screen-space X.
func (p Projection) WorldToScreenX(worldX float64) float64 {
	return (worldX - p.CamX) * p.ScaleX
}

// WorldToScreenY converts a world-space Y to virtual-screen-space Y.
func (p Projection) WorldToScreenY(worldY float64) float64 {
	return worldY * p.ScaleY
}

// WorldToScreenW scales a world-space width to virtual-screen-space.
func (p Projection) WorldToScreenW(w float64) float64 {
	return w * p.ScaleX
}

// WorldToScreenH scales a world-space height to virtual-screen-space.
func (p Projection) WorldToScreenH(h float64) float64 {
	return h * p.ScaleY
}
