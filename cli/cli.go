package cli

import (
	"errors"
	"html/template"
	"math"
	"strconv"
	"strings"

	"github.com/dh1tw/gorigctl/radio"
	sbRadio "github.com/dh1tw/gorigctl/sb_radio"
	"github.com/dh1tw/gorigctl/utils"
)

type CliCmd struct {
	Cmd         func(r radio.Radio, args []string)
	Name        string
	Shortcut    string
	Parameters  string
	Description string
	Example     string
}

func PopulateCliCmds() []CliCmd {

	cliCmds := make([]CliCmd, 0, 40)

	cliSetFrequency := CliCmd{
		Cmd:         setFrequency,
		Name:        "set_freq",
		Shortcut:    "F",
		Parameters:  "Frequency [kHz]",
		Description: "Set Frequency for the current VFO",
		Example:     "F 14250.000",
	}

	cliCmds = append(cliCmds, cliSetFrequency)

	cliGetFrequency := CliCmd{
		Cmd:         getFrequency,
		Name:        "get_freq",
		Shortcut:    "f",
		Description: "Frequency [kHz] of current VFO",
	}

	cliCmds = append(cliCmds, cliGetFrequency)

	cliSetMode := CliCmd{
		Cmd:         setMode,
		Name:        "set_mode",
		Shortcut:    "M",
		Parameters:  "Mode and optionally Filter bandwidth [Hz]",
		Description: "Set Mode and optionally Filter Bandwidth for the current VFO)",
		Example:     "M USB 2400",
	}

	cliCmds = append(cliCmds, cliSetMode)

	cliGetMode := CliCmd{
		Cmd:         getMode,
		Name:        "get_mode",
		Shortcut:    "m",
		Description: "Get Mode",
	}

	cliCmds = append(cliCmds, cliGetMode)

	cliSetVfo := CliCmd{
		Cmd:         setVfo,
		Name:        "set_vfo",
		Shortcut:    "V",
		Parameters:  "VFO Name",
		Description: "Change to another VFO",
		Example:     "V VFOB",
	}

	cliCmds = append(cliCmds, cliSetVfo)

	cliGetVfo := CliCmd{
		Cmd:         getVfo,
		Name:        "get_vfo",
		Shortcut:    "v",
		Description: "Get Vfo",
	}

	cliCmds = append(cliCmds, cliGetVfo)

	cliSetRit := CliCmd{
		Cmd:         setRit,
		Name:        "set_rit",
		Shortcut:    "J",
		Parameters:  "RX Offset [Hz]",
		Description: "Set RX Offset (0 = Off)",
		Example:     "J -500",
	}

	cliCmds = append(cliCmds, cliSetRit)

	cliGetRit := CliCmd{
		Cmd:         getRit,
		Name:        "get_rit",
		Shortcut:    "j",
		Description: "Get Rit [Hz]",
	}

	cliCmds = append(cliCmds, cliGetRit)

	cliSetXit := CliCmd{
		Cmd:         setXit,
		Name:        "set_xit",
		Shortcut:    "Z",
		Description: "Set TX Offset (0 = Off)",
		Parameters:  "TX Offset [Hz]",
		Example:     "Z -500",
	}

	cliCmds = append(cliCmds, cliSetXit)

	cliGetXit := CliCmd{
		Cmd:         getXit,
		Name:        "get_xit",
		Shortcut:    "z",
		Description: "Get Xit [Hz]",
	}

	cliCmds = append(cliCmds, cliGetXit)

	cliSetAnt := CliCmd{
		Cmd:         setAnt,
		Name:        "set_ant",
		Shortcut:    "y",
		Parameters:  "Antenna",
		Description: "Set Antenna",
		Example:     "Y 2",
	}

	cliCmds = append(cliCmds, cliSetAnt)

	cliGetAnt := CliCmd{
		Cmd:         getAnt,
		Name:        "get_ant",
		Shortcut:    "y",
		Description: "Get Antenna",
	}

	cliCmds = append(cliCmds, cliGetAnt)

	cliSetPtt := CliCmd{
		Cmd:         setPtt,
		Name:        "set_ptt",
		Shortcut:    "t",
		Parameters:  "Ptt [true, t, 1, false, f, 0]",
		Description: "Set Transmit on/off",
		Example:     "t 1",
	}

	cliCmds = append(cliCmds, cliSetPtt)

	cliGetPtt := CliCmd{
		Cmd:         getPtt,
		Name:        "get_ptt",
		Shortcut:    "y",
		Description: "Get Ptt",
	}

	cliCmds = append(cliCmds, cliGetPtt)

	cliExecVfoOp := CliCmd{
		Cmd:         execVfoOp,
		Name:        "vfo_op",
		Shortcut:    "G",
		Parameters:  "VFO Operation",
		Description: "Execute a VFO Operation",
		Example:     "G XCHG",
	}

	cliCmds = append(cliCmds, cliExecVfoOp)

	cliSetFunction := CliCmd{
		Cmd:         setFunction,
		Name:        "set_func",
		Shortcut:    "U",
		Parameters:  "Function [string], Value [bool]",
		Description: "Set a Rig function",
		Example:     "U NB 1",
	}

	cliCmds = append(cliCmds, cliSetFunction)

	// cliGetFunction := CliCmd{
	// 	Cmd:         getFunctionsPlain,
	// 	Name:        "get_func",
	// 	Shortcut:    "u",
	// 	Description: "List the activated functions",
	// }

	// cliCmds = append(cliCmds, cliGetFunction)

	cliSetLevel := CliCmd{
		Cmd:         setLevel,
		Name:        "set_level",
		Shortcut:    "L",
		Parameters:  "Level & Value",
		Description: "Set a Level",
		Example:     "L CWPITCH 500",
	}

	cliCmds = append(cliCmds, cliSetLevel)

	// cliGetLevelsPlain := CliCmd{
	// 	Cmd:         getLevelsPlain,
	// 	Name:        "get_level",
	// 	Shortcut:    "l",
	// 	Description: "Lists all available levels",
	// }

	// cliCmds = append(cliCmds, cliGetLevelsPlain)

	cliSetTuningStep := CliCmd{
		Cmd:         setTuningStep,
		Name:        "set_ts",
		Shortcut:    "N",
		Parameters:  "Tuning Step [Hz]",
		Description: "Set the tuning step of the radio",
		Example:     "N 1000",
	}

	cliCmds = append(cliCmds, cliSetTuningStep)

	cliGetTuningStep := CliCmd{
		Cmd:         getTuningStep,
		Name:        "get_ts",
		Shortcut:    "n",
		Description: "Get the current tuning step [Hz]",
	}

	cliCmds = append(cliCmds, cliGetTuningStep)

	cliSetPowerStat := CliCmd{
		Cmd:         setPowerStat,
		Name:        "set_powerstat",
		Shortcut:    "",
		Parameters:  "Rig Power Status [true, t, 1, false, f, 0]",
		Description: "Turn the radio on/off",
		Example:     "set_powerstat 1",
	}

	cliCmds = append(cliCmds, cliSetPowerStat)

	cliGetPowerStat := CliCmd{
		Cmd:         getPowerStat,
		Name:        "get_powerstat",
		Shortcut:    "",
		Description: "Get the power status of the radio (On/Off)",
	}

	cliCmds = append(cliCmds, cliGetPowerStat)

	cliSetSplit := CliCmd{
		Cmd:         setSplitVfo,
		Name:        "set_split",
		Shortcut:    "S",
		Parameters:  "Split VFO [true, t, 1, false, f, 0]",
		Description: "Turn Split On/Off for a VFO",
		Example:     "S 1 VFOB",
	}

	cliCmds = append(cliCmds, cliSetSplit)

	cliGetSplit := CliCmd{
		Cmd:         getSplit,
		Name:        "get_split",
		Shortcut:    "s",
		Description: "Get the split status (if enabled: VFO, Frequency, Mode, Filter)",
	}

	cliCmds = append(cliCmds, cliGetSplit)

	cliSetSplitFrequency := CliCmd{
		Cmd:         setSplitFreq,
		Name:        "set_split_freq",
		Shortcut:    "I",
		Parameters:  "TX Frequency [kHz]",
		Description: "Set the TX Split Frequency (the Split VFO will be determined automatically)",
		Example:     "I 14205000",
	}

	cliCmds = append(cliCmds, cliSetSplitFrequency)

	cliSetSplitMode := CliCmd{
		Cmd:         setSplitMode,
		Name:        "set_split_mode",
		Shortcut:    "X",
		Parameters:  "TX Mode and optionally Filter bandwidth [Hz]",
		Description: "Set the TX Split Mode (optionally with Bandwidth [Hz])",
		Example:     "X CW 200",
	}

	cliCmds = append(cliCmds, cliSetSplitMode)

	cliSetSplitFreqMode := CliCmd{
		Cmd:         setSplitFreqMode,
		Name:        "set_split_freq_mode",
		Shortcut:    "K",
		Parameters:  "TX Frequency [kHz], TX Mode and optionally Filter BW [Hz]",
		Description: "Set the Split Tx Frequency, Mode (optionally with Bandwidth [Hz])",
		Example:     "K 7170000 AM 6000",
	}

	cliCmds = append(cliCmds, cliSetSplitFreqMode)

	// cliSetPollingInterval := CliCmd{
	// 	Cmd:         setPollingInterval,
	// 	Name:        "set_polling_interval",
	// 	Shortcut:    "",
	// 	Parameters:  "Polling rate [ms]",
	// 	Description: "Set the polling interval for updating the meter values (SWR, ALC, Field Strength...)",
	// 	Example:     "set_polling_interval 50",
	// }

	// cliCmds = append(cliCmds, cliSetPollingInterval)

	// cliGetPollingInterval := CliCmd{
	// 	Cmd:         getPollingInterval,
	// 	Name:        "get_polling_interval",
	// 	Shortcut:    "",
	// 	Description: "Get the polling interval for updating the meter values (SWR, ALC, Field Strength...)",
	// }

	// cliCmds = append(cliCmds, cliGetPollingInterval)

	// cliGetSyncInterval := CliCmd{
	// 	Cmd:         getSyncInterval,
	// 	Name:        "get_sync_interval",
	// 	Shortcut:    "",
	// 	Description: "Get the interval for synchronizing all radio values",
	// }

	// cliCmds = append(cliCmds, cliGetSyncInterval)

	// cliSetSyncInterval := CliCmd{
	// 	Cmd:         setSyncInterval,
	// 	Name:        "set_sync_interval",
	// 	Shortcut:    "",
	// 	Parameters:  "Sync rate [s]",
	// 	Description: "Set the interval for synchronizing all radio values",
	// 	Example:     "set_sync_interval 5",
	// }

	// cliCmds = append(cliCmds, cliSetSyncInterval)

	// cliSetPrintUpdates := CliCmd{
	// 	Cmd:         setPrintRigUpdates,
	// 	Name:        "set_print_rig_updates",
	// 	Parameters:  "[true, t, 1, false, f, 0]",
	// 	Shortcut:    "",
	// 	Description: "Print rig values which have changed",
	// }

	// cliCmds = append(cliCmds, cliSetPrintUpdates)

	return cliCmds
}

