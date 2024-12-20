package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/user-web/api"
	"shop_api/user-web/global"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")

	global.GetSugar().Debug("初始化用户相关url")
	{
		UserRouter.GET("list", api.GetUserList)
	}

}
