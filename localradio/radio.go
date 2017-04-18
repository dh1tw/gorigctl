package localradio

import (
	"errors"
	"log"

	hl "github.com/dh1tw/goHamlib"
	sbRadio "github.com/dh1tw/gorigctl/sb_radio"
	"github.com/dh1tw/gorigctl/utils"
)

func (r *LocalRadio) GetCaps() (sbRadio.Capabilities, error) {

	caps := sbRadio.Capabilities{}
	caps.Vfos = r.rig.Caps.Vfos
	caps.Modes = r.rig.Caps.Modes
	caps.VfoOps = r.rig.Caps.Operations
	caps.GetFunctions = r.rig.Caps.GetFunctions
	caps.SetFunctions = r.rig.Caps.SetFunctions
	caps.GetLevels = utils.HlValuesToPbValues(r.rig.Caps.GetLevels)
	caps.SetLevels = utils.HlValuesToPbValues(r.rig.Caps.SetLevels)
	caps.GetParameters = utils.HlValuesToPbValues(r.rig.Caps.GetParameters)
	caps.SetParameters = utils.HlValuesToPbValues(r.rig.Caps.SetParameters)
	caps.MaxRit = int32(r.rig.Caps.MaxRit)
	caps.MaxXit = int32(r.rig.Caps.MaxXit)
	caps.MaxIfShift = int32(r.rig.Caps.MaxIfShift)
	caps.Filters = utils.HlMapToPbMap(r.rig.Caps.Filters)
	caps.TuningSteps = utils.HlMapToPbMap(r.rig.Caps.TuningSteps)
	caps.Preamps = utils.IntListToint32List(r.rig.Caps.Preamps)
	caps.Attenuators = utils.IntListToint32List(r.rig.Caps.Attenuators)
	caps.RigModel = int32(r.rig.Caps.RigModel)
	caps.ModelName = r.rig.Caps.ModelName
	caps.Version = r.rig.Caps.Version
	caps.MfgName = r.rig.Caps.MfgName
	caps.HasPowerstat = r.rig.Caps.HasGetPowerStat
	caps.HasPtt = r.rig.Caps.HasGetPtt
	caps.HasRit = r.rig.Caps.HasGetRit
	caps.HasXit = r.rig.Caps.HasGetXit
	caps.HasSplit = r.rig.Caps.HasGetSplitVfo
	caps.HasTs = r.rig.Caps.HasGetTs
	caps.HasAnt = r.rig.Caps.HasGetAnt

	return caps, nil
}

func (r *LocalRadio) GetState() (sbRadio.State, error) {
	return r.queryVfo()
}

func (r *LocalRadio) GetFrequency() (float64, error) {
	freq, err := r.rig.GetFreq(r.vfo)
	if err != nil {
		return 0, err
	}

	return freq, nil
}

func (r *LocalRadio) SetFrequency(freq float64) error {
	return r.rig.SetFreq(r.vfo, freq)
}

func (r *LocalRadio) GetMode() (string, int, error) {
	m, pbWidth, err := r.rig.GetMode(r.vfo)
	if err != nil {
		return "", 0, err
	}

	mode, ok := hl.ModeName[m]
	if !ok {
		return "", 0, err
	}

	return mode, pbWidth, nil
}

func (r *LocalRadio) SetMode(mode string, pbWidth int) error {
	m, ok := hl.ModeValue[mode]
	if !ok {
		return errors.New("unkown mode")
	}

	err := r.rig.SetMode(r.vfo, m, pbWidth)
	if err != nil {
		return err
	}

	return nil
}

func (r *LocalRadio) GetVfo() (string, error) {
	v, err := r.rig.GetVfo()
	if err != nil {
		return "", err
	}

	vfo, ok := hl.VfoName[v]
	if !ok {
		return "", errors.New("unknown vfo")
	}

	return vfo, nil
}

func (r *LocalRadio) SetVfo(vfo string) error {
	v, ok := hl.VfoValue[vfo]
	if !ok {
		return errors.New("unknown vfo")
	}

	err := r.rig.SetVfo(v)
	if err != nil {
		return err
	}

	r.vfo = v
	return nil
}

func (r *LocalRadio) GetRit() (int, error) {
	rit, err := r.rig.GetRit(r.vfo)
	if err != nil {
		return 0, err
	}
	return rit, err
}

func (r *LocalRadio) SetRit(rit int) error {
	return r.rig.SetRit(r.vfo, rit)
}

func (r *LocalRadio) GetXit() (int, error) {
	xit, err := r.rig.GetXit(r.vfo)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return xit, nil
}

func (r *LocalRadio) SetXit(xit int) error {
	return r.rig.SetXit(r.vfo, xit)
}

func (r *LocalRadio) GetAntenna() (int, error) {
	ant, err := r.rig.GetAnt(r.vfo)
	if err != nil {
		return 0, err
	}
	return ant, nil
}

func (r *LocalRadio) SetAntenna(ant int) error {
	return r.rig.SetAnt(r.vfo, ant)
}

