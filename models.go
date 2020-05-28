package main

import (
	"math"
)

type coor struct {
	Lat  float64
	Long float64
}

type coor64 struct {
	Lat  float64
	Long float64
}

type point struct {
	X float64
	Y float64
}

type pointsLine struct {
	a point
	b point
}

type diffResult struct {
	AvgClosestPointDiffA float64
	AvgClosestPointDiffB float64
	minArrayA            []float64
	minArrayB            []float64
}

func (s *diffResult) avg() float64 {
	return (s.AvgClosestPointDiffA + s.AvgClosestPointDiffB) / 2
}

func (s *diffResult) worst() float64 {
	return math.Max(s.AvgClosestPointDiffA, s.AvgClosestPointDiffB)
}

type swirlDecision struct {
	Score     int
	IsSimilar bool
}

type user struct {
	Name     string
	Password string
}

type run struct {
	Coors    []coor
	Distance float64
}

type runMainCoors struct {
	boxCoors  runBoxCoors
	startCoor coor
	endCoor   coor
}

type runBoxCoors struct {
	southCoor coor
	northCoor coor
	eastCoor  coor
	westCoor  coor
}

func (box runBoxCoors) centerCoor() coor {
	lat := (box.southCoor.Lat + box.southCoor.Lat) / 2
	long := (box.westCoor.Long + box.eastCoor.Long) / 2
	return coor{Lat: lat, Long: long}
}

type runBoxCoorsLoader struct {
	boxCoors runBoxCoors
	hasAdded bool
}

func newBoxCoorsLoader() runBoxCoorsLoader {
	new := runBoxCoorsLoader{hasAdded: false}
	return new
}

func (box *runBoxCoorsLoader) check(c coor) {
	if box.hasAdded {
		if c.Lat < box.boxCoors.southCoor.Lat {
			box.boxCoors.southCoor = c
		}
		if c.Lat > box.boxCoors.northCoor.Lat {
			box.boxCoors.northCoor = c
		}
		if c.Long < box.boxCoors.westCoor.Long {
			box.boxCoors.westCoor = c
		}
		if c.Long > box.boxCoors.eastCoor.Long {
			box.boxCoors.eastCoor = c
		}
	} else {
		box.boxCoors.southCoor = c
		box.boxCoors.northCoor = c
		box.boxCoors.westCoor = c
		box.boxCoors.eastCoor = c
		box.hasAdded = true
	}
}

type runWrapper struct {
	run       run
	mainCoors runMainCoors
}

type runInput struct {
	RunA run
	RunB run
}

type request struct {
	User     user
	Option   string
	RunInput runInput
}

func (c coor) to64() coor64 {
	return coor64{float64(c.Lat), float64(c.Long)}
}

func (c coor) asPoint() point {
	return point{c.Long, c.Lat}
}
