package main

import "math"

type transformMatrix struct {
	iHat point
	jHat point
}

// this is nice because A -> B or A->B->A can occur mostly on an axsis, which will yeild very little
// area. By tilting the shape a bit, this problem can be overcome
func tilt(shapeA []point, shapeB []point) ([]point, []point) {
	// try a rotation
	rotation := math.Pi / 4
	tiltedShapeA := getTiltedShape(shapeA, rotation)
	tiltedShapeB := getTiltedShape(shapeB, rotation)

	// get the area
	areaA := leftRightArea(shapeA).aggrArea()
	areaB := leftRightArea(shapeB).aggrArea()

	// get the tiltedArea
	tiltedAreaA := leftRightArea(tiltedShapeA).aggrArea()
	tiltedAreaB := leftRightArea(tiltedShapeB).aggrArea()

	// avg area
	avgNormalArea := (areaA + areaB) / 2
	avgTiltedArea := (tiltedAreaA + tiltedAreaB) / 2

	// if the tilted area is bigger than the normal one, return that instead
	if avgTiltedArea > avgNormalArea {
		return tiltedShapeA, tiltedShapeB
	}

	return shapeA, shapeB
}

func getTiltedShape(shape []point, rotation float64) []point {
	count := len(shape)
	matrix := createTransformMatrixForRotation(rotation)
	newShape := make([]point, count)
	for i := 0; i < count; i++ {
		point := shape[i]
		newPoint := transformVector(point, matrix)
		newShape[i] = newPoint
	}
	return newShape
}

func createTransformMatrixForRotation(rotation float64) transformMatrix {
	iHat := point{math.Cos(rotation), math.Sin(rotation)}
	iToJ := math.Pi / 2
	jHat := point{math.Cos(rotation + iToJ), math.Sin(rotation + iToJ)}
	return transformMatrix{iHat, jHat}
}

func transformVector(v point, matrix transformMatrix) point {
	x := (v.X * matrix.iHat.X) + (v.Y * matrix.jHat.X)
	y := (v.X * matrix.iHat.Y) + (v.Y * matrix.jHat.Y)
	return point{x, y}
}
