package cli

import (
	"errors"
	"html/template"
	"log"
	"math"
	"strconv"
	"strings"

	"bytes"

	"github.com/dh1tw/gorigctl/radio"
	sbRadio "github.com/dh1tw/gorigctl/sb_radio"
	"github.com/dh1tw/gorigctl/utils"
	"github.com/olekukonko/tablewriter"
)

type CliCmd struct {
	Cmd         func(r radio.Radio, log *log.Logger, args []string)
	Name        string
	Shortcut    string
	Parameters  string
	Description string
	Example     string
}

type Help struct {
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
		Shortcut:    "Y",
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
		Shortcut:    "T",
		Parameters:  "Ptt [true, t, 1, false, f, 0]",
		Description: "Set Transmit on/off",
		Example:     "T 1",
	}

	cliCmds = append(cliCmds, cliSetPtt)

	cliGetPtt := CliCmd{
		Cmd:         getPtt,
		Name:        "get_ptt",
		Shortcut:    "t",
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

	cliGetFunction := CliCmd{
		Cmd:         getFunction,
		Name:        "get_func",
		Shortcut:    "u",
		Description: "Get the value for a particular level",
	}

	cliCmds = append(cliCmds, cliGetFunction)

	cliGetFunctions := CliCmd{
		Cmd:         getFunctions,
		Name:        "get_funcs",
		Description: "List the functions and their value",
	}

	cliCmds = append(cliCmds, cliGetFunctions)

	cliSetLevel := CliCmd{
		Cmd:         setLevel,
		Name:        "set_level",
		Shortcut:    "L",
		Parameters:  "Level & Value",
		Description: "Set a Level",
		Example:     "L CWPITCH 500",
	}

	cliCmds = append(cliCmds, cliSetLevel)

	cliGetLevel := CliCmd{
		Cmd:         getLevel,
		Name:        "get_level",
		Shortcut:    "l",
		Description: "Get the value for particular level",
	}

	cliCmds = append(cliCmds, cliGetLevel)

	cliGetLevels := CliCmd{
		Cmd:         getLevels,
		Name:        "get_levels",
		Description: "Lists all available levels",
	}

	cliCmds = append(cliCmds, cliGetLevels)

	cliSetParm := CliCmd{
		Cmd:         setParameter,
		Name:        "set_parm",
		Shortcut:    "P",
		Parameters:  "Parameter & Value",
		Description: "Set a Parameter",
		Example:     "P BACKLIGHT 1",
	}

	cliCmds = append(cliCmds, cliSetParm)

	cliGetParms := CliCmd{
		Cmd:         getParameters,
		Name:        "get_parms",
		Description: "Lists all available parameters",
	}

	cliCmds = append(cliCmds, cliGetParms)

	cliGetParm := CliCmd{
		Cmd:         getParameter,
		Name:        "get_parm",
		Shortcut:    "p",
		Description: "Get the value for particular parameter",
	}

	cliCmds = append(cliCmds, cliGetParm)

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
		Name:        "set_split_vfo",
		Shortcut:    "S",
		Parameters:  "Split VFO [true, t, 1, false, f, 0]",
		Description: "Turn Split On/Off for a VFO",
		Example:     "S 1 VFOB",
	}

	cliCmds = append(cliCmds, cliSetSplit)

	cliGetSplitVfo := CliCmd{
		Cmd:         getSplitVfo,
		Name:        "get_split_vfo",
		Shortcut:    "s",
		Description: "Get the split status (if enabled: VFO, Frequency, Mode, Filter)",
	}

	cliCmds = append(cliCmds, cliGetSplitVfo)

	cliSetSplitFrequency := CliCmd{
		Cmd:         setSplitFreq,
		Name:        "set_split_freq",
		Shortcut:    "I",
		Parameters:  "TX Frequency [kHz]",
		Description: "Set the TX Split Frequency (the Split VFO will be determined automatically)",
		Example:     "I 14205",
	}

	cliCmds = append(cliCmds, cliSetSplitFrequency)

	cliGetSplitFrequency := CliCmd{
		Cmd:         getSplitFreq,
		Name:        "get_split_freq",
		Shortcut:    "i",
		Description: "Get the TX Split Frequency (the Split VFO will be determined automatically)",
		Example:     "I 14205",
	}

	cliCmds = append(cliCmds, cliGetSplitFrequency)

	cliSetSplitMode := CliCmd{
		Cmd:         setSplitMode,
		Name:        "set_split_mode",
		Shortcut:    "X",
		Parameters:  "TX Mode and optionally Filter bandwidth [Hz]",
		Description: "Set the TX Split Mode (optionally with Bandwidth [Hz])",
		Example:     "X CW 200",
	}

	cliCmds = append(cliCmds, cliSetSplitMode)

	cliGetSplitMode := CliCmd{
		Cmd:         getSplitMode,
		Name:        "get_split_mode",
		Shortcut:    "x",
		Description: "Get the TX Split Mode and Filter",
		Example:     "X CW 200",
	}

	cliCmds = append(cliCmds, cliGetSplitMode)

	cliSetSplitFreqMode := CliCmd{
		Cmd:         setSplitFreqMode,
		Name:        "set_split_freq_mode",
		Shortcut:    "K",
		Parameters:  "TX Frequency [kHz], TX Mode and optionally Filter BW [Hz]",
		Description: "Set the Split Tx Frequency, Mode (optionally with Bandwidth [Hz])",
		Example:     "K 7170000 AM 6000",
	}
	cliCmds = append(cliCmds, cliSetSplitFreqMode)

	cliGetSplitFreqMode := CliCmd{
		Cmd:         getSplitFreqMode,
		Name:        "get_split_freq_mode",
		Shortcut:    "k",
		Description: "Get the TX Split Frequency, Mode and Filter Bandwidth",
		Example:     "X CW 200",
	}

	cliCmds = append(cliCmds, cliGetSplitFreqMode)

	cliDumpCaps := CliCmd{
		Cmd:         dumpCaps,
		Name:        "dump_caps",
		Shortcut:    "1",
		Description: "Print the capabilities of the radio",
	}

	cliCmds = append(cliCmds, cliDumpCaps)

	cliDumpState := CliCmd{
		Cmd:         dumpState,
		Name:        "dump_state",
		Shortcut:    "5",
		Description: "Print the complete state of the radio",
	}

	cliCmds = append(cliCmds, cliDumpState)

	return cliCmds
}

