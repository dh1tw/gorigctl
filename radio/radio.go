package radio

import sbRadio "github.com/dh1tw/gorigctl/sb_radio"

// Radio is the interface which is implemented by localRadio and remoteRadio
type Radio interface {
	GetCaps() sbRadio.Capabilities
	GetState() sbRadio.State
	GetFrequency() float64
	SetFrequency(freq float64) error
	GetMode() (string, int)
	SetMode(mode string, pbWidth int) error
	GetVfo() string
	SetVfo(string) error
	GetRit() int
	SetRit(rit int) error
	GetXit() int
	SetXit(xit int) error
	GetAntenna() int
	SetAntenna(int) error
	GetPtt() bool
	SetPtt(bool) error
	ExecVfoOps([]string) error
	GetTuningStep() int
	SetTuningStep(int) error
	GetPowerstat() bool
	SetPowerstat(bool) error
	GetSplitVfo() (string, bool)
	SetSplitVfo(string, bool) error
	GetSplitFrequency() float64
	SetSplitFrequency(float64) error
	GetSplitMode() (string, int)
	SetSplitMode(string, int) error
	GetSplitPbWidth() int
	SetSplitPbWidth(int) error
	SetSplitFrequencyMode(float64, string, int) error
	GetSplitFrequencyMode() (float64, string, int)
	GetFunction(string) (bool, error)
	SetFunction(string, bool) error
	GetLevel(string) (float32, error)
	SetLevel(string, float32) error
	GetParameter(string) (float32, error)
	SetParameter(string, float32) error
}