func getFrequency(r radio.Radio, args []string) {
	r.Printf("Frequency: %.3f kHz\n", r.GetFrequency()/1000)
}

func setFrequency(r radio.Radio, args []string) {

	if err := checkArgs(args, 1); err != nil {
		r.Print(err)
		return
	}

	freq, err := strconv.ParseFloat(args[0], 10)
	if err != nil {
		r.Print("ERROR: frequency [kHz] must be float")
		return
	}

	// multiply frequency with 1000 since we enter kHz
	// but send over the wire Hz
	if err := r.SetFrequency(freq * 1000); err != nil {
		r.Print(err)
	}
}

func getMode(r radio.Radio, args []string) {
	mode, pbWidth := r.GetMode()
	r.Print("Mode:", mode)
	r.Printf("Filter: %dHz\n", pbWidth)
}

func setMode(r radio.Radio, args []string) {

	if len(args) < 1 || len(args) > 2 {
		r.Print("ERROR: wrong number of arguments")
		return
	}

	mode := strings.ToUpper(args[0])

	caps := r.GetCaps()
	pbWidth := 0

	if ok := utils.StringInSlice(mode, caps.Modes); !ok {
		r.Print("ERROR: unsupported mode")
		return
	}

	if len(args) == 2 {

		pbWidth, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			r.Print("ERROR: filter width [Hz] must be integer")
			pbWidth = 0
		}

		filters, ok := caps.Filters[mode]
		if !ok {
			r.Print("WARN: no filters found for this mode in rig caps")
		} else {
			if ok := utils.Int32InSlice(int32(pbWidth), filters.Value); !ok {
				r.Print("WARN: unspported passband width")
			}
		}
	}

	if err := r.SetMode(mode, pbWidth); err != nil {
		r.Print(err)
	}
}

