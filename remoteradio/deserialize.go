package remoteradio

import (
	"reflect"

	"github.com/dh1tw/gorigctl/events"
	sbRadio "github.com/dh1tw/gorigctl/sb_radio"
	sbStatus "github.com/dh1tw/gorigctl/sb_status"
)

func (r *RemoteRadio) DeserializeRadioStatus(data []byte) error {

	rStatus := sbStatus.Status{}
	if err := rStatus.Unmarshal(data); err != nil {
		return err
	}

	if r.radioOnline != rStatus.Online {
		r.radioOnline = rStatus.Online
		r.events.Pub(rStatus.Online, events.RadioOnline)
	}

	return nil
}

func (r *RemoteRadio) DeserializeCaps(msg []byte) error {

	caps := sbRadio.Capabilities{}
	err := caps.Unmarshal(msg)
	if err != nil {
		return err
	}

	r.caps = caps

	return nil
}

func (r *RemoteRadio) DeserializeCatResponse(msg []byte) error {

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
					r.logger.Println(err)
				}
			}
		}

		if ns.Vfo.GetTuningStep() != r.state.Vfo.TuningStep {
			r.state.Vfo.TuningStep = ns.Vfo.GetTuningStep()
			if r.printRigUpdates {
				r.logger.Printf("Updated Tuning Step: %dHz\n", r.state.Vfo.TuningStep)
			}
		}

		if !reflect.DeepEqual(ns.GetVfo().GetFunctions(), r.state.Vfo.Functions) {
			if err := r.updateFunctions(ns.Vfo.GetFunctions()); err != nil {
				r.logger.Println(err)
			}
		}

		if !reflect.DeepEqual(ns.GetVfo().GetLevels(), r.state.Vfo.Levels) {
			if err := r.updateLevels(ns.Vfo.GetLevels()); err != nil {
				r.logger.Println(err)
			}
		}

		if !reflect.DeepEqual(ns.GetVfo().GetParameters(), r.state.Vfo.Parameters) {
			if err := r.updateParams(ns.Vfo.GetParameters()); err != nil {
				r.logger.Println(err)
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

	if ns.GetSyncInterval() != r.state.SyncInterval {
		r.state.SyncInterval = ns.GetSyncInterval()
		if r.printRigUpdates {
			r.logger.Printf("Updated rig sync interval: %ds\n", r.state.SyncInterval)
		}
	}
	return nil
}

func (r *RemoteRadio) updateSplit(newSplit *sbRadio.Split) error {

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

func (r *RemoteRadio) updateFunctions(newFuncs map[string]bool) error {

	r.state.Vfo.Functions = newFuncs
	if r.printRigUpdates {
		r.logger.Println("Updated functions:")
		for name, value := range r.state.Vfo.Functions {
			r.logger.Printf("%v: %v", name, value)
		}
	}

	// vfo, _ := hl.VfoValue[r.State.CurrentVfo]

	// functions to be enabled
	// diff := utils.SliceDiff(newFuncs, r.State.Vfo.Functions)
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
	// diff = utils.SliceDiff(r.State.Vfo.Functions, newFuncs)
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

func (r *RemoteRadio) updateLevels(newLevels map[string]float32) error {

	r.state.Vfo.Levels = newLevels

	if r.printRigUpdates {
		r.logger.Println("Updated levels:")
		for name, value := range r.state.Vfo.Levels {
			r.logger.Printf("%v: %v", name, value)
		}
	}

	// vfo, _ := hl.VfoValue[r.State.CurrentVfo]

	// for k, v := range newLevels {
	// 	levelValue, ok := hl.LevelValue[k]
	// 	if !ok {
	// 		return errors.New("unknown Level")
	// 	}
	// 	if _, ok := r.State.Vfo.Levels[k]; !ok {
	// 		return errors.New("unsupported Level for this rig")
	// 	}

	// if r.State.Vfo.Levels[k] != v {
	// 	err := r.rig.SetLevel(vfo, levelValue, v)
	// 	if err != nil {
	// 		return nil
	// 	}
	// }
	// }

	return nil
}

func (r *RemoteRadio) updateParams(newParams map[string]float32) error {

	r.state.Vfo.Parameters = newParams

	if r.printRigUpdates {
		r.logger.Println("Updated parameters:")
		for name, value := range r.state.Vfo.Parameters {
			r.logger.Printf("%v: %v", name, value)
		}
	}

	// vfo, _ := hl.VfoValue[r.State.CurrentVfo]

	// for k, v := range newParams {
	// 	paramValue, ok := hl.ParmValue[k]
	// 	if !ok {
	// 		return errors.New("unknown Parameter")
	// 	}
	// if _, ok := r.State.Vfo.Parameters[k]; !ok {
	// 	return errors.New("unsupported Parameter for this rig")
	// }
	// if r.State.Vfo.Levels[k] != v {
	// 	err := r.rig.SetLevel(vfo, paramValue, v)
	// 	if err != nil {
	// 		return nil
	// 	}
	// }
	// }

	return nil
}
