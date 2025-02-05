package middlewares

import (
	"github.com/gin-gonic/gin"
	"snowdream.tech/gserver/pkg/tools"
)

// Empty Empty
func Empty() gin.HandlerFunc {
	tools.DebugPrintF("[INFO] Starting Middleware %s", "Empty")

	return func(c *gin.Context) {
		c.Next()
	}
}
