package game

import (
	"math/rand"
	"time"
)

const MapWidth = 5
const MapHeight = 6

func RandMap() GameMap {
	rand.Seed(time.Now().UnixNano())

	res := GameMap{
		Width:  MapWidth,
		Height: MapHeight,
		Grids:  make([]MapGrid, MapWidth*MapHeight),
	}

	res.Grids[0] = MapGrid{
		Type:     "R",
		Soldiers: new(int),
	}
	res.Grids[MapWidth*MapHeight-1] = MapGrid{
		Type:     "B",
		Soldiers: new(int),
	}

	for i := 1; i < MapWidth*MapHeight-1; i++ {
		res.Grids[i] = randGrid()
	}

	return res
}

func randGrid() MapGrid {
	grid := MapGrid{
		Type:     "",
		Soldiers: nil,
	}
	rnd := rand.Float64()
	if rnd < 0.8 {
		grid.Type = "V"
	} else if rnd < 0.85 {
		grid.Type = "M"
	} else {
		grid.Type = "C"
		grid.Soldiers = new(int)
		*grid.Soldiers = rand.Int()%30 + 10
	}
	return grid
}
