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
		GoodsRouter.GET("list", goods.GetGoodsList) // 商品列表
	}

	GoodsRouter.Use(middlewares.JWTAuth(), middlewares.IsAdminAuth())
	{
		GoodsRouter.POST("new", goods.New)
		GoodsRouter.GET("/:id", goods.Detail) // 获取商品详情
		GoodsRouter.DELETE("/:id", goods.Delete)
		GoodsRouter.GET("/:id/stocks", goods.Stocks) // 获取商品的

		GoodsRouter.PUT("/:id", goods.Update)
		GoodsRouter.PATCH("/:id", goods.UpdateStatus)
	}

}
