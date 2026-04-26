package physics

import "math"

const MaxBounceAngle = math.Pi / 3 // 60 degrees

// BounceResult holds the outcome of a bounce computation:
// the new velocity and the displacement needed to push the ball
// out of the obstacle.
type BounceResult struct {
	VelX, VelY   float64
	PushX, PushY float64
}

// WallBounce reflects a velocity vector off a surface with the given normal,
// and returns the minimum translation to resolve the overlap.
func WallBounce(velX, velY, normalX, normalY, mtvX, mtvY float64) BounceResult {
	// Reflect: v' = v - 2(v·n)n
	dot := velX*normalX + velY*normalY
	return BounceResult{
		VelX:  velX - 2*dot*normalX,
		VelY:  velY - 2*dot*normalY,
		PushX: mtvX,
		PushY: mtvY,
	}
}

// PaddleBounce computes angle-based deflection. The ball's exit angle depends
// on where it hit the paddle face (center → straight, edge → steep angle).
//
// Parameters:
//   - velX, velY: incoming ball velocity
//   - ballCenterY: Y center of the ball
//   - paddleCenterY: Y center of the paddle
//   - paddleHalfH: half the paddle's height
//   - mtvX, mtvY: minimum translation vector to push ball out of paddle
func PaddleBounce(velX, velY, ballCenterY, paddleCenterY, paddleHalfH, mtvX, mtvY float64) BounceResult {
	relativeIntersect := (paddleCenterY - ballCenterY) / paddleHalfH
	angle := relativeIntersect * MaxBounceAngle

	speed := math.Sqrt(velX*velX + velY*velY)
	dir := 1.0
	if velX > 0 {
		dir = -1.0
	}

	return BounceResult{
		VelX:  dir * speed * math.Cos(angle),
		VelY:  -speed * math.Sin(angle),
		PushX: mtvX,
		PushY: mtvY,
	}
}

// ScoreResult indicates whether a point was scored and for which side.
type ScoreResult struct {
	Scored bool
	// LeftScored is true when the ball left the right edge (left player scores).
	LeftScored bool
	// RightScored is true when the ball left the left edge (right player scores).
	RightScored bool
}

// CheckScore determines whether the ball has left the screen horizontally.
//   - ballMinX, ballMaxX: horizontal bounds of the ball
//   - screenW: width of the play area
func CheckScore(ballMinX, ballMaxX, screenW float64) ScoreResult {
	if ballMaxX < 0 {
		return ScoreResult{Scored: true, RightScored: true}
	}
	if ballMinX > screenW {
		return ScoreResult{Scored: true, LeftScored: true}
	}
	return ScoreResult{}
}

// ResetVelocity returns the velocity to assign after a score reset.
// The ball launches toward the side that just scored.
func ResetVelocity(currentVelX, ballSpeed float64) (velX, velY float64) {
	dir := 1.0
	if currentVelX > 0 {
		dir = -1.0
	}
	return dir * ballSpeed, ballSpeed
}
