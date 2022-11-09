package services

import (
	"log"
	"valyria-dc/game"
	"valyria-dc/model"
)

func HandleGameEnd(process game.GameProcess, result game.GameResult) {
	g := model.Game{}
	err := db.Where("id=?", process.ID).First(&g).Error
	if err != nil {
		log.Printf("Error querying game: %v\n", err)
		return
	}
	g.Ticks = game.GameTicks{
		Ticks: process.Ticks,
	}
	g.Result = result
	g.Finished = true

	db.Save(&g)
}
