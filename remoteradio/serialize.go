package remoteradio

import (
	"errors"

	"github.com/dh1tw/gorigctl/comms"
	sbRadio "github.com/dh1tw/gorigctl/sb_radio"
)

func (r *RemoteRadio) GetCaps() sbRadio.Capabilities {
	return r.caps
}

func (r *RemoteRadio) GetState() sbRadio.State {
	return r.state
}

func (r *RemoteRadio) GetFrequency() float64 {
	return r.state.Vfo.Frequency
}

func (r *RemoteRadio) SetFrequency(freq float64) error {
	req := r.initSetState()
	req.Vfo.Frequency = freq
	req.Md.HasFrequency = true
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetMode() (string, int) {
	return r.state.Vfo.Mode, int(r.state.Vfo.PbWidth)
}

func (r *RemoteRadio) SetMode(mode string, pbWidth int) error {
	req := r.initSetState()
	req.Md.HasMode = true
	req.Vfo.Mode = mode
	req.Vfo.PbWidth = int32(pbWidth)

	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetVfo() string {
	return r.state.CurrentVfo
}

func (r *RemoteRadio) SetVfo(vfo string) error {
	req := r.initSetState()
	req.CurrentVfo = vfo
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetRit() int {
	return int(r.state.Vfo.Rit)
}

func (r *RemoteRadio) SetRit(rit int) error {
	req := r.initSetState()
	req.Md.HasRit = true
	req.Vfo.Rit = int32(rit)
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetXit() int {
	return int(r.state.Vfo.Xit)
}

func (r *RemoteRadio) SetXit(xit int) error {
	req := r.initSetState()
	req.Md.HasXit = true
	req.Vfo.Xit = int32(xit)
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetAntenna() int {
	return int(r.state.Vfo.Ant)
}

func (r *RemoteRadio) SetAntenna(ant int) error {
	req := r.initSetState()
	req.Md.HasAnt = true
	req.Vfo.Ant = int32(ant)
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetPtt() bool {
	return r.state.Ptt
}

func (r *RemoteRadio) SetPtt(ptt bool) error {
	req := r.initSetState()
	req.Md.HasPtt = true
	req.Ptt = ptt
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetTuningStep() int {
	return int(r.state.Vfo.TuningStep)
}

func (r *RemoteRadio) SetTuningStep(ts int) error {
	req := r.initSetState()
	req.Md.HasTuningStep = true
	req.Vfo.TuningStep = int32(ts)
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetPowerstat() bool {
	return r.state.RadioOn
}

func (r *RemoteRadio) SetPowerstat(ps bool) error {
	req := r.initSetState()
	req.Md.HasRadioOn = true
	req.RadioOn = ps
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) ExecVfoOps(ops []string) error {
	req := r.initSetState()
	req.VfoOperations = ops
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetSplitVfo() (string, bool) {
	return r.state.Vfo.Split.Vfo, r.state.Vfo.Split.Enabled
}

func (r *RemoteRadio) SetSplitVfo(vfo string, enabled bool) error {
	req := r.initSetState()
	req.Md.HasSplit = true
	req.Vfo.Split.Enabled = enabled
	req.Vfo.Split.Vfo = vfo
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetSplitFrequency() float64 {
	return r.state.Vfo.Split.Frequency
}

func (r *RemoteRadio) SetSplitFrequency(freq float64) error {
	req := r.initSetState()
	req.Md.HasSplit = true
	req.Vfo.Split.Enabled = r.state.Vfo.Split.Enabled
	req.Vfo.Split.Vfo = r.state.Vfo.Split.Vfo
	req.Vfo.Split.Frequency = freq
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetSplitMode() (string, int) {
	return r.state.Vfo.Split.Mode, int(r.state.Vfo.Split.PbWidth)
}

func (r *RemoteRadio) SetSplitMode(mode string, pbWidth int) error {
	req := r.initSetState()
	req.Md.HasSplit = true
	req.Vfo.Split.Enabled = r.state.Vfo.Split.Enabled
	req.Vfo.Split.Mode = mode
	req.Vfo.Split.PbWidth = int32(pbWidth)
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetSplitPbWidth() int {
	return int(r.state.Vfo.Split.PbWidth)
}

func (r *RemoteRadio) SetSplitPbWidth(pbWidth int) error {
	req := r.initSetState()
	req.Md.HasSplit = true
	req.Vfo.Split.Enabled = r.state.Vfo.Split.Enabled
	req.Vfo.Split.Mode = r.state.Vfo.Split.Mode
	req.Vfo.Split.PbWidth = int32(pbWidth)
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetSplitFrequencyMode() (float64, string, int) {
	return r.state.Vfo.Split.Frequency, r.state.Vfo.Split.Mode, int(r.state.Vfo.Split.PbWidth)
}

func (r *RemoteRadio) SetSplitFrequencyMode(freq float64, mode string, pbWidth int) error {
	req := r.initSetState()
	req.Md.HasSplit = true
	req.Vfo.Split.Enabled = r.state.Vfo.Split.Enabled
	req.Vfo.Split.Frequency = freq
	req.Vfo.Split.Mode = mode
	req.Vfo.Split.PbWidth = int32(pbWidth)
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetFunction(function string) (bool, error) {
	value, ok := r.state.Vfo.Functions[function]
	if !ok {
		return false, errors.New("unsupported function")
	}
	return value, nil
}

func (r *RemoteRadio) SetFunction(function string, value bool) error {
	req := r.initSetState()
	req.Md.HasFunctions = true
	req.Vfo.Functions = make(map[string]bool)
	req.Vfo.Functions[function] = value
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetLevel(level string) (float32, error) {
	value, ok := r.state.Vfo.Levels[level]
	if !ok {
		return 0, errors.New("unsupported level")
	}
	return value, nil
}

func (r *RemoteRadio) SetLevel(level string, value float32) error {
	req := r.initSetState()
	req.Md.HasLevels = true
	req.Vfo.Levels = make(map[string]float32)
	req.Vfo.Levels[level] = value
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) GetParameter(parm string) (float32, error) {
	value, ok := r.state.Vfo.Parameters[parm]
	if !ok {
		return 0, errors.New("unsupported parameter")
	}
	return value, nil
}

func (r *RemoteRadio) SetParameter(parm string, value float32) error {
	req := r.initSetState()
	req.Md.HasParameters = true
	req.Vfo.Parameters = make(map[string]float32)
	req.Vfo.Parameters[parm] = value
	return r.sendCatRequest(req)
}

func (r *RemoteRadio) sendCatRequest(req sbRadio.SetState) error {
	data, err := req.Marshal()
	if err != nil {
		return err
	}

	msg := comms.IOMsg{}
	msg.Data = data
	msg.Topic = r.catRequestTopic
	msg.Retain = false
	msg.Qos = 0

	r.toWireCh <- msg

	return nil
}
