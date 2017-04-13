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
	Logger           *log.Logger
}

type radio struct {
	rig           hl.Rig
	state         sbRadio.State
	settings      *RadioSettings
	pollingTicker *time.Ticker
	logger        *log.Logger
	lastUpdate    time.Time
	updateTicker  *time.Ticker
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
	r.logger = rs.Logger

	r.state.PollingInterval = int32(r.settings.PollingInterval.Nanoseconds() / 1000000)

	err := r.rig.Init(rs.RigModel)
	if err != nil {
		r.logger.Println(err)
		return
	}

	r.rig.SetDebugLevel(rs.HlDebugLevel)

	err = r.rig.SetPort(rs.Port)
	if err != nil {
		// if we can not set the port, we shut down
		r.logger.Println(err)
		r.settings.Events.Pub(true, events.Shutdown)
		return
	}

	if err := r.rig.Open(); err != nil {
		// if we can not open the port, we shut down
		r.logger.Println(err)
		r.settings.Events.Pub(true, events.Shutdown)
		return
	}

	// publish the radio's capabilities
	if err := r.sendCaps(); err != nil {
		r.logger.Println("Couldn't get all capabilities:", err)
	}

	if err := r.queryVfo(); err != nil {
		fmt.Println(err)
	}

	// publish the radio's state
	if err := r.sendState(); err != nil {
		r.logger.Println(err)
	}

	r.pollingTicker = time.NewTicker(r.settings.PollingInterval)
	r.updateTicker = time.NewTicker(time.Second)

	for {
		select {
		case msg := <-rs.CatRequestCh:
			r.deserializeCatRequest(msg)
			r.sendState()

		case <-shutdownCh:
			r.logger.Println("Disconnecting from Radio")
			// maybe we have to check if the connection is really open
			r.rig.Close()
			r.rig.Cleanup()
			return

		case <-r.pollingTicker.C:
			r.updateMeter()

		case <-r.updateTicker.C:
			// r.queryVfo()
		}
	}
}

func (r *radio) queryVfo() error {

	if r.rig.Caps.HasGetPowerStat {
		if pwrOn, err := r.rig.GetPowerStat(); err != nil {
			r.logger.Println(err)
			// if the radio doesn't respond, lets assume that the radio if off
			r.state.RadioOn = false
		} else {
			if pwrOn == hl.RIG_POWER_ON {
				r.state.RadioOn = true
			} else {
				r.state.RadioOn = false
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
				r.logger.Print(err)
			} else {
				r.state.CurrentVfo = hl.VfoName[vfo]
			}
		} else {
			r.state.CurrentVfo = "CURR"
		}

		if r.rig.Caps.HasGetFreq {
			freq, err := r.rig.GetFreq(vfo)
			if err != nil {
				r.logger.Println(err)
			} else {
				r.state.Vfo.Frequency = freq
			}
		}

		if r.rig.Caps.HasGetMode {
			mode, pbWidth, err := r.rig.GetMode(vfo)
			if err != nil {
				r.logger.Println(err)
			} else {
				if modeName, ok := hl.ModeName[mode]; ok {
					r.state.Vfo.Mode = modeName
				} else {
					r.logger.Println("unknown mode:", mode)
				}
				r.state.Vfo.PbWidth = int32(pbWidth)
			}
		}

		if r.rig.Caps.HasGetAnt {
			ant, err := r.rig.GetAnt(vfo)
			if err != nil {
				r.logger.Println(err)
			} else {
				r.state.Vfo.Ant = int32(ant)
			}
		}

		if r.rig.Caps.HasGetRit {
			rit, err := r.rig.GetRit(vfo)
			if err != nil {
				r.logger.Println(err)
			} else {
				r.state.Vfo.Rit = int32(rit)
			}
		}

		if r.rig.Caps.HasGetRit {
			xit, err := r.rig.GetXit(vfo)
			if err != nil {
				r.logger.Println(err)
			} else {
				r.state.Vfo.Xit = int32(xit)
			}
		}

		split := sbRadio.Split{}

		if r.rig.Caps.HasGetSplitVfo {
			splitOn, txVfo, err := r.rig.GetSplit(vfo)
			if err != nil {
				r.logger.Println(err)
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
						r.logger.Println(err)
					} else {
						split.Frequency = txFreq
					}
					// }

					// if r.rig.Caps.HasGetSplitMode {
					txMode, txPbWidth, err := r.rig.GetSplitMode(txVfo)
					if err != nil {
						r.logger.Println(err)
					} else {
						if txModeName, ok := hl.ModeName[txMode]; ok {
							split.Mode = txModeName
						} else {
							r.logger.Println("unknown Tx Mode")
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
				r.logger.Println(err)
			} else {
				r.state.Vfo.TuningStep = int32(tStep)
			}
		}

		r.state.Vfo.Functions = make(map[string]bool)

		for _, f := range r.rig.Caps.GetFunctions {
			fValue, err := r.rig.GetFunc(vfo, hl.FuncValue[f])
			if err != nil {
				r.logger.Println(err)
			}
			r.state.Vfo.Functions[f] = fValue
		}

		r.state.Vfo.Levels = make(map[string]float32)
		for _, level := range r.rig.Caps.GetLevels {
			lValue, err := r.rig.GetLevel(vfo, hl.LevelValue[level.Name])
			if err != nil {
				r.logger.Println("Warning:", level.Name, "-", err)
			}
			r.state.Vfo.Levels[level.Name] = lValue
		}

		r.state.Vfo.Parameters = make(map[string]float32)
		for _, param := range r.rig.Caps.GetParameters {
			pValue, err := r.rig.GetParm(vfo, hl.ParmValue[param.Name])
			if err != nil {
				r.logger.Println(err)
			}
			r.state.Vfo.Parameters[param.Name] = pValue
		}

		r.logger.Println("Functions: ", r.state.Vfo.Functions)
		r.logger.Println("Levels:", r.state.Vfo.Levels)
		r.logger.Println("Parameters", r.state.Vfo.Parameters)
	}

	r.lastUpdate = time.Now()

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
		r.logger.Println(err)
	}

	return nil
}

func (r *radio) updateMeter() error {

	// Only update the meter when we can be sure that the radio is
	// actually turned on. If the rig does not provide the powerstat
	// we quit to avoid sending messages to the radio which will be
	// continously rejected
	if !r.rig.Caps.HasGetPowerStat && !r.rig.Caps.HasSetPowerStat {
		return nil
	}

	if !r.state.RadioOn {
		return nil
	}

	vfo := hl.VfoValue[r.state.CurrentVfo]
	newValueAvailable := false

	if r.state.Ptt {

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