func getFrequency(r radio.Radio, log *log.Logger, args []string) {
	freq, err := r.GetFrequency()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Frequency: %.3f kHz\n", freq/1000)
}

func setFrequency(r radio.Radio, log *log.Logger, args []string) {

	if err := CheckArgs(args, 1); err != nil {
		log.Println(err)
		return
	}

	freq, err := strconv.ParseFloat(args[0], 10)
	if err != nil {
		log.Println("ERROR: frequency [kHz] must be float")
		return
	}

	// multiply frequency with 1000 since we enter kHz
	// but send over the wire Hz
	if err := r.SetFrequency(freq * 1000); err != nil {
		log.Println(err)
	}
}

func getMode(r radio.Radio, log *log.Logger, args []string) {
	mode, pbWidth, err := r.GetMode()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Mode:", mode)
	log.Printf("Filter: %d Hz", pbWidth)
}

func setMode(r radio.Radio, log *log.Logger, args []string) {

	if len(args) < 1 || len(args) > 2 {
		log.Println("ERROR: wrong number of arguments")
		return
	}

	mode := strings.ToUpper(args[0])

	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}
	pbWidth := 0

	if ok := utils.StringInSlice(mode, caps.Modes); !ok {
		log.Println("ERROR: unsupported mode")
		return
	}

	if len(args) == 2 {

		width, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			log.Println("ERROR: filter width [Hz] must be integer")
			return
		}

		filters, ok := caps.Filters[mode]
		if !ok {
			log.Println("WARN: no filters found for this mode in rig caps")
		} else {
			if ok := utils.Int32InSlice(int32(width), filters.Value); !ok {
				log.Println("WARN: unspported passband width")
			}
		}

		pbWidth = int(width)
	}

	if err := r.SetMode(mode, pbWidth); err != nil {
		log.Println(err)
	}
}

func getVfo(r radio.Radio, log *log.Logger, args []string) {
	vfo, err := r.GetVfo()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Current Vfo:", vfo)
}

func setVfo(r radio.Radio, log *log.Logger, args []string) {
	if err := CheckArgs(args, 1); err != nil {
		log.Println(err)
		return
	}

	vfo := strings.ToUpper(args[0])

	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}

	if ok := utils.StringInSlice(vfo, caps.Vfos); !ok {
		log.Println("ERROR: unsupported vfo")
		return
	}

	if err := r.SetVfo(vfo); err != nil {
		log.Println(err)
	}
}

func getRit(r radio.Radio, log *log.Logger, args []string) {
	rit, err := r.GetRit()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Rit: %d Hz\n", rit)
}

