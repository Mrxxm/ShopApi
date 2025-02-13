package goods

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shop_api/goods-web/global"
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

	ctx.JSON(http.StatusOK, gin.H{
		"data":     0,
		"total":    0,
		"page":     1,
		"pageSize": 5,
	})
}
