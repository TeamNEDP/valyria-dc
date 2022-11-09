package services

import (
	"github.com/gin-gonic/gin"
	"log"
	"valyria-dc/game"
	"valyria-dc/model"
)

func handleGameEnd(process game.GameProcess, result game.GameResult) {
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

func gameEndpoints(r *gin.RouterGroup) {
	g := r.Group("", AuthRequired())

	g.GET("/:id/live", func(ctx *gin.Context) {
		game.ServeLive(ctx.Param("id"), ctx.Writer, ctx.Request)
	})
}