func setRit(r radio.Radio, log *log.Logger, args []string) {

	if err := CheckArgs(args, 1); err != nil {
		log.Println(err)
		return
	}

	rit, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		log.Println("ERROR: rit value [Hz] must be integer")
		return
	}

	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}

	if math.Abs(float64(rit)) > float64(caps.MaxRit) {
		log.Println("WARN: rit value larger than supported by rig")
	}

	if err := r.SetRit(int(rit)); err != nil {
		log.Println(err)
	}
}

func getXit(r radio.Radio, log *log.Logger, args []string) {
	xit, err := r.GetXit()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Xit: %d Hz\n", xit)
}

func setXit(r radio.Radio, log *log.Logger, args []string) {

	if err := CheckArgs(args, 1); err != nil {
		log.Println(err)
		return
	}

	xit, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		log.Println("ERROR: xit value [Hz] must be integer")
		return
	}

	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}

	if math.Abs(float64(xit)) > float64(caps.MaxXit) {
		log.Println("WARN: xit value larger than supported by rig")
	}

	if err := r.SetXit(int(xit)); err != nil {
		log.Println(err)
	}
}

func getAnt(r radio.Radio, log *log.Logger, args []string) {
	ant, err := r.GetAntenna()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Antenna:", ant)
}

func setAnt(r radio.Radio, log *log.Logger, args []string) {

	if err := CheckArgs(args, 1); err != nil {
		log.Println(err)
		return
	}

	ant, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		log.Println("ERROR: antenna value must be integer")
		return
	}

	if err := r.SetAntenna(int(ant)); err != nil {
		log.Println(err)
	}
}

func getPowerStat(r radio.Radio, log *log.Logger, args []string) {
	ps, err := r.GetPowerstat()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Power On:", ps)
}

func setPowerStat(r radio.Radio, log *log.Logger, args []string) {

	if err := CheckArgs(args, 1); err != nil {
		log.Println(err)
		return
	}

	power, err := strconv.ParseBool(args[0])
	if err != nil {
		log.Println("ERROR: power value must be of type bool (1,t,true / 0,f,false)")
		return
	}

	if err := r.SetPowerstat(power); err != nil {
		log.Println(err)
	}
}

func getPtt(r radio.Radio, log *log.Logger, args []string) {
	ptt, err := r.GetPtt()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("PTT On:", ptt)
}

func setPtt(r radio.Radio, log *log.Logger, args []string) {

	if err := CheckArgs(args, 1); err != nil {
		log.Println(err)
		return
	}

	ptt, err := strconv.ParseBool(args[0])
	if err != nil {
		log.Println("ERROR: ptt value must be of type bool (1,t,true / 0,f,false)")
		return
	}

	if err := r.SetPtt(ptt); err != nil {
		log.Println(err)
	}
}

func getLevel(r radio.Radio, log *log.Logger, args []string) {

	if err := CheckArgs(args, 1); err != nil {
		log.Println(err)
		return
	}

	level := args[0]

	value, err := r.GetLevel(level)
	if err != nil {
		log.Printf("ERROR: unsupported level")
		return
	}
	log.Printf("Level %s: %.2f", level, value)
}

func getLevels(r radio.Radio, log *log.Logger, args []string) {

	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Levels:")

	for _, level := range caps.GetLevels {
		value, err := r.GetLevel(level.Name)
		if err != nil {
			continue
		}
		log.Printf(" %s: %.2f", level.Name, value)
	}
}

func setLevel(r radio.Radio, log *log.Logger, args []string) {
	if err := CheckArgs(args, 2); err != nil {
		log.Println(err)
		return
	}

	levelName := strings.ToUpper(args[0])

	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}

	if !valueInValueList(levelName, caps.SetLevels) {
		log.Println("ERROR: unsupported set level")
	}

	levelValue, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		log.Println("ERROR: level value must be of type float")
		return
	}

	if err := r.SetLevel(levelName, float32(levelValue)); err != nil {
		log.Println(err)
	}
}

func getFunction(r radio.Radio, log *log.Logger, args []string) {

	if err := CheckArgs(args, 1); err != nil {
		log.Println(err)
		return
	}

	function := args[0]

	value, err := r.GetFunction(function)
	if err != nil {
		log.Printf("ERROR: unsupported level")
		return
	}
	log.Printf("Level %s: %v", function, value)
}

func getFunctions(r radio.Radio, log *log.Logger, args []string) {

	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Functions:")

	for _, function := range caps.GetFunctions {
		value, err := r.GetFunction(function)
		if err != nil {
			continue
		}
		log.Printf(" %s: %v", function, value)
	}
}

