package gui

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/cskr/pubsub"
	"github.com/dh1tw/gorigctl/cli"
	"github.com/dh1tw/gorigctl/comms"
	"github.com/dh1tw/gorigctl/events"
	"github.com/dh1tw/gorigctl/remoteradio"
	sbLog "github.com/dh1tw/gorigctl/sb_log"
	"github.com/dh1tw/gorigctl/utils"
	ui "github.com/gizak/termui"
	"github.com/olekukonko/tablewriter"
)

type GuiSettings struct {
	CatResponseCh   chan []byte
	RadioStatusCh   chan []byte
	CatRequestTopic string
	RadioLogCh      chan []byte
	PongCh          chan []int64
	ToWireCh        chan comms.IOMsg
	CapabilitiesCh  chan []byte
	WaitGroup       *sync.WaitGroup
	Events          *pubsub.PubSub
	UserID          string
}

type gui struct {
	cliCmds       []cli.CliCmd
	remoteCliCmds []remoteradio.RemoteCliCmd
	radio         remoteradio.RemoteRadio
	settings      GuiSettings
	logger        *log.Logger
}

func StartGui(gs GuiSettings) {
	defer gs.WaitGroup.Done()

	shutdownCh := gs.Events.Sub(events.Shutdown)

	gui := gui{}

	logger := utils.NewChLogger(gs.Events, events.AppLog, "")
	gui.logger = logger

	gui.radio = remoteradio.NewRemoteRadio(gs.CatRequestTopic, gs.UserID, gs.ToWireCh, logger, gs.Events)
	gui.settings = gs
	gui.cliCmds = cli.PopulateCliCmds()
	gui.remoteCliCmds = remoteradio.GetRemoteCliCmds()

	loggingCh := gs.Events.Sub(events.AppLog)

	// rs.Events.Pub(true, events.ForwardCat)

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	cliInputCh := gs.Events.Sub(events.CliInput)
	pongCh := gs.Events.Sub(events.Pong)
	radioOnlineCh := gs.Events.Sub(events.RadioOnline)

	caps, _ := gui.radio.GetCaps()

	go guiLoop(caps, gs.Events)

	for {
		select {
		case msg := <-gs.CapabilitiesCh:
			gui.radio.DeserializeCaps(msg)
			caps, _ := gui.radio.GetCaps()
			ui.SendCustomEvt("/radio/caps", caps)

		case msg := <-gs.CatResponseCh:
			// r.printRigUpdates = true
			err := gui.radio.DeserializeCatResponse(msg)
			if err != nil {
				ui.SendCustomEvt("/log/msg", err.Error())
			}
			state, _ := gui.radio.GetState()
			ui.SendCustomEvt("/radio/state", state)

		case msg := <-gs.RadioStatusCh:
			gui.radio.DeserializeRadioStatus(msg)

		case msg := <-cliInputCh:
			gui.parseCli(msg.([]string))

		case msg := <-gs.RadioLogCh:
			deserializeRadioLogMsg(msg)

		case msg := <-loggingCh:
			// forward to GUI event handler to be shown in the
			// approriate window
			ui.SendCustomEvt("/log/msg", msg)

		case msg := <-radioOnlineCh:
			if msg.(bool) {
				logger.Println("radio is online")
			} else {
				logger.Println("radio is offline")
			}
			ui.SendCustomEvt("/radio/status", msg.(bool))

		case msg := <-pongCh:
			ui.SendCustomEvt("/network/latency", msg)

		case <-shutdownCh:
			log.Println("disconnecting from radio")
			return
		}
	}
}

func deserializeRadioLogMsg(ba []byte) {

	radioLogMsg := sbLog.LogMsg{}
	err := radioLogMsg.Unmarshal(ba)
	if err != nil {
		fmt.Println("could not unmarshal radio log message")
		return
	}

	ui.SendCustomEvt("/log/msg", radioLogMsg.Msg)
}

func (rg *gui) parseCli(cliInput []string) {

	found := false
	for _, cmd := range rg.cliCmds {
		if cmd.Name == cliInput[0] || cmd.Shortcut == cliInput[0] {
			cmd.Cmd(&rg.radio, rg.logger, cliInput[1:])
			found = true
		}
	}

	for _, cmd := range rg.remoteCliCmds {
		if cmd.Name == cliInput[0] || cmd.Shortcut == cliInput[0] {
			cmd.Cmd(&rg.radio, rg.logger, cliInput[1:])
			found = true
		}
	}

	if cliInput[0] == "help" || cliInput[0] == "?" {
		rg.PrintHelp(rg.logger)
		found = true
	}

	if !found {
		rg.logger.Println("unknown command")
	}
}

func (rg *gui) PrintHelp(log *log.Logger) {

	buf := bytes.Buffer{}

	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Command", "Shortcut", "Parameter"})
	table.SetCenterSeparator("|")
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(50)

	for _, el := range rg.cliCmds {
		table.Append([]string{el.Name, el.Shortcut, el.Parameters})
	}

	for _, el := range rg.remoteCliCmds {
		table.Append([]string{el.Name, el.Shortcut, el.Parameters})
	}

	table.Render()

	lines := strings.Split(buf.String(), "\n")

	for _, line := range lines {
		log.Println(line)
	}
}
