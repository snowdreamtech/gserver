package log

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/snowdreamtech/gserver/pkg/configs"
	"github.com/snowdreamtech/gserver/pkg/tools"
)

// InitLoggerConfig Init Logger Config
func InitLoggerConfig() {
	r := configs.GetAppConfig()
	logDir := r.LogDir

	if _, err := os.Stat(logDir); err != nil {
		err = os.MkdirAll(logDir, 0640)

		if err != nil {
			return
		}
	}

	// Set gin.log
	ginlog := &lumberjack.Logger{
		Filename:   logDir + "/gin.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	tools.DefaultGinWriter = io.MultiWriter(ginlog, os.Stdout)
	log.SetOutput(tools.DefaultGinWriter)
	log.SetPrefix("[LOG-debug] ")

	// gin.DefaultWriter = io.MultiWriter(accesslog)
	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultWriter = tools.DefaultGinWriter

	// Set error.log
	errorlog := &lumberjack.Logger{
		Filename:   logDir + "/error.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	// gin.DefaultErrorWriter = io.MultiWriter(errorlog)
	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultErrorWriter = io.MultiWriter(errorlog, os.Stderr)

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		tools.DebugPrintF("[INFO] %-6s %-25s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}
