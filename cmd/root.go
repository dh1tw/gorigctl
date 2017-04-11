package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gorigctl",
	Short: "An 'opinionated' drop-in replacement for hamlib's gorigctl",
	Long: `gorigctl is an opinionated drop-in replacement for hamlib's gorigctl(d)
	
gorigctl allows you to connect to a local or remote radio, either through a 
command line interface (CLI) or a cli based graphical user interface (GUI).

gorigctl provides also a server (daemon) which makes the radio available on the 
network.

The user experience depends heavily on the hamlib implementations for each 
particular radio. Radios with a 'stable' backend provide the best user experience."
`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gorigctl.[yaml|toml|json])")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".gorigctl") // name of config file (without extension)
		viper.AddConfigPath("$HOME")     // adding home directory as first search path
	}

	viper.AutomaticEnv() // read in environment variables that match
}
