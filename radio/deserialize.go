package radio

import (
	"errors"
	"reflect"

	"time"

	"math"

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
				r.logger.Println(err)
			} else {
				if r.state.RadioOn {
					r.queryVfo()
				}
			}
		}
	}

	// if r.state.RadioOn {

	if ns.CurrentVfo != r.state.CurrentVfo {
		r.logger.Println("updating CurrentVFO to", ns.CurrentVfo)
		if err := r.updateCurrentVfo(ns.CurrentVfo); err != nil {
			r.logger.Println(err)
		}
	}

	if len(ns.VfoOperations) > 0 {
		if err := r.execVfoOperations(ns.GetVfoOperations()); err != nil {
			r.logger.Println(err)
		}
	}

	if ns.Md.HasFrequency {
		if ns.Vfo.GetFrequency() != r.state.Vfo.Frequency {
			r.logger.Println("updating Frequency to", ns.Vfo.Frequency)
			if err := r.updateFrequency(ns.Vfo.GetFrequency()); err != nil {
				r.logger.Println(err)
			}
		}
	}

	if ns.Md.HasMode {
		if ns.Vfo.GetMode() != r.state.Vfo.Mode {
			r.logger.Println("updating Mode to", ns.Vfo.Mode)
			if err := r.updateMode(ns.Vfo.GetMode(), ns.Vfo.GetPbWidth()); err != nil {
				r.logger.Println(err)
			}
		}
	}

	if ns.Md.HasPbWidth {
		if ns.Vfo.GetPbWidth() != r.state.Vfo.PbWidth {
			r.logger.Println("updating PbWidth to", ns.Vfo.PbWidth)
			if err := r.updatePbWidth(ns.Vfo.GetPbWidth()); err != nil {
				r.logger.Println(err)
			}
		}
	}

	if ns.Md.HasAnt {
		if ns.Vfo.GetAnt() != r.state.Vfo.Ant {
			r.logger.Println("updating Antenna to", ns.Vfo.Ant)
			if err := r.updateAntenna(ns.Vfo.GetAnt()); err != nil {
				r.logger.Println(err)
			}
		}
	}

	if ns.Md.HasRit {
		if ns.Vfo.GetRit() != r.state.Vfo.Rit {
			r.logger.Println("updating Rit to", ns.Vfo.Rit)
			if err := r.updateRit(ns.Vfo.GetRit()); err != nil {
				r.logger.Println(err)
			}
		}
	}

	if ns.Md.HasXit {
		if ns.Vfo.GetXit() != r.state.Vfo.Xit {
			r.logger.Println("updating Xit to", ns.Vfo.Xit)
			if err := r.updateXit(ns.Vfo.GetXit()); err != nil {
				r.logger.Println(err)
			}
		}
	}

	if ns.Md.HasSplit {
		if ns.Vfo.Split != nil {
			if !reflect.DeepEqual(ns.Vfo.Split, r.state.Vfo.Split) {
				r.logger.Println("updating Split to", ns.Vfo.Split)
				if err := r.updateSplit(ns.Vfo.Split); err != nil {
					r.logger.Println(err)
				}
			}
		}
	}

	if ns.Md.HasTuningStep {
		if ns.Vfo.GetTuningStep() != r.state.Vfo.TuningStep {
			r.logger.Println("updating TS to", ns.Vfo.TuningStep)
			if err := r.updateTs(ns.Vfo.GetTuningStep()); err != nil {
				r.logger.Println(err)
			}
		}
	}

	if ns.Md.HasFunctions {
		if ns.Vfo.Functions != nil {
			if !reflect.DeepEqual(ns.Vfo.Functions, r.state.Vfo.Functions) {
				if err := r.updateFunctions(ns.Vfo.GetFunctions()); err != nil {
					r.logger.Println(err)
				}
			}
		}
	}

	if ns.Md.HasLevels {
		if ns.Vfo.Levels != nil {
			if !reflect.DeepEqual(ns.Vfo.Levels, r.state.Vfo.Levels) {
				if err := r.updateLevels(ns.Vfo.GetLevels()); err != nil {
					r.logger.Println(err)
				}
			}
		}
	}

	if ns.Md.HasParameters {
		if ns.Vfo.Parameters != nil {
			if !reflect.DeepEqual(ns.Vfo.Parameters, r.state.Vfo.Parameters) {
				if err := r.updateParams(ns.Vfo.GetParameters()); err != nil {
					r.logger.Println(err)
				}
			}
		}
	}
	// }

	if ns.Md.HasPtt {
		if ns.GetPtt() != r.state.Ptt {
			r.logger.Println("updating PTT to", ns.Ptt)
			if err := r.updatePtt(ns.GetPtt()); err != nil {
				r.logger.Println(err)
			}
		}
	}

	if ns.Md.HasPollingInterval {
		if ns.GetPollingInterval() != r.state.PollingInterval {
			if ns.GetPollingInterval() > 0 {
				newPollingInterval := time.Millisecond * time.Duration(ns.GetPollingInterval())
				r.pollingTicker.Stop()
				r.pollingTicker = time.NewTicker(newPollingInterval)
				r.state.PollingInterval = ns.GetPollingInterval()
			} else {
				r.pollingTicker.Stop()
				r.state.PollingInterval = 0
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
	err := r.rig.SetFreq(vfo, newFreq)
	if err != nil {
		return err
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

			r.queryVfo()
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

	ts := 0
	if r.rig.Caps.HasGetTs {
		ts, err = r.rig.GetTs(vfo)
		if err != nil {
			return err
		}
	}

	r.state.Vfo.TuningStep = int32(ts)
	r.state.Vfo.Mode = hl.ModeName[mode]
	r.state.Vfo.PbWidth = int32(pbWidth)

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

	ts := 0
	if r.rig.Caps.HasGetTs {
		ts, err = r.rig.GetTs(vfo)
		if err != nil {
			return err
		}
	}

	r.state.Vfo.TuningStep = int32(ts)
	r.state.Vfo.Mode = hl.ModeName[mode]
	r.state.Vfo.PbWidth = int32(pbWidth)

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
	rit, err := r.rig.GetRit(vfo)
	if err != nil {
		return err
	}

	xit, err := r.rig.GetXit(vfo)
	if err != nil {
		return err
	}

	freq, err := r.rig.GetFreq(vfo)
	if err != nil {
		return err
	}

	r.state.Vfo.Frequency = freq
	r.state.Vfo.Rit = int32(rit)
	r.state.Vfo.Xit = int32(xit)

	return nil
}

func (r *radio) updateXit(newXit int32) error {
	vfo, _ := hl.VfoValue[r.state.CurrentVfo]
	err := r.rig.SetXit(vfo, int(newXit))
	if err != nil {
		return err
	}

	xit, err := r.rig.GetXit(vfo)
	if err != nil {
		return err
	}

	rit, err := r.rig.GetRit(vfo)
	if err != nil {
		return err
	}

	r.state.Vfo.Xit = int32(xit)
	r.state.Vfo.Rit = int32(rit)

	return nil
}

func (r *radio) updateSplit(newSplit *sbRadio.Split) error {

	vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	if !r.rig.Caps.HasGetSplitVfo || !r.rig.Caps.HasSetSplitVfo {
		return errors.New("radio doesn't support split")
	}

	if newSplit.GetEnabled() != r.state.Vfo.Split.Enabled {

		txVfo, ok := hl.VfoValue[newSplit.GetVfo()]
		if !ok {
			return errors.New("unknown tx (split) vfo")
		}

		err := r.rig.SetSplitVfo(vfo, utils.Btoi(newSplit.GetEnabled()), txVfo)
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

	if !r.state.Vfo.Split.Enabled {
		// unset split frequency, mode, pbWidth
		r.state.Vfo.Split.Frequency = 0
		r.state.Vfo.Split.Mode = ""
		r.state.Vfo.Split.Vfo = ""
		r.state.Vfo.Split.PbWidth = 0
		return nil
	}

	if r.rig.Caps.HasGetSplitFreq && r.rig.Caps.HasSetSplitFreq {
		if newSplit.GetFrequency() != r.state.Vfo.Split.Frequency &&
			newSplit.GetFrequency() > 0 {

			txVfo, ok := hl.VfoValue[r.state.Vfo.Split.Vfo]
			if !ok {
				return errors.New("unknown tx (split) vfo")
			}

			err := r.rig.SetSplitFreq(txVfo, newSplit.GetFrequency())
			if err != nil {
				return err
			}
		}

		// verify and copy the radio's split frequency
		txVfo, ok := hl.VfoValue[r.state.Vfo.Split.Vfo]
		if !ok {
			return errors.New("unknown tx (split) vfo")
		}
		txFreq, err := r.rig.GetSplitFreq(txVfo)
		if err != nil {
			return err
		}
		r.state.Vfo.Split.Frequency = txFreq
	}

	// this check should be performed as some of the rigs (e.g. TS-480)
	// don't work well with the fallback functions

	// if r.rig.Caps.HasGetSplitMode && r.rig.Caps.HasSetSplitMode {
	if newSplit.GetMode() != r.state.Vfo.Split.Mode &&
		len(newSplit.GetMode()) > 0 {

		txVfo, ok := hl.VfoValue[r.state.Vfo.Split.Vfo]
		if !ok {
			return errors.New("unknown tx (split) vfo")
		}

		newSplitModeValue, ok := hl.ModeValue[newSplit.GetMode()]
		if !ok {
			return errors.New("unknown tx (split) mode")
		}

		pbWidth := r.state.Vfo.Split.PbWidth
		if newSplit.GetPbWidth() > 0 {
			pbWidth = newSplit.GetPbWidth()
		}

		err := r.rig.SetSplitMode(txVfo, newSplitModeValue, int(pbWidth))
		if err != nil {
			// get the standard filter width for this radio
			pbNormal, err := r.rig.GetPbNormal(newSplitModeValue)
			if err != nil {
				return err
			}
			// try again with the standard filter width
			err = r.rig.SetSplitMode(txVfo, newSplitModeValue, pbNormal)
			if err != nil {
				return err
			}
		}
	}

	// verify and copy the radio's txMode and txPbWidth
	txVfo, ok := hl.VfoValue[r.state.Vfo.Split.Vfo]
	if !ok {
		return errors.New("unknown VFO")
	}
	txMode, txPbWidth, err := r.rig.GetSplitMode(txVfo)
	if err != nil {
		return err
	}
	r.state.Vfo.Split.Mode = hl.ModeName[txMode]
	r.state.Vfo.Split.PbWidth = int32(txPbWidth)

	return nil
	// }

	// we only reach this code if the mode is the same, but we want
	// to update the filter width
	if r.rig.Caps.HasGetSplitMode && r.rig.Caps.HasSetSplitMode {
		if newSplit.GetPbWidth() != r.state.Vfo.Split.PbWidth &&
			len(newSplit.GetMode()) > 0 {

			txVfo, ok := hl.VfoValue[r.state.Vfo.Split.Vfo]
			if !ok {
				return errors.New("unknown VFO")
			}

			splitModeValue := hl.ModeValue[newSplit.GetMode()]
			err := r.rig.SetSplitMode(txVfo, splitModeValue, int(newSplit.GetPbWidth()))
			if err != nil {
				return err
			}
		}

		// verify and copy the radio's mode and txPbWidth
		txVfo, ok = hl.VfoValue[r.state.Vfo.Split.Vfo]
		if !ok {
			return errors.New("unknown VFO")
		}
		txMode, txPbWidth, err := r.rig.GetSplitMode(txVfo)
		if err != nil {
			return err
		}
		r.state.Vfo.Split.Mode = hl.ModeName[txMode]
		r.state.Vfo.Split.PbWidth = int32(txPbWidth)
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

func (r *radio) updateFunctions(newFuncs []string) error {
	vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	// functions to be enabled
	diff := utils.SliceDiff(newFuncs, r.state.Vfo.Functions)
	for _, f := range diff {
		funcValue, ok := hl.FuncValue[f]
		if !ok {
			return errors.New("unknown function")
		}
		err := r.rig.SetFunc(vfo, funcValue, true)
		if err != nil {
			return err
		}

		r.state.Vfo.Functions = append(r.state.Vfo.Functions, f)
	}

	// functions to be disabled
	diff = utils.SliceDiff(r.state.Vfo.Functions, newFuncs)
	for _, f := range diff {
		funcValue, ok := hl.FuncValue[f]
		if !ok {
			return errors.New("unknown function")
		}
		err := r.rig.SetFunc(vfo, funcValue, false)
		if err != nil {
			return err
		}

		r.state.Vfo.Functions = utils.RemoveStringFromSlice(f, r.state.Vfo.Functions)
	}

	return nil
}

func (r *radio) updateLevels(newLevels map[string]float32) error {
	vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	for k, v := range newLevels {
		levelValue, ok := hl.LevelValue[k]
		if !ok {
			return errors.New("unknown Level")
		}
		if _, ok := r.state.Vfo.Levels[k]; !ok {
			return errors.New("unsupported Level for this rig")
		}

		if r.state.Vfo.Levels[k] != v {
			err := r.rig.SetLevel(vfo, levelValue, v)
			if err != nil {
				return nil
			}

			// lets verify if the value has been set
			gl, err := r.rig.HasGetLevel(levelValue)
			if err != nil {
				return err
			}

			if gl == levelValue {
				gv, err := r.rig.GetLevel(vfo, levelValue)
				if err != nil {
					return err
				}
				// allow some rounding when comparing floats
				if math.Abs(float64(gv))-math.Abs(float64(v)) > 0.1 {
					r.logger.Printf("WARN: Level %s could not be set", k)
				}
				r.state.Vfo.Levels[k] = gv
			}

		}
	}

	return nil
}

func (r *radio) updateParams(newParams map[string]float32) error {
	vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	for k, v := range newParams {
		paramValue, ok := hl.ParmValue[k]
		if !ok {
			return errors.New("unknown Parameter")
		}
		if _, ok := r.state.Vfo.Parameters[k]; !ok {
			return errors.New("unsupported Parameter for this rig")
		}
		if r.state.Vfo.Levels[k] != v {
			err := r.rig.SetLevel(vfo, paramValue, v)
			if err != nil {
				return nil
			}
		}
	}

	return nil
}

func (r *radio) updatePowerOn(pwrOn bool) error {

	var pwrStat int
	if pwrOn {
		pwrStat = hl.RIG_POWER_ON
	} else {
		pwrStat = hl.RIG_POWER_OFF
	}

	err := r.rig.SetPowerStat(pwrStat)
	if err != nil {
		return err
	}

	r.state.RadioOn = pwrOn

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

	r.logger.Println("ptt should be:", ptt)
	r.logger.Println("ptt is:", p)

	if p == hl.RIG_PTT_ON {
		r.state.Ptt = true
	} else {
		r.state.Ptt = false
	}

	return nil
}
