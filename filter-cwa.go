package main

import (
	"fmt"
	"math"
)

func filterCWA(runA, runB run) float64 {
	const tippintPoint = 120
	var bouns float64 = 0
	refCoor := runA.Coors[0]
	numLargeChunksGoal := 50
	largeChunksShapeA := producePointsFilter(runA, refCoor, numLargeChunksGoal)
	largeChunksShapeB := producePointsFilter(runB, refCoor, numLargeChunksGoal)

	// fold them
	foldedLargeChunksShapeA, foldedLargeChunksShapeB := fold(largeChunksShapeA, largeChunksShapeB)

	// tilt them
	foldedTiltedLargeChunksShapeA, foldedTiltedLargeChunksShapeB := tilt(foldedLargeChunksShapeA, foldedLargeChunksShapeB)

	// distance
	foldedTiltedDistanceA := distanceForShape(foldedTiltedLargeChunksShapeA)
	foldedTiltedDistanceB := distanceForShape(foldedTiltedLargeChunksShapeB)

	// left right area
	leftRightAreaA := leftRightArea(foldedTiltedLargeChunksShapeA)
	leftRightAreaB := leftRightArea(foldedTiltedLargeChunksShapeB)

	// clockwise area
	clockwiseAreaA := leftRightAreaA.clockwiseArea()
	clockwiseAreaB := leftRightAreaB.clockwiseArea()

	// aggr area
	aggrAreaA := leftRightAreaA.aggrArea()
	aggrAreaB := leftRightAreaB.aggrArea()

	// i should only care about differences in clockwise
	// area if it is large enough on both
	shortSideRectangleA := smallSideOfRectangle(math.Abs(leftRightAreaA.unfiltered().clockwiseArea()), foldedTiltedDistanceA)
	shortSideRectangleB := smallSideOfRectangle(math.Abs(leftRightAreaB.unfiltered().clockwiseArea()), foldedTiltedDistanceB)
	shortestSide := shortSideRectangleA
	longestSide := shortSideRectangleB
	if shortSideRectangleB < shortSideRectangleA {
		shortestSide = shortSideRectangleB
		longestSide = shortSideRectangleA
	}

	if shortestSide < 2 && longestSide < 4 {
		// cannot compare the clockwise area
	} else {
		// can compare the clockwise area
		clockwiseFactorDiff := getFactorDiff(clockwiseAreaA, clockwiseAreaB)
		if clockwiseFactorDiff > 1.2 {
			bouns = tippintPoint + 1
		}
	}

	aggrAreaDiff := getFactorDiff(aggrAreaA, aggrAreaB)
	if aggrAreaDiff > 1.4 {
		bouns = tippintPoint + 1
	}

	distanceDiff := getFactorDiff(foldedTiltedDistanceA, foldedTiltedDistanceB)
	if distanceDiff > 1.4 {
		bouns = tippintPoint + 1
	}

	areaPointsA := areaPointsFromPoints(foldedTiltedLargeChunksShapeA)
	areaPointsB := areaPointsFromPoints(foldedTiltedLargeChunksShapeB)
	areaDiff := areaPointDiff(areaPointsA, areaPointsB)
	avgAreaDiff := areaDiff.avg()
	worstAreaDiff := areaDiff.worst()
	fmt.Printf("%+v\n", areaDiff)
	fmt.Println("Avg: ", avgAreaDiff)
	fmt.Println("Worst: ", worstAreaDiff)

	return worstAreaDiff + bouns
}

// they should be able to be negative as well
func getFactorDiff(a, b float64) float64 {
	var smallestAbs float64
	if math.Abs(a) < math.Abs(b) {
		smallestAbs = math.Abs(a)
	} else {
		smallestAbs = math.Abs(b)
	}
	diff := math.Abs(a - b)
	fraction := (diff / smallestAbs) + 1
	return fraction
}
