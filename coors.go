package main

import (
	"fmt"
)

type rollingAvgCoor struct {
	avg           coor
	num           int
	hasAddedFirst bool
}

func newRollingAvg() *rollingAvgCoor {
	r := rollingAvgCoor{hasAddedFirst: false}
	return &r
}

func (r *rollingAvgCoor) add(c coor) {
	if r.hasAddedFirst {
		numFloat := float64(r.num)
		sumLat := (r.avg.Lat * numFloat) + c.Lat
		sumLong := (r.avg.Long * numFloat) + c.Long
		r.num++
		numFloat = float64(r.num)
		newLat := sumLat / numFloat
		newLong := sumLong / numFloat
		r.avg = coor{newLat, newLong}
	} else {
		r.avg = c
		r.num = 1
		r.hasAddedFirst = true
	}
}

// const goalNumPoints int = 100

func producePoints(run run, refCoor coor, goalNumPoints int) []point {
	if len(run.Coors) < 100 {
		return producePointsNoFilter(run, refCoor)
	}
	points := producePointsFilter(run, refCoor, goalNumPoints)
	//fmt.Println("Points")
	//fmt.Println(points)
	return points
}

// should propably make it so that it check that the resulting points are not too few, in which case it should
// decrease the required distane and do it again
func producePointsFilter(run run, refCoor coor, goalNumPoints int) []point {
	var result []point

	reqDistance := run.Distance / (float64(goalNumPoints) * 1.1)
	basePoint := run.Coors[0]
	count := len(run.Coors)

	var avgPoint = newRollingAvg()
	avgPoint.add(basePoint)

	for i := 1; i < count; i++ {
		coor := run.Coors[i]
		d := coorDistance(coor, basePoint)
		if d > reqDistance {
			point := pointForCoor(avgPoint.avg, refCoor)
			result = append(result, point)
			basePoint = coor
			avgPoint = newRollingAvg()
			avgPoint.add(basePoint)
		}
		avgPoint.add(coor)
	}

	// always add the last (unfinnished) avg point
	point := pointForCoor(avgPoint.avg, refCoor)
	result = append(result, point)
	return result
}

func producePointsNoFilter(run run, refCoor coor) []point {
	result := make([]point, len(run.Coors))
	for index, coor := range run.Coors {
		point := pointForCoor(coor, refCoor)
		result[index] = point
	}
	return result
}

func getMainCoorsFromRun(run run) runMainCoors {
	boxLoader := newBoxCoorsLoader()
	mainCoors := runMainCoors{}
	count := len(run.Coors)
	mainCoors.startCoor = run.Coors[0]
	mainCoors.endCoor = run.Coors[count-1]
	for _, c := range run.Coors {
		boxLoader.check(c)
	}
	mainCoors.boxCoors = boxLoader.boxCoors
	return mainCoors
}

// the maximum average distance that is allowed is 300
// each single point mus be within 400 meters
func mainCoorsSimilarity(mainCoorsA, mainCoorsB runMainCoors, max float64) (bool, float64) {
	//fmt.Printf("%+v\n", mainCoorsA)
	//fmt.Printf("%+v\n", mainCoorsB)
	// southDistance := coorDistance(mainCoorsA.boxCoors.southCoor, mainCoorsB.boxCoors.southCoor)
	// northDistance := coorDistance(mainCoorsA.boxCoors.northCoor, mainCoorsB.boxCoors.northCoor)
	// eastDistance := coorDistance(mainCoorsA.boxCoors.eastCoor, mainCoorsB.boxCoors.eastCoor)
	// westDistance := coorDistance(mainCoorsA.boxCoors.westCoor, mainCoorsB.boxCoors.westCoor)
	centerDistance := coorDistance(mainCoorsA.boxCoors.centerCoor(), mainCoorsB.boxCoors.centerCoor())
	// avgBoundingBoxDistance := (southDistance + northDistance + eastDistance + westDistance) / 4
	fmt.Println("centerDistance ", centerDistance)

	startDistance := coorDistance(mainCoorsA.startCoor, mainCoorsB.startCoor)
	endDistance := coorDistance(mainCoorsA.endCoor, mainCoorsB.endCoor)

	if (centerDistance > max) ||
		(startDistance > max) ||
		(endDistance > max) {
		fmt.Println("FoundOneDistanceAboveAllowedLimit")
		return false, max
	}
	return true, centerDistance

}

func pointForCoor(c coor, refCoor coor) point {
	xCoor := coor{refCoor.Lat, c.Long}
	yCoor := coor{c.Lat, refCoor.Long}

	xPointVal := coorDistance(refCoor, xCoor)
	yPointVal := coorDistance(refCoor, yCoor)

	if c.Long-refCoor.Long < 0 {
		xPointVal = -xPointVal
	}
	if c.Lat-refCoor.Lat < 0 {
		yPointVal = -yPointVal
	}

	return point{xPointVal, yPointVal}
}
