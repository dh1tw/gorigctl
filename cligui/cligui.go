package cligui

import (
	"log"
	"reflect"
	"sync"

	"github.com/cskr/pubsub"
	"github.com/dh1tw/gorigctl/comms"
	"github.com/dh1tw/gorigctl/events"
	sbRadio "github.com/dh1tw/gorigctl/sb_radio"
	sbStatus "github.com/dh1tw/gorigctl/sb_status"
	"github.com/dh1tw/gorigctl/utils"
	ui "github.com/gizak/termui"
	"github.com/spf13/viper"
)

type RemoteRadioSettings struct {
	CatResponseCh   chan []byte
	RadioStatusCh   chan []byte
	CatRequestTopic string
	PongCh          chan []int64
	ToWireCh        chan comms.IOMsg
	CapabilitiesCh  chan []byte
	WaitGroup       *sync.WaitGroup
	Events          *pubsub.PubSub
}

type remoteRadio struct {
	state           sbRadio.State
	newState        sbRadio.SetState
	caps            sbRadio.Capabilities
	settings        RemoteRadioSettings
	cliCmds         []cliCmd
	printRigUpdates bool
	userID          string
	radioOnline     bool
	logger          *log.Logger
}

type cliCmd struct {
	Cmd         func(r *remoteRadio, args []string)
	Name        string
	Shortcut    string
	Parameters  string
	Description string
	Example     string
}

func HandleRemoteRadio(rs RemoteRadioSettings) {
	defer rs.WaitGroup.Done()

	shutdownCh := rs.Events.Sub(events.Shutdown)

	r := remoteRadio{}
	r.state.Vfo = &sbRadio.Vfo{}
	r.state.Vfo.Functions = make(map[string]bool)
	r.state.Vfo.Levels = make(map[string]float32)
	r.state.Vfo.Parameters = make(map[string]float32)
	r.state.Vfo.Split = &sbRadio.Split{}

	r.settings = rs

	r.cliCmds = make([]cliCmd, 0, 30)
	r.populateCliCmds()

	if viper.IsSet("general.user_id") {
		r.userID = viper.GetString("general.user_id")
	} else {
		r.userID = "unknown_" + utils.RandStringRunes(5)
	}

	logger := utils.NewChLogger(rs.Events, events.AppLog, "")
	r.logger = logger

	loggingCh := rs.Events.Sub(events.AppLog)

	// rs.Events.Pub(true, events.ForwardCat)

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	cliInputCh := rs.Events.Sub(events.CliInput)
	pongCh := rs.Events.Sub(events.Pong)
	serverStatusCh := rs.Events.Sub(events.ServerOnline)

	go guiLoop(r.caps, r.settings.Events)

	for {
		select {
		case msg := <-rs.CapabilitiesCh:
			r.deserializeCaps(msg)
			ui.SendCustomEvt("/radio/caps", r.caps)

		case msg := <-rs.CatResponseCh:
			r.deserializeCatResponse(msg)
			ui.SendCustomEvt("/radio/state", r.state)

		case msg := <-rs.RadioStatusCh:
			r.deserializeRadioStatus(msg)

		case msg := <-cliInputCh:
			r.parseCli(msg.([]string))

		case msg := <-loggingCh:
			// forward to GUI event handler to be shown in the
			// approriate window
			ui.SendCustomEvt("/log/msg", msg)

		case msg := <-serverStatusCh:
			if msg.(bool) {
				logger.Println("Server Online")
			} else {
				logger.Println("Server Offline")
			}
			ui.SendCustomEvt("/radio/status", msg.(bool))

		case msg := <-pongCh:
			ui.SendCustomEvt("/network/latency", msg)

		case <-shutdownCh:
			log.Println("Disconnecting from Radio")
			return
		}
	}
}

func (r *remoteRadio) deserializeRadioStatus(data []byte) error {

	rStatus := sbStatus.Status{}
	if err := rStatus.Unmarshal(data); err != nil {
		return err
	}

	if r.radioOnline != rStatus.GetOnline() {
		r.radioOnline = rStatus.GetOnline()
		r.logger.Println("Radio Online:", r.radioOnline)
	}

	return nil
}

func (r *remoteRadio) sendCatRequest(req sbRadio.SetState) error {
	data, err := req.Marshal()
	if err != nil {
		return err
	}

	msg := comms.IOMsg{}
	msg.Data = data
	msg.Topic = r.settings.CatRequestTopic
	msg.Retain = false
	msg.Qos = 0

	r.settings.ToWireCh <- msg

	return nil
}

