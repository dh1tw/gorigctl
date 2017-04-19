package server

import (
	"errors"
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
	SyncInterval     time.Duration
	RadioLogger      *log.Logger
	AppLogger        *log.Logger
}

type localRadio struct {
	rig            hl.Rig
	state          sbRadio.State
	settings       *RadioSettings
	pollingTicker  *time.Ticker
	radioLogger    *log.Logger
	appLogger      *log.Logger
	lastUpdateSent time.Time
	lastCmdRecvd   time.Time
	syncTicker     *time.Ticker
}

func StartRadioServer(rs RadioSettings) {

	defer rs.WaitGroup.Done()

	prepareShutdownCh := rs.Events.Sub(events.PrepareShutdown)
	shutdownCh := rs.Events.Sub(events.Shutdown)

	r := localRadio{}
	r.rig = hl.Rig{}
	r.state = sbRadio.State{}
	r.state.Vfo = &sbRadio.Vfo{}
	r.state.Vfo.Split = &sbRadio.Split{}
	r.state.Vfo.Levels = make(map[string]float32)
	r.state.Vfo.Parameters = make(map[string]float32)
	r.state.Vfo.Functions = make(map[string]bool)
	r.state.Channel = &sbRadio.Channel{}
	r.settings = &rs
	r.radioLogger = rs.RadioLogger
	r.appLogger = rs.AppLogger

	r.state.PollingInterval = int32(r.settings.PollingInterval.Nanoseconds() / 1000000)
	r.state.SyncInterval = int32(r.settings.SyncInterval.Seconds())

	r.rig.SetDebugLevel(rs.HlDebugLevel)

	err := r.rig.Init(rs.RigModel)
	if err != nil {
		log.Println(err)
		return
	}

	if rs.RigModel != 1 {
		err = r.rig.SetPort(rs.Port)
		if err != nil {
			// if we can not set the port, we shut down
			log.Println(err)
			r.settings.Events.Pub(true, events.Shutdown)
			return
		}
	}

	if err := r.rig.Open(); err != nil {
		// if we can not open the port, we shut down
		log.Println(err)
		r.settings.Events.Pub(true, events.Shutdown)
		return
	}

	// publish the radio's capabilities
	if err := r.sendCaps(); err != nil {
		r.radioLogger.Println("Couldn't get all capabilities:", err)
	}

	if err := r.queryVfo(); err != nil {
		r.radioLogger.Println(err)
	}

	// publish the radio's state
	if err := r.sendState(); err != nil {
		r.radioLogger.Println(err)
	}

	if r.settings.PollingInterval > 0 {
		r.pollingTicker = time.NewTicker(r.settings.PollingInterval)
	} else {
		r.pollingTicker = time.NewTicker(time.Second * 100)
		r.pollingTicker.Stop()
	}
	if r.settings.SyncInterval > 0 {
		r.syncTicker = time.NewTicker(r.settings.SyncInterval)
	} else {
		r.syncTicker = time.NewTicker(time.Second * 100)
		r.syncTicker.Stop()
	}

	for {
		select {
		case msg := <-rs.CatRequestCh:
			r.deserializeCatRequest(msg)
			r.sendState()
			r.lastCmdRecvd = time.Now()

		case <-prepareShutdownCh:
			r.pollingTicker.Stop()
			r.syncTicker.Stop()
			r.sendClearState()
			time.Sleep(time.Millisecond * 100)
			r.sendClearCaps()

		case <-shutdownCh:
			r.appLogger.Println("Disconnecting from Radio")
			// maybe we have to check if the connection is really open
			r.rig.Close()
			r.rig.Cleanup()
			return

		case <-r.pollingTicker.C:
			r.updateMeter()

		case <-r.syncTicker.C:
			// make sure we don't interrupt while receiving data
			// as updating takes a few hundred milliseconds
			if time.Since(r.lastCmdRecvd) < time.Second*3 {
				continue
			}

			r.queryVfo()

			if (r.rig.Caps.HasGetPowerStat && r.state.RadioOn) || !r.rig.Caps.HasGetPowerStat {

				if err := r.sendState(); err != nil {
					r.radioLogger.Println(err)
				}
			}
		}
	}
}

