package system

import (
	"github.com/soockee/terminal-games/game-template/component"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// UpdateCamera tracks a target entity horizontally.
// Customise the tracking logic for your game (e.g. follow player, scroll).
func UpdateCamera(e *ecs.ECS) {
	_, ok := component.Camera.First(e.World)
	if !ok {
		return
	}
	// TODO: implement camera tracking for your game.
	// Example: track a player entity at screenW/4 from the left.
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