func getVfo(r radio.Radio, args []string) {
	r.Print("Current Vfo:", r.GetVfo())
}

func setVfo(r radio.Radio, args []string) {
	if err := checkArgs(args, 1); err != nil {
		r.Print(err)
		return
	}

	vfo := strings.ToUpper(args[0])

	caps := r.GetCaps()

	if ok := utils.StringInSlice(vfo, caps.Vfos); !ok {
		r.Print("ERROR: unsupported vfo")
		return
	}

	if err := r.SetVfo(vfo); err != nil {
		r.Print(err)
	}
}

func getRit(r radio.Radio, args []string) {
	r.Printf("Rit: %d Hz\n", r.GetRit())
}

func setRit(r radio.Radio, args []string) {

	if err := checkArgs(args, 1); err != nil {
		r.Print(err)
		return
	}

	rit, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		r.Print("ERROR: rit value [Hz] must be integer")
		return
	}

	caps := r.GetCaps()

	if math.Abs(float64(rit)) > float64(caps.MaxRit) {
		r.Print("WARN: rit value larger than supported by rig")
	}

	if err := r.SetRit(int(rit)); err != nil {
		r.Print(err)
	}
}

func getXit(r radio.Radio, args []string) {
	r.Printf("Xit: %d Hz\n", r.GetXit())
}