func (r *remoteRadio) deserializeCaps(msg []byte) error {

	caps := sbRadio.Capabilities{}
	err := caps.Unmarshal(msg)
	if err != nil {
		return err
	}

	r.caps = caps

	return nil
}

func (r *remoteRadio) deserializeCatResponse(msg []byte) error {

	ns := sbRadio.State{}
	err := ns.Unmarshal(msg)
	if err != nil {
		return err
	}

	if ns.CurrentVfo != r.state.CurrentVfo {
		r.state.CurrentVfo = ns.CurrentVfo
		if r.printRigUpdates {
			r.logger.Println("Updated Current Vfo:", r.state.CurrentVfo)
		}
	}

	if ns.Vfo != nil {

		if ns.Vfo.GetFrequency() != r.state.Vfo.Frequency {
			r.state.Vfo.Frequency = ns.Vfo.GetFrequency()
			if r.printRigUpdates {
				r.logger.Printf("Updated Frequency: %.3fkHz\n", r.state.Vfo.Frequency/1000)
			}
		}

		if ns.Vfo.GetMode() != r.state.Vfo.Mode {
			r.state.Vfo.Mode = ns.Vfo.GetMode()
			if r.printRigUpdates {
				r.logger.Println("Updated Mode:", r.state.Vfo.Mode)
			}
		}

		if ns.Vfo.GetPbWidth() != r.state.Vfo.PbWidth {
			r.state.Vfo.PbWidth = ns.Vfo.GetPbWidth()
			if r.printRigUpdates {
				r.logger.Printf("Updated Filter: %dHz\n", r.state.Vfo.PbWidth)
			}
		}

		if ns.Vfo.GetAnt() != r.state.Vfo.Ant {
			r.state.Vfo.Ant = ns.Vfo.GetAnt()
			if r.printRigUpdates {
				r.logger.Println("Updated Antenna:", r.state.Vfo.Ant)
			}
		}

		if ns.Vfo.GetRit() != r.state.Vfo.Rit {
			r.state.Vfo.Rit = ns.Vfo.GetRit()
			if r.printRigUpdates {
				r.logger.Printf("Updated Rit: %dHz\n", r.state.Vfo.Rit)
			}
		}

		if ns.Vfo.GetXit() != r.state.Vfo.Xit {
			r.state.Vfo.Xit = ns.Vfo.GetXit()
			if r.printRigUpdates {
				r.logger.Printf("Updated Xit: %dHz\n", r.state.Vfo.Xit)
			}
		}

		if ns.Vfo.GetSplit() != nil {
			if !reflect.DeepEqual(ns.Vfo.GetSplit(), r.state.Vfo.Split) {
				if err := r.updateSplit(ns.Vfo.Split); err != nil {
					log.Println(err)
				}
			}
		}

		if ns.Vfo.GetTuningStep() != r.state.Vfo.TuningStep {
			r.state.Vfo.TuningStep = ns.Vfo.GetTuningStep()
			if r.printRigUpdates {
				r.logger.Printf("Updated Tuning Step: %dHz\n", r.state.Vfo.TuningStep)
			}
		}

		if ns.Vfo.Functions != nil {
			if !reflect.DeepEqual(ns.Vfo.Functions, r.state.Vfo.Functions) {
				if err := r.updateFunctions(ns.Vfo.GetFunctions()); err != nil {
					log.Println(err)
				}
			}
		}

		if ns.Vfo.Levels != nil {
			if !reflect.DeepEqual(ns.Vfo.Levels, r.state.Vfo.Levels) {
				if err := r.updateLevels(ns.Vfo.GetLevels()); err != nil {
					log.Println(err)
				}
			}
		}

		if ns.Vfo.Parameters != nil {
			if !reflect.DeepEqual(ns.Vfo.Parameters, r.state.Vfo.Parameters) {
				if err := r.updateParams(ns.Vfo.GetParameters()); err != nil {
					log.Println(err)
				}
			}
		}

	}

	if ns.GetRadioOn() != r.state.RadioOn {
		r.state.RadioOn = ns.GetRadioOn()
		if r.printRigUpdates {
			r.logger.Println("Updated Radio Power On:", r.state.RadioOn)
		}
	}

	if ns.GetPtt() != r.state.Ptt {
		r.state.Ptt = ns.GetPtt()
		if r.printRigUpdates {
			r.logger.Println("Updated PTT On:", r.state.Ptt)
		}
	}

	if ns.GetPollingInterval() != r.state.PollingInterval {
		r.state.PollingInterval = ns.GetPollingInterval()
		if r.printRigUpdates {
			r.logger.Printf("Updated rig polling interval: %dms\n", r.state.PollingInterval)
		}
	}

	return nil
}

