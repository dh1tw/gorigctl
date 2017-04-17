package radioServer

import (
	"log"
	"sync"
	"time"

	"github.com/cskr/pubsub"
	hl "github.com/dh1tw/goHamlib"
	"github.com/dh1tw/gorigctl/comms"
)

type RadioServerSettings struct {
	RigModel         int
	Port             hl.Port
	HlDebugLevel     int
	CatRequestCh     chan []byte
	ToWireCh         chan comms.IOMsg
	CatResponseTopic string
	CapsTopic        string
	WaitGroup        *sync.WaitGroup
	Events           *pubsub.PubSub
	PollingInterval  time.Duration
	SyncInterval     time.Duration
	RadioLogger      *log.Logger
	AppLogger        *log.Logger
}

func RunRadioServer() {

}
