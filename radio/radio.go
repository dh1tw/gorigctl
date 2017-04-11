package radio

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"time"

	"github.com/cskr/pubsub"
	hl "github.com/dh1tw/goHamlib"
	"github.com/dh1tw/gorigctl/comms"
	"github.com/dh1tw/gorigctl/events"
	sbRadio "github.com/dh1tw/gorigctl/sb_radio"
)

type RadioSettings struct {
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
}

type radio struct {
	rig           hl.Rig
	state         sbRadio.State
	settings      *RadioSettings
	pollingTicker *time.Ticker
}

func HandleRadio(rs RadioSettings) {

	defer rs.WaitGroup.Done()

	shutdownCh := rs.Events.Sub(events.Shutdown)

	r := radio{}
	r.rig = hl.Rig{}
	r.state = sbRadio.State{}
	r.state.Vfo = &sbRadio.Vfo{}
	r.state.Channel = &sbRadio.Channel{}
	r.settings = &rs

	r.state.PollingInterval = int32(r.settings.PollingInterval.Nanoseconds() / 1000000)

	err := r.rig.Init(rs.RigModel)
	if err != nil {
		log.Println(err)
		return
	}

	r.rig.SetDebugLevel(rs.HlDebugLevel)

	err = r.rig.SetPort(rs.Port)
	if err != nil {
		// if we can not set the port, we shut down
		log.Println(err)
		r.settings.Events.Pub(true, events.Shutdown)
		return
	}

	if rs.RigModel != 1 { // exception for Dummy Rig
		if err := r.rig.Open(); err != nil {
			// if we can not open the port, we shut down
			log.Println(err)
			r.settings.Events.Pub(true, events.Shutdown)
			return
		}
	}

	// publish the radio's capabilities
	if err := r.sendCaps(); err != nil {
		log.Println(err)
	}

	// // check if the radio is turned on and query its state
	// rigOn, err := r.rig.GetPowerStat()
	// if err != nil {
	// 	log.Println(err)
	// 	// we should check if the rig might no have the
	// 	// ability to be turn on/off through CapsTopic

	// 	// let's hope for the best and query it
	// 	if err := r.queryVfo(); err != nil {
	// 		log.Println(err)
	// 	}
	// } else {
	// 	// no error and the rig is on so we can query it
	// 	if rigOn == hl.RIG_POWER_ON {
	// 		if err := r.queryVfo(); err != nil {
	// 			log.Println(err)
	// 		}
	// 	} else {
	// 		r.state.RadioOn = false
	// 	}
	// }

	if err := r.queryVfo(); err != nil {
		fmt.Println(err)
	}

	// if the rig supports fast_commands, then we will use it
	token, err := r.rig.GetConf("fast_commands_token")
	if err != nil {
		log.Println(err)
	}
	if len(token) > 0 {
		err = r.rig.SetConf("fast_commands_token", "1")
		if err != nil {
			log.Println(err)
		}
	}

	// publish the radio's state
	if err := r.sendState(); err != nil {
		log.Println(err)
	}

	r.pollingTicker = time.NewTicker(r.settings.PollingInterval)

	for {
		select {
		case msg := <-rs.CatRequestCh:
			r.deserializeCatRequest(msg)
			r.sendState()

		case <-shutdownCh:
			log.Println("Disconnecting from Radio")
			// maybe we have to check if the connection is really open
			r.rig.Close()
			r.rig.Cleanup()
			return

		case <-r.pollingTicker.C:
			r.updateMeter()
		}
	}
}

