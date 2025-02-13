package goods

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shop_api/goods-web/global"
	"shop_api/goods-web/proto"
	"strconv"
)

func HandleValidatorError(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"msg": errs.Translate(global.Trans),
	})
}

func HandleGrpcErrorToHttp(err error, ctx *gin.Context) {
	if err != nil {
		zap.S().Errorw("[HandleGrpcErrorToHttp] 处理grpc错误", "msg", err.Error())
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, gin.H{"msg": e.Message()})
			case codes.InvalidArgument:
				ctx.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
			case codes.Unavailable:
				ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "服务不可用"})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "内部错误"})
			}
		}
		return
	}
}

func GetGoodsList(ctx *gin.Context) {
	zap.S().Debug("获取商品列表页接口")

	request := &proto.GoodsFilterRequest{}

	priceMin := ctx.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)

	priceMax := ctx.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)

	isHot := ctx.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}

	isNew := ctx.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}

	isTab := ctx.DefaultQuery("it", "0")
	if isTab == "1" {
		request.IsTab = true
	}

	categoryId := ctx.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	request.TopCategory = int32(categoryIdInt)

	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	keywords := ctx.DefaultQuery("q", "")
	request.KeyWords = keywords

	brandId := ctx.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	request.Brand = int32(brandIdInt)

	//请求商品的service服务、负载均衡
	//parent, _ := ctx.Get("parentSpan")
	//opentracing.ContextWithSpan(context.Background(), parent.(opentracing.Span))

	//e, b := sentinel.Entry("goods-list", sentinel.WithTrafficType(base.Inbound))
	//if b != nil {
	//	ctx.JSON(http.StatusTooManyRequests, gin.H{
	//		"msg": "请求过于频繁，请稍后重试",
	//	})
	//	return
	//}
	r, err := global.GoodsSrvClient.GoodsList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("[List] 查询 【商品列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := map[string]interface{}{
		"total":     r.Total,
		"page":      pagesInt,
		"page_size": perNumsInt,
	}

	goodsList := make([]interface{}, 0)
	for _, value := range r.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.IsFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}
	reMap["data"] = goodsList

	ctx.JSON(http.StatusOK, reMap)
}
