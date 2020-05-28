package main

import (
	"math"
	"sort"
)

func distanceArrayForShape(shape []point) []float64 {
	count := len(shape)
	result := make([]float64, count)
	var accumDistance float64 = 0
	result[0] = accumDistance
	for i := 1; i < count; i++ {
		lastPoint := shape[i-1]
		point := shape[i]
		d := distance(lastPoint, point)
		accumDistance += d
		result[i] = accumDistance
	}
	return result
}

func getPointsFromArrays(arrX, arrY []float64) []point {
	result := make([]point, len(arrX))
	for i := range arrX {
		x := arrX[i]
		y := arrY[i]
		newPoint := point{x, y}
		result[i] = newPoint
	}
	return result
}

type pointDiffDetourFinder struct {
	baseWeight         float64
	detoursStartAtDiff float64
	scoreLevels        []pointDiffLevel
}

func newPointDiffDetourFinder() *pointDiffDetourFinder {
	obj := pointDiffDetourFinder{baseWeight: 1, detoursStartAtDiff: 100}
	obj.setScoreLevels()
	return &obj
}

func (s pointDiffDetourFinder) analyse(curve []point) float64 {
	detours := []float64{}
	currentArea := newMaybeFloat64()
	collectPoints := make([]point, len(curve))
	for i := 0; i < len(curve)-1; i++ {
		p := curve[i]
		nextP := curve[i+1]
		line := pointsLine{p, nextP}
		area := s.scoreLine(line)
		if area <= 0 {
			if !currentArea.isNil() {
				// end of an area section
				detours = append(detours, currentArea.value)
				currentArea.setNil()
			}
		} else {
			if currentArea.isNil() {
				// begging of an area section
				currentArea.set(area)
			} else {
				// continue of an area section
				currentArea.value += area
			}
		}
		newCollectPoint := point{p.X, currentArea.value}
		collectPoints[i] = newCollectPoint
	}
	if len(detours) < 0 {
		return 0
	}
	sort.Float64s(detours)
	// printJSON(collectPoints)
	return detours[len(detours)-1]
}

func (s pointDiffDetourFinder) scoreLine(line pointsLine) float64 {
	diffVal := line.a.Y
	if s.detoursStartAtDiff > line.a.Y {
		return 0
	}
	deltaX := (line.b.X - line.a.X)
	defaultArea := deltaX * s.detoursStartAtDiff
	level := s.getLevelForDiff(diffVal)
	area := line.edgeArea() - defaultArea
	area = math.Sqrt(area)
	return area * level.weight
}

func (s *pointDiffDetourFinder) getLevelForDiff(val float64) pointDiffLevel {
	if val < s.scoreLevels[0].limit {
		return s.scoreLevels[0]
	}

	for _, level := range s.scoreLevels {
		if val > level.limit {
			return level
		}
	}

	return s.scoreLevels[0]
}

func (s *pointDiffDetourFinder) setScoreLevels() {
	L1 := pointDiffLevel{100, s.baseWeight * 1.00}
	L2 := pointDiffLevel{125, s.baseWeight * 1.00}
	L3 := pointDiffLevel{150, s.baseWeight * 1.25}
	L4 := pointDiffLevel{175, s.baseWeight * 1.25}
	L5 := pointDiffLevel{200, s.baseWeight * 1.50}
	L6 := pointDiffLevel{250, s.baseWeight * 1.50}
	L7 := pointDiffLevel{300, s.baseWeight * 1.75}
	L8 := pointDiffLevel{350, s.baseWeight * 1.75}
	L9 := pointDiffLevel{400, s.baseWeight * 2.00}
	L10 := pointDiffLevel{500, s.baseWeight * 3.00}
	L11 := pointDiffLevel{600, s.baseWeight * 4.00}
	L12 := pointDiffLevel{700, s.baseWeight * 5.00}
	L13 := pointDiffLevel{800, s.baseWeight * 6.00}
	L14 := pointDiffLevel{900, s.baseWeight * 7.00}
	L15 := pointDiffLevel{1000, s.baseWeight * 8.00}
	s.scoreLevels = []pointDiffLevel{L1, L2, L3, L4, L5, L6, L7, L8, L9, L10, L11, L12, L13, L14, L15}
}

type pointDiffLevel struct {
	limit  float64
	weight float64
}

type pointDiffDetour struct {
	area float64
}

type maybeFloat64 struct {
	hasValue bool
	value    float64
}

func newMaybeFloat64() *maybeFloat64 {
	return &maybeFloat64{false, 0}
}

func (m *maybeFloat64) get() (bool, float64) {
	if m.hasValue {
		return true, m.value
	}
	return false, m.value
}

func (m *maybeFloat64) set(value float64) {
	m.value = value
	m.hasValue = true
}

func (m *maybeFloat64) setNil() {
	m.value = 0
	m.hasValue = false
}

func (m *maybeFloat64) isNil() bool {
	return m.hasValue
}
