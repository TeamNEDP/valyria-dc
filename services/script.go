package services

import (
	"github.com/gin-gonic/gin"
	"valyria-dc/model"
)

type ScriptForm struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func scriptEndpoints(r *gin.RouterGroup) {
	g := r.Group("", AuthRequired())

	g.POST("", createScript)
	g.GET("/", listScripts)
	g.PATCH("/:name", queryScript, patchScript)
	g.DELETE("/:name", queryScript, deleteScript)
}

func queryScript(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)
	name := ctx.Param("name")

	script := model.UserScript{}
	err := db.Where("user_id=?", user.ID).Where("name=?", name).First(&script).Error
	if err != nil {
		ctx.JSON(notFound("script not found"))
		ctx.Abort()
		return
	}

	ctx.Set("script", script)
}

func createScript(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)
	form := ScriptForm{}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(invalidParams("malformed script form"))
		return
	}
	script := model.UserScript{
		UserID: user.ID,
		Name:   form.Name,
		Code:   form.Code,
	}
	err := db.Save(&script).Error
	if err != nil {
		ctx.JSON(internalError(err.Error()))
		return
	}
	ctx.JSON(resOk(nil))
}

func listScripts(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)
	var scripts []model.UserScript
	db.Where("user_id=?", user.ID).Find(&scripts)

	res := make([]model.ScriptInfo, 0, len(scripts))
	for _, v := range scripts {
		res = append(res, v.Info())
	}
	ctx.JSON(resOk(res))
}

func patchScript(ctx *gin.Context) {
	script := ctx.MustGet("script").(model.UserScript)
	form := ScriptForm{}

	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(invalidParams("malformed script form"))
		return
	}

	script.Name = form.Name
	script.Code = form.Code

	err := db.Save(&script).Error

	if err != nil {
		ctx.JSON(internalError(err.Error()))
		return
	}

	ctx.JSON(resOk(nil))
}

func deleteScript(ctx *gin.Context) {
	script := ctx.MustGet("script").(model.UserScript)
	script.Name = randCode(16)
	db.Save(&script)
	db.Delete(&script)
	ctx.JSON(resOk(nil))
}
