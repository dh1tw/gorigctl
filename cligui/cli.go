package cligui

import (
	"math"
	"strconv"
	"strings"

	sbRadio "github.com/dh1tw/gorigctl/sb_radio"
	"github.com/dh1tw/gorigctl/utils"
)

func (r *remoteRadio) populateCliCmds() {

	cliSetFrequency := cliCmd{
		Cmd:         setFrequency,
		Name:        "set_freq",
		Shortcut:    "F",
		Parameters:  "Frequency [kHz]",
		Description: "Set Frequency for the current VFO",
		Example:     "F 14250.000",
	}

	r.cliCmds = append(r.cliCmds, cliSetFrequency)

	cliSetMode := cliCmd{
		Cmd:         setMode,
		Name:        "set_mode",
		Shortcut:    "M",
		Parameters:  "Mode and optionally Filter bandwidth [Hz]",
		Description: "Set Mode and optionally Filter Bandwidth for the current VFO)",
		Example:     "M USB 2400",
	}

	r.cliCmds = append(r.cliCmds, cliSetMode)

	cliSetVfo := cliCmd{
		Cmd:         setVfo,
		Name:        "set_vfo",
		Shortcut:    "V",
		Parameters:  "VFO Name",
		Description: "Change to another VFO",
		Example:     "V VFOB",
	}

	r.cliCmds = append(r.cliCmds, cliSetVfo)

	cliSetRit := cliCmd{
		Cmd:         setRit,
		Name:        "set_rit",
		Shortcut:    "J",
		Parameters:  "RX Offset [Hz]",
		Description: "Set RX Offset (0 = Off)",
		Example:     "J -500",
	}

	r.cliCmds = append(r.cliCmds, cliSetRit)

	cliSetXit := cliCmd{
		Cmd:         setXit,
		Name:        "set_xit",
		Shortcut:    "Z",
		Description: "Set TX Offset (0 = Off)",
		Parameters:  "TX Offset [Hz]",
		Example:     "Z -500",
	}

	r.cliCmds = append(r.cliCmds, cliSetXit)

	cliSetAnt := cliCmd{
		Cmd:         setAnt,
		Name:        "set_ant",
		Shortcut:    "y",
		Parameters:  "Antenna",
		Description: "Set Antenna",
		Example:     "Y 2",
	}

	r.cliCmds = append(r.cliCmds, cliSetAnt)

	cliSetPtt := cliCmd{
		Cmd:         setPtt,
		Name:        "set_ptt",
		Shortcut:    "t",
		Parameters:  "Ptt [true, t, 1, false, f, 0]",
		Description: "Set Transmit on/off",
		Example:     "t 1",
	}

	r.cliCmds = append(r.cliCmds, cliSetPtt)

	cliExecVfoOp := cliCmd{
		Cmd:         execVfoOp,
		Name:        "vfo_op",
		Shortcut:    "G",
		Parameters:  "VFO Operation",
		Description: "Execute a VFO Operation",
		Example:     "G XCHG",
	}

	r.cliCmds = append(r.cliCmds, cliExecVfoOp)

	cliSetFunction := cliCmd{
		Cmd:         setFunction,
		Name:        "set_func",
		Shortcut:    "U",
		Parameters:  "Function",
		Description: "Set a Rig function",
		Example:     "U NB true",
	}

	r.cliCmds = append(r.cliCmds, cliSetFunction)

	cliSetLevel := cliCmd{
		Cmd:         setLevel,
		Name:        "set_level",
		Shortcut:    "L",
		Parameters:  "Level & Value",
		Description: "Set a Level",
		Example:     "L CWPITCH 500",
	}

	r.cliCmds = append(r.cliCmds, cliSetLevel)

	cliSetTuningStep := cliCmd{
		Cmd:         setTuningStep,
		Name:        "set_ts",
		Shortcut:    "N",
		Parameters:  "Tuning Step [Hz]",
		Description: "Set the tuning step of the radio",
		Example:     "N 1000",
	}

	r.cliCmds = append(r.cliCmds, cliSetTuningStep)

	cliSetPowerStat := cliCmd{
		Cmd:         setPowerStat,
		Name:        "set_powerstat",
		Shortcut:    "",
		Parameters:  "Rig Power Status [true, t, 1, false, f, 0]",
		Description: "Turn the radio on/off",
		Example:     "set_powerstat 1",
	}

	r.cliCmds = append(r.cliCmds, cliSetPowerStat)

	cliSetSplit := cliCmd{
		Cmd:         setSplitVfo,
		Name:        "set_split_vfo",
		Shortcut:    "S",
		Parameters:  "Split [true, t, 1, false, f, 0], VFO",
		Description: "Turn Split On/Off",
		Example:     "S 1",
	}

	r.cliCmds = append(r.cliCmds, cliSetSplit)

	cliSetSplitFrequency := cliCmd{
		Cmd:         setSplitFreq,
		Name:        "set_split_freq",
		Shortcut:    "I",
		Parameters:  "TX Frequency [kHz]",
		Description: "Set the TX Split Frequency (the Split VFO will be determined automatically)",
		Example:     "I 14205000",
	}

	r.cliCmds = append(r.cliCmds, cliSetSplitFrequency)

	cliSetSplitMode := cliCmd{
		Cmd:         setSplitMode,
		Name:        "set_split_mode",
		Shortcut:    "X",
		Parameters:  "TX Mode and optionally Filter bandwidth [Hz]",
		Description: "Set the TX Split Mode (optionally with Bandwidth [Hz])",
		Example:     "X CW 200",
	}

	r.cliCmds = append(r.cliCmds, cliSetSplitMode)

	cliSetSplitFreqMode := cliCmd{
		Cmd:         setSplitFreqMode,
		Name:        "set_split_freq_mode",
		Shortcut:    "K",
		Parameters:  "TX Frequency [kHz], TX Mode and optionally Filter BW [Hz]",
		Description: "Set the Split Tx Frequency, Mode (optionally with Bandwidth [Hz])",
		Example:     "K 7170000 AM 6000",
	}

	r.cliCmds = append(r.cliCmds, cliSetSplitFreqMode)

	cliSetPollingInterval := cliCmd{
		Cmd:         setPollingInterval,
		Name:        "set_polling_interval",
		Shortcut:    "",
		Parameters:  "Polling rate [ms]",
		Description: "Set the radios polling Rate for updating the meters (SWR, ALC, Field Strength...)",
		Example:     "set_polling_rate 50",
	}

	r.cliCmds = append(r.cliCmds, cliSetPollingInterval)

	cliSetPrintUpdates := cliCmd{
		Cmd:         setPrintRigUpdates,
		Name:        "set_print_rig_updates",
		Parameters:  "[true, t, 1, false, f, 0]",
		Shortcut:    "",
		Description: "Print rig values which have changed",
	}

	r.cliCmds = append(r.cliCmds, cliSetPrintUpdates)

}