func (r *localRadio) queryVfo() error {

	if r.rig.Caps.HasGetPowerStat {
		if pwrOn, err := r.rig.GetPowerStat(); err != nil {
			r.radioLogger.Println(err)
			// if the radio doesn't respond, lets assume that the radio if off
			r.state.RadioOn = false
		} else {
			if pwrOn == hl.RIG_POWER_ON {
				r.state.RadioOn = true
			} else {
				r.state.RadioOn = false
				// announce that the radio has ben turned off
				if err := r.sendState(); err != nil {
					return err
				}
				return nil
			}
		}
	}

	// Only query radio if Power is On or if Radio has now PowerStat function
	// in this case we will assume that the radio is turned on
	if (r.rig.Caps.HasGetPowerStat && r.state.RadioOn) || !r.rig.Caps.HasGetPowerStat {

		vfo := hl.VfoValue["CURR"]

		if r.rig.Caps.HasGetVfo {
			vfo, err := r.rig.GetVfo()
			if err != nil {
				r.radioLogger.Print(err)
			} else {
				r.state.CurrentVfo = hl.VfoName[vfo]
			}
		} else {
			r.state.CurrentVfo = "CURR"
		}

		if r.rig.Caps.HasGetFreq {
			freq, err := r.rig.GetFreq(vfo)
			if err != nil {
				r.radioLogger.Println(err)
			} else {
				r.state.Vfo.Frequency = freq
			}
		}

		if r.rig.Caps.HasGetMode {
			mode, pbWidth, err := r.rig.GetMode(vfo)
			if err != nil {
				r.radioLogger.Println(err)
			} else {
				if modeName, ok := hl.ModeName[mode]; ok {
					r.state.Vfo.Mode = modeName
				} else {
					r.radioLogger.Println("unknown mode:", mode)
				}
				r.state.Vfo.PbWidth = int32(pbWidth)
			}
		}

		if r.rig.Caps.HasGetAnt {
			ant, err := r.rig.GetAnt(vfo)
			if err != nil {
				r.radioLogger.Println(err)
			} else {
				r.state.Vfo.Ant = int32(ant)
			}
		}

		if r.rig.Caps.HasGetRit {
			rit, err := r.rig.GetRit(vfo)
			if err != nil {
				r.radioLogger.Println(err)
			} else {
				r.state.Vfo.Rit = int32(rit)
			}
		}

		if r.rig.Caps.HasGetRit {
			xit, err := r.rig.GetXit(vfo)
			if err != nil {
				r.radioLogger.Println(err)
			} else {
				r.state.Vfo.Xit = int32(xit)
			}
		}

		split := sbRadio.Split{}

		if r.rig.Caps.HasGetSplitVfo {
			splitOn, txVfo, err := r.rig.GetSplit(vfo)
			if err != nil {
				r.radioLogger.Println(err)
			} else {
				if splitOn == hl.RIG_SPLIT_ON {
					split.Enabled = true
				} else {
					split.Enabled = false
				}
				if txVfoName, ok := hl.VfoName[txVfo]; ok {
					split.Vfo = txVfoName
				} else {
					return errors.New("unknown Vfo Name")
				}

				if splitOn == hl.RIG_SPLIT_ON {

					// these checks should be enabled, but most of the
					// backends don't have these functions implemented
					// therefore they use the emulated functions which
					// unfortunately don't work everywhere well (e.g. TS-480)
					// if r.rig.Caps.HasGetSplitFreq {
					txFreq, err := r.rig.GetSplitFreq(txVfo)
					if err != nil {
						r.radioLogger.Println(err)
					} else {
						split.Frequency = txFreq
					}
					// }

					// if r.rig.Caps.HasGetSplitMode {
					txMode, txPbWidth, err := r.rig.GetSplitMode(txVfo)
					if err != nil {
						r.radioLogger.Println(err)
					} else {
						if txModeName, ok := hl.ModeName[txMode]; ok {
							split.Mode = txModeName
						} else {
							r.radioLogger.Println("unknown Tx Mode")
						}
						split.PbWidth = int32(txPbWidth)
					}
					// }
				}
			}
		}

		r.state.Vfo.Split = &split

		if r.rig.Caps.HasGetTs {
			tStep, err := r.rig.GetTs(vfo)
			if err != nil {
				r.radioLogger.Println(err)
			} else {
				r.state.Vfo.TuningStep = int32(tStep)
			}
		}

		for _, f := range r.rig.Caps.GetFunctions {
			fValue, err := r.rig.GetFunc(vfo, hl.FuncValue[f])
			if err != nil {
				r.radioLogger.Println(err)
			}
			r.state.Vfo.Functions[f] = fValue
		}

		for _, level := range r.rig.Caps.GetLevels {
			lValue, err := r.rig.GetLevel(vfo, hl.LevelValue[level.Name])
			if err != nil {
				r.radioLogger.Println("Warning:", level.Name, "-", err)
			}
			r.state.Vfo.Levels[level.Name] = lValue
		}

		for _, param := range r.rig.Caps.GetParameters {
			pValue, err := r.rig.GetParm(vfo, hl.ParmValue[param.Name])
			if err != nil {
				r.radioLogger.Println(err)
			}
			r.state.Vfo.Parameters[param.Name] = pValue
		}
	}

	r.lastUpdateSent = time.Now()

	return nil
}

