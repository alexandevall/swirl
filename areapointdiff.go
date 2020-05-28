package main

import "math"

type areaPoint struct {
	point     point
	leftArea  float64
	rightArea float64
}

func (ap *areaPoint) clockwiseArea() float64 {
	return (ap.rightArea - ap.leftArea)
}

func (ap areaPoint) aggrArea() float64 {
	return ap.rightArea + ap.leftArea
}

type areaPointDiffResult struct {
	avgDiffA float64
	avgDiffB float64
}

type areaPointWorkerInput struct {
	pointA areaPoint
	shapeB []areaPoint
}

func areaPointsFromPoints(shape []point) []areaPoint {
	var leftArea, rightArea float64
	count := len(shape)
	result := make([]areaPoint, count)
	result[0] = areaPoint{point: shape[0], leftArea: 0, rightArea: 0}
	for i := 1; i < count; i++ {
		lastPoint := shape[i-1]
		point := shape[i]
		clockwiseArea := sumEdge(lastPoint, point)
		if clockwiseArea < 0 {
			leftArea += -(clockwiseArea)
		} else {
			rightArea += clockwiseArea
		}
		newAreaPoint := areaPoint{point: point, leftArea: leftArea, rightArea: rightArea}
		result[i] = newAreaPoint
	}
	return result
}

func areaPointDiffForPoints(shapeA, shapeB []point) diffResult {
	return areaPointDiff(areaPointsFromPoints(shapeB), areaPointsFromPoints(shapeA))
}

func areaPointDiff(shapeA []areaPoint, shapeB []areaPoint) diffResult {
	countA := len(shapeA)
	countB := len(shapeB)

	minArrayB := make([]float64, countB)
	for i := range minArrayB {
		minArrayB[i] = math.MaxFloat64
	}

	var totalDiffA float64

	jobs := make(chan areaPointWorkerInput, countA)
	minDistanceA := make(chan float64, countA)
	totalLoops := countA * countB
	distanceB := make(chan minDistanceResult, totalLoops)

	go areaPointWorker(jobs, minDistanceA, distanceB)
	go areaPointWorker(jobs, minDistanceA, distanceB)
	go areaPointWorker(jobs, minDistanceA, distanceB)
	go areaPointWorker(jobs, minDistanceA, distanceB)
	go areaPointWorker(jobs, minDistanceA, distanceB)
	go areaPointWorker(jobs, minDistanceA, distanceB)
	go areaPointWorker(jobs, minDistanceA, distanceB)
	go areaPointWorker(jobs, minDistanceA, distanceB)
	go areaPointWorker(jobs, minDistanceA, distanceB)

	for i := 0; i < countA; i++ {
		pointA := shapeA[i]
		jobs <- areaPointWorkerInput{pointA, shapeB}
	}

	for i := 0; i < countA+totalLoops; i++ {
		select {
		case m := <-minDistanceA:
			totalDiffA += m
		case d := <-distanceB:
			minForShape(minArrayB, d.d, d.index)
		}
	}

	avgDiffA := totalDiffA / float64(countA)
	totalDiffB := sum(minArrayB)
	avgDiffB := totalDiffB / float64(countB)

	return diffResult{avgDiffA, avgDiffB, make([]float64, 0), make([]float64, 0)}
}

func areaPointWorker(jobs <-chan areaPointWorkerInput, results chan<- float64, distanceB chan<- minDistanceResult) {
	for n := range jobs {
		results <- minAreaPointDistance(n.pointA, n.shapeB, distanceB)
	}
}

func minAreaPointDistance(pointA areaPoint, shapeB []areaPoint, distanceB chan<- minDistanceResult) float64 {
	count := len(shapeB)
	min := math.MaxFloat64
	for i := 0; i < count; i++ {
		pointB := shapeB[i]
		d := areaPointDistance3D(pointA, pointB)
		if d < min {
			min = d
		}
		distanceB <- minDistanceResult{d, i}
	}
	return min
}

func clockwiseAreaScore(a areaPoint) float64 {
	// return a.clockwiseArea()

	clockwiseArea := a.clockwiseArea()
	sign := 1.0
	if clockwiseArea < 0 {
		sign = -1.0
		clockwiseArea = clockwiseArea * sign
	}
	adjusted := (math.Sqrt(clockwiseArea/2) * sign) * 10
	// adjusted := clockwiseArea / 2 * sign
	return adjusted

}

func aggrLeftRightAreaScore(a areaPoint) float64 {
	aggrArea := a.aggrArea()
	sign := 1.0
	if aggrArea < 0 {
		sign = -1.0
		aggrArea = aggrArea * sign
	}
	adjusted := math.Sqrt(aggrArea/2) * sign
	// adjusted := aggrArea / 2 * sign
	return adjusted
}

func areaPointDistance3D(a areaPoint, b areaPoint) float64 {
	clockwiseAreaScoreA := clockwiseAreaScore(a)
	clockwiseAreeScoreB := clockwiseAreaScore(b)
	pow2 := math.Pow(a.point.X-b.point.X, 2) + math.Pow(a.point.Y-a.point.Y, 2) + math.Pow(clockwiseAreaScoreA-clockwiseAreeScoreB, 2)
	return math.Sqrt(pow2)
}

func areaPointDistance4D(a areaPoint, b areaPoint) float64 {
	clockwiseAreaScoreA := clockwiseAreaScore(a)
	clockwiseAreeScoreB := clockwiseAreaScore(b)
	aggrAreaA := aggrLeftRightAreaScore(a)
	aggrAreaB := aggrLeftRightAreaScore(b)
	pow2 := math.Pow(a.point.X-b.point.X, 2) + math.Pow(a.point.Y-a.point.Y, 2) + math.Pow(clockwiseAreaScoreA-clockwiseAreeScoreB, 2) + math.Pow(aggrAreaA-aggrAreaB, 2)
	return math.Sqrt(pow2)
}