func (r *remoteRadio) parseCli(cliCmd []string) {

	for _, cmd := range r.cliCmds {
		if cmd.Name == cliCmd[0] || cmd.Shortcut == cliCmd[0] {
			cmd.Cmd(r, cliCmd[1:])
		}
	}
}

func getFrequency(r *remoteRadio, args []string) {
	r.logger.Printf("Frequency: %.3fkHz\n", r.state.Vfo.Frequency/1000)
}

func setFrequency(r *remoteRadio, args []string) {

	if ok := r.checkArgs(args, 1); !ok {
		return
	}

	freq, err := strconv.ParseFloat(args[0], 10)
	if err != nil {
		r.logger.Println("ERROR: frequency [kHz] must be float")
		return
	}

	// req := r.deepCopyState()
	req := r.initSetState()
	req.Vfo.Frequency = freq * 1000
	req.Md.HasFrequency = true
	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println("ERROR:", err)
	}
}

func getMode(r *remoteRadio, args []string) {
	r.logger.Println("Mode:", r.state.Vfo.Mode)
	r.logger.Printf("Filter: %dHz\n", r.state.Vfo.PbWidth)
}

func setMode(r *remoteRadio, args []string) {

	if len(args) < 1 || len(args) > 2 {
		r.logger.Println("ERROR: wrong number of arguments")
		return
	}

	mode := strings.ToUpper(args[0])

	if ok := utils.StringInSlice(mode, r.caps.Modes); !ok {
		r.logger.Println("ERROR: unsupported mode")
		return
	}

	req := r.initSetState()
	req.Vfo.Mode = mode
	req.Md.HasMode = true

	if len(args) == 2 {

		pbWidth, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			r.logger.Println("ERROR: filter width [Hz] must be integer")
			pbWidth = 0
		}

		filters, ok := r.caps.Filters[mode]
		if !ok {
			r.logger.Println("WARN: no filters found for this mode in rig caps")
		} else {
			if ok := utils.Int32InSlice(int32(pbWidth), filters.Value); !ok {
				r.logger.Println("WARN: unspported passband width")
			}
		}
		req.Vfo.PbWidth = int32(pbWidth)
		req.Md.HasPbWidth = true
	}

	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println("ERROR:", err)
	}
}

