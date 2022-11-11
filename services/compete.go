package services

import (
	"github.com/gin-gonic/gin"
	"valyria-dc/model"
)

func competeEndpoints(r *gin.RouterGroup) {
	g := r.Group("", AuthRequired())

	g.GET("", userCompetitionStatus)
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
