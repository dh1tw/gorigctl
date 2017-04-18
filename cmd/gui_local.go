package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cskr/pubsub"
	hl "github.com/dh1tw/goHamlib"
	"github.com/dh1tw/gorigctl/cli"
	"github.com/dh1tw/gorigctl/events"
	"github.com/dh1tw/gorigctl/gui"
	"github.com/dh1tw/gorigctl/localradio"
	"github.com/dh1tw/gorigctl/utils"
	ui "github.com/gizak/termui"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var guiLocalCmd = &cobra.Command{
	Use:   "local",
	Short: "GUI client for a local radio",
	Long:  `GUI client for a local radio`,
	Run:   runLocalGui,
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
	guiLocalCmd.Flags().IntP("hl-debug-level", "D", 0, "Hamlib Debug Level (0=ERROR, 5=TRACE")
}

type localGui struct {
	cliCmds []cli.CliCmd
	radio   *localradio.LocalRadio
}

func runLocalGui(cmd *cobra.Command, args []string) {

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.BindPFlag("radio.rig-model", cmd.Flags().Lookup("rig-model"))
	viper.BindPFlag("radio.baudrate", cmd.Flags().Lookup("baudrate"))
	viper.BindPFlag("radio.portname", cmd.Flags().Lookup("portname"))
	viper.BindPFlag("radio.databits", cmd.Flags().Lookup("databits"))
	viper.BindPFlag("radio.stopbits", cmd.Flags().Lookup("stopbits"))
	viper.BindPFlag("radio.parity", cmd.Flags().Lookup("parity"))
	viper.BindPFlag("radio.handshake", cmd.Flags().Lookup("handshake"))
	viper.BindPFlag("radio.hl-debug-level", cmd.Flags().Lookup("hl-debug-level"))

	rigModel := viper.GetInt("radio.rig-model")
	debugLevel := viper.GetInt("radio.hl-debug-level")

	port := hl.Port{}
	port.Baudrate = viper.GetInt("radio.baudrate")
	port.Databits = viper.GetInt("radio.databits")
	port.Stopbits = viper.GetInt("radio.stopbits")
	port.Portname = viper.GetString("radio.portname")
	port.RigPortType = hl.RIG_PORT_SERIAL
	switch viper.GetString("radio.parity") {
	case "none":
		port.Parity = hl.N
	case "even":
		port.Parity = hl.E
	case "odd":
		port.Parity = hl.O
	default:
		port.Parity = hl.N
	}

	switch viper.GetString("radio.handshake") {
	case "none":
		port.Handshake = hl.NO_HANDSHAKE
	case "RTSCTS":
		port.Handshake = hl.RTSCTS_HANDSHAKE
	default:
		port.Handshake = hl.NO_HANDSHAKE
	}

	evPS := pubsub.New(10000)

	logger := utils.NewChLogger(evPS, events.AppLog, "")

	lr, err := localradio.NewLocalRadio(rigModel, debugLevel, port, logger)
	if err != nil {
		fmt.Println("Unable to initialize radio:", err)
		os.Exit(-1)
	}

	lGui := localGui{
		radio:   lr,
		cliCmds: cli.PopulateCliCmds(),
	}

	// prepareShutdownCh := evPS.Sub(events.PrepareShutdown)
	shutdownCh := evPS.Sub(events.Shutdown)
	cliInputCh := evPS.Sub(events.CliInput)
	loggingCh := evPS.Sub(events.AppLog)

	caps, err := lGui.radio.GetCaps()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	go gui.Loop(evPS)

	state, err := lGui.radio.GetState()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	ui.SendCustomEvt("/radio/status", true)
	ui.SendCustomEvt("/radio/caps", caps)
	ui.SendCustomEvt("/radio/state", state)

	for {
		select {
		// shutdown the application gracefully
		case <-shutdownCh:
			//force exit after 1 sec
			exitTimeout := time.NewTimer(time.Second)
			go func() {
				<-exitTimeout.C
				os.Exit(-1)
			}()
			os.Exit(0)

		// case msg := <-toDeserializeCatResponseCh:
		// 	// r.printRigUpdates = true
		// 	state, _ := rGui.radio.GetState()
		// 	ui.SendCustomEvt("/radio/state", state)

		case msg := <-cliInputCh:
			lGui.parseCli(logger, msg.([]string))

		case msg := <-loggingCh:
			// forward to GUI event handler to be shown in the
			// approriate window
			ui.SendCustomEvt("/log/msg", msg)

		case <-shutdownCh:
			log.Println("disconnecting from radio")
			return
		}
	}

}

func (lcli *localGui) parseCli(logger *log.Logger, cliInput []string) {

	found := false

	if len(cliInput) == 0 {
		return
	}

	for _, cmd := range lcli.cliCmds {
		if cmd.Name == cliInput[0] || cmd.Shortcut == cliInput[0] {
			cmd.Cmd(lcli.radio, logger, cliInput[1:])
			found = true
		}
	}

	if cliInput[0] == "help" || cliInput[0] == "?" {
		lcli.PrintHelp(logger)
		found = true
	}

	if !found {
		log.Println("unknown command")
	}
}

func (lGui *localGui) PrintHelp(log *log.Logger) {

	buf := bytes.Buffer{}

	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Command", "Shortcut", "Parameter"})
	table.SetCenterSeparator("|")
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(50)

	for _, el := range lGui.cliCmds {
		table.Append([]string{el.Name, el.Shortcut, el.Parameters})
	}

	table.Render()

	lines := strings.Split(buf.String(), "\n")

	for _, line := range lines {
		log.Println(line)
	}
}
