package services

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
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
	r.POST("/avatar", AuthRequired(), userChangeAvatar)
	r.GET("/:id/avatar", queryUser, userAvatar)
	r.GET("/info", AuthRequired(), func(ctx *gin.Context) {
		user := ctx.MustGet("user").(model.User)
		ctx.JSON(resOk(user.Info()))
	})
}

func queryUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user := model.User{}
	err := db.Where("id=?", id).First(&user).Error
	if err != nil {
		ctx.JSON(notFound("user not found"))
		ctx.Abort()
		return
	}
	ctx.Set("pathUser", user)
}

func userRegister(ctx *gin.Context) {
	form := UserRegisterForm{}
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(invalidParams("malformed register form"))
		return
	}

	hpw, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(internalError(err.Error()))
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
		ctx.JSON(internalError(err.Error()))
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
		ctx.JSON(internalError(err.Error()))
		return
	}
	ctx.SetCookie("session_id", session.SessionID, 86400*30, "/", "", false, true)
	ctx.JSON(resOk(nil))
}

func userLogout(ctx *gin.Context) {
	ctx.SetCookie("session_id", "", 0, "/", "", false, false)
	ctx.JSON(resOk(nil))
}

func userChangeAvatar(ctx *gin.Context) {
	user := ctx.MustGet("user").(model.User)
	ff, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(invalidParams("malformed multipart form"))
		return
	}
	if ff.Size > 2*1024*1024 {
		ctx.JSON(invalidParams("file size too large"))
		return
	}
	f, err := ff.Open()
	if err != nil {
		ctx.JSON(invalidParams("failed to open form file"))
		return
	}
	defer f.Close()
	conf, _, err := image.DecodeConfig(f)

	if conf.Width > 2000 || conf.Height > 2000 || conf.Width != conf.Height {
		ctx.JSON(invalidParams("image size too large or aspect ratio is not 1:1"))
		return
	}

	_, err = f.Seek(0, io.SeekStart)

	if err != nil {
		ctx.JSON(internalError("failed to seek"))
		return
	}

	img, _, err := image.Decode(f)
	if err != nil {
		ctx.JSON(internalError("failed to decode img"))
		return
	}

	w := new(bytes.Buffer)
	err = jpeg.Encode(w, img, nil)
	if err != nil {
		ctx.JSON(internalError("failed to encode jpeg"))
		return
	}

	user.Avatar = w.Bytes()
	err = db.Save(&user).Error
	if err != nil {
		ctx.JSON(internalError(err.Error()))
		return
	}

	ctx.JSON(resOk(nil))
}

func userAvatar(ctx *gin.Context) {
	user := ctx.MustGet("pathUser").(model.User)
	ctx.Data(http.StatusOK, "image/jpeg", user.Avatar)
}
