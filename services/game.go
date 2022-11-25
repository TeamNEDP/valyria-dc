package services

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"math"
	"strconv"
	"time"
	"valyria-dc/game"
	"valyria-dc/model"
)

func handleGameEnd(process *game.GameProcess, result game.GameResult) {
	g := model.Game{}
	err := db.
		Preload("RScript.User").
		Preload("BScript.User").
		Where("id=?", process.ID).First(&g).Error
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

	if g.Official {
		r := g.RScript.User
		b := g.BScript.User
		sa := 0.0
		sb := 0.0
		if g.Result.Winner == "R" {
			sa = 1.0
		} else if g.Result.Winner == "B" {
			sb = 1.0
		} else {
			sa = 0.5
			sb = 0.5
		}
		ea := 1.0 / (1 + math.Pow(10, float64(r.Rating-b.Rating)/32.0))
		eb := 1.0 / (1 + math.Pow(10, float64(b.Rating-r.Rating)/32.0))

		r.Rating = int(math.Round(float64(r.Rating) + 32.0*(sa-ea)))
		b.Rating = int(math.Round(float64(b.Rating) + 32.0*(sb-eb)))

		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Save(&r).Error; err != nil {
				return err
			}
			return tx.Save(&b).Error
		})
		if err != nil {
			log.Printf("Failed to save rating change: %v\n", err)
		}
	}
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

type CustomGameForm struct {
	RScript string `json:"r_script"`
	BScript string `json:"b_script"`
}

func gameEndpoints(r *gin.RouterGroup) {
	var games []model.Game
	db.Where("NOT finished").Find(&games)

	for _, g := range games {
		log.Printf("Restarting game %s\n", g.ID)
		go game.StartGame(g.ID, g.Setting)
	}

	g := r.Group("", AuthRequired())

	g.GET("/", listGames)
	each := g.Group("/:id", queryGame)
	each.GET("/details", gameDetails)
	each.GET("/live", func(ctx *gin.Context) {
		game.ServeLive(ctx.Param("id"), ctx.Writer, ctx.Request)
	})

	g.POST("/custom", func(ctx *gin.Context) {
		user := ctx.MustGet("user").(model.User)
		form := CustomGameForm{}
		if err := ctx.ShouldBindJSON(&form); err != nil {
			ctx.JSON(invalidParams("malformed custom form"))
			return
		}

		rScript := model.UserScript{}
		bScript := model.UserScript{}

		if err := db.Where("user_id=?", user.ID).Where("name=?", form.RScript).First(&rScript).Error; err != nil {
			ctx.JSON(notFound("r_script not found"))
			return
		}

		if err := db.Where("user_id=?", user.ID).Where("name=?", form.BScript).First(&bScript).Error; err != nil {
			ctx.JSON(notFound("b_script not found"))
			return
		}

		gameSetting := game.GameSetting{
			Map: game.RandMap(),
			Users: map[string]game.GameUser{
				"r": {
					ID: user.ID,
					Script: game.UserScript{
						Type:    "javascript",
						Content: &rScript.Code,
					},
				},
				"b": {
					ID: user.ID,
					Script: game.UserScript{
						Type:    "javascript",
						Content: &bScript.Code,
					},
				},
			},
		}

		gameId := randCode(8)
		g := model.Game{
			ID:        gameId,
			Finished:  false,
			RScriptID: rScript.ID,
			BScriptID: bScript.ID,
			Official:  false,
			Setting:   gameSetting,
			Ticks:     game.GameTicks{},
			Result:    game.GameResult{},
			CreatedAt: time.Now(),
		}

		db.Save(&g)

		go game.StartGame(gameId, gameSetting)

		ctx.JSON(resOk(nil))
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

	limit := 20
	offset := 0

	if lim, err := strconv.Atoi(ctx.Query("limit")); err == nil {
		limit = lim
	}
	if limit > 20 {
		ctx.JSON(invalidParams("limit too large"))
	}
	if off, err := strconv.Atoi(ctx.Query("offset")); err == nil {
		offset = off
	}

	var games []model.Game
	db.
		Preload("RScript").
		Preload("BScript").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Where("r_script_id IN (?) OR b_script_id IN (?)",
			db.Model(&model.UserScript{}).Where("user_id=?", user.ID).Select("id"),
			db.Model(&model.UserScript{}).Where("user_id=?", user.ID).Select("id"),
		).
		Omit("ticks").
		Find(&games)

	res := make([]UserGameEntry, 0, len(games))
	for _, v := range games {
		entry := UserGameEntry{
			ID:       v.ID,
			Date:     v.CreatedAt.Unix(),
			Official: v.Official,
		}
		if v.RScript.UserID == user.ID {
			entry.Role = "R"
		} else {
			entry.Role = "B"
		}
		if v.Finished {
			result := v.Result
			entry.Status = "finished"
			entry.Result = &result
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
