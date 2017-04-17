package radio

import sbRadio "github.com/dh1tw/gorigctl/sb_radio"

// Radio is the interface which is implemented by localRadio and remoteRadio
type Radio interface {
	GetCaps() (sbRadio.Capabilities, error)
	GetState() (sbRadio.State, error)
	GetFrequency() (float64, error)
	SetFrequency(freq float64) error
	GetMode() (string, int, error)
	SetMode(mode string, pbWidth int) error
	GetVfo() (string, error)
	SetVfo(string) error
	GetRit() (int, error)
	SetRit(rit int) error
	GetXit() (int, error)
	SetXit(xit int) error
	GetAntenna() (int, error)
	SetAntenna(int) error
	GetPtt() (bool, error)
	SetPtt(bool) error
	ExecVfoOps([]string) error
	GetTuningStep() (int, error)
	SetTuningStep(int) error
	GetPowerstat() (bool, error)
	SetPowerstat(bool) error
	GetSplitVfo() (string, bool, error)
	SetSplitVfo(string, bool) error
	GetSplitFrequency() (float64, error)
	SetSplitFrequency(float64) error
	GetSplitMode() (string, int, error)
	SetSplitMode(string, int) error
	GetSplitPbWidth() (int, error)
	SetSplitPbWidth(int) error
	SetSplitFrequencyMode(float64, string, int) error
	GetSplitFrequencyMode() (float64, string, int, error)
	GetFunction(string) (bool, error)
	SetFunction(string, bool) error
	GetLevel(string) (float32, error)
	SetLevel(string, float32) error
	GetParameter(string) (float32, error)
	SetParameter(string, float32) error
}
