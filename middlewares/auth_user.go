package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online-judge/helper"
)

func AuthUserCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		userClaim, err := helper.AnalyseToken(auth)
		if err != nil {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized Authorization",
			})
			return
		}
		if userClaim == nil {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized Admin",
			})
			return
		}
		ctx.Set("user_claims", userClaim)
		ctx.Next()
	}
}
