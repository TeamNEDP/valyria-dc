package services

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
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

type UserGameEntry struct {
	ID       string           `json:"id"`
	Role     string           `json:"role"`
	Date     int64            `json:"date"`
	Status   string           `json:"status"`
	Official bool             `json:"official"`
	Result   *game.GameResult `json:"result,omitempty"`
}

type GameDetails struct {
	Map    game.GameMap    `json:"map"`
	Ticks  []game.GameTick `json:"ticks"`
	Result game.GameResult `json:"result"`
}

func gameEndpoints(r *gin.RouterGroup) {
	g := r.Group("", AuthRequired())

	g.GET("/", listGames)
	each := g.Group("/:id", queryGame)
	each.GET("/details", gameDetails)
	each.GET("/live", func(ctx *gin.Context) {
		game.ServeLive(ctx.Param("id"), ctx.Writer, ctx.Request)
	})
}

func queryGame(ctx *gin.Context) {
	id := ctx.Param("id")
	g := model.Game{}
	err := db.Where("id=?", id).First(&g).Error
	if err != nil {
		ctx.JSON(notFound("game not found"))
		ctx.Abort()
		return
	}
	ctx.Set("game", g)
}

func listGames(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)

	limit := 50
	offset := 0

	if lim, err := strconv.Atoi(ctx.Query("limit")); err != nil {
		limit = lim
	}
	if off, err := strconv.Atoi(ctx.Query("offset")); err != nil {
		offset = off
	}

	var games []model.Game
	db.
		Preload("RScript").
		Preload("BScript").
		Limit(limit).
		Offset(offset).
		Where("r_script_id IN (?) OR b_script_id IN (?)",
			db.Model(&model.UserScript{}).Where("user_id=?", user.ID).Select("id"),
			db.Model(&model.UserScript{}).Where("user_id=?", user.ID).Select("id"),
		).
		Find(&games)

	res := make([]UserGameEntry, 0, len(games))
	for _, v := range games {
		entry := UserGameEntry{
			ID:       v.ID,
			Date:     v.CreatedAt.Unix(),
			Official: false,
		}
		if v.RScript.UserID == user.ID {
			entry.Role = "R"
		} else {
			entry.Role = "B"
		}
		if v.Finished {
			entry.Status = "finished"
			entry.Result = &v.Result
		} else {
			if game.IsRunning(v.ID) {
				entry.Status = "running"
			} else {
				entry.Status = "queue"
			}
		}
		res = append(res, entry)
	}

	ctx.JSON(resOk(res))
}

func gameDetails(ctx *gin.Context) {
	g := ctx.MustGet("game").(model.Game)
	if !g.Finished {
		ctx.JSON(invalidParams("game not finished"))
		return
	}
	ctx.JSON(resOk(GameDetails{
		Map:    g.Setting.Map,
		Ticks:  g.Ticks.Ticks,
		Result: g.Result,
	}))
}
