package middlewares

import (
	"regexp"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/snowdreamtech/gserver/pkg/tools"
)

const (
	localhostRegexString = "[http://|https://]?[localhost|127\\.0\\.0\\.1]:?(6553[0-5]|655[0-2][0-9]|65[0-4][0-9]{2}|6[0-4][0-9]{3}|[1-5][0-9]{4}|[1-9][0-9]{0,3}|0)?/?$"
)

var (
	localhostRegex = regexp.MustCompile(localhostRegexString)
)

// Cors cors
func Cors() gin.HandlerFunc {
	tools.DebugPrintF("[INFO] Starting Middleware %s", "Cors")

	return cors.New(cors.Config{
		AllowOrigins:     []string{},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", " Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return localhostRegex.MatchString(origin)
		},
		MaxAge: 24 * time.Hour,
	})
}