func setFunction(r radio.Radio, log *log.Logger, args []string) {
	if err := CheckArgs(args, 2); err != nil {
		log.Println(err)
		return
	}

	funcName := args[0]

	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}

	if !utils.StringInSlice(funcName, caps.SetFunctions) {
		log.Println("unsupported set function")
	}

	value, err := strconv.ParseBool(args[1])
	if err != nil {
		log.Println("ERROR: function value must be of type bool (1,t,true / 0,f,false)")
		return
	}

	if err := r.SetFunction(funcName, value); err != nil {
		log.Println(err)
	}
}

func getParameter(r radio.Radio, log *log.Logger, args []string) {

	if err := CheckArgs(args, 1); err != nil {
		log.Println(err)
		return
	}

	parm := args[0]

	value, err := r.GetParameter(parm)
	if err != nil {
		log.Printf("ERROR: unsupported level")
		return
	}
	log.Printf("Parameter %s: %.2f", parm, value)
}

func getParameters(r radio.Radio, log *log.Logger, args []string) {

	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Parameters:")

	for _, parm := range caps.GetParameters {
		value, err := r.GetParameter(parm.Name)
		if err != nil {
			continue
		}
		log.Printf(" %s: %.2f", parm.Name, value)
	}
}

func setParameter(r radio.Radio, log *log.Logger, args []string) {
	if err := CheckArgs(args, 2); err != nil {
		log.Println(err)
		return
	}

	parmName := strings.ToUpper(args[0])

	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}

	if !valueInValueList(parmName, caps.SetParameters) {
		log.Println("ERROR: unsupported set parameter")
	}

	parmValue, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		log.Println("ERROR: parameter value must be of type float")
		return
	}

	if err := r.SetParameter(parmName, float32(parmValue)); err != nil {
		log.Println(err)
	}
}

func getSplitVfo(r radio.Radio, log *log.Logger, args []string) {
	vfo, enabled, err := r.GetSplitVfo()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Split Enabled:", enabled)
	if enabled {
		log.Println("Split Vfo:", vfo)
	}
}

func getSplitMode(r radio.Radio, log *log.Logger, args []string) {
	_, enabled, err := r.GetSplitVfo()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Split Enabled:", enabled)
	if enabled {
		mode, pbWidth, err := r.GetSplitMode()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Split Mode:", mode)
		log.Printf("Split PbWidth: %d Hz\n", pbWidth)
	}
}

func getSplitFreqMode(r radio.Radio, log *log.Logger, args []string) {
	vfo, enabled, err := r.GetSplitVfo()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Split Enabled:", enabled)
	if enabled {
		freq, mode, pbWidth, err := r.GetSplitFrequencyMode()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Split Vfo:", vfo)
		log.Printf("Split Freq: %.3f kHz\n", freq)
		log.Println("Split Mode:", mode)
		log.Printf("Split PbWidth: %d Hz\n", pbWidth)
	}
}

func setSplitVfo(r radio.Radio, log *log.Logger, args []string) {

	if err := CheckArgs(args, 2); err != nil {
		log.Println(err)
		return
	}

	splitEnabled, err := strconv.ParseBool(args[0])
	if err != nil {
		log.Println("ERROR: split enable/disable value must be of type bool (1,t,true / 0,f,false)")
		return
	}

	txVfo := args[1]

	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}

	if !utils.StringInSlice(txVfo, caps.Vfos) {
		log.Println("ERROR: Vfo not supported by this radio")
	}

	if err := r.SetSplitVfo(txVfo, splitEnabled); err != nil {
		log.Println(err)
	}
}

func setSplitFreq(r radio.Radio, log *log.Logger, args []string) {

	if err := CheckArgs(args, 1); err != nil {
		log.Println(err)
		return
	}

	freq, err := strconv.ParseFloat(args[0], 10)
	if err != nil {
		log.Println("ERROR: frequency [kHz] must be float")
		return
	}

	if err := r.SetSplitFrequency(freq); err != nil {
		log.Println(err)
	}
}

func getSplitFreq(r radio.Radio, log *log.Logger, args []string) {
	_, enabled, err := r.GetSplitVfo()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Split Enabled:", enabled)
	if enabled {
		freq, err := r.GetSplitFrequency()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Split Freq: %.3f kHz\n", freq)
	}
}

