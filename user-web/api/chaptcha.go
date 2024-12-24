package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha() (ctx *gin.Context) {
	base64Captcha.NewDriverDigit(240, 80, 5, 0.7, 80)

}
