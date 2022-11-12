package services

import (
	"github.com/gin-gonic/gin"
	"valyria-dc/model"
)

func topsEndpoint(r *gin.RouterGroup) {
	r.GET("", func(ctx *gin.Context) {
		var users []model.User
		db.Order("rating DESC").Limit(10).Find(&users)
		res := make([]model.UserInfo, 0, len(users))
		for _, v := range users {
			res = append(res, v.Info())
		}
		ctx.JSON(resOk(res))
	})
}
