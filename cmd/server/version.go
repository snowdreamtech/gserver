package server

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/snowdreamtech/gserver/pkg/env"
	"github.com/spf13/cobra"
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
		BuildVersion := fmt.Sprintf("%s version %s-%s %s\n", env.ProjectName, env.GitTag, env.CommitHash, OSArch)
		CopyrightDetail := fmt.Sprintf("%s\n", env.COPYRIGHT)
		LicenseDetail := fmt.Sprintf("License: %s\n", env.LICENSE)
		AuthorDetail := fmt.Sprintf("Written by %s", env.Author)
		BuildDetail := fmt.Sprintf("Built at %s", env.BuildTime)

		var builder strings.Builder
		builder.WriteString(BuildVersion)
		builder.WriteString(CopyrightDetail)
		builder.WriteString(LicenseDetail)

		builder.WriteString("\n")
		builder.WriteString(AuthorDetail)
		builder.WriteString("\n")
		builder.WriteString(BuildDetail)

		fmt.Println(builder.String())
	},
}
