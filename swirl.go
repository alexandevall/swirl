package main

import (
	"fmt"
	"math"
)

const swirlScoreTippingPoint = 50.0

/*
// if runs are compared
func swirlRuns(runA run, runB run) swirlDecision {
	// refCoor := runA.Coors[0]
	refCoor := getLowestCoorForRuns(runA, runB)
	shapeA := producePoints(runA, refCoor)
	shapeB := producePoints(runB, refCoor)
	return swirlDeep(shapeA, shapeB)
} */

// try to eliminate the possibility of similarity as fast as it can
/*
func swirlRunsFlat(runA, runB run) swirlDecision {


} */

// always do all calculations to get a propper score
func swirlRunsDeep(runA, runB run) swirlDecision {

	var scoreFilter1, scoreFilter2, scoreFilter3, scoreFilter4 float64

	const filter1Weight = 0.01
	const filter2Weight = 0.70
	const filter3Weight = 0.05
	const filter4Weight = 0.24

	fmt.Println("")
	fmt.Println("___Filter1___")
	// ***
	// FILTER 1: main coors
	//
	// begin by looking at the main coors of the runs
	mainCoorA := getMainCoorsFromRun(runA)
	mainCoorB := getMainCoorsFromRun(runB)
	maxMainCoorDistance := 400.0
	_, mainCoorsScore := mainCoorsSimilarity(mainCoorA, mainCoorB, maxMainCoorDistance)
	fmt.Println("MainCoorsScore ", mainCoorsScore)
	scoreFilter1 = suqshToTippingPoint(maxMainCoorDistance, mainCoorsScore)
	if scoreFilter1 >= swirlScoreTippingPoint {
		fmt.Println(">>>IsNotSimilar<<<")
		// they are not similar
		// I have to be sure that score increases with at least 50 (weighted)
		scoreFilter1 = exagerateScoreOverTippingPoint(scoreFilter1, filter1Weight)
	}

	// to begin the other steps
	// refCoor := getLowestCoorForRuns(runA, runB)
	refCoor := runA.Coors[0]

	fmt.Println("")
	fmt.Println("___Filter2___")
	// ***
	// FILTER 2: pointDiff
	//
	numChunksGoal := 500
	shapeA := producePoints(runA, refCoor, numChunksGoal)
	shapeB := producePoints(runB, refCoor, numChunksGoal)
	pointDiff := pointdiff(shapeA, shapeB)
	avgPointDiff := pointDiff.avg()

	// fmt.Printf("%+v\n", pointDiff)
	fmt.Println("Avg: ", avgPointDiff)

	scoreFilter2 = suqshToTippingPoint(50, avgPointDiff)
	if scoreFilter2 >= swirlScoreTippingPoint {
		// not similar
		fmt.Println(">>>IsNotSimilar<<<")
		scoreFilter2 = exagerateScoreOverTippingPoint(scoreFilter2, filter2Weight)
	}

	fmt.Println("")
	fmt.Println("___Filter3___")
	// ***
	// FILTER 3: detours
	shapeDistanceArrA := distanceArrayForShape(shapeA)
	diffCurveA := getPointsFromArrays(shapeDistanceArrA, pointDiff.minArrayA)
	shapeDistanceArrB := distanceArrayForShape(shapeB)
	diffCurveB := getPointsFromArrays(shapeDistanceArrB, pointDiff.minArrayB)

	detourFinder := newPointDiffDetourFinder()
	detourAreaA := detourFinder.analyse(diffCurveA)
	detourAreaB := detourFinder.analyse(diffCurveB)
	maxDetourArea := math.Max(detourAreaA, detourAreaB)

	fmt.Println("DetourAreaA")
	fmt.Println(detourAreaA)
	fmt.Println("DetourAreaB")
	fmt.Println(detourAreaB)
	fmt.Println("MaxDetourArea")
	fmt.Println(maxDetourArea)

	scoreFilter3 = suqshToTippingPoint(800, maxDetourArea)
	if detourAreaA > swirlScoreTippingPoint {
		scoreFilter3 = exagerateScoreOverTippingPoint(scoreFilter3, filter3Weight)
	}

	fmt.Println("")
	fmt.Println("___Filter4___")
	// ***
	// FILTER 4: clockwise area

	filter4Value := filterCWA(runA, runB)
	scoreFilter4 = suqshToTippingPoint(120, filter4Value)
	if scoreFilter4 >= swirlScoreTippingPoint {
		// it is not similar
		fmt.Println(">>>IsNotSimilar<<<")
		scoreFilter4 = exagerateScoreOverTippingPoint(scoreFilter4, filter4Weight)
	}

	fmt.Println("")
	fmt.Println("___Final___")
	// if it can reach this point with a score lower than lets say 50, it is similar enough
	score := (scoreFilter1 * filter1Weight) + (scoreFilter2 * filter2Weight) + (scoreFilter3 * filter3Weight) + (scoreFilter4 + filter4Weight)
	fmt.Println("Score ", score)
	decision := swirlDecision{}
	if score < swirlScoreTippingPoint {
		decision = swirlDecision{Score: int(score), IsSimilar: true}
	} else {
		decision = swirlDecision{Score: int(score), IsSimilar: false}
	}

	fmt.Printf("%+v\n", decision)
	return decision
}

