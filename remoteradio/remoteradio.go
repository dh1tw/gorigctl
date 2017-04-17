package remoteradio

import (
	"log"

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
}

func NewRemoteRadio(topic, userID string, toWire chan comms.IOMsg, logger *log.Logger) RemoteRadio {
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

func (r *RemoteRadio) Print(v ...interface{}) {
	r.logger.Println(v)
}

func (r *RemoteRadio) Printf(format string, v ...interface{}) {
	r.logger.Printf(format, v)
}
