package server

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/snowdreamtech/gserver/pkg/env"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of  " + env.ProjectName,
	Long:  "All software has versions. This is  " + env.ProjectName + "'s",
	Run: func(cmd *cobra.Command, args []string) {
		OSArch := runtime.GOOS + "/" + runtime.GOARCH
		BuildVersion := fmt.Sprintf("%s version %s %s\n", env.ProjectName, env.GitTag, OSArch)
		LicenseDetail := fmt.Sprintf("%s\n", env.LICENSE)
		AuthorDetail := fmt.Sprintf("Written by  %s", env.Author)

		var builder strings.Builder
		builder.WriteString(BuildVersion)
		builder.WriteString(LicenseDetail)

		builder.WriteString("\n")
		builder.WriteString(AuthorDetail)

		fmt.Println(builder.String())
	},
}
