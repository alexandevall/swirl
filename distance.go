package main

import (
	"math"
)

// from https://gist.github.com/cdipaolo/d3f8db3848278b49db68

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func coorDistance(coor1 coor, coor2 coor) float64 {

	// changes this if i want to keep 64
	coor64v1 := coor1.to64()
	coor64v2 := coor2.to64()

	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = coor64v1.Lat * math.Pi / 180
	lo1 = coor64v1.Long * math.Pi / 180
	la2 = coor64v2.Lat * math.Pi / 180
	lo2 = coor64v2.Long * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)
	result := 2 * r * math.Asin(math.Sqrt(h))
	return float64(result)
}
