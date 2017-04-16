package radio

import (
	"errors"
	"reflect"

	"time"

	hl "github.com/dh1tw/goHamlib"
	sbRadio "github.com/dh1tw/gorigctl/sb_radio"
	"github.com/dh1tw/gorigctl/utils"
)

func (r *radio) deserializeCatRequest(request []byte) error {

	ns := sbRadio.SetState{}
	if err := ns.Unmarshal(request); err != nil {
		return err
	}

	if ns.Md.HasRadioOn {
		if ns.GetRadioOn() != r.state.RadioOn {
			if err := r.updatePowerOn(ns.GetRadioOn()); err != nil {
				r.radioLogger.Println(err)
			}

			return nil
		}
	}

	// if r.state.RadioOn {

	if ns.CurrentVfo != r.state.CurrentVfo {
		r.appLogger.Println("updating vfo to", ns.CurrentVfo)
		if err := r.updateCurrentVfo(ns.CurrentVfo); err != nil {
			r.radioLogger.Println(err)
		}
	}

	if len(ns.VfoOperations) > 0 {
		r.appLogger.Println("executing vfo operation(s)", ns.VfoOperations)
		if err := r.execVfoOperations(ns.GetVfoOperations()); err != nil {
			r.radioLogger.Println(err)
		}
	}

	if ns.Md.HasFrequency {
		if ns.Vfo.Frequency != r.state.Vfo.Frequency {
			r.appLogger.Printf("updating frequency to %.0f Hz\n", ns.Vfo.Frequency)
			if err := r.updateFrequency(ns.Vfo.Frequency); err != nil {
				r.radioLogger.Println(err)
			}
		}
	}

	if ns.Md.HasMode {
		if ns.Vfo.GetMode() != r.state.Vfo.Mode {
			r.appLogger.Println("updating mode to", ns.Vfo.Mode)
			if err := r.updateMode(ns.Vfo.GetMode(), ns.Vfo.GetPbWidth()); err != nil {
				r.radioLogger.Println(err)
			}
		}
	}

	if ns.Md.HasPbWidth {
		if ns.Vfo.GetPbWidth() != r.state.Vfo.PbWidth {
			r.appLogger.Printf("updating pbwidth to %d Hz\n", ns.Vfo.PbWidth)
			if err := r.updatePbWidth(ns.Vfo.GetPbWidth()); err != nil {
				r.radioLogger.Println(err)
			}
		}
	}

	if ns.Md.HasAnt {
		if ns.Vfo.GetAnt() != r.state.Vfo.Ant {
			r.appLogger.Println("updating antenna to", ns.Vfo.Ant)
			if err := r.updateAntenna(ns.Vfo.GetAnt()); err != nil {
				r.radioLogger.Println(err)
			}
		}
	}

	if ns.Md.HasRit {
		if ns.Vfo.GetRit() != r.state.Vfo.Rit {
			r.appLogger.Printf("updating rit to %d Hz\n", ns.Vfo.Rit)
			if err := r.updateRit(ns.Vfo.GetRit()); err != nil {
				r.radioLogger.Println(err)
			}
		}
	}

	if ns.Md.HasXit {
		if ns.Vfo.GetXit() != r.state.Vfo.Xit {
			r.appLogger.Printf("updating xit to %d Hz\n", ns.Vfo.Xit)
			if err := r.updateXit(ns.Vfo.GetXit()); err != nil {
				r.radioLogger.Println(err)
			}
		}
	}

	if ns.Md.HasSplit {
		if ns.Vfo.Split != nil {
			if !reflect.DeepEqual(ns.Vfo.Split, r.state.Vfo.Split) {
				r.appLogger.Println("updating split to", ns.Vfo.Split)
				if err := r.updateSplit(ns.Vfo.Split); err != nil {
					r.radioLogger.Println(err)
				}
			}
		}
	}

	if ns.Md.HasTuningStep {
		if ns.Vfo.GetTuningStep() != r.state.Vfo.TuningStep {
			r.appLogger.Printf("updating tuning step to %d Hz\n", ns.Vfo.TuningStep)
			if err := r.updateTs(ns.Vfo.GetTuningStep()); err != nil {
				r.radioLogger.Println(err)
			}
		}
	}

	if ns.Md.HasFunctions {
		if ns.Vfo.Functions != nil {
			if !reflect.DeepEqual(ns.Vfo.Functions, r.state.Vfo.Functions) {
				r.appLogger.Println("updating one or more functions")
				if err := r.updateFunctions(ns.Vfo.GetFunctions()); err != nil {
					r.radioLogger.Println(err)
				}
			}
		}
	}

	if ns.Md.HasLevels {
		if ns.Vfo.Levels != nil {
			if !reflect.DeepEqual(ns.Vfo.Levels, r.state.Vfo.Levels) {
				r.appLogger.Println("updating one or more levels")
				if err := r.updateLevels(ns.Vfo.GetLevels()); err != nil {
					r.radioLogger.Println(err)
				}
			}
		}
	}

	if ns.Md.HasParameters {
		if ns.Vfo.Parameters != nil {
			if !reflect.DeepEqual(ns.Vfo.Parameters, r.state.Vfo.Parameters) {
				r.appLogger.Println("updating one or more parameters")
				if err := r.updateParams(ns.Vfo.GetParameters()); err != nil {
					r.radioLogger.Println(err)
				}
			}
		}
	}
	// }

	if ns.Md.HasPtt {
		if ns.GetPtt() != r.state.Ptt {
			r.radioLogger.Println("updating ptt to", ns.Ptt)
			if err := r.updatePtt(ns.GetPtt()); err != nil {
				r.radioLogger.Println(err)
			}
		}
	}

	if ns.Md.HasPollingInterval {
		if ns.GetPollingInterval() != r.state.PollingInterval {
			if ns.GetPollingInterval() > 0 {
				r.radioLogger.Printf("updating rig polling interval to %dms\n", ns.PollingInterval)
				newPollingInterval := time.Millisecond * time.Duration(ns.GetPollingInterval())
				r.pollingTicker.Stop()
				r.pollingTicker = time.NewTicker(newPollingInterval)
				r.state.PollingInterval = ns.GetPollingInterval()
			} else {
				r.radioLogger.Println("stopped rig polling")
				r.pollingTicker.Stop()
				r.state.PollingInterval = 0
			}
		}
	}

	if ns.Md.HasSyncInterval {
		if ns.GetSyncInterval() != r.state.SyncInterval {
			if ns.GetSyncInterval() > 0 {
				r.radioLogger.Printf("updating rig sync interval to %ds\n", ns.SyncInterval)
				newSyncInterval := time.Second * time.Duration(ns.GetSyncInterval())
				r.syncTicker.Stop()
				r.syncTicker = time.NewTicker(newSyncInterval)
				r.state.SyncInterval = ns.GetSyncInterval()
			} else {
				r.radioLogger.Println("stopped rig sync")
				r.syncTicker.Stop()
				r.state.SyncInterval = 0
			}
		}
	}

	return nil
}