// what i like to test atm
func swirly(shapeA, shapeB []point) swirlDecision {
	leftRightAreaA := leftRightArea(shapeA)
	leftRightAreaB := leftRightArea(shapeB)
	fmt.Println("NO FOLD")
	fmt.Println("AREA-A")
	leftRightAreaA.print()
	fmt.Println("")
	fmt.Println("AREA-B")
	leftRightAreaB.print()
	fmt.Println("")

	fmt.Println("WITH FOLD")
	foldedA, foldedB := fold(shapeA, shapeB)
	leftRightAreaFoldA := leftRightArea(foldedA)
	leftRightAreaFoldB := leftRightArea(foldedB)
	fmt.Println("AREA-A")
	leftRightAreaFoldA.print()
	fmt.Println("")
	fmt.Println("AREA-B")
	leftRightAreaFoldB.print()
	fmt.Println("")

	curveDiff := clockwiseCurveDiff(areaPointsFromPoints(shapeA), areaPointsFromPoints(shapeB))
	fmt.Println("curveDiff:")
	fmt.Printf("%+v\n", curveDiff)
	fmt.Println("Score")
	fmt.Println(curveDiff.score())
	fmt.Println("")

	diffPoint := areaPointDiffForPoints(shapeA, shapeB)
	fmt.Println("-- RESULT --")
	fmt.Printf("%+v\n", diffPoint)
	avgDiffPoint := diffPoint.avg()
	areaDiff := scoreDiffArea(leftRightAreaA.clockwiseArea(), leftRightAreaB.clockwiseArea())

	/*

		if avgDiffPoint > 50 || areaDiff > 1.3 {
			// fucked
			score := normalizeScore(avgDiffPoint, areaDiff, 30.0) // 30 so that even perfect area are unsimilar
			return swirlDecision{Score: int(score), IsSimilar: false}
		}

		// _____
		// can be similar

		// have to see its not a no-clockwisearea type
		distanceA := distanceForShape(shapeA)
		distanceB := distanceForShape(shapeB)


		if isOneNoAreaShape(clockwiseAreaA, distanceA, clockwiseAreaB, distanceB) {
			diffAreaPoint := areaPointDiffForPoints(shapeA, shapeB)
			avgDiffAreaPoint := diffAreaPoint.avg()

			if avgDiffAreaPoint > 100 {
				// if this is not good enough, return
				score := normalizeScore(avgDiffPoint, areaDiff, avgDiffAreaPoint)
				return swirlDecision{Score: int(score), IsSimilar: false}
			}
		}
	*/

	score := normalizeScore(avgDiffPoint, areaDiff, 0.0)
	isSimilar := false
	return swirlDecision{Score: int(score), IsSimilar: isSimilar}
}

// the key entry point to the algo
/*
func swirlShapes(shapeA []point, shapeB []point) swirlDecision {
	diff := pointdiff(shapeA, shapeB)
	fmt.Println("")
	fmt.Println("-- RESULT --")
	fmt.Printf("%+v\n", diff)

	avgPointDiff := diff.avg()
	areaDiff := scoreDiffArea(diff.ClockwiseAreaA, diff.ClockwiseAreaB)
	fmt.Println("AvgPointDiff")
	fmt.Println(avgPointDiff)
	fmt.Println("AreaDiff")
	fmt.Println(areaDiff)
	var score float64
	var isSimilar bool

	// a score below 80 is simliar (50 + 100 * (1.3-1))
	if avgPointDiff > 50 || areaDiff > 1.3 {
		score = normalizeScore(avgPointDiff, areaDiff, 30.0) // 30 so that even perfect area are unsimilar
		isSimilar = false
	} else {
		score = normalizeScore(avgPointDiff, areaDiff, 0.0)
		isSimilar = true
	}

	decision := swirlDecision{Score: int(score), IsSimilar: isSimilar}
	fmt.Println("DECISION:")
	fmt.Printf("%+v\n", decision)
	return decision
} */

func isPointDifferenceSmall(diff diffResult) bool {
	avgDiff := (diff.AvgClosestPointDiffA + diff.AvgClosestPointDiffB) / 2
	if avgDiff < 50 {
		return true
	}
	return false
}

func scoreDiffArea(areaA float64, areaB float64) float64 {
	var smallestAbs float64
	if math.Abs(areaA) < math.Abs(areaB) {
		smallestAbs = math.Abs(areaA)
	} else {
		smallestAbs = math.Abs(areaB)
	}
	areaDifference := math.Abs(areaA - areaB)
	fraction := (areaDifference / smallestAbs) + 1
	return fraction
}

func normalizeScore(pointDiff float64, areaDiff float64, constant float64) float64 {
	return pointDiff + ((areaDiff - 1) * 100) + constant
}

// lets say that i want to adda value of 40 to score, because it
// is less than som threashold (say 50).
// lets also say that i do not allow this value to be (in 'score') lower than 30.
// Now i need to convert the value of 40 to some value below 30.
// (it the value has been 50, i would want to return 30)
func squshScoreInMax(score, scoreIsBelow, max float64) float64 {
	factor := max / scoreIsBelow
	return score * factor
}

type squshScore struct {
	tippingPoint                     float64
	tippingPointCorrespondingToValue float64
	weight                           float64
}

func newScoreSqusher50(tippingPointCorrespondingToValue, weight float64) squshScore {
	return squshScore{50, tippingPointCorrespondingToValue, weight}
}

func suqshToTippingPoint(tippingPointCorrespondingToValue, val float64) float64 {
	factor := swirlScoreTippingPoint / tippingPointCorrespondingToValue
	return factor * val
}

func (s squshScore) squshAndWeight(x float64) float64 {
	return s.sqush(x) * s.weight
}

func (s squshScore) sqush(x float64) float64 {
	factor := s.tippingPoint / s.tippingPointCorrespondingToValue
	return factor * x
}

func (s squshScore) squshAndWeightIsUnsimilar(x float64) float64 {
	weightedScore := s.squshAndWeight(x)
	if weightedScore < s.tippingPoint {
		return s.tippingPoint
	}
	return weightedScore
}

func exagerateScoreOverTippingPoint(score, weight float64) float64 {
	add := swirlScoreTippingPoint / weight
	return add + score
}
