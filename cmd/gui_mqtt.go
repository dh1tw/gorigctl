package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cskr/pubsub"
	"github.com/dh1tw/gorigctl/cli"
	"github.com/dh1tw/gorigctl/comms"
	"github.com/dh1tw/gorigctl/events"
	"github.com/dh1tw/gorigctl/gui"
	"github.com/dh1tw/gorigctl/ping"
	"github.com/dh1tw/gorigctl/remoteradio"
	sbLog "github.com/dh1tw/gorigctl/sb_log"
	"github.com/dh1tw/gorigctl/utils"
	ui "github.com/gizak/termui"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"net/http"
	_ "net/http/pprof"
)

// mqttCmd represents the mqtt command
var guiMqttCmd = &cobra.Command{
	Use:   "mqtt",
	Short: "GUI client which connects via MQTT to a remote Radio",
	Long: `GUI client which connects via MQTT to a remote Radio
	
The MQTT Topics follow the Shackbus convention and must match on the
Server and the Client.

The parameters in "<>" can be set through flags or in the config file:
<station>/radios/<radio>/cat

`,
	Run: guiCliClient,
}

func init() {
	guiCmd.AddCommand(guiMqttCmd)
	guiMqttCmd.Flags().StringP("broker-url", "u", "localhost", "Broker URL")
	guiMqttCmd.Flags().IntP("broker-port", "p", 1883, "Broker Port")
	guiMqttCmd.Flags().StringP("station", "X", "mystation", "Your station callsign")
	guiMqttCmd.Flags().StringP("radio", "Y", "myradio", "Radio ID")
}

type remoteGui struct {
	cliCmds       []cli.CliCmd
	remoteCliCmds []remoteradio.RemoteCliCmd
	radio         remoteradio.RemoteRadio
	logger        *log.Logger
}