func (r *LocalRadio) GetPtt() (bool, error) {
	ptt, err := r.rig.GetPtt(r.vfo)
	if err != nil {
		return false, err
	}
	if ptt == hl.RIG_PTT_ON {
		return true, nil
	}
	return false, nil
}

func (r *LocalRadio) SetPtt(ptt bool) error {
	p := hl.RIG_PTT_OFF
	if ptt {
		p = hl.RIG_PTT_ON
	}
	return r.rig.SetPtt(r.vfo, p)
}

func (r *LocalRadio) GetTuningStep() (int, error) {
	ts, err := r.rig.GetTs(r.vfo)
	if err != nil {
		return 0, err
	}
	return ts, nil
}

func (r *LocalRadio) SetTuningStep(ts int) error {
	return r.rig.SetTs(r.vfo, ts)
}

func (r *LocalRadio) GetPowerstat() (bool, error) {
	ps, err := r.rig.GetPowerStat()
	if err != nil {
		return false, err
	}
	if ps == hl.RIG_POWER_ON {
		return true, nil
	}
	return false, nil
}

func (r *LocalRadio) SetPowerstat(ps bool) error {
	p := hl.RIG_POWER_OFF
	if ps {
		p = hl.RIG_POWER_ON
	}

	return r.rig.SetPowerStat(p)
}

