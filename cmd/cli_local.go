package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"sync"

	"github.com/cskr/pubsub"
	hl "github.com/dh1tw/goHamlib"
	"github.com/dh1tw/gorigctl/cli"
	"github.com/dh1tw/gorigctl/events"
	"github.com/dh1tw/gorigctl/localradio"
	"github.com/dh1tw/gorigctl/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cliLocalCmd = &cobra.Command{
	Use:   "local",
	Short: "command line client for a local radio",
	Long:  `command line client for a local radio`,
	Run:   runLocalCli,
}

func init() {
	cliCmd.AddCommand(cliLocalCmd)
	cliLocalCmd.Flags().IntP("rig-model", "m", 1, "Hamlib Rig Model ID")
	cliLocalCmd.Flags().IntP("baudrate", "b", 38400, "Baudrate")
	cliLocalCmd.Flags().StringP("portname", "o", "/dev/mhux/cat", "Portname / Device path")
	cliLocalCmd.Flags().IntP("databits", "d", 8, "Databits")
	cliLocalCmd.Flags().IntP("stopbits", "s", 1, "Stopbits")
	cliLocalCmd.Flags().StringP("parity", "r", "none", "Parity")
	cliLocalCmd.Flags().StringP("handshake", "a", "none", "Handshake")
	cliLocalCmd.Flags().IntP("hl-debug-level", "D", 0, "Hamlib Debug Level (0=ERROR,..., 5=TRACE)")
}

type localCli struct {
	cliCmds []cli.CliCmd
	radio   *localradio.LocalRadio
}

func runLocalCli(cmd *cobra.Command, args []string) {

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

	logger := utils.NewStdLogger("", 0)

	lr, err := localradio.NewLocalRadio(rigModel, debugLevel, port, logger)
	if err != nil {
		fmt.Println("Unable to initialize radio:", err)
		os.Exit(-1)
	}

	lcli := localCli{
		radio:   lr,
		cliCmds: cli.PopulateCliCmds(),
	}

	evPS := pubsub.New(10)

	wg := sync.WaitGroup{}

	// SystemEvents
	wg.Add(1)

	go events.WatchSystemEvents(evPS, &wg)
	go events.CaptureKeyboard(evPS)

	prepareShutdownCh := evPS.Sub(events.PrepareShutdown)
	shutdownCh := evPS.Sub(events.Shutdown)
	cliInputCh := evPS.Sub(events.CliInput)

	fmt.Println()
	fmt.Printf("rig command: ")

	for {
		select {
		case <-prepareShutdownCh:
			evPS.Pub(true, events.Shutdown)

		case <-shutdownCh:
			exitTicker := time.NewTicker(time.Second)
			go func() {
				<-exitTicker.C
				os.Exit(-1)
			}()
			wg.Wait()
			os.Exit(0)

		case msg := <-cliInputCh:
			lcli.parseCli(logger, msg.([]string))
		}
	}
}

func (lcli *localCli) parseCli(logger *log.Logger, cliInput []string) {

	found := false

	if len(cliInput) == 0 {
		fmt.Printf("rig command: ")
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
		fmt.Println("unknown command")
	}

	fmt.Println()
	fmt.Printf("rig command: ")
}

func (lcli *localCli) PrintHelp(log *log.Logger) {

	buf := bytes.Buffer{}

	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Command", "Shortcut", "Parameter"})
	table.SetCenterSeparator("|")
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(50)

	for _, el := range lcli.cliCmds {
		table.Append([]string{el.Name, el.Shortcut, el.Parameters})
	}

	table.Render()

	lines := strings.Split(buf.String(), "\n")

	for _, line := range lines {
		log.Println(line)
	}
}