func setXit(r radio.Radio, args []string) {

	if err := checkArgs(args, 1); err != nil {
		r.Print(err)
		return
	}

	xit, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		r.Print("ERROR: xit value [Hz] must be integer")
		return
	}

	caps := r.GetCaps()

	if math.Abs(float64(xit)) > float64(caps.MaxXit) {
		r.Print("WARN: xit value larger than supported by rig")
	}

	if err := r.SetXit(int(xit)); err != nil {
		r.Print(err)
	}
}

func getAnt(r radio.Radio, args []string) {
	r.Print("Antenna:", r.GetAntenna())
}

func setAnt(r radio.Radio, args []string) {

	if err := checkArgs(args, 1); err != nil {
		r.Print(err)
		return
	}

	ant, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		r.Print("ERROR: antenna value must be integer")
		return
	}

	if err := r.SetAntenna(int(ant)); err != nil {
		r.Print(err)
	}
}

func getPowerStat(r radio.Radio, args []string) {
	r.Print("Power On:", r.GetPowerstat())
}

func setPowerStat(r radio.Radio, args []string) {

	if err := checkArgs(args, 1); err != nil {
		r.Print(err)
		return
	}

	power, err := strconv.ParseBool(args[0])
	if err != nil {
		r.Print("ERROR: power value must be of type bool (1,t,true / 0,f,false)")
		return
	}

	if err := r.SetPowerstat(power); err != nil {
		r.Print(err)
	}
}

func getPtt(r radio.Radio, args []string) {
	r.Print("PTT On:", r.GetPtt())
}

func setPtt(r radio.Radio, args []string) {

	if err := checkArgs(args, 1); err != nil {
		r.Print(err)
		return
	}

	ptt, err := strconv.ParseBool(args[0])
	if err != nil {
		r.Print("ERROR: ptt value must be of type bool (1,t,true / 0,f,false)")
		return
	}

	if err := r.SetPtt(ptt); err != nil {
		r.Print(err)
	}
}

// func getLevelsPlain(r radio.Radio, args []string) {
// 	r.printLevelsPlain()
// }

func setLevel(r radio.Radio, args []string) {
	if err := checkArgs(args, 2); err != nil {
		r.Print(err)
		return
	}

	levelName := strings.ToUpper(args[0])

	caps := r.GetCaps()

	if !valueInValueList(levelName, caps.SetLevels) {
		r.Print("ERROR: unknown level")
	}

	levelValue, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		r.Print("ERROR: level value must be of type float")
		return
	}

	if err := r.SetLevel(levelName, float32(levelValue)); err != nil {
		r.Print(err)
	}
}

