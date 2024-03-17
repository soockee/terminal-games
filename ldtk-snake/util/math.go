package util

import (
	"math"
	"math/rand"

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

func EuclideanDistance(vec1, vec2 resolv.Vector) float64 {
	sum := (vec1.X-vec2.X)*2 + (vec1.Y-vec2.Y)*2
	return math.Sqrt(sum)
}

func DirectionVector(vec1, vec2 resolv.Vector) resolv.Vector {
	direction := resolv.NewVector(vec2.X-vec1.X, vec2.Y-vec1.Y)
	return direction
}

func RandomPointInBounds(xMin int, yMin int, xMax int, yMax int) (int, int) {
	// Ensure xMax and yMax are greater than or equal to xMin and yMin, respectively
	if xMax <= xMin || yMax <= yMin {
		panic("Invalid bounds: xMax and yMax must be greater than or equal to xMin and yMin, respectively")
	}

	// Calculate the range for x and y values
	xRange := xMax - xMin
	yRange := yMax - yMin

	// Generate random integers within the calculated ranges, starting from the lower bounds
	randomX := xMin + rand.Intn(xRange)
	randomY := yMin + rand.Intn(yRange)

	return randomX, randomY
}
