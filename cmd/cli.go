package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "command line interface for a (local/remote) radio",
	Long: `Run a command line interface.
	
Run a CLI which can be attached either to a local radio or connected through
a specific transportation protocol to a remote radio.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`Please specify if you want to run the CLI with a local radio or connect
to a remote radio. For a remote radio you have to specify the transportation protocol (--help for available options)`)
	},
}

func init() {
	RootCmd.AddCommand(cliCmd)
}
