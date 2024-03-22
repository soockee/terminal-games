package util

import (
	"math"
	"math/rand"

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

func CalculateHighscore(score float64, time float64) int {
	score = score * 100

	// Define a base weight for the score (can be adjusted)
	weightScore := 0.8

	// Define a weight for the time (can be adjusted, negative for lower time is better)
	weightTime := -0.2

	// Normalize the time (optional, adjust range based on your expected times)
	normalizedTime := 1.0 / (1 + time/100) // Example normalization (scales time between 0 and 1)

	// Calculate highscore using weighted sum and rounding to int
	highscore := int(math.Floor(weightScore*score + weightTime*normalizedTime))

	// Ensure non-negative highscore (optional)
	if highscore < 0 {
		highscore = 0
	}

	return highscore
}
