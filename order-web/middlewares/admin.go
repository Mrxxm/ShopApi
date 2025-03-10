package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_api/order-web/models"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		curUser := claims.(*models.CustomClaims)

		if curUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
