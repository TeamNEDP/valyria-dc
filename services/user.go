package services

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"valyria-dc/model"
)

type UserRegisterForm struct {
	Name     string ` json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func userEndpoints(r *gin.RouterGroup) {
	r.POST("", userRegister)
	r.POST("/login", userLogin)
	r.POST("/logout", AuthRequired(), userLogout)
}

func userRegister(ctx *gin.Context) {
	form := UserRegisterForm{}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(invalidParams("malformed register form"))
		return
	}

	hpw, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(internalError(err))
		return
	}

	user := model.User{
		ID:       randCode(8),
		Name:     form.Name,
		Email:    form.Email,
		Avatar:   nil,
		Rating:   0,
		Password: hpw,
	}
	err = db.Save(&user).Error
	if err != nil {
		ctx.JSON(internalError(err))
		return
	}
	ctx.JSON(resOk(nil))
}

func userLogin(ctx *gin.Context) {
	form := UserLoginForm{}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(invalidParams("malformed login form"))
		return
	}
	var user model.User
	err := db.Where("email=?", form.Email).First(&user).Error
	if err != nil {
		ctx.JSON(notFound("user not found"))
		return
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(form.Password))
	if err != nil {
		ctx.JSON(invalidParams("invalid password"))
		return
	}
	session := model.UserSession{
		UserID:    user.ID,
		UserAgent: ctx.GetHeader("user-agent"),
		SessionID: randCode(16),
	}
	err = db.Save(&session).Error
	if err != nil {
		ctx.JSON(internalError(err))
		return
	}
	ctx.SetCookie("session_id", session.SessionID, 86400*30, "/", "", false, true)
	ctx.JSON(resOk(nil))
}

func userLogout(ctx *gin.Context) {
	ctx.SetCookie("session_id", "", 0, "/", "", false, false)
	ctx.JSON(resOk(nil))
}
