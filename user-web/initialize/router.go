package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_api/user-web/middlewares"
	"shop_api/user-web/router"
)

func Routers() *gin.Engine {
	// 1.初始化router
	r := gin.Default()
	// 2.使用中间件，解决跨域问题
	r.Use(middlewares.Cors())
	// 3.健康检测
	r.GET("health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})
	// 4.调用Group创建路由分组
	ApiV1Group := r.Group("/u/v1")
	// 5.初始化路由信息
	router.InitBaseRouter(ApiV1Group)
	router.InitUserRouter(ApiV1Group)
	// 6.返回路由
	return r
}