// func getFunctionsPlain(r radio.Radio, args []string) {
// 	r.printFunctionsPlain()
// }

func setFunction(r radio.Radio, args []string) {
	if err := checkArgs(args, 2); err != nil {
		r.Print(err)
		return
	}

	funcName := args[0]

	caps := r.GetCaps()

	if !utils.StringInSlice(funcName, caps.SetFunctions) {
		r.Print("unknown function")
	}

	value, err := strconv.ParseBool(args[1])
	if err != nil {
		r.Print("ERROR: function value must be of type bool (1,t,true / 0,f,false)")
		return
	}

	if err := r.SetFunction(funcName, value); err != nil {
		r.Print(err)
	}
}

func getSplit(r radio.Radio, args []string) {
	vfo, enabled := r.GetSplitVfo()
	r.Print("Split Enabled:", enabled)
	if enabled {
		freq, mode, pbWidth := r.GetSplitFrequencyMode()
		r.Print("Split Vfo:", vfo)
		r.Printf("Split Freq: %.3f kHz\n", freq)
		r.Print("Split Mode:", mode)
		r.Printf("Split PbWidth: %d Hz\n", pbWidth)
	}
}

func setSplitVfo(r radio.Radio, args []string) {

	if err := checkArgs(args, 2); err != nil {
		r.Print(err)
		return
	}

	splitEnabled, err := strconv.ParseBool(args[0])
	if err != nil {
		r.Print("ERROR: split enable/disable value must be of type bool (1,t,true / 0,f,false)")
		return
	}

	txVfo := args[1]

	caps := r.GetCaps()

	if !utils.StringInSlice(txVfo, caps.Vfos) {
		r.Print("ERROR: Vfo not supported by this radio")
	}

	if err := r.SetSplitVfo(txVfo, splitEnabled); err != nil {
		r.Print(err)
	}
}

func setSplitFreq(r radio.Radio, args []string) {

	if err := checkArgs(args, 1); err != nil {
		r.Print(err)
		return
	}

	freq, err := strconv.ParseFloat(args[0], 10)
	if err != nil {
		r.Print("ERROR: frequency [kHz] must be float")
		return
	}

	if err := r.SetSplitFrequency(freq); err != nil {
		r.Print(err)
	}
}

func setSplitMode(r radio.Radio, args []string) {
	if len(args) < 1 || len(args) > 2 {
		r.Print("ERROR: wrong number of arguments")
		return
	}

	caps := r.GetCaps()

	mode := args[0]

	if ok := utils.StringInSlice(mode, caps.Modes); !ok {
		r.Print("ERROR: unsupported mode")
		return
	}

	pbWidth := 0

	if len(args) == 2 {

		pbWidth, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			r.Print("ERROR: filter width [Hz] must be integer")
		}

		filters, ok := caps.Filters[args[0]]
		if !ok {
			r.Print("WARN: no filters found for this mode in rig caps")
		}
		if ok := utils.Int32InSlice(int32(pbWidth), filters.Value); !ok {
			r.Print("WARN: unspported filter width")
		}
	}

	if err := r.SetSplitMode(mode, pbWidth); err != nil {
		r.Print(err)
	}
}

func setSplitFreqMode(r radio.Radio, args []string) {
	if len(args) < 2 || len(args) > 3 {
		r.Print("ERROR: wrong number of arguments")
		return
	}

	freq, err := strconv.ParseFloat(args[0], 10)
	if err != nil {
		r.Print("ERROR: frequency [Hz] must be float")
		return
	}

	caps := r.GetCaps()
	mode := args[1]
	pbWidth := 0

	if ok := utils.StringInSlice(mode, caps.Modes); !ok {
		r.Print("ERROR: unsupported mode")
		return
	}

	if len(args) == 3 {

		pbWidth, err := strconv.ParseInt(args[2], 10, 32)
		if err != nil {
			r.Print("ERROR: filter width [Hz] must be integer")
		}

		filters, ok := caps.Filters[args[2]]
		if !ok {
			r.Print("WARN: no filters found for this mode in rig caps")
		}
		if ok := utils.Int32InSlice(int32(pbWidth), filters.Value); !ok {
			r.Print("WARN: unspported filter width")
		}
	}

	if err := r.SetSplitFrequencyMode(freq, mode, pbWidth); err != nil {
		r.Print(err)
	}
}

