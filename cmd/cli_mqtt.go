package cmd

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cskr/pubsub"
	"github.com/dh1tw/gorigctl/cliclient"
	"github.com/dh1tw/gorigctl/comms"
	"github.com/dh1tw/gorigctl/events"
	"github.com/dh1tw/gorigctl/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// mqttCmd represents the mqtt command
var clientMqttCmd = &cobra.Command{
	Use:   "mqtt",
	Short: "command line client which connects via MQTT to a remote radio",
	Long:  `command line client which connects via MQTT to a remote radio`,
	Run:   mqttCliClient,
}

func init() {
	cliCmd.AddCommand(clientMqttCmd)
	clientMqttCmd.Flags().StringP("broker-url", "u", "localhost", "Broker URL")
	clientMqttCmd.Flags().IntP("broker-port", "p", 1883, "Broker Port")
	clientMqttCmd.Flags().StringP("station", "X", "mystation", "Your station callsign")
	clientMqttCmd.Flags().StringP("radio", "Y", "myradio", "Radio ID")
}

func mqttCliClient(cmd *cobra.Command, args []string) {

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// bind the pflags to viper settings
	viper.BindPFlag("mqtt.broker_url", cmd.Flags().Lookup("broker-url"))
	viper.BindPFlag("mqtt.broker_port", cmd.Flags().Lookup("broker-port"))
	viper.BindPFlag("mqtt.station", cmd.Flags().Lookup("station"))
	viper.BindPFlag("mqtt.radio", cmd.Flags().Lookup("radio"))

	if viper.IsSet("general.user_id") {
		viper.Set("general.user_id", utils.RandStringRunes(5))
	} else {
		viper.Set("general.user_id", "unknown_"+utils.RandStringRunes(5))
	}

	mqttBrokerURL := viper.GetString("mqtt.broker_url")
	mqttBrokerPort := viper.GetInt("mqtt.broker_port")
	mqttClientID := viper.GetString("general.user_id")

	baseTopic := viper.GetString("mqtt.station") +
		"/radios/" + viper.GetString("mqtt.radio") +
		"/cat"

	serverCatRequestTopic := baseTopic + "/setstate"
	serverStatusTopic := baseTopic + "/status"
	//	serverPingTopic := baseTopic + "/ping"
	// errorTopic := baseTopic + "/error"

	// tx topics
	serverCatResponseTopic := baseTopic + "/state"
	serverCapsTopic := baseTopic + "/caps"
	serverPongTopic := baseTopic + "/pong"

	mqttRxTopics := []string{serverCatResponseTopic, serverCapsTopic, serverPongTopic, serverStatusTopic}

	toWireCh := make(chan comms.IOMsg, 20)
	toDeserializeCatResponseCh := make(chan []byte, 10)
	toDeserializePingResponseCh := make(chan []byte, 10)
	toDeserializeCapsCh := make(chan []byte, 5)
	toDeserializeStatusCh := make(chan []byte, 5)

	// Event PubSub
	evPS := pubsub.New(1)

	// WaitGroup to coordinate a graceful shutdown
	var wg sync.WaitGroup

	logger := utils.NewChLogger(evPS, events.AppLog, "")

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
		ToWire:                      toWireCh,
		Events:                      evPS,
		LastWill:                    nil,
		Logger:                      logger,
	}

	remoteRadioSettings := cliClient.RemoteRadioSettings{
		CatResponseCh:   toDeserializeCatResponseCh,
		RadioStatusCh:   toDeserializeStatusCh,
		CapabilitiesCh:  toDeserializeCapsCh,
		ToWireCh:        toWireCh,
		CatRequestTopic: serverCatRequestTopic,
		Events:          evPS,
		WaitGroup:       &wg,
	}

	wg.Add(3) //MQTT + RemoteRadio + SysEvents

	connectionStatusCh := evPS.Sub(events.MqttConnStatus)
	prepareShutdownCh := evPS.Sub(events.PrepareShutdown)
	shutdownCh := evPS.Sub(events.Shutdown)

	go events.WatchSystemEvents(evPS, &wg)
	go cliClient.HandleRemoteRadio(remoteRadioSettings)
	time.Sleep(200 * time.Millisecond)
	go comms.MqttClient(mqttSettings)
	go events.CaptureKeyboard(evPS)

	for {
		select {

		// CTRL-C has been pressed; let's prepare the shutdown
		case <-prepareShutdownCh:
			// advice that we are going offline
			time.Sleep(time.Millisecond * 200)
			evPS.Pub(true, events.Shutdown)

		// shutdown the application gracefully
		case <-shutdownCh:
			//force exit after 1 sec
			exitTicker := time.NewTicker(time.Second)
			go func() {
				<-exitTicker.C
				os.Exit(0)
			}()
			wg.Wait()
			os.Exit(0)

		case ev := <-connectionStatusCh:
			connStatus := ev.(int)
			if connStatus == comms.CONNECTED {
			}
		}
	}
}
