package middleware

import (
	"ChildrenMath/pkg/settings"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"log"
)

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(settings.ServerIp + ":" + settings.ServerPort)
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     settings.ServerIp + ":" + settings.ServerPort,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}