func execVfoOp(r radio.Radio, args []string) {

	caps := r.GetCaps()

	for _, vfoOp := range args {
		if !utils.StringInSlice(vfoOp, caps.VfoOps) {
			r.Print("ERROR: unknown vfo operation:", vfoOp)
			return
		}
	}

	if err := r.ExecVfoOps(args); err != nil {
		r.Print(err)
	}

}

func getTuningStep(r radio.Radio, args []string) {
	r.Printf("Tuning Step: %d Hz\n", r.GetTuningStep())
}

func setTuningStep(r radio.Radio, args []string) {

	if err := checkArgs(args, 1); err != nil {
		r.Print(err)
		return
	}

	ts, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		r.Print("ERROR: tuning step [Hz] must be integer")
		return
	}

	caps := r.GetCaps()

	// check if the given tuning step is supported by the rig
	mode, _ := r.GetMode()
	supportedTs, ok := caps.TuningSteps[mode]
	if !ok {
		r.Print("WARN: No tuning step values registered for this mode")
	}
	if ok := utils.Int32InSlice(int32(ts), supportedTs.Value); !ok {
		r.Print("WARN: tuning step not supported for this mode")
	}

	if err := r.SetTuningStep(int(ts)); err != nil {
		r.Print(err)
	}
}

// func getPollingInterval(r radio.Radio, args []string) {
// 	r.Printf("Rig polling interval: %dms\n", r.state.PollingInterval)
// }

// func setPollingInterval(r radio.Radio, args []string) {
// 	if !r.checkArgs(args, 1) {
// 		return
// 	}

// 	req := r.initSetState()

// 	ur, err := strconv.ParseInt(args[0], 10, 32)
// 	if err != nil {
// 		r.Print("ERROR: polling interval must be integer [ms]")
// 		return
// 	}

// 	req.PollingInterval = int32(ur)
// 	req.Md.HasPollingInterval = true

// 	if err := r.sendCatRequest(req); err != nil {
// 		r.Print("ERROR:", err)
// 	}
// }

// func getSyncInterval(r radio.Radio, args []string) {
// 	r.Printf("Rig sync interval: %ds\n", r.state.SyncInterval)
// }

// func setSyncInterval(r radio.Radio, args []string) {
// 	if !r.checkArgs(args, 1) {
// 		return
// 	}

// 	req := r.initSetState()

// 	ur, err := strconv.ParseInt(args[0], 10, 32)
// 	if err != nil {
// 		r.Print("ERROR: polling interval must be integer [s]")
// 		return
// 	}

// 	req.SyncInterval = int32(ur)
// 	req.Md.HasSyncInterval = true

// 	if err := r.sendCatRequest(req); err != nil {
// 		r.Print("ERROR:", err)
// 	}
// }

// func setPrintRigUpdates(r radio.Radio, args []string) {
// 	if !r.checkArgs(args, 1) {
// 		return
// 	}

// 	ru, err := strconv.ParseBool(args[0])
// 	if err != nil {
// 		r.Print("ERROR: value must be of type bool (1,t,true / 0,f,false)")
// 		return
// 	}

// 	r.printRigUpdates = ru
// }

func checkArgs(args []string, length int) error {
	if len(args) != length {
		return errors.New("ERROR: wrong number of arguments")
	}
	return nil
}

func valueInValueList(vName string, vList []*sbRadio.Value) bool {
	for _, value := range vList {
		if value.Name == vName {
			return true
		}
	}
	return false
}

var levelsTmpl = template.Must(template.New("").Parse(
	`
Levels: {{range $name, $val := .}}
    {{$name}}: {{$val}} {{end}}
`,
))

// func (r *radio.Radio) printLevelsPlain() {

// 	r.Print("Levels:")
// 	for levelName, levelValue := range r.state.Vfo.Levels {
// 		r.Printf(" %s: %.3f", levelName, levelValue)
// 	}
// }

// func (r radio.Radio) printFunctionsPlain() {

// 	r.Print("Functions:")
// 	for funcName, funcValue := range r.state.Vfo.Functions {
// 		r.Printf(" %s: %v", funcName, funcValue)
// 	}
// }
