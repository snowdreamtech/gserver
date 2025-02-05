package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/snowdreamtech/gserver/pkg/configs"
	"github.com/snowdreamtech/gserver/pkg/tools"
)

// Referer Referer
func Referer() gin.HandlerFunc {
	tools.DebugPrintF("[INFO] Starting Middleware %s", "Referer")

	app := configs.GetAppConfig()

	if !app.RefererLimiter {
		return Empty()
	}

	return func(c *gin.Context) {
		referer := c.Request.Referer()

		if referer == "" || localhostRegex.MatchString(referer) {
			c.Next()

			return
		}

		if strings.HasPrefix(referer, "http://"+c.Request.Host) || strings.HasPrefix(referer, "https://"+c.Request.Host) {
			c.Next()

			return
		}

		c.AbortWithStatus(http.StatusForbidden)
	}
}
