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
	"github.com/dh1tw/gorigctl/remoteradio"
	"github.com/dh1tw/gorigctl/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// mqttCmd represents the mqtt command
var clientMqttCmd = &cobra.Command{
	Use:   "mqtt",
	Short: "command line client which connects via MQTT to a remote radio",
	Long: `command line client which connects via MQTT to a remote radio
	
The MQTT Topics follow the Shackbus convention and must match on the
Server and the Client.

The parameters in "<>" can be set through flags or in the config file:
<station>/radios/<radio>/cat

`,
	Run: mqttCliClient,
}

func init() {
	cliCmd.AddCommand(clientMqttCmd)
	clientMqttCmd.Flags().StringP("broker-url", "u", "test.mosquitto.org", "MQTT Broker URL")
	clientMqttCmd.Flags().IntP("broker-port", "p", 1883, "MQTT Broker Port")
	clientMqttCmd.Flags().StringP("username", "U", "", "MQTT Username")
	clientMqttCmd.Flags().StringP("password", "P", "", "MQTT Password")
	clientMqttCmd.Flags().StringP("client-id", "C", "gorigctl-cli", "MQTT ClientID")
	clientMqttCmd.Flags().StringP("station", "X", "mystation", "remote station callsign")
	clientMqttCmd.Flags().StringP("radio", "Y", "myradio", "Radio ID")
}

type remoteCli struct {
	cliCmds       []cli.CliCmd
	remoteCliCmds []remoteradio.RemoteCliCmd
	radio         remoteradio.RemoteRadio
}

func mqttCliClient(cmd *cobra.Command, args []string) {

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// bind the pflags to viper settings
	viper.BindPFlag("mqtt.broker-url", cmd.Flags().Lookup("broker-url"))
	viper.BindPFlag("mqtt.broker-port", cmd.Flags().Lookup("broker-port"))
	viper.BindPFlag("mqtt.station", cmd.Flags().Lookup("station"))
	viper.BindPFlag("mqtt.radio", cmd.Flags().Lookup("radio"))
	viper.BindPFlag("mqtt.username", cmd.Flags().Lookup("username"))
	viper.BindPFlag("mqtt.password", cmd.Flags().Lookup("password"))
	viper.BindPFlag("mqtt.client-id", cmd.Flags().Lookup("client-id"))

	mqttBrokerURL := viper.GetString("mqtt.broker-url")
	mqttBrokerPort := viper.GetInt("mqtt.broker-port")
	mqttUsername := viper.GetString("mqtt.username")
	mqttPassword := viper.GetString("mqtt.password")
	mqttClientID := viper.GetString("mqtt.client-id")

	if mqttClientID == "gorigctl-cli" {
		mqttClientID = mqttClientID + "-" + utils.RandStringRunes(5)
	}

	baseTopic := viper.GetString("mqtt.station") +
		"/radios/" + viper.GetString("mqtt.radio") +
		"/cat"

	serverCatRequestTopic := baseTopic + "/setstate"
	serverStatusTopic := baseTopic + "/status"

	// tx topics
	serverCatResponseTopic := baseTopic + "/state"
	serverCapsTopic := baseTopic + "/caps"

	mqttRxTopics := []string{serverCatResponseTopic, serverCapsTopic, serverStatusTopic}

	toWireCh := make(chan comms.IOMsg, 20)
	toDeserializeCatResponseCh := make(chan []byte, 10)
	toDeserializePingResponseCh := make(chan []byte, 10)
	toDeserializeCapsCh := make(chan []byte, 5)
	toDeserializeStatusCh := make(chan []byte, 5)

	// Event PubSub
	evPS := pubsub.New(1)

	// WaitGroup to coordinate a graceful shutdown
	var wg sync.WaitGroup

	// logger := utils.NewChLogger(evPS, events.AppLog, "")
	logger := utils.NewStdLogger("", 0)

	mqttSettings := comms.MqttSettings{
		WaitGroup:  &wg,
		Transport:  "tcp",
		BrokerURL:  mqttBrokerURL,
		BrokerPort: mqttBrokerPort,
		ClientID:   mqttClientID,
		Username:   mqttUsername,
		Password:   mqttPassword,
		Topics:     mqttRxTopics,
		ToDeserializeCatResponseCh:  toDeserializeCatResponseCh,
		ToDeserializeCatRequestCh:   toDeserializePingResponseCh,
		ToDeserializeCapabilitiesCh: toDeserializeCapsCh,
		ToDeserializeStatusCh:       toDeserializeStatusCh,
		ToWire:                      toWireCh,
		Events:                      evPS,
		LastWill:                    nil,
		Logger:                      logger,
	}

	wg.Add(3) //MQTT + SysEvents

	connectionStatusCh := evPS.Sub(events.MqttConnStatus)
	prepareShutdownCh := evPS.Sub(events.PrepareShutdown)
	shutdownCh := evPS.Sub(events.Shutdown)
	cliInputCh := evPS.Sub(events.CliInput)
	radioOnlineCh := evPS.Sub(events.RadioOnline)

	rcli := remoteCli{}
	rcli.radio = remoteradio.NewRemoteRadio(serverCatRequestTopic, mqttClientID, toWireCh, logger, evPS)
	rcli.cliCmds = cli.PopulateCliCmds()
	rcli.remoteCliCmds = remoteradio.GetRemoteCliCmds()

	go events.WatchSystemEvents(evPS, &wg)
	time.Sleep(200 * time.Millisecond)
	go comms.MqttClient(mqttSettings)
	go events.CaptureKeyboard(evPS)

	for {
		select {

		// CTRL-C has been pressed; let's prepare the shutdown
		case <-prepareShutdownCh:
			// advice that we are going offline
			evPS.Pub(true, events.Shutdown)

		// shutdown the application gracefully
		case <-shutdownCh:
			//force exit after 1 sec
			exitTicker := time.NewTicker(time.Second)
			go func() {
				<-exitTicker.C
				os.Exit(-1)
			}()
			wg.Wait()
			os.Exit(0)

		case msg := <-toDeserializeCapsCh:
			if err := rcli.radio.DeserializeCaps(msg); err != nil {
				logger.Println(err)
			}

		case msg := <-toDeserializeCatResponseCh:
			if err := rcli.radio.DeserializeCatResponse(msg); err != nil {
				logger.Println(err)
			}

		case msg := <-toDeserializeStatusCh:
			if err := rcli.radio.DeserializeRadioStatus(msg); err != nil {
				logger.Println(err)
			}

		case msg := <-cliInputCh:
			rcli.parseCli(logger, msg.([]string))

		case ev := <-connectionStatusCh:
			connStatus := ev.(int)
			if connStatus == comms.CONNECTED {

			}
		case ev := <-radioOnlineCh:
			radioOnline := ev.(bool)
			if radioOnline {
				fmt.Println("radio is online")
				fmt.Println()
				fmt.Printf("rig command: ")
			} else {
				logger.Println("radio is offline")
			}
		}
	}
}