func guiCliClient(cmd *cobra.Command, args []string) {

	// profiling server can be enabled through a hidden pflag
	go func() {
		log.Println(http.ListenAndServe("localhost:6061", nil))
	}()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// bind the pflags to viper settings
	viper.BindPFlag("mqtt.broker_url", cmd.Flags().Lookup("broker-url"))
	viper.BindPFlag("mqtt.broker_port", cmd.Flags().Lookup("broker-port"))
	viper.BindPFlag("mqtt.station", cmd.Flags().Lookup("station"))
	viper.BindPFlag("mqtt.radio", cmd.Flags().Lookup("radio"))

	userID := "unknown_" + utils.RandStringRunes(5)
	mqttClientID := "unknown_" + utils.RandStringRunes(5)

	if viper.IsSet("general.user_id") {
		userID = viper.GetString("general.user_id")
		mqttClientID = viper.GetString("general.user_id") + utils.RandStringRunes(5)
	}

	mqttBrokerURL := viper.GetString("mqtt.broker_url")
	mqttBrokerPort := viper.GetInt("mqtt.broker_port")

	baseTopic := viper.GetString("mqtt.station") +
		"/radios/" + viper.GetString("mqtt.radio") +
		"/cat"

	serverCatRequestTopic := baseTopic + "/setstate"
	serverStatusTopic := baseTopic + "/status"
	serverPingTopic := baseTopic + "/ping"
	serverLogTopic := baseTopic + "/log"

	// tx topics
	serverCatResponseTopic := baseTopic + "/state"
	serverCapsTopic := baseTopic + "/caps"
	serverPongTopic := baseTopic + "/pong"

	mqttRxTopics := []string{
		serverCatResponseTopic,
		serverCapsTopic,
		serverPongTopic,
		serverStatusTopic,
		serverLogTopic,
	}

	toWireCh := make(chan comms.IOMsg, 20)
	toDeserializeCatResponseCh := make(chan []byte, 50)
	toDeserializePingResponseCh := make(chan []byte, 50)
	toDeserializeCapsCh := make(chan []byte, 5)
	toDeserializeStatusCh := make(chan []byte, 5)
	toDeserializeLogCh := make(chan []byte, 10)

	// Event PubSub
	evPS := pubsub.New(10000)

	// WaitGroup to coordinate a graceful shutdown
	var wg sync.WaitGroup

	pingSettings := ping.Settings{
		ToWireCh:  toWireCh,
		PingTopic: serverPingTopic,
		PongCh:    toDeserializePingResponseCh,
		UserID:    userID,
		WaitGroup: &wg,
		Events:    evPS,
	}

	appLogger := utils.NewChLogger(evPS, events.AppLog, "")

	mqttSettings := comms.MqttSettings{
		WaitGroup:  &wg,
		Transport:  "tcp",
		BrokerURL:  mqttBrokerURL,
		BrokerPort: mqttBrokerPort,
		ClientID:   mqttClientID,
		Topics:     mqttRxTopics,
		ToDeserializeCatResponseCh:  toDeserializeCatResponseCh,
		ToDeserializeCatRequestCh:   toDeserializePingResponseCh,
		ToDeserializeCapabilitiesCh: toDeserializeCapsCh,
		ToDeserializeStatusCh:       toDeserializeStatusCh,
		ToDeserializePingResponseCh: toDeserializePingResponseCh,
		ToDeserializeLogCh:          toDeserializeLogCh,
		ToWire:                      toWireCh,
		Events:                      evPS,
		LastWill:                    nil,
		Logger:                      appLogger,
	}

	wg.Add(2) //MQTT + ping

	rGui := remoteGui{}

	shutdownCh := evPS.Sub(events.Shutdown)
	cliInputCh := evPS.Sub(events.CliInput)
	pongCh := evPS.Sub(events.Pong)
	radioOnlineCh := evPS.Sub(events.RadioOnline)
	loggingCh := evPS.Sub(events.AppLog)

	logger := utils.NewChLogger(evPS, events.AppLog, "")
	rGui.logger = logger

	rGui.radio = remoteradio.NewRemoteRadio(serverCatRequestTopic, userID, toWireCh, logger, evPS)
	rGui.cliCmds = cli.PopulateCliCmds()
	rGui.remoteCliCmds = remoteradio.GetRemoteCliCmds()
	rGui.logger = logger

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	go ping.CheckLatency(pingSettings)
	time.Sleep(200 * time.Millisecond)
	go comms.MqttClient(mqttSettings)
	go gui.Loop(evPS)

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
			wg.Wait()
			os.Exit(0)
		case msg := <-toDeserializeCapsCh:
			rGui.radio.DeserializeCaps(msg)
			caps, _ := rGui.radio.GetCaps()
			ui.SendCustomEvt("/radio/caps", caps)

		case msg := <-toDeserializeCatResponseCh:
			// r.printRigUpdates = true
			err := rGui.radio.DeserializeCatResponse(msg)
			if err != nil {
				ui.SendCustomEvt("/log/msg", err.Error())
			}
			state, _ := rGui.radio.GetState()
			ui.SendCustomEvt("/radio/state", state)

		case msg := <-toDeserializeStatusCh:
			rGui.radio.DeserializeRadioStatus(msg)

		case msg := <-cliInputCh:
			rGui.parseCli(msg.([]string))

		case msg := <-toDeserializeLogCh:
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

func (rGui *remoteGui) parseCli(cliInput []string) {

	found := false
	for _, cmd := range rGui.cliCmds {
		if cmd.Name == cliInput[0] || cmd.Shortcut == cliInput[0] {
			cmd.Cmd(&rGui.radio, rGui.logger, cliInput[1:])
			found = true
		}
	}

	for _, cmd := range rGui.remoteCliCmds {
		if cmd.Name == cliInput[0] || cmd.Shortcut == cliInput[0] {
			cmd.Cmd(&rGui.radio, rGui.logger, cliInput[1:])
			found = true
		}
	}

	if cliInput[0] == "help" || cliInput[0] == "?" {
		rGui.PrintHelp(rGui.logger)
		found = true
	}

	if !found {
		rGui.logger.Println("unknown command")
	}
}

func (rGui *remoteGui) PrintHelp(log *log.Logger) {

	buf := bytes.Buffer{}

	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Command", "Shortcut", "Parameter"})
	table.SetCenterSeparator("|")
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(50)

	for _, el := range rGui.cliCmds {
		table.Append([]string{el.Name, el.Shortcut, el.Parameters})
	}

	for _, el := range rGui.remoteCliCmds {
		table.Append([]string{el.Name, el.Shortcut, el.Parameters})
	}

	table.Render()

	lines := strings.Split(buf.String(), "\n")

	for _, line := range lines {
		rGui.logger.Println(line)
	}
}
