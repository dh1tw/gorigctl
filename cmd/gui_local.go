package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var guiLocalCmd = &cobra.Command{
	Use:   "local",
	Short: "GUI client for a local radio",
	Long:  `GUI client for a local radio`,
	Run:   localGui,
}

func init() {
	guiCmd.AddCommand(guiLocalCmd)
	guiLocalCmd.Flags().IntP("rig-model", "m", 0, "Hamlib Rig Model ID")
	guiLocalCmd.Flags().IntP("baudrate", "b", 38400, "Baudrate")
	guiLocalCmd.Flags().StringP("portname", "o", "/dev/mhux/cat", "Portname (e.g. COM1)")
	guiLocalCmd.Flags().IntP("databits", "d", 8, "Databits")
	guiLocalCmd.Flags().IntP("stopbits", "s", 1, "Stopbits")
	guiLocalCmd.Flags().StringP("parity", "r", "none", "Parity")
	guiLocalCmd.Flags().StringP("handshake", "a", "none", "Handshake")
}

func localGui(cmd *cobra.Command, args []string) {
	viper.BindPFlag("radio.rig-model", cmd.Flags().Lookup("rig-model"))
	viper.BindPFlag("radio.baudrate", cmd.Flags().Lookup("baudrate"))
	viper.BindPFlag("radio.portname", cmd.Flags().Lookup("portname"))
	viper.BindPFlag("radio.databits", cmd.Flags().Lookup("databits"))
	viper.BindPFlag("radio.stopbits", cmd.Flags().Lookup("stopbits"))
	viper.BindPFlag("radio.parity", cmd.Flags().Lookup("parity"))
	viper.BindPFlag("radio.handshake", cmd.Flags().Lookup("handshake"))

	fmt.Println("local gui called")
}
