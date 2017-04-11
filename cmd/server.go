package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// serverCmdrepresents the serve command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server which makes the radio available on the network",
	Long: `Run a gorigctl server.

Start a gorigctl server using a specific transportation protocol.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please select a transportation protocol (--help for available options)")
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
