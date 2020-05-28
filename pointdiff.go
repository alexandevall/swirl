package main

import "math"

type counter struct {
	c int
}

func (c *counter) add() {
	c.c++
}

func pointdiff(shapeA []point, shapeB []point) diffResult {
	countA := len(shapeA)
	countB := len(shapeB)

	minArrayB := make([]float64, countB)
	for i := range minArrayB {
		minArrayB[i] = math.MaxFloat64
	}
	minArrayA := make([]float64, countA)

	var totalDiffA float64
	var totalCount counter

	// A tricky thing here is that all points on one shape, should be compared with all
	// lines (edges) of the other. As there are fewer points than lines, if i loop through all points
	// on A, there will be one time when i shouldnt compare a line on A with a point on B; conversly,
	// if i loop through points on B (which is done in 'minDistance' function) there will be one too
	// many lines. This is why there are two special cases where the first point on A is compared with
	// B or the first point in B is compared with A.

	firstPointA := shapeA[0]
	firstMinDistanceA := minDistanceNoShape(firstPointA, shapeB, totalCount)
	totalDiffA += firstMinDistanceA

	jobs := make(chan loopShapeBInput, countA-1)
	minDistanceA := make(chan minDistanceResult, countA-1)
	totalLoops := (countA - 1) * (countB - 1)
	distanceB := make(chan minDistanceResult, totalLoops)

	/*
		jobs := make(chan loopShapeBInput, countA )
		minDistanceA := make(chan float64)
		// totalLoops := (countA - 1) * (countB - 1)
		distanceB := make(chan distanceForBPoint)
	*/
	go shapeWorkerB(jobs, minDistanceA, distanceB)
	go shapeWorkerB(jobs, minDistanceA, distanceB)
	go shapeWorkerB(jobs, minDistanceA, distanceB)
	go shapeWorkerB(jobs, minDistanceA, distanceB)
	go shapeWorkerB(jobs, minDistanceA, distanceB)
	go shapeWorkerB(jobs, minDistanceA, distanceB)
	go shapeWorkerB(jobs, minDistanceA, distanceB)
	go shapeWorkerB(jobs, minDistanceA, distanceB)
	go shapeWorkerB(jobs, minDistanceA, distanceB)
	go shapeWorkerB(jobs, minDistanceA, distanceB)

	// loop through all lines on A
	for i := 1; i < countA; i++ {
		lastPointA := shapeA[i-1]
		pointA := shapeA[i]
		lineA := pointsLine{lastPointA, pointA}
		// go minDistance(pointA, lineA, shapeB, minArrayB, &totalCount)
		jobs <- loopShapeBInput{pointA, lineA, i, shapeB}
		// totalDiffA += min
	}

	for i := 0; i < (countA-1)+totalLoops; i++ {
		select {
		case m := <-minDistanceA:
			totalDiffA += m.d
			minArrayA[m.index] = m.d
		case d := <-distanceB:
			minForShape(minArrayB, d.d, d.index)
		}
	}

	/*
		// all the distances for B
		for i := 0; i < totalLoops; i++ {
			r := <-distanceB
			minForShape(minArrayB, r.d, r.index)
		}

		// get the results
		for i := 1; i < countA; i++ {
			r := <-minDistanceA
			//fmt.Println(r)
			totalDiffA += r
		}
	*/

	avgDiffA := totalDiffA / float64(countA)

	totalDiffB := sum(minArrayB)
	avgDiffB := totalDiffB / float64(countB)

	return diffResult{avgDiffA, avgDiffB, minArrayA, minArrayB}
}

type loopShapeBInput struct {
	pointA point
	lineA  pointsLine
	indexA int
	shapeB []point
}

type minDistanceResult struct {
	d     float64
	index int
}

func shapeWorkerB(jobs <-chan loopShapeBInput, minDistanceA chan<- minDistanceResult, distanceB chan<- minDistanceResult) {
	for n := range jobs {
		minD := minDistance(n.pointA, n.lineA, n.shapeB, distanceB)
		minDistanceA <- minDistanceResult{minD, n.indexA}
	}
}

func sum(array []float64) float64 {
	var result float64
	for _, val := range array {
		result += val
	}
	return result
}

