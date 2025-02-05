package server

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/snowdreamtech/gserver/pkg/env"
)

func init() {
	rootCmd.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Print " + env.ProjectName + " version and environment info",
	Long:  "All software has versions. This is " + env.ProjectName + "'s",
	Run: func(cmd *cobra.Command, args []string) {
		OSArch := runtime.GOOS + "/" + runtime.GOARCH
		BuildVersion := fmt.Sprintf("%s version %s %s\n", env.ProjectName, env.GitTag, OSArch)
		GOOS := fmt.Sprintf("GOOS=%s\n", runtime.GOOS)
		GOARCH := fmt.Sprintf("GOARCH=%s\n", runtime.GOARCH)
		GOVERSION := fmt.Sprintf("GOVERSION=%s\n", runtime.Version())

		var builder strings.Builder
		builder.WriteString(BuildVersion)
		builder.WriteString(GOOS)
		builder.WriteString(GOARCH)
		builder.WriteString(GOVERSION)

		fmt.Println(builder.String())
	},
}
