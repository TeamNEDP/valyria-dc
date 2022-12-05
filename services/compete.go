package services

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"time"
	"valyria-dc/game"
	"valyria-dc/model"
)

func competeEndpoints(r *gin.RouterGroup) {
	go func() {
		left := false
		for {
			var cp []model.UserCompetition
			db.
				Preload("UserScript").
				Joins("left join users on user_competitions.user_id=users.id").
				Order("users.rating").
				Find(&cp)
			if (len(cp) & 1) == 0 {
				for i := 0; (i + 1) < len(cp); i += 2 {
					startCompetition(cp[i], cp[i+1])
				}
			} else {
				i := 1
				if left {
					i = 0
				}
				for ; (i + 1) < len(cp); i += 2 {
					startCompetition(cp[i], cp[i+1])
				}
				left = !left
			}
			time.Sleep(time.Minute * 10)
		}
	}()

	g := r.Group("", AuthRequired())

	g.GET("", userCompetitionStatus)
	g.POST("", userCompetitionSet)
}

type CompetitionStatus struct {
	Involved   bool    `json:"involved"`
	ScriptName *string `json:"scriptName,omitempty"`
}

func userCompetitionStatus(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)
	c := model.UserCompetition{}
	if err := db.Preload("UserScript").Where("user_id=?", user.ID).Error; err != nil {
		ctx.JSON(resOk(CompetitionStatus{Involved: false}))
		return
	}
	ctx.JSON(resOk(CompetitionStatus{
		Involved:   true,
		ScriptName: &c.UserScript.Name,
	}))
}

func userCompetitionSet(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)

	var name *string
	if err := ctx.ShouldBindJSON(&name); err != nil {
		ctx.JSON(invalidParams("invalid competition arguments"))
		return
	}

	if name == nil {
		err := db.Unscoped().Where("user_id=?", user.ID).Delete(&model.UserCompetition{}).Error
		if err != nil {
			ctx.JSON(internalError(err.Error()))
			return
		}
		ctx.JSON(resOk(nil))
		return
	}

	script := model.UserScript{}
	if err := db.Where("user_id=?", user.ID).Where("name=?", name).First(&script).Error; err != nil {
		ctx.JSON(notFound("script not found"))
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		tx.Unscoped().Where("user_id=?", user.ID).Delete(&model.UserCompetition{})
		return db.Save(&model.UserCompetition{
			UserID:       user.ID,
			UserScriptID: script.ID,
		}).Error
	})

	if err != nil {
		ctx.JSON(internalError(err.Error()))
		return
	}

	ctx.JSON(resOk(nil))
}

func startCompetition(r model.UserCompetition, b model.UserCompetition) {
	gameSetting := game.GameSetting{
		Map: game.RandMap(),
		Users: map[string]game.GameUser{
			"r": {
				ID: r.UserID,
				Script: game.UserScript{
					Type:    "javascript",
					Content: &r.UserScript.Code,
				},
			},
			"b": {
				ID: b.UserID,
				Script: game.UserScript{
					Type:    "javascript",
					Content: &b.UserScript.Code,
				},
			},
		},
	}

	g := model.Game{
		ID:        randCode(8),
		Finished:  false,
		RScriptID: r.UserScriptID,
		BScriptID: b.UserScriptID,
		Official:  true,
		Setting:   gameSetting,
		Ticks:     game.GameTicks{},
		Result:    game.GameResult{},
		CreatedAt: time.Now(),
	}
	if err := db.Save(&g).Error; err != nil {
		log.Printf("Failed to start competition: %v\n", err)
		return
	}
	game.StartGame(g.ID, gameSetting)
}
