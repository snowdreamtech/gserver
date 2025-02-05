package middlewares

import (
	"github.com/gin-gonic/gin"
	"snowdream.tech/gserver/pkg/tools"
)

// Header Header
func Header() gin.HandlerFunc {
	tools.DebugPrintF("[INFO] Starting Middleware %s", "Header")

	return func(c *gin.Context) {
		c.Writer.Header().Set("server", "SnowdreamTech Static HTTP Server/0.1")

		c.Next()

	}
}
