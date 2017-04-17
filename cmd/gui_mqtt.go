package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cskr/pubsub"
	"github.com/dh1tw/gorigctl/comms"
	"github.com/dh1tw/gorigctl/events"
	"github.com/dh1tw/gorigctl/gui"
	"github.com/dh1tw/gorigctl/ping"
	"github.com/dh1tw/gorigctl/utils"
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

	guiSettings := gui.GuiSettings{
		CatResponseCh:   toDeserializeCatResponseCh,
		RadioStatusCh:   toDeserializeStatusCh,
		CapabilitiesCh:  toDeserializeCapsCh,
		RadioLogCh:      toDeserializeLogCh,
		ToWireCh:        toWireCh,
		CatRequestTopic: serverCatRequestTopic,
		Events:          evPS,
		WaitGroup:       &wg,
		UserID:          userID,
	}

	// serverStatusSettings := serverstatus.Settings{
	// 	Waitgroup:      &wg,
	// 	ServerStatusCh: toDeserializeStatusCh,
	// 	Events:         evPS,
	// 	Logger:         appLogger,
	// }

	wg.Add(4) //MQTT + ping + cligui + MonitorServerStatus

	shutdownCh := evPS.Sub(events.Shutdown)

	go ping.CheckLatency(pingSettings)
	go gui.StartGui(guiSettings)
	// go serverstatus.MonitorServerStatus(serverStatusSettings)
	go time.Sleep(200 * time.Millisecond)
	go comms.MqttClient(mqttSettings)

	for {
		select {
		// shutdown the application gracefully
		case <-shutdownCh:
			//force exit after 1 sec
			exitTimeout := time.NewTimer(time.Second)
			go func() {
				<-exitTimeout.C
				os.Exit(0)
			}()
			wg.Wait()
			os.Exit(0)
		}
	}
}
