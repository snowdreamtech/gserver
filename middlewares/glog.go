package middlewares

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/snowdreamtech/gserver/pkg/configs"
	"github.com/snowdreamtech/gserver/pkg/tools"
)

// LoggerWithFormatter instance a Logger middleware with the specified log format function.
func LoggerWithFormatter() gin.HandlerFunc {
	tools.DebugPrintF("[INFO] Starting Middleware %s", "Logger")

	r := configs.GetAppConfig()
	logDir := r.LogDir

	// Set access.log
	accesslog := &lumberjack.Logger{
		Filename:   logDir + "/access.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	tools.DefaultAccessWriter = io.MultiWriter(accesslog, os.Stdout)

	// Set access.log middleware
	accessLogFormatter := func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] %s %s %s %s %d %s \"%s\" %s \n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Request.Header.Get("X-Request-ID"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}

	accessLogConfig := gin.LoggerConfig{
		Formatter: accessLogFormatter,
		Output:    tools.DefaultAccessWriter,
	}

	return gin.LoggerWithConfig(accessLogConfig)
}
