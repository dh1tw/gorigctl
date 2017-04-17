package gui

import (
	"fmt"
	"log"
	"sync"

	"github.com/cskr/pubsub"
	"github.com/dh1tw/gorigctl/cli"
	"github.com/dh1tw/gorigctl/comms"
	"github.com/dh1tw/gorigctl/events"
	"github.com/dh1tw/gorigctl/remoteradio"
	sbLog "github.com/dh1tw/gorigctl/sb_log"
	"github.com/dh1tw/gorigctl/utils"
	ui "github.com/gizak/termui"
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
	cliCmds  []cli.CliCmd
	radio    remoteradio.RemoteRadio
	settings GuiSettings
	logger   *log.Logger
}

func StartGui(gs GuiSettings) {
	defer gs.WaitGroup.Done()

	shutdownCh := gs.Events.Sub(events.Shutdown)

	gui := gui{}

	logger := utils.NewChLogger(gs.Events, events.AppLog, "")
	gui.logger = logger

	gui.radio = remoteradio.NewRemoteRadio(gs.CatRequestTopic, gs.UserID, gs.ToWireCh, logger)
	gui.settings = gs
	gui.cliCmds = cli.PopulateCliCmds()

	loggingCh := gs.Events.Sub(events.AppLog)

	// rs.Events.Pub(true, events.ForwardCat)

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	cliInputCh := gs.Events.Sub(events.CliInput)
	pongCh := gs.Events.Sub(events.Pong)
	serverStatusCh := gs.Events.Sub(events.ServerOnline)

	go guiLoop(gui.radio.GetCaps(), gs.Events)

	for {
		select {
		case msg := <-gs.CapabilitiesCh:
			gui.radio.DeserializeCaps(msg)
			ui.SendCustomEvt("/radio/caps", gui.radio.GetCaps())

		case msg := <-gs.CatResponseCh:
			// r.printRigUpdates = true
			err := gui.radio.DeserializeCatResponse(msg)
			if err != nil {
				ui.SendCustomEvt("/log/msg", err.Error())
			}
			ui.SendCustomEvt("/radio/state", gui.radio.GetState())

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

		case msg := <-serverStatusCh:
			if msg.(bool) {
				logger.Println("Radio Online")
			} else {
				logger.Println("Radio Offline")
			}
			ui.SendCustomEvt("/radio/status", msg.(bool))

		case msg := <-pongCh:
			ui.SendCustomEvt("/network/latency", msg)

		case <-shutdownCh:
			log.Println("Disconnecting from Radio")
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
			cmd.Cmd(&rg.radio, cliInput[1:])
			found = true
		}
	}
	if !found {
		rg.radio.Print("unknown command")
	}
}
