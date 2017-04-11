package ping

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cskr/pubsub"
	"github.com/dh1tw/gorigctl/comms"
	"github.com/dh1tw/gorigctl/events"
	sbPing "github.com/dh1tw/gorigctl/sb_ping"
)

type Settings struct {
	ToWireCh  chan comms.IOMsg
	PongCh    chan []byte
	PingCh    chan []byte
	PingTopic string
	PongTopic string
	UserID    string
	WaitGroup *sync.WaitGroup
	Events    *pubsub.PubSub
}

// CheckLatency sends out a ping every second to the server
// to determine the system latency.  This Function is
// typically executed as a goroutine in client applications
func CheckLatency(ps Settings) {

	defer ps.WaitGroup.Done()

	shutdownCh := ps.Events.Sub(events.Shutdown)

	connectionStatusCh := ps.Events.Sub(events.MqttConnStatus)

	connectionStatus := comms.DISCONNECTED

	pingTicker := time.NewTicker(time.Second)

	for {
		select {
		case <-shutdownCh:
			return

		case <-pingTicker.C:
			if connectionStatus == comms.CONNECTED {
				sendPing(ps.UserID, ps.PingTopic, ps.ToWireCh)
			}

		case msg := <-ps.PongCh:
			pong, err := deserializePong(msg, ps.UserID)
			if err == nil {
				ps.Events.Pub(pong, events.Pong)
			}

		case ev := <-connectionStatusCh:
			connectionStatus = ev.(int)
		}
	}
}

// EchoPing receives a Ping Request and sends it back (Pong). This Function is
// typically executed as a goroutine on server applications
func EchoPing(ps Settings) {

	defer ps.WaitGroup.Done()

	shutdownCh := ps.Events.Sub(events.Shutdown)

	for {
		select {
		case <-shutdownCh:
			return

		case msg := <-ps.PongCh:
			pong := comms.IOMsg{}
			pong.Data = msg
			pong.Topic = ps.PongTopic
			ps.ToWireCh <- pong
		}
	}
}

func sendPing(userID, topic string, toWireCh chan comms.IOMsg) {
	now := time.Now().UnixNano()

	req := sbPing.Ping{}
	req.UserId = userID
	req.Timestamp = now

	data, err := req.Marshal()
	if err != nil {
		fmt.Println(err)
	} else {
		wireMsg := comms.IOMsg{
			Topic: topic,
			Data:  data,
		}
		toWireCh <- wireMsg
	}
}

// deserialize Pong (Ping reply) and return the passed Duration (in Nanoseconds)
func deserializePong(msg []byte, myUserID string) (int64, error) {
	pong := sbPing.Ping{}
	err := pong.Unmarshal(msg)
	if err != nil {
		return 0, err
	}

	if myUserID != pong.UserId {
		return 0, errors.New("not determined for this user")
	}

	pingTimestamp := time.Unix(0, pong.Timestamp)
	delta := time.Since(pingTimestamp)
	return delta.Nanoseconds(), nil
}
