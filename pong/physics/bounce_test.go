package physics_test

import (
	"math"
	"testing"

	"github.com/soockee/terminal-games/pong/physics"
)

func approxEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

const eps = 1e-9

// --- WallBounce ---

func TestWallBounce_HorizontalNormal(t *testing.T) {
	// Ball moving right hits a wall with left-facing normal.
	r := physics.WallBounce(5, 0, -1, 0, -2, 0)

	if !approxEqual(r.VelX, -5, eps) {
		t.Errorf("VelX = %f, want -5", r.VelX)
	}
	if !approxEqual(r.VelY, 0, eps) {
		t.Errorf("VelY = %f, want 0", r.VelY)
	}
	if r.PushX != -2 || r.PushY != 0 {
		t.Errorf("Push = (%f,%f), want (-2,0)", r.PushX, r.PushY)
	}
}

func TestWallBounce_VerticalNormal(t *testing.T) {
	// Ball moving down hits a floor with upward normal.
	r := physics.WallBounce(3, 4, 0, -1, 0, -1)

	if !approxEqual(r.VelX, 3, eps) {
		t.Errorf("VelX = %f, want 3", r.VelX)
	}
	if !approxEqual(r.VelY, -4, eps) {
		t.Errorf("VelY = %f, want -4", r.VelY)
	}
}

func TestWallBounce_DiagonalNormal(t *testing.T) {
	// Ball moving right (5, 0) hits a 45-degree wall (normal = (-√2/2, -√2/2)).
	// Reflection: v' = v - 2(v·n)n
	// v·n = 5*(-√2/2) = -5√2/2
	// v' = (5,0) - 2*(-5√2/2)*(-√2/2, -√2/2) = (5,0) - (5,5) = (0, -5)
	n := -math.Sqrt2 / 2
	r := physics.WallBounce(5, 0, n, n, -1, -1)

	if !approxEqual(r.VelX, 0, eps) {
		t.Errorf("VelX = %f, want ~0", r.VelX)
	}
	if !approxEqual(r.VelY, -5, eps) {
		t.Errorf("VelY = %f, want ~-5", r.VelY)
	}
}

// --- PaddleBounce ---

func TestPaddleBounce_CenterHit(t *testing.T) {
	// Ball hits the center of the paddle → angle 0 → straight horizontal.
	r := physics.PaddleBounce(5, 3, 100, 100, 30, -2, 0)

	speed := math.Sqrt(5*5 + 3*3)

	// Ball was moving right (velX > 0), so dir = -1 → exits left.
	if !approxEqual(r.VelX, -speed, eps) {
		t.Errorf("VelX = %f, want %f", r.VelX, -speed)
	}
	// sin(0) = 0, so VelY should be 0.
	if !approxEqual(r.VelY, 0, eps) {
		t.Errorf("VelY = %f, want 0", r.VelY)
	}
}

func TestPaddleBounce_EdgeHit(t *testing.T) {
	// Ball hits the very top edge of the paddle → relativeIntersect = 1 → max angle.
	paddleCenterY := 100.0
	paddleHalfH := 30.0
	ballCenterY := paddleCenterY - paddleHalfH // top edge

	r := physics.PaddleBounce(-4, 3, ballCenterY, paddleCenterY, paddleHalfH, 2, 0)

	speed := math.Sqrt(16 + 9)
	expectedAngle := physics.MaxBounceAngle

	// Ball was moving left (velX < 0), so dir = 1 → exits right.
	expectedVelX := speed * math.Cos(expectedAngle)
	expectedVelY := -speed * math.Sin(expectedAngle)

	if !approxEqual(r.VelX, expectedVelX, eps) {
		t.Errorf("VelX = %f, want %f", r.VelX, expectedVelX)
	}
	if !approxEqual(r.VelY, expectedVelY, eps) {
		t.Errorf("VelY = %f, want %f", r.VelY, expectedVelY)
	}
}

func TestPaddleBounce_PreservesSpeed(t *testing.T) {
	// Speed should be preserved regardless of where the ball hits.
	inSpeed := math.Sqrt(4*4 + 3*3)
	r := physics.PaddleBounce(4, 3, 90, 100, 30, -1, 0)

	outSpeed := math.Sqrt(r.VelX*r.VelX + r.VelY*r.VelY)
	if !approxEqual(outSpeed, inSpeed, eps) {
		t.Errorf("out speed = %f, want %f", outSpeed, inSpeed)
	}
}

// --- CheckScore ---

func TestCheckScore_NoScore(t *testing.T) {
	r := physics.CheckScore(50, 60, 320)
	if r.Scored {
		t.Error("should not have scored")
	}
}

func TestCheckScore_BallOffLeft(t *testing.T) {
	// Ball fully past left edge → right player scores.
	r := physics.CheckScore(-10, -1, 320)
	if !r.Scored || !r.RightScored || r.LeftScored {
		t.Errorf("got %+v, want Scored+RightScored", r)
	}
}

func TestCheckScore_BallOffRight(t *testing.T) {
	// Ball fully past right edge → left player scores.
	r := physics.CheckScore(321, 331, 320)
	if !r.Scored || !r.LeftScored || r.RightScored {
		t.Errorf("got %+v, want Scored+LeftScored", r)
	}
}

func TestCheckScore_BallPartiallyOffscreen(t *testing.T) {
	// Ball partially off left (minX < 0 but maxX >= 0) → no score yet.
	r := physics.CheckScore(-5, 5, 320)
	if r.Scored {
		t.Error("should not score when ball is only partially off-screen")
	}
}

// --- ResetVelocity ---

func TestResetVelocity_MovingRight(t *testing.T) {
	// Ball was moving right → reset launches left.
	vx, vy := physics.ResetVelocity(4, 3)
	if !approxEqual(vx, -3, eps) {
		t.Errorf("VelX = %f, want -3", vx)
	}
	if !approxEqual(vy, 3, eps) {
		t.Errorf("VelY = %f, want 3", vy)
	}
}

func TestResetVelocity_MovingLeft(t *testing.T) {
	// Ball was moving left → reset launches right.
	vx, vy := physics.ResetVelocity(-4, 3)
	if !approxEqual(vx, 3, eps) {
		t.Errorf("VelX = %f, want 3", vx)
	}
	if !approxEqual(vy, 3, eps) {
		t.Errorf("VelY = %f, want 3", vy)
	}
}