func getVfo(r *remoteRadio, args []string) {
	r.logger.Println("Current Vfo:", r.state.CurrentVfo)
}

func setVfo(r *remoteRadio, args []string) {
	if ok := r.checkArgs(args, 1); !ok {
		return
	}

	vfo := strings.ToUpper(args[0])

	if ok := utils.StringInSlice(vfo, r.caps.Vfos); !ok {
		r.logger.Println("ERROR: unsupported vfo")
		return
	}

	req := r.initSetState()
	req.CurrentVfo = vfo

	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println(err)
	}
}

func getRit(r *remoteRadio, args []string) {
	r.logger.Printf("Rit: %dHz\n", r.state.Vfo.Rit)
}

func setRit(r *remoteRadio, args []string) {

	if ok := r.checkArgs(args, 1); !ok {
		return
	}

	rit, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		r.logger.Println("ERROR: rit value [Hz] must be integer")
		return
	}

	if math.Abs(float64(rit)) > float64(r.caps.MaxRit) {
		r.logger.Println("WARN: rit value larger than supported by rig")
	}

	req := r.initSetState()
	req.Vfo.Rit = int32(rit)
	req.Md.HasRit = true

	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println("ERROR:", err)
	}
}

func getXit(r *remoteRadio, args []string) {
	r.logger.Printf("Xit: %dHz\n", r.state.Vfo.Xit)
}

func setXit(r *remoteRadio, args []string) {

	if !r.checkArgs(args, 1) {
		return
	}

	xit, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		r.logger.Println("ERROR: xit value [Hz] must be integer")
		return
	}

	if math.Abs(float64(xit)) > float64(r.caps.MaxXit) {
		r.logger.Println("WARN: xit value larger than supported by rig")
	}

	req := r.initSetState()

	req.Vfo.Xit = int32(xit)
	req.Md.HasXit = true

	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println("ERROR:", err)
	}
}

func getAnt(r *remoteRadio, args []string) {
	r.logger.Println("Antenna:", r.state.Vfo.Ant)
}

func setAnt(r *remoteRadio, args []string) {
	if !r.checkArgs(args, 1) {
		return
	}

	ant, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		r.logger.Println("ERROR: antenna value must be integer")
		return
	}

	// check Antenna in CAPS
	req := r.initSetState()
	req.Vfo.Ant = int32(ant)
	req.Md.HasAnt = true

	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println("ERROR:", err)
	}
}

func getPowerStat(r *remoteRadio, args []string) {
	r.logger.Println("Power On:", r.state.RadioOn)
}

func setPowerStat(r *remoteRadio, args []string) {
	if !r.checkArgs(args, 1) {
		return
	}

	power, err := strconv.ParseBool(args[0])
	if err != nil {
		r.logger.Println("ERROR: power value must be of type bool (1,t,true / 0,f,false")
		return
	}

	req := r.initSetState()
	req.RadioOn = power
	req.Md.HasRadioOn = true

	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println("ERROR:", err)
	}
}

func getPtt(r *remoteRadio, args []string) {
	r.logger.Println("PTT On:", r.state.Ptt)
}

func setPtt(r *remoteRadio, args []string) {
	if !r.checkArgs(args, 1) {
		return
	}

	ptt, err := strconv.ParseBool(args[0])
	if err != nil {
		r.logger.Println("ERROR: ptt value must be of type bool (1,t,true / 0,f,false")
		return
	}

	req := r.initSetState()
	req.Ptt = ptt
	req.Md.HasPtt = true

	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println("ERROR:", err)
	}
}

