package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/user-web/api"
	"shop_api/user-web/global"
	"shop_api/user-web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")

	UserRouter.POST("password_login", api.PasswordLogin)

	global.GetSugar().Debug("初始化用户相关url")
	{
		UserRouter.Use(middlewares.JWTAuth()).GET("list", api.GetUserList)
	}

}
