package main

func distanceForShape(shape []point) float64 {
	var daccum float64
	for i := 1; i < len(shape); i++ {
		lastPoint := shape[i-1]
		point := shape[i]
		daccum += distance(lastPoint, point)
	}
	return daccum
}