func setLevel(r *remoteRadio, args []string) {
	if !r.checkArgs(args, 2) {
		return
	}

	levelName := strings.ToUpper(args[0])

	if !valueInValueList(levelName, r.caps.SetLevels) {
		r.logger.Println("ERROR: unknown level")
	}

	levelValue, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		r.logger.Println("ERROR: level value must be of type float")
		return
	}

	levelMap := make(map[string]float32)

	levelMap[levelName] = float32(levelValue)

	req := r.initSetState()

	req.Vfo.Levels = levelMap
	req.Md.HasLevels = true
	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println("ERROR:", err)
	}
}

func getFunction(r *remoteRadio, args []string) {
	r.logger.Println("Functions:", r.state.Vfo.Functions)
}

func setFunction(r *remoteRadio, args []string) {
	if !r.checkArgs(args, 2) {
		return
	}

	funcName := args[0]
	if !utils.StringInSlice(funcName, r.caps.SetFunctions) {
		r.logger.Println("unknown function")
	}

	value, err := strconv.ParseBool(args[1])
	if err != nil {
		r.logger.Println("ERROR: function value must be of type bool (1,t,true / 0,f,false")
		return
	}

	req := r.initSetState()
	req.Md.HasFunctions = true
	req.Vfo.Functions = make(map[string]bool)
	req.Vfo.Functions[funcName] = value

	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println("ERROR:", err)
	}
}

func getSplit(r *remoteRadio, args []string) {
	r.logger.Println("Split Enabled:", r.state.Vfo.Split.Enabled)
	if r.state.Vfo.Split.Enabled {
		r.logger.Println("Split Vfo:", r.state.Vfo.Split.Vfo)
		r.logger.Printf("Split Freq: %.3fkHz\n", r.state.Vfo.Split.Frequency)
		r.logger.Println("Split Mode:", r.state.Vfo.Split.Mode)
		r.logger.Printf("Split PbWidth: %dHz\n", r.state.Vfo.Split.PbWidth)
	}
}

func setSplitVfo(r *remoteRadio, args []string) {
	if !r.checkArgs(args, 2) {
		return
	}

	splitEnabled, err := strconv.ParseBool(args[0])
	if err != nil {
		r.logger.Println("ERROR: split enable/disable value must be of type bool (1,t,true / 0,f,false")
		return
	}

	txVfo := args[1]
	if !utils.StringInSlice(txVfo, r.caps.Vfos) {
		r.logger.Println("ERROR: Vfo not supported by this radio")
	}

	req := r.initSetState()
	req.Md.HasSplit = true
	req.Vfo.Split.Enabled = splitEnabled
	req.Vfo.Split.Vfo = txVfo

	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println("ERROR:", err)
	}
}

func setSplitFreq(r *remoteRadio, args []string) {
	if !r.checkArgs(args, 1) {
		return
	}

	freq, err := strconv.ParseFloat(args[0], 10)
	if err != nil {
		r.logger.Println("ERROR: frequency [kHz] must be float")
		return
	}

	req := r.initSetState()
	req.Vfo.Split.Enabled = true
	req.Vfo.Split.Vfo = r.state.Vfo.Split.Vfo
	req.Vfo.Split.Frequency = freq * 1000
	req.Md.HasSplit = true

	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println(err)
	}
}

func setSplitMode(r *remoteRadio, args []string) {
	if len(args) < 1 || len(args) > 2 {
		r.logger.Println("ERROR: wrong number of arguments")
		return
	}

	if ok := utils.StringInSlice(args[0], r.caps.Modes); !ok {
		r.logger.Println("ERROR: unsupported mode")
		return
	}

	req := r.initSetState()
	req.Vfo.Split.Mode = args[0]
	req.Md.HasSplit = true

	if len(args) == 2 {

		pbWidth, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			r.logger.Println("ERROR: filter width [Hz] must be integer")
		}

		filters, ok := r.caps.Filters[args[0]]
		if !ok {
			r.logger.Println("WARN: no filters found for this mode in rig caps")
		}
		if ok := utils.Int32InSlice(int32(pbWidth), filters.Value); !ok {
			r.logger.Println("WARN: unspported filter width")
		}
		req.Vfo.Split.PbWidth = int32(pbWidth)
	}

	req.Vfo.Split.Enabled = true
	req.Vfo.Split.Vfo = r.state.Vfo.Split.Vfo

	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println("ERROR:", err)
	}
}

