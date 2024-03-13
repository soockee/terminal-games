package util

import (
	"math"

	"github.com/solarlune/resolv"
)

func radToDegre(rad float64) float64 {
	return rad * (180 / math.Pi)
}

func CalculateAngle(direction resolv.Vector) float64 {
	if direction.X == 0 {
		if direction.Y == 0 {
			return 0.0
		} else if direction.Y > 0 {
			return 90.0
		} else {
			return 270.0
		}
	}

	rad := math.Atan2(direction.Y, direction.X)
	degree := radToDegre(rad)

	return degree
}