func (r *radio) queryVfo() error {
	vfo, err := r.rig.GetVfo()
	if err != nil {
		return err
	}
	r.state.CurrentVfo = hl.VfoName[vfo]

	if pwrOn, err := r.rig.GetPowerStat(); err != nil {
		return err
	} else {
		if pwrOn == hl.RIG_POWER_ON {
			r.state.RadioOn = true
		} else {
			r.state.RadioOn = false
		}
	}

	freq, err := r.rig.GetFreq(vfo)
	if err != nil {
		return err
	}
	r.state.Vfo.Frequency = freq

	mode, pbWidth, err := r.rig.GetMode(vfo)
	if err != nil {
		return err
	}
	if modeName, ok := hl.ModeName[mode]; ok {
		r.state.Vfo.Mode = modeName
	} else {
		return errors.New("unknown mode")
	}

	r.state.Vfo.PbWidth = int32(pbWidth)

	ant, err := r.rig.GetAnt(vfo)
	if err != nil {
		return err
	}
	r.state.Vfo.Ant = int32(ant)

	rit, err := r.rig.GetRit(vfo)
	if err != nil {
		return err
	}
	r.state.Vfo.Rit = int32(rit)

	xit, err := r.rig.GetXit(vfo)
	if err != nil {
		return err
	}
	r.state.Vfo.Xit = int32(xit)

	split := sbRadio.Split{}

	splitOn, txVfo, err := r.rig.GetSplit(vfo)
	if err != nil {
		return err
	}

	if splitOn == hl.RIG_SPLIT_ON {
		split.Enabled = true
	} else {
		split.Enabled = false
	}

	if splitOn == hl.RIG_SPLIT_ON {

		txFreq, err := r.rig.GetSplitFreq(txVfo)
		if err != nil {
			return err
		}

		txMode, txPbWidth, err := r.rig.GetSplitMode(txVfo)
		if err != nil {
			return err
		}
		split.Frequency = txFreq
		if txVfoName, ok := hl.VfoName[txVfo]; ok {
			split.Vfo = txVfoName
		} else {
			return errors.New("unknown Vfo Name")
		}

		if txModeName, ok := hl.ModeName[txMode]; ok {
			split.Mode = txModeName
		} else {
			return errors.New("unknown Mode")
		}

		split.PbWidth = int32(txPbWidth)

	}

	r.state.Vfo.Split = &split

	tStep, err := r.rig.GetTs(vfo)
	if err != nil {
		return err
	}
	r.state.Vfo.TuningStep = int32(tStep)

	r.state.Vfo.Functions = make([]string, 0, len(hl.FuncName))

	for _, f := range r.rig.Caps.GetFunctions {
		fValue, err := r.rig.GetFunc(vfo, hl.FuncValue[f])
		if err != nil {
			return err
		}
		if fValue {
			r.state.Vfo.Functions = append(r.state.Vfo.Functions, f)
		}
	}

	r.state.Vfo.Levels = make(map[string]float32)
	for _, level := range r.rig.Caps.GetLevels {
		lValue, err := r.rig.GetLevel(vfo, hl.LevelValue[level.Name])
		if err != nil {
			// return err
			log.Println("Warning:", level.Name, "-", err)
		}
		r.state.Vfo.Levels[level.Name] = lValue
	}

	r.state.Vfo.Parameters = make(map[string]float32)
	for _, param := range r.rig.Caps.GetParameters {
		pValue, err := r.rig.GetParm(vfo, hl.ParmValue[param.Name])
		if err != nil {
			return err
		}
		r.state.Vfo.Parameters[param.Name] = pValue
	}

	return nil
}

func (r *radio) sendState() error {

	if state, err := r.state.Marshal(); err == nil {
		stateMsg := comms.IOMsg{}
		stateMsg.Data = state
		stateMsg.Retain = true
		stateMsg.Topic = r.settings.CatResponseTopic
		r.settings.ToWireCh <- stateMsg
	} else {
		return err
	}

	return nil
}

func (r *radio) sendCaps() error {

	if caps, err := r.serializeCaps(); err == nil {
		capsMsg := comms.IOMsg{}
		capsMsg.Data = caps
		capsMsg.Retain = true
		capsMsg.Topic = r.settings.CapsTopic
		r.settings.ToWireCh <- capsMsg
	} else {
		log.Println(err)
	}

	return nil
}

func (r *radio) updateMeter() error {

	if !r.state.RadioOn {
		return nil
	}

	vfo := hl.VfoValue[r.state.CurrentVfo]

	newValueAvailable := false

	if r.state.Ptt {
		swr, err := r.rig.GetLevel(vfo, hl.RIG_LEVEL_SWR)
		if err != nil {
			return err
		}

		alc, err := r.rig.GetLevel(vfo, hl.RIG_LEVEL_ALC)
		if err != nil {
			return err
		}

		if r.state.Vfo.Levels["SWR"] != swr {
			r.state.Vfo.Levels["SWR"] = swr
			newValueAvailable = true
		}

		if r.state.Vfo.Levels["ALC"] != alc {
			r.state.Vfo.Levels["ALC"] = alc
			newValueAvailable = true
		}
	} else {
		strength, err := r.rig.GetLevel(vfo, hl.RIG_LEVEL_STRENGTH)
		if err != nil {
			return err
		}

		if r.state.Vfo.Levels["STRENGTH"] != strength {
			r.state.Vfo.Levels["STRENGTH"] = strength
			newValueAvailable = true
		}
	}

	if newValueAvailable {
		err := r.sendState()
		if err != nil {
			return err
		}
	}

	return nil
}