func (rcli *remoteCli) parseCli(logger *log.Logger, cliInput []string) {

	found := false

	if len(cliInput) == 0 {
		fmt.Printf("rig command: ")
		return
	}

	for _, cmd := range rcli.cliCmds {
		if cmd.Name == cliInput[0] || cmd.Shortcut == cliInput[0] {
			cmd.Cmd(&rcli.radio, logger, cliInput[1:])
			found = true
		}
	}

	for _, cmd := range rcli.remoteCliCmds {
		if cmd.Name == cliInput[0] || cmd.Shortcut == cliInput[0] {
			cmd.Cmd(&rcli.radio, logger, cliInput[1:])
			found = true
		}
	}

	if cliInput[0] == "help" || cliInput[0] == "?" {
		rcli.PrintHelp(logger)
		found = true
	}

	if !found {
		fmt.Println("unknown command")
	}

	fmt.Println()
	fmt.Printf("rig command: ")
}

func (rcli *remoteCli) PrintHelp(log *log.Logger) {

	buf := bytes.Buffer{}

	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Command", "Shortcut", "Parameter"})
	table.SetCenterSeparator("|")
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(50)

	for _, el := range rcli.cliCmds {
		table.Append([]string{el.Name, el.Shortcut, el.Parameters})
	}

	for _, el := range rcli.remoteCliCmds {
		table.Append([]string{el.Name, el.Shortcut, el.Parameters})
	}

	table.Render()

	lines := strings.Split(buf.String(), "\n")

	for _, line := range lines {
		log.Println(line)
	}
}
