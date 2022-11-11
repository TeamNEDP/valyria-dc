package services

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"valyria-dc/model"
)

func competeEndpoints(r *gin.RouterGroup) {
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
	script := model.UserScript{}
	if err := db.Where("user_id=?", user.ID).Where("name=?", name).First(&script).Error; err != nil {
		ctx.JSON(notFound("script not found"))
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		tx.Where("user_id=?", user.ID).Delete(&model.UserCompetition{})
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
