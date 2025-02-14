package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop_api/goods-web/api/goods"
	"shop_api/goods-web/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	//
	zap.S().Debug("初始化商品相关url")
	{
		GoodsRouter.GET("list", goods.GetGoodsList)
	}

	GoodsRouter.Use(middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		GoodsRouter.POST("new", goods.New)
		GoodsRouter.GET("/:id", goods.Detail)
		GoodsRouter.DELETE("/:id", goods.Delete)
	}

}
