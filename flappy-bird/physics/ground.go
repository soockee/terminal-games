package physics

// ClampGround returns the displacement needed to keep boundsMaxY at or
// above groundY, and whether clamping occurred.
func ClampGround(boundsMaxY, groundY float64) (pushY float64, clamped bool) {
	if boundsMaxY > groundY {
		return groundY - boundsMaxY, true
	}
	return 0, false
}

// Overlaps returns true when two axis-aligned rects overlap.
func Overlaps(aMinX, aMinY, aMaxX, aMaxY, bMinX, bMinY, bMaxX, bMaxY float64) bool {
	return aMinX < bMaxX && aMaxX > bMinX && aMinY < bMaxY && aMaxY > bMinY
}
