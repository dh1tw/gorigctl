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
	"github.com/dh1tw/gorigctl/comms"
	"github.com/dh1tw/gorigctl/events"
	"github.com/dh1tw/gorigctl/gui"
	"github.com/dh1tw/gorigctl/remoteradio"
	"github.com/dh1tw/gorigctl/server"
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
	guiLocalCmd.Flags().IntP("rig-model", "m", 1, "Hamlib Rig Model ID")
	guiLocalCmd.Flags().IntP("baudrate", "b", 38400, "Baudrate")
	guiLocalCmd.Flags().StringP("portname", "o", "/dev/mhux/cat", "Portname / Device path")
	guiLocalCmd.Flags().IntP("databits", "d", 8, "Databits")
	guiLocalCmd.Flags().IntP("stopbits", "s", 1, "Stopbits")
	guiLocalCmd.Flags().StringP("parity", "r", "none", "Parity")
	guiLocalCmd.Flags().StringP("handshake", "a", "none", "Handshake")
	guiLocalCmd.Flags().DurationP("polling-interval", "t", time.Duration(time.Millisecond*100), "Timer for polling the rig's meter values [ms] (0 = disabled)")
	guiLocalCmd.Flags().DurationP("sync-interval", "k", time.Duration(time.Second*3), "Timer for syncing all values with the rig [s] (0 = disabled)")

}

type localGui struct {
	cliCmds       []cli.CliCmd
	remoteCliCmds []remoteradio.RemoteCliCmd
	radio         remoteradio.RemoteRadio
	logger        *log.Logger
}

func runLocalGui(cmd *cobra.Command, args []string) {

	// profiling server can be enabled through a hidden pflag
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6061", nil))
	// }()

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
	viper.BindPFlag("radio.polling-interval", cmd.Flags().Lookup("polling-interval"))
	viper.BindPFlag("radio.sync-interval", cmd.Flags().Lookup("sync-interval"))

	rigModel := viper.GetInt("radio.rig-model")
	debugLevel := 0 // off
	pollingInterval := viper.GetDuration("radio.polling-interval")
	syncInterval := viper.GetDuration("radio.sync-interval")

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

	toServerCh := make(chan comms.IOMsg, 1000)
	toClientCh := make(chan comms.IOMsg, 1000)
	fromClientCh := make(chan []byte, 1000)

	logger := utils.NewChLogger(evPS, events.AppLog, "")
	nullLogger := utils.NewNullLogger()

	userID := ""
	serverCatRequestTopic := "toServer"

	remRadio := remoteradio.NewRemoteRadio(serverCatRequestTopic, userID, toServerCh, logger, evPS)

	lGui := localGui{
		radio:         remRadio,
		cliCmds:       cli.PopulateCliCmds(),
		remoteCliCmds: remoteradio.GetRemoteCliCmds(),
	}

	wg := sync.WaitGroup{}

	rs := server.RadioSettings{
		RigModel:         rigModel,
		Port:             port,
		HlDebugLevel:     debugLevel,
		CatRequestCh:     fromClientCh,
		CatResponseTopic: "state",
		ToWireCh:         toClientCh,
		CapsTopic:        "caps",
		WaitGroup:        &wg,
		Events:           evPS,
		PollingInterval:  pollingInterval,
		SyncInterval:     syncInterval,
		RadioLogger:      logger,
		AppLogger:        nullLogger,
	}

	wg.Add(1) // radioServer

	// prepareShutdownCh := evPS.Sub(events.PrepareShutdown)
	shutdownCh := evPS.Sub(events.Shutdown)
	cliInputCh := evPS.Sub(events.CliInput)
	loggingCh := evPS.Sub(events.AppLog)

	go server.StartRadioServer(rs)

	// give a few milliseconds to check if radio
	// produces an error before we initialize the gui
	time.Sleep(time.Millisecond * 300)

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	go gui.Loop(evPS)
	lGui.radio.SetOnline(true)
	ui.SendCustomEvt("/radio/status", true)

	for {
		select {
		// shutdown the application gracefully
		case <-shutdownCh:
			//force exit after 1 sec
			exitTimeout := time.NewTimer(time.Second)
			ui.Close()
			go func() {
				<-exitTimeout.C
				os.Exit(-1)
			}()
			os.Exit(0)

		case msg := <-toClientCh:
			ioMsg := comms.IOMsg(msg)
			switch ioMsg.Topic {
			case "state":
				if err := lGui.radio.DeserializeCatResponse(ioMsg.Data); err != nil {
					ui.SendCustomEvt("/log/msg", err.Error())
					continue
				}
				state, err := lGui.radio.GetState()
				if err != nil {
					ui.SendCustomEvt("/log/msg", err.Error())
					continue
				}
				ui.SendCustomEvt("/radio/state", state)

			case "caps":
				lGui.radio.DeserializeCaps(ioMsg.Data)
				caps, err := lGui.radio.GetCaps()
				if err != nil {
					ui.SendCustomEvt("/log/msg", err.Error())
					continue
				}
				ui.SendCustomEvt("/radio/caps", caps)
			}
		case msg := <-toServerCh:
			ioMsg := comms.IOMsg(msg)
			if ioMsg.Topic == "toServer" {
				fromClientCh <- ioMsg.Data
			}

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

func (lGui *localGui) parseCli(logger *log.Logger, cliInput []string) {

	found := false

	if len(cliInput) == 0 {
		return
	}

	for _, cmd := range lGui.cliCmds {
		if cmd.Name == cliInput[0] || cmd.Shortcut == cliInput[0] {
			cmd.Cmd(&lGui.radio, logger, cliInput[1:])
			found = true
		}
	}

	for _, cmd := range lGui.remoteCliCmds {
		if cmd.Name == cliInput[0] || cmd.Shortcut == cliInput[0] {
			cmd.Cmd(&lGui.radio, logger, cliInput[1:])
			found = true
		}
	}

	if cliInput[0] == "help" || cliInput[0] == "?" {
		lGui.printHelp(logger)
		found = true
	}

	if !found {
		logger.Println("unknown command")
	}
}

func (lGui *localGui) printHelp(log *log.Logger) {

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

	for _, el := range lGui.remoteCliCmds {
		table.Append([]string{el.Name, el.Shortcut, el.Parameters})
	}

	table.Render()

	lines := strings.Split(buf.String(), "\n")

	for _, line := range lines {
		log.Println(line)
	}
}