func (r *remoteRadio) updateSplit(newSplit *sbRadio.Split) error {

	if newSplit.GetEnabled() != r.state.Vfo.Split.Enabled {
		r.state.Vfo.Split.Enabled = newSplit.GetEnabled()
		if r.printRigUpdates {
			r.logger.Println("Updated Split Enabled:", r.state.Vfo.Split.Enabled)
		}
	}

	if newSplit.GetFrequency() != r.state.Vfo.Split.Frequency {
		r.state.Vfo.Split.Frequency = newSplit.GetFrequency()
		if r.printRigUpdates {
			r.logger.Printf("Updated TX (Split) Frequency: %.3fkHz\n", r.state.Vfo.Split.Frequency/1000)
		}
	}

	if newSplit.GetVfo() != r.state.Vfo.Split.Vfo {
		r.state.Vfo.Split.Vfo = newSplit.GetVfo()
		if r.printRigUpdates {
			r.logger.Println("Updated TX (Split) Vfo:", r.state.Vfo.Split.Vfo)
		}
	}

	if newSplit.GetMode() != r.state.Vfo.Split.Mode {
		r.state.Vfo.Split.Mode = newSplit.GetMode()
		if r.printRigUpdates {
			r.logger.Println("Updated TX (Split) Mode:", r.state.Vfo.Split.Mode)
		}
	}

	if newSplit.GetPbWidth() != r.state.Vfo.Split.PbWidth {

		r.state.Vfo.Split.PbWidth = newSplit.GetPbWidth()
		if r.printRigUpdates {
			r.logger.Printf("Split PbWidth: %dHz\n", r.state.Vfo.Split.PbWidth)
		}
	}

	return nil
}

func (r *remoteRadio) updateFunctions(newFuncs map[string]bool) error {

	r.state.Vfo.Functions = newFuncs

	// vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	// functions to be enabled
	// diff := utils.SliceDiff(newFuncs, r.state.Vfo.Functions)
	// for _, f := range diff {
	// 	funcValue, ok := hl.FuncValue[f]
	// 	if !ok {
	// 		return errors.New("unknown function")
	// 	}
	// 	// err := r.rig.SetFunc(vfo, funcValue, true)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// }

	// // functions to be disabled
	// diff = utils.SliceDiff(r.state.Vfo.Functions, newFuncs)
	// for _, f := range diff {
	// 	funcValue, ok := hl.FuncValue[f]
	// 	if !ok {
	// 		return errors.New("unknown function")
	// 	}
	// 	// err := r.rig.SetFunc(vfo, funcValue, false)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// }

	return nil
}

func (r *remoteRadio) updateLevels(newLevels map[string]float32) error {

	r.state.Vfo.Levels = newLevels

	// vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	// for k, v := range newLevels {
	// 	levelValue, ok := hl.LevelValue[k]
	// 	if !ok {
	// 		return errors.New("unknown Level")
	// 	}
	// 	if _, ok := r.state.Vfo.Levels[k]; !ok {
	// 		return errors.New("unsupported Level for this rig")
	// 	}

	// if r.state.Vfo.Levels[k] != v {
	// 	err := r.rig.SetLevel(vfo, levelValue, v)
	// 	if err != nil {
	// 		return nil
	// 	}
	// }
	// }

	return nil
}

func (r *remoteRadio) updateParams(newParams map[string]float32) error {

	r.state.Vfo.Parameters = newParams

	// vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	// for k, v := range newParams {
	// 	paramValue, ok := hl.ParmValue[k]
	// 	if !ok {
	// 		return errors.New("unknown Parameter")
	// 	}
	// if _, ok := r.state.Vfo.Parameters[k]; !ok {
	// 	return errors.New("unsupported Parameter for this rig")
	// }
	// if r.state.Vfo.Levels[k] != v {
	// 	err := r.rig.SetLevel(vfo, paramValue, v)
	// 	if err != nil {
	// 		return nil
	// 	}
	// }
	// }

	return nil
}

func (r *remoteRadio) initSetState() sbRadio.SetState {
	request := sbRadio.SetState{}

	request.CurrentVfo = r.state.CurrentVfo
	request.Vfo = &sbRadio.Vfo{}
	request.Vfo.Split = &sbRadio.Split{}
	request.Md = &sbRadio.MetaData{}
	request.UserId = r.userID

	return request
}
