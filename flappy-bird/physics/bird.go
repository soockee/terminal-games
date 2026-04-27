package physics

// ApplyGravity adds gravitational acceleration to vertical velocity.
func ApplyGravity(velY, gravity float64) float64 {
	return velY + gravity
}

// ClampTop returns the displacement needed to keep an object's top edge
// at or below y=0, and whether clamping occurred.
func ClampTop(boundsMinY float64) (pushY float64, clamped bool) {
	if boundsMinY < 0 {
		return -boundsMinY, true
	}
	return 0, false
}
