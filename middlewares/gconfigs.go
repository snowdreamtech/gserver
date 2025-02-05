package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/snowdreamtech/gserver/pkg/configs"
	"github.com/snowdreamtech/gserver/pkg/tools"
)

// Configs  Configs with Viper
// Viper is a complete configuration solution for Go applications including 12-Factor apps.
// It is designed to work within an application, and can handle all types of configuration needs and formats. It supports:
// setting defaults
// reading from JSON, TOML, YAML, HCL, envfile and Java properties config files
// live watching and re-reading of config files (optional)
// reading from environment variables
// reading from remote config systems (etcd or Consul), and watching changes
// reading from command line flags
// reading from buffer
// setting explicit values
// Viper can be thought of as a registry for all of your applications configuration needs.
func Configs(conf *configs.Configs) gin.HandlerFunc {
	tools.DebugPrintF("[INFO] Starting Middleware %s", "Configs")

	return func(c *gin.Context) {
		c.Set(configs.ConfigKey, conf)
		c.Next()
	}
}
