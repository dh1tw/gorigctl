package remoteradio

import (
	"log"
	"strconv"

	"github.com/cskr/pubsub"
	"github.com/dh1tw/gorigctl/cli"
	"github.com/dh1tw/gorigctl/comms"
	sbRadio "github.com/dh1tw/gorigctl/sb_radio"
)

type RemoteRadio struct {
	state           sbRadio.State
	caps            sbRadio.Capabilities
	printRigUpdates bool
	userID          string
	radioOnline     bool
	logger          *log.Logger
	catRequestTopic string
	toWireCh        chan comms.IOMsg
	events          *pubsub.PubSub
}

type RemoteCliCmd struct {
	Cmd         func(r *RemoteRadio, log *log.Logger, args []string)
	Name        string
	Shortcut    string
	Parameters  string
	Description string
	Example     string
}

func NewRemoteRadio(topic, userID string, toWire chan comms.IOMsg, logger *log.Logger, events *pubsub.PubSub) RemoteRadio {
	r := RemoteRadio{}
	r.state = sbRadio.State{}
	r.state.Vfo = &sbRadio.Vfo{}
	r.state.Vfo.Split = &sbRadio.Split{}
	r.state.Vfo.Functions = make(map[string]bool)
	r.state.Vfo.Levels = make(map[string]float32)
	r.state.Vfo.Parameters = make(map[string]float32)
	r.state.Channel = &sbRadio.Channel{}
	r.caps = sbRadio.Capabilities{}
	r.toWireCh = toWire
	r.userID = userID
	r.catRequestTopic = topic
	r.logger = logger
	r.events = events

	return r
}

func (r *RemoteRadio) initSetState() sbRadio.SetState {
	request := sbRadio.SetState{}

	request.CurrentVfo = r.state.CurrentVfo
	request.Vfo = &sbRadio.Vfo{}
	request.Vfo.Split = &sbRadio.Split{}
	request.Md = &sbRadio.MetaData{}
	request.UserId = r.userID

	return request
}

func GetPollingInterval(r *RemoteRadio, log *log.Logger, args []string) {
	log.Printf("Rig polling interval: %dms\n", r.state.PollingInterval)
}

func SetPollingInterval(r *RemoteRadio, log *log.Logger, args []string) {
	if err := cli.CheckArgs(args, 1); err != nil {
		log.Println(err)
		return
	}

	req := r.initSetState()

	ur, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		log.Println("ERROR: polling interval must be integer [ms]")
		return
	}

	req.PollingInterval = int32(ur)
	req.Md.HasPollingInterval = true

	if err := r.sendCatRequest(req); err != nil {
		log.Println("ERROR:", err)
	}
}

func GetSyncInterval(r *RemoteRadio, log *log.Logger, args []string) {
	log.Printf("Rig sync interval: %ds\n", r.state.SyncInterval)
}

func SetSyncInterval(r *RemoteRadio, log *log.Logger, args []string) {
	if err := cli.CheckArgs(args, 1); err != nil {
		log.Println(err)
	}

	req := r.initSetState()

	ur, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		log.Println("ERROR: polling interval must be integer [s]")
		return
	}

	req.SyncInterval = int32(ur)
	req.Md.HasSyncInterval = true

	if err := r.sendCatRequest(req); err != nil {
		log.Println("ERROR:", err)
	}
}

func SetPrintRigUpdates(r *RemoteRadio, log *log.Logger, args []string) {
	if err := cli.CheckArgs(args, 1); err != nil {
		log.Println(err)
		return
	}

	ru, err := strconv.ParseBool(args[0])
	if err != nil {
		log.Println("ERROR: value must be of type bool (1,t,true / 0,f,false)")
		return
	}

	r.printRigUpdates = ru
}

func GetPrintRigUpdates(r *RemoteRadio, log *log.Logger, args []string) {
	log.Printf("Print rig updates: %v", r.printRigUpdates)
}

func GetRemoteCliCmds() []RemoteCliCmd {

	cliCmds := make([]RemoteCliCmd, 0, 40)

	cliSetPollingInterval := RemoteCliCmd{
		Cmd:         SetPollingInterval,
		Name:        "set_polling_interval",
		Shortcut:    "",
		Parameters:  "Polling rate [ms]",
		Description: "Set the polling interval for updating the meter values (SWR, ALC, Field Strength...)",
		Example:     "set_polling_interval 50",
	}

	cliCmds = append(cliCmds, cliSetPollingInterval)

	cliGetPollingInterval := RemoteCliCmd{
		Cmd:         GetPollingInterval,
		Name:        "get_polling_interval",
		Shortcut:    "",
		Description: "Get the polling interval for updating the meter values (SWR, ALC, Field Strength...)",
	}

	cliCmds = append(cliCmds, cliGetPollingInterval)

	cliGetSyncInterval := RemoteCliCmd{
		Cmd:         GetSyncInterval,
		Name:        "get_sync_interval",
		Shortcut:    "",
		Description: "Get the interval for synchronizing all radio values",
	}

	cliCmds = append(cliCmds, cliGetSyncInterval)

	cliSetSyncInterval := RemoteCliCmd{
		Cmd:         SetSyncInterval,
		Name:        "set_sync_interval",
		Shortcut:    "",
		Parameters:  "Sync rate [s]",
		Description: "Set the interval for synchronizing all radio values",
		Example:     "set_sync_interval 5",
	}

	cliCmds = append(cliCmds, cliSetSyncInterval)

	cliSetPrintUpdates := RemoteCliCmd{
		Cmd:         SetPrintRigUpdates,
		Name:        "set_print_rig_updates",
		Parameters:  "[true, t, 1, false, f, 0]",
		Shortcut:    "",
		Description: "Print rig values which have changed",
	}

	cliCmds = append(cliCmds, cliSetPrintUpdates)

	cliGetPrintUpdates := RemoteCliCmd{
		Cmd:         GetPrintRigUpdates,
		Name:        "get_print_rig_updates",
		Shortcut:    "",
		Description: "If rig values shall be printed when they have changed",
	}

	cliCmds = append(cliCmds, cliGetPrintUpdates)

	return cliCmds

}