func minDistance(pointA point, lineA pointsLine, shapeB []point, distanceB chan<- minDistanceResult) float64 {
	// special case - first point on B
	firstPointOnShapeB := shapeB[0]
	firstPointOnShapeBToLineA := lineSegmentToPointDistance(firstPointOnShapeB, lineA)
	distanceB <- minDistanceResult{firstPointOnShapeBToLineA, 0}

	// loop through all lines in B to find the smallest distance to pointA
	count := len(shapeB)
	min := math.MaxFloat64
	for i := 1; i < count; i++ {
		pointOnShapeB := shapeB[i]
		lastPointOnShapeB := shapeB[i-1]
		lineB := pointsLine{lastPointOnShapeB, pointOnShapeB}
		// d := distance(pointA, pointOnShapeB)
		d := lineSegmentToPointDistance(pointA, lineB)
		if d < min {
			min = d
		}

		// the distance from pointB to the line on A
		d = lineSegmentToPointDistance(pointA, lineB)
		distanceB <- minDistanceResult{d, i}
		//fmt.Println(totalCount)
	}
	return min
}

// as i will lopp throuhg all points in shapeA
// i will have more points then lines
// the function above 'minDistance' assumes that there is a line on shapeA
// that points in shapeB can compare with
// one time this will not be the case
// and then this is used
func minDistanceNoShape(pointA point, shapeB []point, totalCount counter) float64 {
	// loop through all lines
	count := len(shapeB)
	min := math.MaxFloat64
	for i := 1; i < count; i++ {
		pointOnShapeB := shapeB[i]
		lastPointOnShapeB := shapeB[i-1]
		lineB := pointsLine{lastPointOnShapeB, pointOnShapeB}
		// d := distance(pointA, pointOnShapeB)
		d := lineSegmentToPointDistance(pointA, lineB)
		// fmt.Println("")
		// fmt.Println(pointA)
		// fmt.Println(pointOnShapeB)
		// fmt.Println(d)

		if d < min {
			min = d
		}
	}
	return min
}

func minForShape(minArray []float64, distance float64, index int) {
	if minArray[index] > distance {
		minArray[index] = distance
	}
}

func distance(a point, b point) float64 {
	a2 := (a.X - b.X) * (a.X - b.X)
	b2 := (a.Y - b.Y) * (a.Y - b.Y)
	answer := math.Sqrt(float64(a2) + float64(b2))
	return float64(answer)
}

type linearEquation struct {
	k float64
	m float64
}

func easyPointToLineDistance(p point, line pointsLine) float64 {
	return distance(p, line.b)
}

func lineSegmentToPointDistance(p point, line pointsLine) float64 {
	closestPointOnLine := closestPointToLine(p, line)
	if isPointInsideSpan(closestPointOnLine, line) {
		return distance(closestPointOnLine, p)
	}
	aDistance := distance(p, line.a)
	bDistance := distance(p, line.b)
	if aDistance < bDistance {
		return aDistance
	}
	return bDistance
}

func isPointInsideSpan(p point, span pointsLine) bool {
	var xMax, xMin, yMax, yMin float64
	if span.a.X > span.b.X {
		xMax = span.a.X
		xMin = span.b.X
	} else {
		xMax = span.b.X
		xMin = span.a.X
	}
	if span.a.Y > span.b.Y {
		yMax = span.a.Y
		yMin = span.b.Y
	} else {
		yMax = span.b.Y
		yMin = span.a.Y
	}
	if (p.X < xMax) && (p.X > xMin) && (p.Y < yMax) && (p.Y > yMin) {
		return true
	}
	return false
}

func closestPointToLine(p point, line pointsLine) point {
	var k1, k2, m1, m2 float64
	// representing the line
	k1 = (line.b.Y - line.a.Y) / (line.b.X - line.a.X)
	m1 = line.a.Y - (line.a.X * k1)

	// the line but twisted 180 degrees
	k2 = (line.a.X - line.b.X) / (line.b.Y - line.a.Y)
	m2 = p.Y - (p.Y * k2)

	eq1 := linearEquation{k1, m1}
	eq2 := linearEquation{k2, m2}

	return solveLinearEquation(eq1, eq2)
}

func solveLinearEquation(a linearEquation, b linearEquation) point {
	var x, y float64
	x = (b.m - a.m) / (a.k - b.k)
	y = (a.k * x) + a.m
	return point{x, y}
}

func (line pointsLine) asLinearEquation() linearEquation {
	var k, m float64
	k = (line.b.Y - line.a.Y) / (line.b.X - line.a.X)
	m = line.a.Y - (line.a.X * k)
	return linearEquation{k, m}
}

func (eq linearEquation) f(x float64) float64 {
	return (eq.k * x) + eq.m
}