func (r *radio) updateCurrentVfo(newVfo string) error {
	if vfo, ok := hl.VfoValue[newVfo]; ok {
		err := r.rig.SetVfo(vfo)
		if err != nil {
			return err
		}
		r.queryVfo()
	} else {
		return errors.New("unknown Vfo")
	}
	return nil
}

func (r *radio) updateFrequency(newFreq float64) error {
	vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	// if the rig supports fast_commands, then we will use it
	hasFastToken := r.rig.HasToken("fast_commands_token")

	if hasFastToken {
		err := r.rig.SetConf("fast_commands_token", "1")
		if err != nil {
			r.radioLogger.Println(err)
		}
	}

	err := r.rig.SetFreq(vfo, newFreq)
	if err != nil {
		return err
	}

	if hasFastToken {
		err = r.rig.SetConf("fast_commands_token", "0")
		if err != nil {
			r.radioLogger.Println(err)
		}
	}

	r.state.Vfo.Frequency = newFreq
	return nil
}

func (r *radio) execVfoOperations(vfoOps []string) error {
	vfo, _ := hl.VfoValue[r.state.CurrentVfo]
	for _, v := range vfoOps {
		vfoOpValue, ok := hl.OperationValue[v]
		if !ok {
			return errors.New("unknown VFO Operation")
		}
		err := r.rig.VfoOp(vfo, vfoOpValue)
		if err != nil {
			return err
		}

		// if there are major changes to the VFO we have
		// to re-read the VFO data from the RIG
		if utils.StringInSlice("XCHG", vfoOps) ||
			utils.StringInSlice("TOGGLE", vfoOps) ||
			utils.StringInSlice("FROM_VFO", vfoOps) ||
			utils.StringInSlice("TO_VFO", vfoOps) ||
			utils.StringInSlice("MCL", vfoOps) {

			if err := r.queryVfo(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *radio) updateMode(newMode string, newPbWidth int32) error {

	if !r.rig.Caps.HasSetMode || !r.rig.Caps.HasGetMode {
		return errors.New("unable to update mode; function not implemented")
	}

	vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	pbWidth := int(r.state.Vfo.PbWidth)
	if newPbWidth > 0 {
		pbWidth = int(newPbWidth)
	}

	newModeValue, ok := hl.ModeValue[newMode]
	if !ok {
		return errors.New("unknown mode")
	}

	err := r.rig.SetMode(vfo, newModeValue, pbWidth)
	if err != nil {
		pbNormal, err := r.rig.GetPbNormal(newModeValue)
		if err != nil {
			return err
		}
		err = r.rig.SetMode(vfo, newModeValue, pbNormal)
		if err != nil {
			return err
		}
	}

	mode, pbWidth, err := r.rig.GetMode(vfo)
	if err != nil {
		return err
	}

	r.state.Vfo.Mode = hl.ModeName[mode]
	r.state.Vfo.PbWidth = int32(pbWidth)

	ts := 0
	if r.rig.Caps.HasGetTs {
		ts, err = r.rig.GetTs(vfo)
		if err != nil {
			return err
		}
	}

	r.state.Vfo.TuningStep = int32(ts)

	return nil
}

func (r *radio) updatePbWidth(newPbWidth int32) error {

	if !r.rig.Caps.HasSetMode || !r.rig.Caps.HasGetMode {
		return errors.New("unable to update mode/filter; function not implemented")
	}

	vfo, _ := hl.VfoValue[r.state.CurrentVfo]
	modeValue := hl.ModeValue[r.state.Vfo.Mode]
	err := r.rig.SetMode(vfo, modeValue, int(newPbWidth))
	if err != nil {
		return err
	}

	mode, pbWidth, err := r.rig.GetMode(vfo)
	if err != nil {
		return err
	}

	r.state.Vfo.Mode = hl.ModeName[mode]
	r.state.Vfo.PbWidth = int32(pbWidth)

	if r.rig.Caps.HasGetTs {
		ts := 0
		ts, err = r.rig.GetTs(vfo)
		if err != nil {
			return err
		}
		r.state.Vfo.TuningStep = int32(ts)
	}

	return nil
}

func (r *radio) updateAntenna(newAnt int32) error {

	vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	err := r.rig.SetAnt(vfo, int(newAnt))
	if err != nil {
		return err
	}
	ant, err := r.rig.GetAnt(vfo)
	if err != nil {
		return err
	}

	r.state.Vfo.Ant = int32(ant)

	return nil
}

func (r *radio) updateRit(newRit int32) error {
	vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	err := r.rig.SetRit(vfo, int(newRit))
	if err != nil {
		return err
	}

	// verify rit
	rit, err := r.rig.GetRit(vfo)
	if err != nil {
		return err
	}
	r.state.Vfo.Rit = int32(rit)

	// on some rigs rit will also update xit
	xit, err := r.rig.GetXit(vfo)
	if err != nil {
		return err
	}
	r.state.Vfo.Xit = int32(xit)

	// vfo frequency might also have changed
	freq, err := r.rig.GetFreq(vfo)
	if err != nil {
		return err
	}
	r.state.Vfo.Frequency = freq

	return nil
}

func (r *radio) updateXit(newXit int32) error {
	vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	err := r.rig.SetXit(vfo, int(newXit))
	if err != nil {
		return err
	}

	// verify rit
	xit, err := r.rig.GetXit(vfo)
	if err != nil {
		return err
	}
	r.state.Vfo.Xit = int32(xit)

	// one some rigs rit will also update rit
	rit, err := r.rig.GetRit(vfo)
	if err != nil {
		return err
	}
	r.state.Vfo.Rit = int32(rit)

	return nil
}

func (r *radio) updateSplit(newSplit *sbRadio.Split) error {

	vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	if !r.rig.Caps.HasGetSplitVfo || !r.rig.Caps.HasSetSplitVfo {
		return errors.New("radio doesn't support split")
	}

	if newSplit.Enabled != r.state.Vfo.Split.Enabled {

		txVfo, ok := hl.VfoValue[newSplit.Vfo]
		if !ok {
			return errors.New("unknown tx (split) vfo")
		}

		err := r.rig.SetSplitVfo(vfo, utils.Btoi(newSplit.Enabled), txVfo)
		if err != nil {
			return err
		}

		// verify and set the radio's split vfo, enabled flag
		checkSplitEnabled, checkTxVfo, err := r.rig.GetSplitVfo(vfo)
		if err != nil {
			return err
		}

		r.state.Vfo.Split.Vfo, _ = hl.VfoName[checkTxVfo]
		r.state.Vfo.Split.Enabled = utils.Itob(checkSplitEnabled)
	}

	// clear if split has been deactivated
	if !r.state.Vfo.Split.Enabled {
		// unset split frequency, mode, pbWidth
		r.state.Vfo.Split.Frequency = 0
		r.state.Vfo.Split.Mode = ""
		r.state.Vfo.Split.Vfo = ""
		r.state.Vfo.Split.PbWidth = 0
		return nil
	}

	if r.rig.Caps.HasGetSplitFreq && r.rig.Caps.HasSetSplitFreq {

		txVfo, ok := hl.VfoValue[r.state.Vfo.Split.Vfo]
		if !ok {
			return errors.New("unknown tx (split) vfo")
		}

		if newSplit.Frequency != r.state.Vfo.Split.Frequency &&
			newSplit.Frequency > 0 {

			err := r.rig.SetSplitFreq(txVfo, newSplit.Frequency)
			if err != nil {
				return err
			}
		}

		// verify the radios split frequency
		txFreq, err := r.rig.GetSplitFreq(txVfo)
		if err != nil {
			return err
		}
		r.state.Vfo.Split.Frequency = txFreq
	}

	// the check below should be performed, but unfortunately, most of the rigs
	// have not implemented specific functions to set/get split mode. Some
	// even don't work well with the fallback functions (e.g. TS-480 which
	// disables the split when querying the split frequency)

	// if r.rig.Caps.HasGetSplitMode && r.rig.Caps.HasSetSplitMode {
	if newSplit.Mode != r.state.Vfo.Split.Mode &&
		len(newSplit.Mode) > 0 {

		txVfo, ok := hl.VfoValue[r.state.Vfo.Split.Vfo]
		if !ok {
			return errors.New("unknown tx (split) vfo")
		}

		newSplitModeValue, ok := hl.ModeValue[newSplit.Mode]
		if !ok {
			return errors.New("unknown tx (split) mode")
		}

		pbWidth := r.state.Vfo.Split.PbWidth
		if newSplit.PbWidth > 0 {
			pbWidth = newSplit.PbWidth
		}

		err := r.rig.SetSplitMode(txVfo, newSplitModeValue, int(pbWidth))
		if err != nil {
			// if this fails, we will try to set again the split
			// mode with the default filter width

			// get the standard filter width for this radio
			pbNormal, err := r.rig.GetPbNormal(newSplitModeValue)
			if err != nil {
				return err
			}

			// try again
			err = r.rig.SetSplitMode(txVfo, newSplitModeValue, pbNormal)
			if err != nil {
				return err
			}
		}
	}

	// verify the radios txMode and txPbWidth
	txVfo, ok := hl.VfoValue[r.state.Vfo.Split.Vfo]
	if !ok {
		return errors.New("unknown vfo")
	}
	txMode, txPbWidth, err := r.rig.GetSplitMode(txVfo)
	if err != nil {
		return err
	}
	r.state.Vfo.Split.Mode = hl.ModeName[txMode]
	r.state.Vfo.Split.PbWidth = int32(txPbWidth)

	// return nil
	// }

	// we only reach this code if the mode is the same, but we want
	// to update the filter width
	if r.rig.Caps.HasGetSplitMode && r.rig.Caps.HasSetSplitMode {
		txVfo, ok = hl.VfoValue[r.state.Vfo.Split.Vfo]
		if !ok {
			return errors.New("unknown vfo")
		}

		if newSplit.GetPbWidth() != r.state.Vfo.Split.PbWidth &&
			len(newSplit.GetMode()) > 0 {

			splitModeValue := hl.ModeValue[newSplit.GetMode()]
			err := r.rig.SetSplitMode(txVfo, splitModeValue, int(newSplit.GetPbWidth()))
			if err != nil {
				return err
			}
		}

		// verify the radios mode and txPbWidth
		txMode, txPbWidth, err := r.rig.GetSplitMode(txVfo)
		if err != nil {
			return err
		}
		r.state.Vfo.Split.Mode = hl.ModeName[txMode]
		r.state.Vfo.Split.PbWidth = int32(txPbWidth)
	}

	// Double check if non of the above functions have disabled split
	// accidentally (as it is the case for the TS-480 as of 15.4.2017)
	checkSplitEnabled, _, err := r.rig.GetSplitVfo(vfo)
	if err != nil {
		return err
	}

	// clear if split has been (acidentally) deactivated
	if !utils.Itob(checkSplitEnabled) {
		r.state.Vfo.Split.Enabled = false
		r.state.Vfo.Split.Frequency = 0
		r.state.Vfo.Split.Mode = ""
		r.state.Vfo.Split.Vfo = ""
		r.state.Vfo.Split.PbWidth = 0
		return errors.New("unable to set split")
	}

	return nil
}

func (r *radio) updateTs(newTs int32) error {
	vfo, _ := hl.VfoValue[r.state.CurrentVfo]
	err := r.rig.SetTs(vfo, int(newTs))
	if err != nil {
		return err
	}

	ts, err := r.rig.GetTs(vfo)
	if err != nil {
		return err
	}

	r.state.Vfo.TuningStep = int32(ts)

	return nil
}

func (r *radio) updateFunctions(newFuncs map[string]bool) error {
	vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	// functions to be enabled
	for funcName, newFuncValue := range newFuncs {

		funcValue, ok := hl.FuncValue[funcName]
		if !ok {
			r.radioLogger.Println("unknown function", funcName)
			continue
		}

		// make sure that the radio can actually set this function
		if utils.StringInSlice(funcName, r.rig.Caps.SetFunctions) {

			err := r.rig.SetFunc(vfo, funcValue, newFuncValue)
			if err != nil {
				r.radioLogger.Println("unable to set function", funcValue, err)
			}
		} else {
			r.radioLogger.Println("radio does not support setting function", funcName)
		}

		// before we can verify the function's value we have to check that
		// the radio can actually get this function
		if utils.StringInSlice(funcName, r.rig.Caps.GetFunctions) {
			cfmFuncValue, err := r.rig.GetFunc(vfo, funcValue)
			if err != nil {
				r.radioLogger.Println("unable to verify function", funcValue, err)
				continue
			}
			r.state.Vfo.Functions[funcName] = cfmFuncValue
		}
	}

	return nil
}

func (r *radio) updateLevels(newLevels map[string]float32) error {
	vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	// iterate over all new Levels in the map
	for levelName, newLevel := range newLevels {

		hlLevel, ok := hl.LevelValue[levelName]
		if !ok {
			r.radioLogger.Println("unknown level", levelName)
			continue
		}

		hasSetLevel := false
		for _, setLevel := range r.rig.Caps.SetLevels {
			if setLevel.Name == levelName {
				hasSetLevel = true
				break
			}
		}
		if !hasSetLevel {
			r.radioLogger.Println("radio does not support setting the level", levelName)
			continue
		}

		if r.state.Vfo.Levels[levelName] != newLevel {

			err := r.rig.SetLevel(vfo, hlLevel, newLevel)
			if err != nil {
				r.radioLogger.Println("unable to set level", levelName)
			}
		}

		// before we can verify the level we have to check that
		// the radio can actually get this level

		hasGetLevel := false
		for _, getLevel := range r.rig.Caps.GetLevels {
			if getLevel.Name == levelName {
				hasGetLevel = true
				break
			}
		}

		if hasGetLevel {
			cfmLevel, err := r.rig.GetLevel(vfo, hlLevel)
			if err != nil {
				r.radioLogger.Println("unable to verify level", levelName)
				continue
			}
			r.state.Vfo.Levels[levelName] = cfmLevel
		}
	}

	return nil
}

func (r *radio) updateParams(newParams map[string]float32) error {
	vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	// iterate over all new Levels in the map
	for parmName, newParm := range newParams {

		hlParam, ok := hl.ParmValue[parmName]
		if !ok {
			r.radioLogger.Println("unknown parameter", parmName)
			continue
		}

		hasSetParm := false
		for _, setParm := range r.rig.Caps.SetParameters {
			if setParm.Name == parmName {
				hasSetParm = true
				break
			}
		}
		if !hasSetParm {
			r.radioLogger.Println("radio does not support setting the parameter", parmName)
			continue
		}

		if r.state.Vfo.Parameters[parmName] != newParm {

			err := r.rig.SetParm(vfo, hlParam, newParm)
			if err != nil {
				r.radioLogger.Println("unable to set parameter", parmName)
			}
		}

		// before we can verify the parameter we have to check that
		// the radio can actually get this parameter

		hasGetParm := false
		for _, getParm := range r.rig.Caps.GetParameters {
			if getParm.Name == parmName {
				hasGetParm = true
				break
			}
		}

		if hasGetParm {
			cfmParm, err := r.rig.GetParm(vfo, hlParam)
			if err != nil {
				r.radioLogger.Println("unable to verify parameter", parmName)
				continue
			}
			r.state.Vfo.Parameters[parmName] = cfmParm
		}
	}

	return nil
}

func (r *radio) updatePowerOn(pwrOn bool) error {

	if !r.rig.Caps.HasSetPowerStat || !r.rig.Caps.HasGetPowerStat {
		return errors.New("radio doesn't support set/get powerstat")
	}

	var pwrStat int
	if pwrOn {
		pwrStat = hl.RIG_POWER_ON
	} else {
		pwrStat = hl.RIG_POWER_OFF
	}

	if err := r.rig.SetPowerStat(pwrStat); err != nil {
		return err
	}

	// give the radio a little bit of time to turn on/off
	time.Sleep(time.Millisecond * 500)

	// verify powerstat
	cfmPwrStat, err := r.rig.GetPowerStat()
	if err != nil {
		return err
	}

	if cfmPwrStat == hl.RIG_POWER_OFF {
		r.state = sbRadio.State{}
		r.state.Vfo = &sbRadio.Vfo{}
		r.state.Channel = &sbRadio.Channel{}
		r.state.Vfo.Split = &sbRadio.Split{}
		r.state.Vfo.Levels = make(map[string]float32)
		r.state.Vfo.Parameters = make(map[string]float32)
		r.state.Vfo.Functions = make(map[string]bool)
	} else if cfmPwrStat == hl.RIG_POWER_ON {
		r.queryVfo()
	} else {
		r.radioLogger.Println("unknown powerstat", cfmPwrStat)
	}

	return nil
}

func (r *radio) updatePtt(ptt bool) error {
	vfo, _ := hl.VfoValue[r.state.CurrentVfo]
	var pttValue int
	if ptt {
		pttValue = hl.RIG_PTT_ON
	} else {
		pttValue = hl.RIG_PTT_OFF
	}

	err := r.rig.SetPtt(vfo, pttValue)
	if err != nil {
		return err
	}

	time.Sleep(30 * time.Millisecond)
	p, err := r.rig.GetPtt(vfo)
	if err != nil {
		return err
	}

	if p == hl.RIG_PTT_ON {
		r.state.Ptt = true
	} else {
		r.state.Ptt = false
	}

	return nil
}
