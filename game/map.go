package game

import (
	"math/rand"
	"time"
)

func RandMap() GameMap {
	rand.Seed(time.Now().UnixNano())

	res := GameMap{
		Width:  13,
		Height: 13,
		Grids:  make([]MapGrid, 13*13),
	}

	res.Grids[0] = MapGrid{
		Type:     "R",
		Soldiers: new(int),
	}
	res.Grids[13*13-1] = MapGrid{
		Type:     "B",
		Soldiers: new(int),
	}

	for i := 1; i < 13*13-1; i++ {
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