func (r *LocalRadio) ExecVfoOps(ops []string) error {
	for op := range ops {
		err := r.rig.VfoOp(r.vfo, op)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *LocalRadio) GetSplitVfo() (string, bool, error) {
	e, v, err := r.rig.GetSplitVfo(r.vfo)
	if err != nil {
		return "", false, err
	}
	enabled := false
	if e == 1 {
		enabled = true
	}

	vfo, ok := hl.VfoName[v]
	if !ok {
		return "", false, errors.New("unknown vfo")
	}

	return vfo, enabled, nil
}

func (r *LocalRadio) SetSplitVfo(vfo string, enabled bool) error {
	v, ok := hl.VfoValue[vfo]
	if !ok {
		return errors.New("unknown vfo")
	}
	return r.rig.SetSplitVfo(r.vfo, utils.Btoi(enabled), v)
}

func (r *LocalRadio) GetSplitFrequency() (float64, error) {
	_, v, err := r.rig.GetSplitVfo(r.vfo)
	if err != nil {
		return 0, err
	}

	freq, err := r.rig.GetSplitFreq(v)
	if err != nil {
		return 0, err
	}
	return freq, nil
}

func (r *LocalRadio) SetSplitFrequency(freq float64) error {
	_, v, err := r.rig.GetSplitVfo(r.vfo)
	if err != nil {
		return err
	}

	return r.rig.SetSplitFreq(v, freq*1000) // Hz
}

func (r *LocalRadio) GetSplitMode() (string, int, error) {

	_, v, err := r.rig.GetSplitVfo(r.vfo)
	if err != nil {
		return "", 0, err
	}

	m, pbWidth, err := r.rig.GetSplitMode(v)
	if err != nil {
		return "", 0, err
	}
	mode, ok := hl.ModeName[m]
	if !ok {
		return "", 0, errors.New("unknown mode")
	}
	return mode, pbWidth, nil
}

func (r *LocalRadio) SetSplitMode(mode string, pbWidth int) error {
	_, v, err := r.rig.GetSplitVfo(r.vfo)
	if err != nil {
		return err
	}

	m, ok := hl.ModeValue[mode]
	if !ok {
		return errors.New("unknown mode")
	}
	return r.rig.SetSplitMode(v, m, pbWidth)
}

func (r *LocalRadio) GetSplitPbWidth() (int, error) {
	_, v, err := r.rig.GetSplitVfo(r.vfo)
	if err != nil {
		return 0, err
	}

	_, pbWidth, err := r.rig.GetSplitMode(v)
	if err != nil {
		return 0, err
	}
	return pbWidth, nil
}

func (r *LocalRadio) SetSplitPbWidth(pbWidth int) error {

	_, v, err := r.rig.GetSplitVfo(r.vfo)
	if err != nil {
		return err
	}

	mode, _, err := r.rig.GetSplitMode(v)
	if err != nil {
		return err
	}

	return r.rig.SetSplitMode(r.vfo, mode, pbWidth)
}

func (r *LocalRadio) GetSplitFrequencyMode() (float64, string, int, error) {

	freq, err := r.GetSplitFrequency()
	if err != nil {
		return 0, "", 0, err
	}

	mode, pbWidth, err := r.GetSplitMode()
	if err != nil {
		return 0, "", 0, err
	}

	return freq, mode, pbWidth, nil

}

func (r *LocalRadio) SetSplitFrequencyMode(freq float64, mode string, pbWidth int) error {

	if err := r.SetSplitFrequency(freq); err != nil {
		return err
	}

	if err := r.SetSplitMode(mode, pbWidth); err != nil {
		return err
	}
	return nil
}

func (r *LocalRadio) GetFunction(function string) (bool, error) {

	f, ok := hl.FuncValue[function]
	if !ok {
		return false, errors.New("unknown function")
	}
	return r.rig.GetFunc(r.vfo, f)
}

func (r *LocalRadio) SetFunction(function string, value bool) error {

	f, ok := hl.FuncValue[function]
	if !ok {
		return errors.New("unknown function")
	}
	return r.rig.SetFunc(r.vfo, f, value)
}

func (r *LocalRadio) GetLevel(level string) (float32, error) {
	l, ok := hl.LevelValue[level]
	if !ok {
		return 0, errors.New("unknown level")
	}
	return r.rig.GetLevel(r.vfo, l)
}

func (r *LocalRadio) SetLevel(level string, value float32) error {
	l, ok := hl.LevelValue[level]
	if !ok {
		return errors.New("unknown level")
	}
	return r.rig.SetLevel(r.vfo, l, value)
}

func (r *LocalRadio) GetParameter(parm string) (float32, error) {
	p, ok := hl.LevelValue[parm]
	if !ok {
		return 0, errors.New("unknown parameter")
	}
	return r.rig.GetParm(r.vfo, p)

}

func (r *LocalRadio) SetParameter(parm string, value float32) error {
	p, ok := hl.LevelValue[parm]
	if !ok {
		return errors.New("unknown parameter")
	}
	return r.rig.SetParm(r.vfo, p, value)
}

func (r *LocalRadio) queryVfo() (sbRadio.State, error) {

	state := sbRadio.State{}
	state.Vfo = &sbRadio.Vfo{}
	state.Vfo.Levels = make(map[string]float32)
	state.Vfo.Parameters = make(map[string]float32)
	state.Vfo.Functions = make(map[string]bool)
	state.Vfo.Split = &sbRadio.Split{}
	state.Channel = &sbRadio.Channel{}

	if r.rig.Caps.HasGetPowerStat {
		pwrOn, err := r.rig.GetPowerStat()
		if err != nil {
			return state, err
		}
		if pwrOn == hl.RIG_POWER_ON {
			state.RadioOn = true
		} else {
			state.RadioOn = false
		}
	}

	// Only query radio if Power is On or if Radio has now PowerStat function
	// in this case we will assume that the radio is turned on
	if (r.rig.Caps.HasGetPowerStat && state.RadioOn) || !r.rig.Caps.HasGetPowerStat {

		vfo := hl.VfoValue["CURR"]

		if r.rig.Caps.HasGetVfo {
			vfo, err := r.GetVfo()
			if err != nil {
				return state, err
			}
			state.CurrentVfo = vfo
		} else {
			state.CurrentVfo = "CURR"
		}

		if r.rig.Caps.HasGetFreq {
			freq, err := r.GetFrequency()
			if err != nil {
				return state, err
			}
			state.Vfo.Frequency = freq
		}

		if r.rig.Caps.HasGetMode {
			mode, pbWidth, err := r.GetMode()
			if err != nil {
				return state, err
			}
			state.Vfo.Mode = mode
			state.Vfo.PbWidth = int32(pbWidth)
		}

		if r.rig.Caps.HasGetAnt {
			ant, err := r.GetAntenna()
			if err != nil {
				return state, err
			}
			state.Vfo.Ant = int32(ant)
		}

		if r.rig.Caps.HasGetRit {
			rit, err := r.rig.GetRit(vfo)
			if err != nil {
				return state, err
			} else {
				state.Vfo.Rit = int32(rit)
			}
		}

		if r.rig.Caps.HasGetRit {
			xit, err := r.rig.GetXit(vfo)
			if err != nil {
				return state, err
			}
			state.Vfo.Xit = int32(xit)
		}

		if r.rig.Caps.HasGetSplitVfo {
			txVfo, splitOn, err := r.GetSplitVfo()
			if err != nil {
				return state, err
			}

			state.Vfo.Split.Enabled = splitOn
			state.Vfo.Split.Vfo = txVfo

			if splitOn {

				// these checks should be enabled, but most of the
				// backends don't have these functions implemented
				// therefore they use the emulated functions which
				// unfortunately don't work everywhere well (e.g. TS-480)
				// if r.rig.Caps.HasGetSplitFreq {
				txFreq, txMode, txPbWidth, err := r.GetSplitFrequencyMode()
				if err != nil {
					return state, err
				}
				state.Vfo.Split.Frequency = txFreq
				state.Vfo.Split.Mode = txMode
				state.Vfo.Split.PbWidth = int32(txPbWidth)
			}
		}

		if r.rig.Caps.HasGetTs {
			tStep, err := r.rig.GetTs(vfo)
			if err != nil {
				return state, err
			}
			state.Vfo.TuningStep = int32(tStep)

		}

		for _, f := range r.rig.Caps.GetFunctions {
			fValue, err := r.GetFunction(f)
			if err != nil {
				return state, err
			}
			state.Vfo.Functions[f] = fValue
		}

		for _, level := range r.rig.Caps.GetLevels {
			lValue, err := r.GetLevel(level.Name)
			if err != nil {
				return state, err
			}
			state.Vfo.Levels[level.Name] = lValue
		}

		for _, param := range r.rig.Caps.GetParameters {
			pValue, err := r.GetParameter(param.Name)
			if err != nil {
				return state, err
			}
			state.Vfo.Parameters[param.Name] = pValue
		}
	}

	return state, nil
}
