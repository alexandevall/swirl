package main

func getLowCoor(coors []coor) coor {
	southLat := coors[0].Lat
	westLong := coors[0].Long
	for _, coor := range coors {
		if coor.Lat < southLat {
			southLat = coor.Lat
		}
		if coor.Long < westLong {
			westLong = coor.Long
		}
	}
	return coor{Lat: southLat, Long: westLong}
}

func getLowestCoorForRuns(runA, runB run) coor {
	lowA := getLowCoor(runA.Coors)
	lowB := getLowCoor(runB.Coors)
	lowest := lowA
	if lowA.Lat > lowB.Lat {
		lowest.Lat = lowB.Lat
	}
	if lowA.Long > lowB.Long {
		lowest.Long = lowB.Long
	}
	return lowest
}
