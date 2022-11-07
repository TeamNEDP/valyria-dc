package services

import (
	"github.com/gin-gonic/gin"
	"valyria-dc/model"
)

func AuthRequired() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sessionId, err := ctx.Cookie("session_id")
		if err != nil {
			ctx.JSON(unauthorized())
			ctx.Abort()
			return
		}
		var session model.UserSession
		if err := db.Preload("User").
			Where("session_id=?", sessionId).
			First(&session).Error; err != nil {
			ctx.JSON(unauthorized())
			ctx.Abort()
			return
		}
		session.UserAgent = ctx.GetHeader("user-agent")
		db.Save(&session)
		ctx.Set("session", session)
		ctx.Set("user", session.User)
	}
}
