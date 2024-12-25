package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/user-web/api"
	"shop_api/user-web/global"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha)
		//BaseRouter.POST("register", api.Register)
		//BaseRouter.POST("login", api.Login)
	}
	global.GetSugar().Debug("初始化基础路由")
}
