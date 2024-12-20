package initialize

import (
	"github.com/gin-gonic/gin"
	"shop_api/user-web/router"
)

func Routers() *gin.Engine {
	// 1.初始化router
	r := gin.Default()
	// 2.调用Group创建路由分组
	ApiV1Group := r.Group("/u/v1")
	// 3.初始化用户路由信息
	router.InitUserRouter(ApiV1Group)
	// 4.返回路由
	return r
}
