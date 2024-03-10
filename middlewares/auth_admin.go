package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online-judge/helper"
)

func AuthAdminCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//TODO:Check if user is admin
		auth := ctx.GetHeader("Authorization")
		userClaim, err := helper.AnalyseToken(auth)
		if err != nil {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "UnAuthorized",
			})
			return
		}

		if userClaim.IsAdmin != 1 {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "UnAuthorized",
			})
			return
		}

		ctx.Next()
	}
}
