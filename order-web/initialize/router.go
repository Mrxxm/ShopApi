package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_api/order-web/middlewares"
	"shop_api/order-web/router"
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
	ApiV1Group := r.Group("/o/v1")
	// 5.初始化路由信息
	router.InitOrderRouter(ApiV1Group)
	router.InitShopCartRouter(ApiV1Group)
	// 6.返回路由
	return r
}
