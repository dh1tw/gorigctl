package serverstatus

import (
	"log"
	"sync"

	"github.com/cskr/pubsub"
	sbStatus "github.com/dh1tw/gorigctl/sb_status"
	"github.com/dh1tw/remoteAudio/events"
)

type Settings struct {
	ServerStatusCh chan []byte
	Events         *pubsub.PubSub
	Waitgroup      *sync.WaitGroup
	Logger         *log.Logger
}

func MonitorServerStatus(s Settings) {

	defer s.Waitgroup.Done()

	shutdownCh := s.Events.Sub(events.Shutdown)

	for {
		select {
		case msg := <-s.ServerStatusCh:
			status := sbStatus.Status{}
			err := status.Unmarshal(msg)
			if err != nil {
				s.Logger.Println("Unable to Unmarshal Server Status Msg", err.Error())
				break
			}
			s.Events.Pub(status.Online, events.ServerOnline)

		case <-shutdownCh:
			return
		}
	}
}
