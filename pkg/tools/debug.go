package tools

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// DefaultGinWriter  Log for gin app
var DefaultGinWriter io.Writer = os.Stdout

// DefaultAccessWriter  Log for gin app
var DefaultAccessWriter io.Writer = os.Stdout

// DefaultErrorWriter  Log for gin app
var DefaultErrorWriter io.Writer = os.Stderr

// DebugPrintF use it to print formatted logs
func DebugPrintF(format string, values ...any) {
	if gin.IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(DefaultGinWriter, "[GIN-debug] "+format, values...)
	}
}

// DebugPrint use it to print logs
func DebugPrint(values ...any) {
	if gin.IsDebugging() {
		fmt.Fprint(DefaultGinWriter, "[GIN-debug] ", values)
	}
}
