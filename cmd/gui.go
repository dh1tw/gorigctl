package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// guiCmd represents the gui command
var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "graphical user interface for a (local/remote) radio",
	Long: `Run a (CLI) graphical user interface.
	
Run a GUI which can be attached either to a local radio or connected through
a specific transportation protocol to a remote radio.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`Please specify if you want to run the GUI with a local radio or connect
to a remote radio. For a remote radio you have to specify the transportation 
protocol (--help for available options)
`)
	},
}

func init() {
	RootCmd.AddCommand(guiCmd)

}
