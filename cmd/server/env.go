package server

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/snowdreamtech/gserver/pkg/env"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Print " + env.ProjectName + " version and environment info",
	Long:  "All software has versions. This is " + env.ProjectName + "'s",
	Run: func(cmd *cobra.Command, args []string) {
		ProjectName := fmt.Sprintf("ProjectName=%s\n", env.ProjectName)
		GOOS := fmt.Sprintf("GOOS=%s\n", runtime.GOOS)
		GOARCH := fmt.Sprintf("GOARCH=%s\n", runtime.GOARCH)
		GOVERSION := fmt.Sprintf("GOVERSION=%s\n", runtime.Version())
		Author := fmt.Sprintf("Author=%s\n", env.Author)
		BuildTime := fmt.Sprintf("BuildTime=%s\n", env.BuildTime)
		GitTag := fmt.Sprintf("GitTag=%s\n", env.GitTag)
		CommitHash := fmt.Sprintf("CommitHash=%s\n", env.CommitHash)
		CommitHashFull := fmt.Sprintf("CommitHashFull=%s\n", env.CommitHashFull)
		COPYRIGHT := fmt.Sprintf("Copyright=%s\n", env.COPYRIGHT)
		LICENSE := fmt.Sprintf("LICENSE=%s\n", env.LICENSE)

		var builder strings.Builder
		builder.WriteString(ProjectName)
		builder.WriteString(Author)
		builder.WriteString(BuildTime)
		builder.WriteString(GitTag)
		builder.WriteString(CommitHash)
		builder.WriteString(CommitHashFull)
		builder.WriteString(GOOS)
		builder.WriteString(GOARCH)
		builder.WriteString(GOVERSION)
		builder.WriteString(COPYRIGHT)
		builder.WriteString(LICENSE)

		fmt.Println(builder.String())
	},
}