func setSplitMode(r radio.Radio, log *log.Logger, args []string) {
	if len(args) < 1 || len(args) > 2 {
		log.Println("ERROR: wrong number of arguments")
		return
	}

	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}

	mode := args[0]

	if ok := utils.StringInSlice(mode, caps.Modes); !ok {
		log.Println("ERROR: unsupported mode")
		return
	}

	pbWidth := 0

	if len(args) == 2 {

		width, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			log.Println("ERROR: filter width [Hz] must be integer")
			return
		}

		filters, ok := caps.Filters[mode]
		if !ok {
			log.Println("WARN: no filters found for this mode in rig caps")
		} else {
			if ok := utils.Int32InSlice(int32(width), filters.Value); !ok {
				log.Println("WARN: unspported filter width")
			}
		}

		pbWidth = int(width)
	}

	if err := r.SetSplitMode(mode, pbWidth); err != nil {
		log.Println(err)
	}
}

func setSplitFreqMode(r radio.Radio, log *log.Logger, args []string) {
	if len(args) < 2 || len(args) > 3 {
		log.Println("ERROR: wrong number of arguments")
		return
	}

	freq, err := strconv.ParseFloat(args[0], 10)
	if err != nil {
		log.Println("ERROR: frequency [Hz] must be float")
		return
	}

	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}

	mode := args[1]
	pbWidth := 0

	if ok := utils.StringInSlice(mode, caps.Modes); !ok {
		log.Println("ERROR: unsupported mode")
		return
	}

	if len(args) == 3 {

		width, err := strconv.ParseInt(args[2], 10, 32)
		if err != nil {
			log.Println("ERROR: filter width [Hz] must be integer")
			return
		}

		filters, ok := caps.Filters[mode]
		if !ok {
			log.Println("WARN: no filters found for this mode in rig caps")
		} else {
			if ok := utils.Int32InSlice(int32(width), filters.Value); !ok {
				log.Println("WARN: unspported filter width")
			}
		}

		pbWidth = int(width)
	}

	if err := r.SetSplitFrequencyMode(freq, mode, pbWidth); err != nil {
		log.Println(err)
	}
}

func execVfoOp(r radio.Radio, log *log.Logger, args []string) {

	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}

	for _, vfoOp := range args {
		if !utils.StringInSlice(vfoOp, caps.VfoOps) {
			log.Println("ERROR: unknown vfo operation:", vfoOp)
			return
		}
	}

	if err := r.ExecVfoOps(args); err != nil {
		log.Println(err)
	}

}

func getTuningStep(r radio.Radio, log *log.Logger, args []string) {
	ts, err := r.GetTuningStep()
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Tuning Step: %d Hz\n", ts)
}

func setTuningStep(r radio.Radio, log *log.Logger, args []string) {

	if err := CheckArgs(args, 1); err != nil {
		log.Println(err)
		return
	}

	ts, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		log.Println("ERROR: tuning step [Hz] must be integer")
		return
	}

	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}

	// check if the given tuning step is supported by the rig
	mode, _, err := r.GetMode()
	if err != nil {
		log.Println(err)
		return
	}

	supportedTs, ok := caps.TuningSteps[mode]
	if !ok {
		log.Println("WARN: No tuning step values registered for this mode")
	}
	if ok := utils.Int32InSlice(int32(ts), supportedTs.Value); !ok {
		log.Println("WARN: tuning step not supported for this mode")
	}

	if err := r.SetTuningStep(int(ts)); err != nil {
		log.Println(err)
	}
}

