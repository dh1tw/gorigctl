package events

import (
	"os"
	"os/signal"

	"sync"

	"github.com/cskr/pubsub"
)

// Event channel names used for event Pubsub

// internal
const (
	MqttConnStatus  = "mqttConnStatus" // int
	ForwardCat      = "forwardAudio"   //bool
	CliInput        = "cliInput"       // []string
	PrepareShutdown = "prepShutdown"   // no type
	Shutdown        = "shutdown"       // no type
	OsExit          = "osExit"         // bool
	AppLog          = "applog"         // string
	RadioLog        = "radiolog"       // string
	RadioOnline     = "radioOnline"    //bool
	Pong            = "pong"           // int64
)

func WatchSystemEvents(evPS *pubsub.PubSub, wg *sync.WaitGroup) {

	defer wg.Done()

	// Channel to handle OS signals
	osSignals := make(chan os.Signal, 1)

	//subscribe to os.Interrupt (CTRL-C signal)
	signal.Notify(osSignals, os.Interrupt)

	select {
	case osSignal := <-osSignals:
		if osSignal == os.Interrupt {
			evPS.Pub(true, PrepareShutdown)
			return
		}
	}
}
