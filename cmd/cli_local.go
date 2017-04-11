package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cliLocalCmd = &cobra.Command{
	Use:   "local",
	Short: "command line client for a local radio",
	Long:  `command line client for a local radio`,
	Run:   localCli,
}

func init() {
	cliCmd.AddCommand(cliLocalCmd)
	cliLocalCmd.Flags().IntP("rig-model", "m", 0, "Hamlib Rig Model ID")
	cliLocalCmd.Flags().IntP("baudrate", "b", 38400, "Baudrate")
	cliLocalCmd.Flags().StringP("portname", "o", "/dev/mhux/cat", "Portname (e.g. COM1)")
	cliLocalCmd.Flags().IntP("databits", "d", 8, "Databits")
	cliLocalCmd.Flags().IntP("stopbits", "s", 1, "Stopbits")
	cliLocalCmd.Flags().StringP("parity", "r", "none", "Parity")
	cliLocalCmd.Flags().StringP("handshake", "a", "none", "Handshake")
}

func localCli(cmd *cobra.Command, args []string) {
	viper.BindPFlag("radio.rig-model", cmd.Flags().Lookup("rig-model"))
	viper.BindPFlag("radio.baudrate", cmd.Flags().Lookup("baudrate"))
	viper.BindPFlag("radio.portname", cmd.Flags().Lookup("portname"))
	viper.BindPFlag("radio.databits", cmd.Flags().Lookup("databits"))
	viper.BindPFlag("radio.stopbits", cmd.Flags().Lookup("stopbits"))
	viper.BindPFlag("radio.parity", cmd.Flags().Lookup("parity"))
	viper.BindPFlag("radio.handshake", cmd.Flags().Lookup("handshake"))

	fmt.Println("local gui called")
}
