package main

import (
	"math"
)

func foldShapeOnShape(shape []point, onShape []point) []point {
	reslut := make([]point, len(shape))
	for index, p := range shape {
		// finding the closest point for the shape
		bestPoint := point{}
		minDistance := math.MaxFloat64
		for _, onShapePoint := range onShape {
			d := distance(p, onShapePoint)
			if d < minDistance {
				bestPoint = onShapePoint
				minDistance = d
			}
		}
		reslut[index] = bestPoint
	}
	return reslut
}

func fold(shapeA, shapeB []point) ([]point, []point) {
	foldedA := foldShapeOnShape(shapeA, shapeB)
	foldedB := foldShapeOnShape(shapeB, foldedA)
	return foldedA, foldedB
}