var capsTmpl = template.Must(template.New("").Parse(
	`
Radio Capabilities:

Manufacturer: {{.MfgName}}
Model Name: {{.ModelName}}
Hamlib Rig Model ID: {{.RigModel}}
Hamlib Rig Version: {{.Version}}
Hamlib Rig Status: {{.Status}}
Supported VFOs:{{range $vfo := .Vfos}}{{$vfo}} {{end}}
Supported Modes: {{range $mode := .Modes}}{{$mode}} {{end}}
Supported VFO Operations: {{range $vfoOp := .VfoOps}}{{$vfoOp}} {{end}}
Supported Functions (Get):{{range $getF := .GetFunctions}}{{$getF}} {{end}}
Supported Functions (Set): {{range $setF := .SetFunctions}}{{$setF}} {{end}}
Supported Levels (Get): {{range $val := .GetLevels}}
  {{$val.Name}} ({{$val.Min}}..{{$val.Max}}/{{$val.Step}}){{end}}
Supported Levels (Set): {{range $val := .SetLevels}}
  {{$val.Name}} ({{$val.Min}}..{{$val.Max}}/{{$val.Step}}){{end}}
Supported Parameters (Get): {{range $val := .GetParameters}}
  {{$val.Name}} ({{$val.Min}}..{{$val.Max}}/{{$val.Step}}){{end}}
Supported Parameters (Set): {{range $val := .SetParameters}}
  {{$val.Name}} ({{$val.Min}}..{{$val.Max}}/{{$val.Step}}){{end}}
Max Rit: +-{{.MaxRit}}Hz
Max Xit: +-{{.MaxXit}}Hz
Max IF Shift: +-{{.MaxIfShift}}Hz
Filters [Hz]: {{range $mode, $pbList := .Filters}}
  {{$mode}}:		{{range $pb := $pbList.Value}}{{$pb}} {{end}} {{end}}
Tuning Steps [Hz]: {{range $mode, $tsList := .TuningSteps}}
  {{$mode}}:		{{range $ts := $tsList.Value}}{{$ts}} {{end}} {{end}}
Preamps: {{range $preamp := .Preamps}}{{$preamp}}dB {{end}}
Attenuators: {{range $att := .Attenuators}}{{$att}}dB {{end}} 
`,
))

func dumpCaps(r radio.Radio, log *log.Logger, args []string) {
	caps, err := r.GetCaps()
	if err != nil {
		log.Println(err)
		return
	}

	buf := bytes.Buffer{}

	err = capsTmpl.Execute(&buf, caps)
	if err != nil {
		log.Println(err)
	}

	lines := strings.Split(buf.String(), "\n")

	for _, line := range lines {
		log.Println(line)
	}
}

var stateTmpl = template.Must(template.New("").Parse(
	`
Current Vfo: {{.CurrentVfo}}
  Frequency: {{.Vfo.Frequency}}Hz
  Mode: {{.Vfo.Mode}}
  PbWidth: {{.Vfo.PbWidth}}
  Antenna: {{.Vfo.Ant}}
  Rit: {{.Vfo.Rit}}
  Xit: {{.Vfo.Xit}}
  Split: 
    Enabled: {{.Vfo.Split.Enabled}}
    Vfo: {{.Vfo.Split.Vfo}}
    Frequency: {{.Vfo.Split.Frequency}}
    Mode: {{.Vfo.Split.Mode}}
    PbWidth: {{.Vfo.Split.PbWidth}}
  Tuning Step: {{.Vfo.TuningStep}}
  Functions: {{range $name, $val := .Vfo.Functions}}
    {{$name}}: {{$val}} {{end}}
  Levels: {{range $name, $val := .Vfo.Levels}}
    {{$name}}: {{$val}} {{end}}
  Parameters: {{range $name, $val := .Vfo.Parameters}}
    {{$name}}: {{$val}} {{end}}
Radio On: {{.RadioOn}}
Ptt: {{.Ptt}}
Update Rate: {{.PollingInterval}}
`,
))

func dumpState(r radio.Radio, log *log.Logger, args []string) {
	state, err := r.GetState()
	if err != nil {
		log.Println(err)
		return
	}

	buf := bytes.Buffer{}

	err = stateTmpl.Execute(&buf, state)
	if err != nil {
		log.Println(err)
	}

	lines := strings.Split(buf.String(), "\n")

	for _, line := range lines {
		log.Println(line)
	}
}

var helpTmpl = template.Must(template.New("").Parse(
	`
Available commands (some may not be available for this radio):

{{range .}}{{.Name}}:
  Shortcut: {{if .Shortcut}}{{.Shortcut}}{{else}}n/a{{end}}
  Description: {{if .Description}}{{.Description}}{{else}}n/a{{end}}
  Example: {{if .Example}}{{.Example}}{{else}}n/a{{end}}

{{end}}

`,
))

// func PrintHelp(c Cmd, log *log.Logger) {

// 	buf := byte.Buffer{}

// 	list := make([]help, 0, len)

// 	err := helpTmpl.Execute(&buf, r.cliCmds)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

func PrintHelp(cmds []CliCmd, log *log.Logger) {

	buf := bytes.Buffer{}

	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"Command", "Shortcut", "Parameter"})
	table.SetCenterSeparator("|")
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(50)

	for _, el := range cmds {
		table.Append([]string{el.Name, el.Shortcut, el.Parameters})
	}

	table.Render()

	lines := strings.Split(buf.String(), "\n")

	for _, line := range lines {
		log.Println(line)
	}
}

func CheckArgs(args []string, length int) error {
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