func setSplitFreqMode(r *remoteRadio, args []string) {
	if len(args) < 2 || len(args) > 3 {
		r.logger.Println("ERROR: wrong number of arguments")
		return
	}

	freq, err := strconv.ParseFloat(args[0], 10)
	if err != nil {
		r.logger.Println("ERROR: frequency [Hz] must be float")
		return
	}

	if ok := utils.StringInSlice(args[1], r.caps.Modes); !ok {
		r.logger.Println("ERROR: unsupported mode")
		return
	}

	req := r.initSetState()
	req.Vfo.Split.Enabled = true
	req.Vfo.Split.Frequency = freq * 1000
	req.Vfo.Split.Mode = args[1]
	req.Md.HasSplit = true

	if len(args) == 3 {

		pbWidth, err := strconv.ParseInt(args[2], 10, 32)
		if err != nil {
			r.logger.Println("ERROR: filter width [Hz] must be integer")
		}

		filters, ok := r.caps.Filters[args[2]]
		if !ok {
			r.logger.Println("WARN: no filters found for this mode in rig caps")
		}
		if ok := utils.Int32InSlice(int32(pbWidth), filters.Value); !ok {
			r.logger.Println("WARN: unspported filter width")
		}
		req.Vfo.Split.PbWidth = int32(pbWidth)
	}

	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println("ERROR:", err)
	}
}

func execVfoOp(r *remoteRadio, args []string) {

	for _, vfoOp := range args {
		if !utils.StringInSlice(vfoOp, r.caps.VfoOps) {
			r.logger.Println("ERROR: unknown vfo operation:", vfoOp)
			return
		}
	}

	req := r.initSetState()
	req.VfoOperations = args

	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println("ERROR:", err)
	}

}

func getTuningStep(r *remoteRadio, args []string) {
	r.logger.Printf("Tuning Step: %dHz\n", r.state.Vfo.TuningStep)
}

func setTuningStep(r *remoteRadio, args []string) {
	if !r.checkArgs(args, 1) {
		return
	}

	req := r.initSetState()

	ts, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		r.logger.Println("ERROR: tuning step [Hz] must be integer")
		return
	}

	// check if the given tuning step is supported by the rig
	supportedTs, ok := r.caps.TuningSteps[r.state.Vfo.Mode]
	if !ok {
		r.logger.Println("WARN: No tuning step values registered for this mode")
	}
	if ok := utils.Int32InSlice(int32(ts), supportedTs.Value); !ok {
		r.logger.Println("WARN: tuning step not supported for this mode")
	}
	req.Vfo.TuningStep = int32(ts)
	req.Md.HasTuningStep = true

	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println("ERROR:", err)
	}
}

func getPollingInterval(r *remoteRadio, args []string) {
	r.logger.Printf("Rig polling interval: %dms\n", r.state.PollingInterval)
}

func setPollingInterval(r *remoteRadio, args []string) {
	if !r.checkArgs(args, 1) {
		return
	}

	req := r.initSetState()

	ur, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		r.logger.Println("ERROR: polling interval must be integer [ms]")
		return
	}

	req.PollingInterval = int32(ur)
	req.Md.HasPollingInterval = true

	if err := r.sendCatRequest(req); err != nil {
		r.logger.Println("ERROR:", err)
	}
}

func setPrintRigUpdates(r *remoteRadio, args []string) {
	if !r.checkArgs(args, 1) {
		return
	}

	ru, err := strconv.ParseBool(args[0])
	if err != nil {
		r.logger.Println("ERROR: value must be of type bool (1,t,true / 0,f,false")
		return
	}

	r.printRigUpdates = ru
}

func (r *remoteRadio) checkArgs(args []string, length int) bool {
	if len(args) != length {
		r.logger.Println("ERROR: wrong number of arguments")
		return false
	}
	return true
}

func valueInValueList(vName string, vList []*sbRadio.Value) bool {
	for _, value := range vList {
		if value.Name == vName {
			return true
		}
	}
	return false
}
