package util

import (
	"math"

	"github.com/solarlune/resolv"
)

func radToDegree(rad float64) float64 {
	return rad * (180 / math.Pi)
}

func CalculateAngleBetweenVectors(referenceVector resolv.Vector, direction resolv.Vector) float64 {
	// Calculate the dot product of the vectors
	dotProduct := referenceVector.Dot(direction)

	// Calculate the cosine of the angle between the vectors
	cosTheta := dotProduct / (referenceVector.Magnitude() * direction.Magnitude())

	// Handle potential numerical errors due to floating-point precision
	if cosTheta > 1.0 {
		cosTheta = 1.0
	} else if cosTheta < -1.0 {
		cosTheta = -1.0
	}

	// Calculate the angle in radians using arccosine
	rad := math.Acos(cosTheta)

	// Convert radians to degrees
	degree := radToDegree(rad)

	return degree
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
	degree := radToDegree(rad)

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
