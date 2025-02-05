package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"snowdream.tech/gserver/pkg/configs"
	"snowdream.tech/gserver/pkg/tools"
)

// Accounts user credential in basic auth.
var accounts = gin.Accounts{
	"admin": "admin",
}

// BasicAuth BasicAuth
// BasicAuth
func BasicAuth() gin.HandlerFunc {
	tools.DebugPrintF("[INFO] Starting Middleware %s", "BasicAuth")

	app := configs.GetAppConfig()

	if !app.Basic || app.User == "" {
		return Empty()
	}

	arr := strings.SplitN(app.User, ":", 2)

	if arr == nil || len(arr) != 2 {
		return Empty()
	}

	accounts = gin.Accounts{
		arr[0]: arr[1],
	}

	return gin.BasicAuth(accounts)
}
