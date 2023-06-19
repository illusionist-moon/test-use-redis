package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("cors")

		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin")
		if origin != "" {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-Custom-Header, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			ctx.Header("Access-Control-Expose-Headers", "Authorization, Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			ctx.Header("Access-Control-Max-Age", "172800")
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Set("content-type", "application/json") // 设置返回格式是json
		}
		if method == "OPTIONS" {
			ctx.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "Options Request!",
			})
		}
		//处理请求
		ctx.Next()
	}
}

//func Cors() gin.HandlerFunc {
//	fmt.Println("cors")
//	return cors.New(cors.Config{
//		//准许跨域请求网站,多个使用,分开,限制使用*
//		AllowOrigins: []string{"*"},
//		//准许使用的请求方式
//		AllowMethods: []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
//		//准许使用的请求表头
//		AllowHeaders: []string{"Authorization", "Content-Length", "X-CSRF-Token", "Token", "session", "X_Requested_With", "Accept", "Origin", "Host", "Connection", "Accept-Encoding", "Accept-Language", "DNT", "X-Custom-Header", "Keep-Alive", "User-Agent", "X-Requested-With", " If-Modified-Since", "Cache-Control", "Content-Type, Pragma"},
//		//显示的请求表头
//		ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control", "Content-Language", "Content-Type", "Expires", "Last-Modified", "Pragma", "FooBar"},
//		//凭证共享,确定共享
//		AllowCredentials: true,
//		//容许跨域的原点网站,可以直接return true就万事大吉了
//		AllowOriginFunc: func(origin string) bool {
//			return true
//		},
//		//超时时间设定
//		MaxAge: 24 * time.Hour,
//	})
//}
