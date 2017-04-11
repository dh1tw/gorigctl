package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var version string
var commitHash string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gorigctl",
	Long:  `All software has versions. This is gorigctl's.`,
	Run: func(cmd *cobra.Command, args []string) {
		printGorigctlVersion()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func printGorigctlVersion() {
	buildDate := time.Now().Format(time.RFC3339)
	fmt.Printf("gorigctl Version: %s, %s/%s, BuildDate: %s, Commit: %s\n",
		version, runtime.GOOS, runtime.GOARCH, buildDate, commitHash)
}
