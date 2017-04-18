package comms

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cskr/pubsub"
	"github.com/dh1tw/gorigctl/events"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttSettings struct {
	WaitGroup                   *sync.WaitGroup
	Transport                   string
	BrokerURL                   string
	BrokerPort                  int
	ClientID                    string
	Topics                      []string
	ToDeserializeCatRequestCh   chan []byte
	ToDeserializeCatResponseCh  chan []byte
	ToDeserializeCapabilitiesCh chan []byte
	ToDeserializeStatusCh       chan []byte
	ToDeserializePingRequestCh  chan []byte
	ToDeserializePingResponseCh chan []byte
	ToDeserializeLogCh          chan []byte
	ToWire                      chan IOMsg
	Events                      *pubsub.PubSub
	LastWill                    *LastWill
	Logger                      *log.Logger
}

// LastWill defines the LastWill for MQTT. The LastWill will be
// submitted to the broker on connection and will be published
// on Disconnect.
type LastWill struct {
	Topic  string
	Data   []byte
	Qos    byte
	Retain bool
}

// IOMsg is a struct used internally which either originates from or
// will be send to the wire
type IOMsg struct {
	Data       []byte
	Raw        []float32
	Topic      string
	Retain     bool
	Qos        byte
	MQTTts     time.Time
	EnqueuedTs time.Time
}

const (
	DISCONNECTED = 0
	CONNECTED    = 1
)

func MqttClient(s MqttSettings) {

	defer s.WaitGroup.Done()

	// mqtt.DEBUG = log.New(os.Stderr, "DEBUG - ", log.LstdFlags)
	// mqtt.CRITICAL = log.New(os.Stderr, "CRITICAL - ", log.LstdFlags)
	// mqtt.WARN = log.New(os.Stderr, "WARN - ", log.LstdFlags)
	// mqtt.ERROR = log.New(os.Stderr, "ERROR - ", log.LstdFlags)

	shutdownCh := s.Events.Sub(events.Shutdown)

	var msgHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

		if strings.Contains(msg.Topic(), "cat/setstate") {

			s.ToDeserializeCatRequestCh <- msg.Payload()[:len(msg.Payload())]

		} else if strings.Contains(msg.Topic(), "cat/state") {

			s.ToDeserializeCatResponseCh <- msg.Payload()[:len(msg.Payload())]

		} else if strings.Contains(msg.Topic(), "cat/caps") {

			s.ToDeserializeCapabilitiesCh <- msg.Payload()[:len(msg.Payload())]

		} else if strings.Contains(msg.Topic(), "cat/status") {

			s.ToDeserializeStatusCh <- msg.Payload()[:len(msg.Payload())]

		} else if strings.Contains(msg.Topic(), "cat/ping") {

			s.ToDeserializePingRequestCh <- msg.Payload()[:len(msg.Payload())]

		} else if strings.Contains(msg.Topic(), "cat/pong") {

			s.ToDeserializePingResponseCh <- msg.Payload()[:len(msg.Payload())]

		} else if strings.Contains(msg.Topic(), "cat/log") {

			s.ToDeserializeLogCh <- msg.Payload()[:len(msg.Payload())]
		}

	}

	var connectionLostHandler = func(client mqtt.Client, err error) {
		s.Logger.Println("Connection lost to MQTT Broker; Reason:", err)
		s.Events.Pub(DISCONNECTED, events.MqttConnStatus)
	}

	// since we use SetCleanSession we have to subscribe on each
	// connect or reconnect to the channels
	var onConnectHandler = func(client mqtt.Client) {
		s.Logger.Printf("Connected to MQTT Broker %s:%d\n", s.BrokerURL, s.BrokerPort)

		// Subscribe to Task Topics
		for _, topic := range s.Topics {
			if token := client.Subscribe(topic, 0, nil); token.Wait() &&
				token.Error() != nil {
				s.Logger.Println(token.Error())
			}
		}
		s.Events.Pub(CONNECTED, events.MqttConnStatus)
	}

	opts := mqtt.NewClientOptions().AddBroker(s.Transport + "://" + s.BrokerURL + ":" + strconv.Itoa(s.BrokerPort))
	opts.SetClientID(s.ClientID)
	opts.SetDefaultPublishHandler(msgHandler)
	opts.SetKeepAlive(time.Second * 5)
	opts.SetMaxReconnectInterval(time.Second)
	opts.SetCleanSession(true)
	opts.SetOnConnectHandler(onConnectHandler)
	opts.SetConnectionLostHandler(connectionLostHandler)
	opts.SetAutoReconnect(true)

	if s.LastWill != nil {
		opts.SetBinaryWill(s.LastWill.Topic, s.LastWill.Data, s.LastWill.Qos, s.LastWill.Retain)
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		s.Logger.Println(token.Error())
	}

	for {
		select {
		case <-shutdownCh:
			s.Logger.Println("Disconnecting from MQTT Broker")
			if client.IsConnected() {
				client.Disconnect(0)
			}
			return
		case msg := <-s.ToWire:
			token := client.Publish(msg.Topic, msg.Qos, msg.Retain, msg.Data)
			token.WaitTimeout(time.Millisecond * 100)
			token.Wait()

			//indicates if cat data should be forwarded for decoding
			// case ev := <-forwardCatRequestCh:
			// 	forwardCat = ev.(bool)
		}
	}
}
