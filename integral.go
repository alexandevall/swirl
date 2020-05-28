package main

import (
	"math"
)

type integralResult struct {
	positiveArea float64
	negativeArea float64
}

type clockwiseCurveDiffResult struct {
	curveDiff                   float64
	aggrLeftRightAreaFactorDiff float64
}

func clockwiseCurveDiff(shapeA, shapeB []areaPoint) clockwiseCurveDiffResult {
	curveA := getClockwiseAreaCurveFrom(shapeA)
	curveB := getClockwiseAreaCurveFrom(shapeB)
	printJSON(curveA)
	maxAggrAreaA := curveA[len(curveA)-1].X
	maxAggrAreaB := curveB[len(curveB)-1].X
	maxAggrSmallest := maxAggrAreaA
	maxAggrLargest := maxAggrAreaB
	if maxAggrAreaA > maxAggrAreaB {
		maxAggrSmallest = maxAggrAreaB
		maxAggrLargest = maxAggrAreaA
	}
	aggrAreaFactor := maxAggrLargest / maxAggrSmallest
	curveDiff := integralDiffBetweenCurves(curveA, curveB)
	return clockwiseCurveDiffResult{curveDiff, aggrAreaFactor}
}

func getClockwiseAreaCurveFrom(shape []areaPoint) []point {
	result := make([]point, len(shape))
	for index, p := range shape {
		aggrArea := p.aggrArea()
		clockwiseArea := p.clockwiseArea()
		clockwiseCurvePoint := point{aggrArea, clockwiseArea}
		result[index] = clockwiseCurvePoint
	}
	return result
}

func (c clockwiseCurveDiffResult) score() float64 {
	return c.curveDiff * c.aggrLeftRightAreaFactorDiff
}

func (ir *integralResult) add(area float64) {
	if area < 0 {
		ir.negativeArea += area
	} else {
		ir.positiveArea += area
	}
}

func (ir integralResult) aggrArea() float64 {
	return ir.negativeArea + ir.positiveArea
}

func integralForPoints(curve []point) integralResult {
	count := len(curve)
	var result integralResult

	splitLineAtXAxsis := func(line pointsLine) (leftLine, rightLine pointsLine) {
		intersection := line.asLinearEquation().solveForY(0)
		leftLine = pointsLine{line.a, intersection}
		rightLine = pointsLine{intersection, line.b}
		return
	}

	for i := 1; i < count; i++ {
		leftPoint := curve[i-1]
		rightPoint := curve[i]
		line := pointsLine{leftPoint, rightPoint}
		if isPositive(leftPoint.Y) != isPositive(rightPoint.Y) {
			// the points cross the x axsis.
			// divide the point into two lines
			leftLine, rightLine := splitLineAtXAxsis(line)
			result.add(leftLine.edgeArea())
			result.add(rightLine.edgeArea())
		} else {
			result.add(line.edgeArea())
		}
	}
	return result
}

func integralDiffBetweenCurves(curveA, curveB []point) float64 {

	var result float64
	var indexA int = 0
	var indexB int = 0
	countA := len(curveA)
	countB := len(curveB)
	leftPointA := curveA[indexA]
	leftPointB := curveB[indexB]
	indexA++
	indexB++

	next := func() (bool, point, point) {
		rightPointA := curveA[indexA]
		rightPointB := curveB[indexB]
		var xOnAIsSmaller bool
		if rightPointA.X < rightPointB.X {
			indexA++
			xOnAIsSmaller = true
			return xOnAIsSmaller, rightPointA, rightPointB
		}
		indexB++
		xOnAIsSmaller = false
		return xOnAIsSmaller, rightPointA, rightPointB
	}

	cuttLineSegment := func(x float64, line pointsLine) pointsLine {
		eq := line.asLinearEquation()
		y := eq.f(x)
		newPoint := point{x, y}
		return pointsLine{line.a, newPoint}
	}

	diffIntersectingLines := func(lineA, lineB pointsLine) float64 {
		intersectionPoint := solveLinearEquation(lineA.asLinearEquation(), lineB.asLinearEquation())

		var leftTriangleHeight, leftTriangleBase, leftTriangleArea float64
		leftTriangleHeight = math.Abs(lineA.a.Y - lineB.a.Y)
		leftTriangleBase = math.Abs(lineA.a.X - intersectionPoint.X)
		leftTriangleArea = (leftTriangleBase * leftTriangleHeight) / 2

		var rightTriangleHeight, rightTriangleBase, rightTriangleArea float64
		rightTriangleHeight = math.Abs(lineA.b.Y - lineB.b.Y)
		rightTriangleBase = math.Abs(intersectionPoint.X - lineA.b.X)
		rightTriangleArea = (rightTriangleBase * rightTriangleHeight) / 2

		return leftTriangleArea + rightTriangleArea
	}

	// they should have the exact same x-values for a and b
	// lineA.a.X == lineB.a.X && lineA.b.X == lineB.b.X
	// lineA.a and lineB.a shoulb be the left point
	diff := func(lineA pointsLine, lineB pointsLine) float64 {
		yDiffLeft := lineA.a.Y - lineB.a.Y
		yDiffRight := lineA.b.Y - lineB.b.Y
		if isPositive(yDiffLeft) != isPositive(yDiffRight) {
			// the lines are intersecting
			return diffIntersectingLines(lineA, lineB)
		}
		// lineA is positve, lineB is negative
		edgeAreaA := lineA.edgeArea()
		edgeAreaB := lineB.edgeArea()
		return math.Abs(edgeAreaA - edgeAreaB)
	}

	for indexA < countA && indexB < countB {
		var lineA, lineB pointsLine
		xOnAIsSmaller, rightPointA, rightPointB := next()
		if xOnAIsSmaller {
			// next is a point on A
			lineA = pointsLine{leftPointA, rightPointA}
			lineB = pointsLine{leftPointB, rightPointB}
			lineB = cuttLineSegment(rightPointA.X, lineB)
		} else {
			lineB = pointsLine{leftPointB, rightPointB}
			lineA = pointsLine{leftPointA, rightPointA}
			lineA = cuttLineSegment(rightPointB.X, lineA)
		}
		result += diff(lineA, lineB)
		leftPointA = lineA.b
		leftPointB = lineB.b
	}
	return result
}

func (line pointsLine) edgeArea() float64 {
	return (line.b.X - line.a.X) * (line.b.Y + line.a.Y) / 2
}

func (eq linearEquation) solveForY(y float64) point {
	x := (y - eq.m) / eq.k
	return point{x, y}
}
func isPositive(x float64) bool {
	if x < 0 {
		return false
	}
	return true
}
