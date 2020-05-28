package main

import (
	"fmt"
	"math"
)

/*
func isOneNoAreaShape(areaA, distanceA, areaB, distanceB float64) bool {
	if isNoAreaShape(areaA, distanceA) || isNoAreaShape(areaB, distanceB) {
		return true
	}
	return false
}

func isNoAreaShape(area float64, distance float64) bool {
	// imagaine a and b are the two side of a square with a cirfumfrance of
	// 'distance' and an area of 'area'
	b := (distance / 4) + math.Sqrt(math.Pow(distance/4, 2)-area)
	a := area / b
	smallest := a
	largest := b
	if b < a {
		smallest = b
		largest = a
	}
	smallestToLargest := smallest / largest

	// if the largest was 100
	smallSide := smallestToLargest * 100

	if smallest < 3 {
		return true
	}
	return false
}
*/

func smallSideOfRectangle(area float64, distance float64) float64 {
	b := (distance / 4) + math.Sqrt(math.Pow(distance/4, 2)-area)
	a := area / b

	// this is quite strange to me, but apperently the area can be bigger relative to the distance,
	// so that the area-distance combination cannot be desctibed as a rectangle
	// (the distance is greater than sqrt(area))
	// in this case lets for the moment just pretend that it is an square
	if math.IsNaN(b) {
		return 100
	}
	smallest := a
	largest := b
	if b < a {
		smallest = b
		largest = a
	}
	smallestToLargest := smallest / largest

	// if the largest was 100
	smallSide := smallestToLargest * 100
	return smallSide
}

func clockwiseArea(shape []point) float64 {
	count := len(shape)
	var accum float64
	for i := 0; i < count-1; i++ {
		point := shape[i]
		nextPoint := shape[i+1]
		s := sumEdge(point, nextPoint)
		accum += s
	}

	// should i really use the last one
	first := shape[0]
	last := shape[count-1]
	finalEdge := sumEdge(last, first)
	accum += finalEdge

	// i want the area some nice number for lon (about 10 km) runs5
	accum = accum / 10000

	return accum
}

type leftRightAreaResult struct {
	leftArea           float64
	rightArea          float64
	finalClockwiseArea float64 // can be negative!
	divVal             float64
}

func newLeftRightAreaResult() leftRightAreaResult {
	return leftRightAreaResult{0, 0, 0, 1}
}

func (a leftRightAreaResult) unfiltered() *leftRightAreaResult {
	inversedDivVal := 1 / a.divVal
	a.divAll(inversedDivVal)
	return &a
}

func (a leftRightAreaResult) clockwiseArea() float64 {
	return a.rightArea - a.leftArea
}

func (a leftRightAreaResult) leftRightWithFinalPart() (leftArea, rightArea float64) {
	leftArea = a.leftArea
	rightArea = a.rightArea
	if a.finalClockwiseArea < 0 {
		leftArea += -a.finalClockwiseArea
	} else {
		rightArea += a.finalClockwiseArea
	}
	return
}

func (a leftRightAreaResult) clockwiseAreaWithFinalPart() float64 {
	leftArea, rightArea := a.leftRightWithFinalPart()
	return rightArea - leftArea
}

func (a leftRightAreaResult) aggrArea() float64 {
	return a.rightArea + a.leftArea
}

func (a leftRightAreaResult) aggrAreaWithFinalPart() float64 {
	return math.Abs(a.finalClockwiseArea) + a.leftArea + a.rightArea
}

func (a *leftRightAreaResult) divAll(div float64) {
	a.leftArea = a.leftArea / div
	a.rightArea = a.rightArea / div
	a.finalClockwiseArea = a.finalClockwiseArea / div
	a.divVal = div
}

func (a *leftRightAreaResult) print() {
	fmt.Println("RightArea")
	fmt.Println(a.rightArea)
	fmt.Println("LeftArea")
	fmt.Println(a.leftArea)
	fmt.Println("ClockwiseArea")
	fmt.Println(a.clockwiseArea())
	fmt.Println("AggrArea")
	fmt.Println(a.aggrArea())

	/*
		fmt.Println("withFinalMart ->")
		fmt.Println("FinalclockwiseArea")
		fmt.Println(a.finalClockwiseArea)
		leftArea, rightArea := a.leftRightWithFinalPart()
		fmt.Println("RightArea")
		fmt.Println(rightArea)
		fmt.Println("LeftArea")
		fmt.Println(leftArea)
		fmt.Println("ClockwiseArea")
		fmt.Println(a.clockwiseAreaWithFinalPart())
		fmt.Println("AggrArea")
		fmt.Println(a.aggrAreaWithFinalPart()) */
}

func leftRightArea(shape []point) leftRightAreaResult {
	result := newLeftRightAreaResult()
	count := len(shape)
	for i := 0; i < count-1; i++ {
		point := shape[i]
		nextPoint := shape[i+1]
		s := sumEdge(point, nextPoint)
		if s > 0 {
			// it is right
			result.rightArea += s
		} else {
			// it is left (i make the value postive)
			result.leftArea += -s
		}
	}

	// should i really use the last one
	first := shape[0]
	last := shape[count-1]
	finalEdge := sumEdge(last, first)
	result.finalClockwiseArea = finalEdge

	// to make it readable (the two because the area is div 2)
	result.divAll(10000 * 2)

	return result
}

func sumEdge(a point, b point) float64 {
	// return (b.X - a.X) * math.Abs(b.Y+a.Y)
	return (b.X - a.X) * (b.Y + a.Y)
}
