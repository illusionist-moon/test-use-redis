package middleware

import (
	"ChildrenMath/pkg/e"
	"ChildrenMath/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthUserCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		userClaim, err := util.AnalyseToken(auth)
		if err != nil {
			if err.Error() == "timeout" {
				ctx.JSON(http.StatusOK, gin.H{
					"code": e.ErrorAuthCheckTokenTimeout,
					"msg":  e.GetMsg(e.ErrorAuthCheckTokenTimeout),
				})
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"code": e.ErrorAuthCheckTokenFail,
					"msg":  e.GetMsg(e.ErrorAuthCheckTokenFail),
				})
			}
			ctx.Abort()
			return
		}
		ctx.Set("userid", userClaim.UserID)
		ctx.Set("email", userClaim.Email)
		ctx.Next()
	}
}