func (r *localRadio) sendState() error {

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

func (r *localRadio) sendCaps() error {

	if caps, err := r.serializeCaps(); err == nil {
		capsMsg := comms.IOMsg{}
		capsMsg.Data = caps
		capsMsg.Retain = true
		capsMsg.Topic = r.settings.CapsTopic
		r.settings.ToWireCh <- capsMsg
	} else {
		r.radioLogger.Println(err)
	}

	return nil
}

func (r *localRadio) updateMeter() error {

	// Only update the meter when we can be sure that the radio is
	// actually turned on. If the rig does not provide the powerstat
	// we quit to avoid sending messages to the radio which will be
	// continously rejected

	if !r.rig.Caps.HasGetPowerStat || !r.rig.Caps.HasSetPowerStat {
		return nil
	}

	if !r.state.RadioOn {
		return nil
	}

	vfo := hl.VfoValue[r.state.CurrentVfo]
	newValueAvailable := false

	if r.rig.Caps.HasGetPtt && r.state.Ptt {

		if swrCurrValue, ok := r.state.Vfo.Levels["SWR"]; ok {
			swrNewValue, err := r.rig.GetLevel(vfo, hl.RIG_LEVEL_SWR)
			if err != nil {
				return err
			}
			if swrNewValue != swrCurrValue {
				r.state.Vfo.Levels["SWR"] = swrNewValue
				newValueAvailable = true
			}
		}

		if alcCurrValue, ok := r.state.Vfo.Levels["ALC"]; ok {
			alcNewValue, err := r.rig.GetLevel(vfo, hl.RIG_LEVEL_ALC)
			if err != nil {
				return err
			}
			if alcNewValue != alcCurrValue {
				r.state.Vfo.Levels["ALC"] = alcNewValue
				newValueAvailable = true
			}
		}

	} else {

		if sCurrValue, ok := r.state.Vfo.Levels["STRENGTH"]; ok {
			sNewValue, err := r.rig.GetLevel(vfo, hl.RIG_LEVEL_STRENGTH)
			if err != nil {
				return err
			}
			if sNewValue != sCurrValue {
				r.state.Vfo.Levels["STRENGTH"] = sNewValue
				newValueAvailable = true
			}
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

func (r *localRadio) sendClearState() error {

	msg := comms.IOMsg{}
	msg.Data = []byte{}
	msg.Retain = true
	msg.Topic = r.settings.CatResponseTopic
	msg.Qos = 0

	r.settings.ToWireCh <- msg

	return nil
}

func (r *localRadio) sendClearCaps() error {

	msg := comms.IOMsg{}
	msg.Data = []byte{}
	msg.Retain = true
	msg.Topic = r.settings.CapsTopic
	msg.Qos = 0

	r.settings.ToWireCh <- msg

	return nil
}
