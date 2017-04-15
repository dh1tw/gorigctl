package cligui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/cskr/pubsub"
	"github.com/dh1tw/gorigctl/events"
	sbRadio "github.com/dh1tw/gorigctl/sb_radio"
	"github.com/dh1tw/gorigctl/utils"
	ui "github.com/gizak/termui"
)

type radioGui struct {
	latencySpark         ui.Sparkline
	latency              *ui.Sparklines
	powerOn              *ui.Par
	ptt                  *ui.Par
	info                 *ui.List
	functions            *ui.List
	functionsData        []GuiFunction
	levels               *ui.List
	levelsData           []GuiLevel
	parameters           *ui.List
	parametersData       []GuiLevel
	frequency            *CharField
	frequencyInitialized bool
	sMeter               *ui.Gauge
	powerMeter           *ui.Gauge
	swrMeter             *ui.Gauge
	antenna              *ui.Par
	attenuator           *ui.Par
	preamp               *ui.Par
	vfo                  *ui.Par
	mode                 *ui.Par
	filter               *ui.Par
	tuningStep           *ui.Par
	rit                  *ui.Par
	xit                  *ui.Par
	split                *ui.Par
	txFrequency          *ui.Par
	txMode               *ui.Par
	txFilter             *ui.Par
	operations           *ui.List
	log                  *ui.List
	cli                  *Input
	state                sbRadio.State
	caps                 sbRadio.Capabilities
	internalFreq         float64
	lastFreqChange       time.Time
	radioOnline          bool
}

// initialize the gui components
func (rg *radioGui) init() {
	rg.state = sbRadio.State{}
	rg.state.Vfo = &sbRadio.Vfo{}
	rg.state.Channel = &sbRadio.Channel{}
	rg.state.Vfo.Split = &sbRadio.Split{}
	rg.internalFreq = 0.0
	rg.lastFreqChange = time.Now()
	rg.radioOnline = false
	rg.functionsData = make([]GuiFunction, 0, 32)

	rg.latencySpark = ui.NewSparkline()
	rg.latencySpark.Title = "Offline"
	rg.latencySpark.Data = []int{}
	rg.latencySpark.LineColor = ui.ColorYellow | ui.AttrBold

	rg.latency = ui.NewSparklines(rg.latencySpark)
	rg.latency.Height = 5
	rg.latency.BorderLabel = "Latency"

	rg.powerOn = ui.NewPar("")
	rg.powerOn.Height = 3
	rg.powerOn.BorderLabel = "Power On"

	rg.ptt = ui.NewPar("")
	rg.ptt.Height = 3
	rg.ptt.BorderLabel = "PTT"

	rg.info = ui.NewList()
	rg.info.Items = []string{"", ""}
	rg.info.BorderLabel = "Info"
	rg.info.Height = 4

	rg.functions = ui.NewList()
	rg.functions.BorderLabel = "Functions"
	rg.functions.Height = 10

	rg.levels = ui.NewList()
	rg.levels.BorderLabel = "Levels"
	rg.levels.Height = 10
	rg.levelsData = make([]GuiLevel, 0, 32)

	rg.parameters = ui.NewList()
	rg.parameters.BorderLabel = "Parameters"
	rg.parameters.Height = 10
	rg.parametersData = make([]GuiLevel, 0, 32)

	rg.frequency = NewCharField("")
	rg.frequency.BorderLabel = "Frequency"
	rg.frequency.Text = "Radio Offline"
	rg.frequency.Height = 9
	rg.frequency.Alignment = ui.AlignRight
	rg.frequency.PaddingTop = 1
	rg.frequency.PaddingLeft = 0

	rg.sMeter = ui.NewGauge()
	rg.sMeter.Percent = 40
	rg.sMeter.Height = 3
	rg.sMeter.BorderLabel = "S-Meter"
	rg.sMeter.BarColor = ui.ColorGreen
	rg.sMeter.Percent = 0
	rg.sMeter.Label = ""

	rg.swrMeter = ui.NewGauge()
	rg.swrMeter.Percent = 40
	rg.swrMeter.Height = 3
	rg.swrMeter.BorderLabel = "SWR"
	rg.swrMeter.BarColor = ui.ColorYellow
	rg.swrMeter.Percent = 0
	rg.swrMeter.Label = ""

	rg.powerMeter = ui.NewGauge()
	rg.powerMeter.Percent = 40
	rg.powerMeter.Height = 3
	rg.powerMeter.BorderLabel = "Power"
	rg.powerMeter.BarColor = ui.ColorRed
	rg.powerMeter.Percent = 0
	rg.powerMeter.Label = ""

	rg.mode = ui.NewPar("")
	rg.mode.Height = 3
	rg.mode.BorderLabel = "Mode"

	rg.vfo = ui.NewPar("")
	rg.vfo.Height = 3
	rg.vfo.BorderLabel = "VFO"

	rg.filter = ui.NewPar("")
	rg.filter.Height = 3
	rg.filter.BorderLabel = "Filter"

	rg.antenna = ui.NewPar("")
	rg.antenna.Height = 3
	rg.antenna.BorderLabel = "Antenna"

	rg.attenuator = ui.NewPar("")
	rg.attenuator.Height = 3
	rg.attenuator.BorderLabel = "Att"

	rg.preamp = ui.NewPar("")
	rg.preamp.Height = 3
	rg.preamp.BorderLabel = "Preamp"

	rg.tuningStep = ui.NewPar("")
	rg.tuningStep.Height = 3
	rg.tuningStep.BorderLabel = "Tuning Step"

	rg.rit = ui.NewPar("")
	rg.rit.Height = 3
	rg.rit.BorderLabel = "RIT"

	rg.xit = ui.NewPar("")
	rg.xit.Height = 3
	rg.xit.BorderLabel = "XIT"

	rg.split = ui.NewPar("")
	rg.split.Height = 3
	rg.split.BorderLabel = "Split"

	rg.txFrequency = ui.NewPar("")
	rg.txFrequency.Height = 3
	rg.txFrequency.BorderLabel = "TX Frequency"

	rg.txMode = ui.NewPar("")
	rg.txMode.Height = 3
	rg.txMode.BorderLabel = "TX Mode"

	rg.txFilter = ui.NewPar("")
	rg.txFilter.Height = 3
	rg.txFilter.BorderLabel = "TX Filter"

	rg.operations = ui.NewList()
	rg.operations.Items = []string{}
	rg.operations.BorderLabel = "Operations"
	rg.operations.Height = 10

	rg.log = ui.NewList()
	rg.log.Items = []string{}
	rg.log.BorderLabel = "Logging"
	rg.log.Height = rg.calcLogWindowHeight()

	if rg.cli != nil {
		if rg.cli.IsCapturing {
			rg.cli.StopCapture()
		}
	}
	rg.cli = NewInput("", false)
	rg.cli.Height = 3
	rg.cli.BorderLabel = "Rig command:"
	rg.cli.StartCapture()

	//clear the grid
	ui.Body.Rows = []*ui.Row{}
	ui.Clear()

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(2, 0, rg.info, rg.latency),
			ui.NewCol(8, 0, rg.frequency),
			ui.NewCol(2, 0, rg.powerMeter, rg.swrMeter, rg.sMeter)),
		ui.NewRow(
			ui.NewCol(2, 0, rg.powerOn),
			ui.NewCol(1, 0, rg.vfo),
			ui.NewCol(1, 0, rg.mode),
			ui.NewCol(2, 0, rg.filter),
			ui.NewCol(1, 0, rg.rit),
			ui.NewCol(1, 0, rg.xit),
			ui.NewCol(1, 0, rg.antenna),
			ui.NewCol(1, 0, rg.attenuator),
			ui.NewCol(1, 0, rg.preamp),
			ui.NewCol(1, 0, rg.tuningStep)),
		ui.NewRow(
			ui.NewCol(2, 0, rg.ptt),
			ui.NewCol(1, 0, rg.split),
			ui.NewCol(2, 0, rg.txFrequency),
			ui.NewCol(1, 0, rg.txMode),
			ui.NewCol(2, 0, rg.txFilter)),
		ui.NewRow(
			ui.NewCol(2, 0, rg.functions, rg.operations),
			ui.NewCol(8, 0, rg.log),
			ui.NewCol(2, 0, rg.levels, rg.parameters)),
		ui.NewRow(
			ui.NewCol(12, 0, rg.cli)),
	)

	// calculate layout
	ui.Body.Align()

	ui.Render(ui.Body)

}

func (rg *radioGui) calcLogWindowHeight() int {

	height := 0

	leftColumn := rg.functions.Height + rg.operations.Height
	rightColumn := rg.levels.Height + rg.parameters.Height
	if leftColumn > rightColumn {
		height = leftColumn
	} else {
		height = rightColumn
	}

	if height < 20 {
		height = 20
	}

	return height
}

func (rg *radioGui) updateCaps(ev ui.Event) {

	rg.caps = ev.Data.(sbRadio.Capabilities)

	rg.parametersData = make([]GuiLevel, 0, 32)
	rg.functionsData = make([]GuiFunction, 0, 32)
	rg.levelsData = make([]GuiLevel, 0, 32)

	// update widgets
	rg.info.Items[0] = rg.caps.MfgName + " " + rg.caps.ModelName
	rg.info.Items[1] = rg.caps.Version + " " + rg.caps.Status

	// update GUI Layout
	rg.operations.Height = 2 + len(rg.caps.VfoOps)
	rg.parameters.Height = 2 + len(rg.caps.GetParameters)

	for _, funcName := range rg.caps.GetFunctions {
		fData := GuiFunction{Label: funcName}
		rg.functionsData = append(rg.functionsData, fData)
	}

	// add the write-only (set)functions
	for _, funcName := range rg.caps.SetFunctions {
		found := false
		for _, f := range rg.functionsData {
			if f.Label == funcName {
				found = true
			}
		}
		if !found {
			fData := GuiFunction{Label: funcName, SetOnly: true}
			rg.functionsData = append(rg.functionsData, fData)
		}
	}

	rg.functions.Height = 2 + len(rg.functionsData)
	rg.functions.Items = SprintFunctions(rg.functionsData)

	for _, level := range rg.caps.GetLevels {
		lData := GuiLevel{Label: level.Name}
		rg.levelsData = append(rg.levelsData, lData)
	}

	// add the write-only (set)levels
	for _, level := range rg.caps.SetLevels {
		found := false
		for _, l := range rg.levelsData {
			if l.Label == level.Name {
				found = true
			}
		}
		if !found {
			lData := GuiLevel{Label: level.Name, SetOnly: true}
			rg.levelsData = append(rg.levelsData, lData)
		}
	}

	rg.levels.Height = 2 + len(rg.levelsData)
	rg.levels.Items = SprintLevels(rg.levelsData)

	rg.operations.Items = rg.caps.VfoOps

	rg.log.Height = rg.calcLogWindowHeight()

	if !rg.caps.HasPowerstat {
		rg.powerOn.Text = "n/a"
	}

	ui.Clear()
	ui.Body.Align()
	ui.Render(ui.Body)
	rg.updateGUI()
}

func (rg *radioGui) addLogEntry(ev ui.Event) {
	msg := ev.Data.(string)
	if len(rg.log.Items) >= rg.log.Height-2 {
		rg.log.Items = rg.log.Items[1:]
	}
	rg.log.Items = append(rg.log.Items, msg)
	ui.Render(rg.log)
}

func (rg *radioGui) updateState(ev ui.Event) {

	rg.state = ev.Data.(sbRadio.State)

	if !rg.radioOnline {
		return
	}

	rg.updateGUI()
}

func (rg *radioGui) updateGUI() {

	// assume radio is powered on
	powerOn := true

	// verify if radio is powered on (if possible)
	if rg.caps.HasPowerstat && !rg.state.RadioOn {
		rg.clear()
		rg.frequency.Text = "RADIO OFF"
		rg.powerOn.Bg = ui.ColorDefault
		ui.Render(rg.frequency, rg.powerOn)
		powerOn = false
	}

	if !rg.frequencyInitialized {
		rg.internalFreq = rg.state.Vfo.Frequency
		rg.frequencyInitialized = true
		rg.frequency.Text = utils.FormatFreq(rg.internalFreq)
	}

	rg.setVfo(powerOn, rg.state.CurrentVfo)
	rg.setMode(powerOn, rg.state.Vfo.Mode)
	rg.setFilter(powerOn, rg.state.Vfo.PbWidth)
	rg.setRit(powerOn, rg.state.Vfo.Rit)
	rg.setXit(powerOn, rg.state.Vfo.Xit)
	splitEnabled := rg.state.Vfo.Split.Enabled
	rg.setSplitVfo(splitEnabled, rg.state.Vfo.Split.Vfo)
	rg.setTxMode(splitEnabled, rg.state.Vfo.Split.Mode)
	rg.setTxFrequency(splitEnabled, rg.state.Vfo.Split.Frequency)
	rg.setTxFilter(splitEnabled, rg.state.Vfo.Split.PbWidth)
	rg.setPtt(powerOn, rg.state.Ptt)
	rg.setAntenna(powerOn, rg.state.Vfo.Ant)
	rg.setPowerOn(rg.caps.HasPowerstat, rg.state.RadioOn)
	rg.setTuningStep(powerOn, rg.state.Vfo.TuningStep)

	ptt := rg.state.Ptt

	if swrValue, ok := rg.state.Vfo.Levels["SWR"]; ok {
		rg.setSwrMeter(ptt, swrValue)
	}

	if powerValue, ok := rg.state.Vfo.Levels["METER"]; ok {
		rg.setPowerMeter(ptt, powerValue)
	}

	if sMeterValue, ok := rg.state.Vfo.Levels["STRENGTH"]; ok {
		rg.setSMeter(ptt, sMeterValue)
	}

	if attValue, ok := rg.state.Vfo.Levels["ATT"]; ok {
		rg.setAttenuator(powerOn, attValue)
	} else {
		rg.setAttenuator(false, 0)
	}

	if preampValue, ok := rg.state.Vfo.Levels["PREAMP"]; ok {
		rg.setPreamp(powerOn, preampValue)
	} else {
		rg.setPreamp(false, 0)
	}

	// fmt.Println(rg.state.Vfo.Functions)
	rg.setFunctions(rg.state.Vfo.Functions)
	rg.setLevels(rg.state.Vfo.Levels)
	rg.setParameters(rg.state.Vfo.Parameters)
}

func (rg *radioGui) clear() {

	rg.internalFreq = 0.0

	for i := range rg.levelsData {
		rg.levelsData[i].Value = 0.0
	}

	for i := range rg.functionsData {
		rg.functionsData[i].Set = false
	}

	for i := range rg.parametersData {
		rg.parametersData[i].Value = 0.0
	}
}

func (rg *radioGui) setVfo(powerOn bool, vfo string) {
	if powerOn {
		rg.vfo.Text = vfo
	} else {
		rg.vfo.Text = ""
	}
	ui.Render(rg.vfo)
}

func (rg *radioGui) setMode(powerOn bool, mode string) {
	if powerOn {
		rg.mode.Text = mode
	} else {
		rg.mode.Text = ""
	}
	ui.Render(rg.mode)
}

func (rg *radioGui) setFilter(powerOn bool, filter int32) {
	if powerOn {
		rg.filter.Text = fmt.Sprintf("%v Hz", filter)
	} else {
		rg.filter.Text = ""
	}
	ui.Render(rg.filter)
}

func (rg *radioGui) setRit(powerOn bool, rit int32) {
	if powerOn {
		rg.rit.Text = fmt.Sprintf("%v Hz", rit)
		if rg.state.Vfo.Rit != 0 {
			rg.rit.TextBgColor = ui.ColorGreen
			rg.rit.Bg = ui.ColorGreen
		} else {
			rg.rit.TextBgColor = ui.ColorDefault
			rg.rit.Bg = ui.ColorDefault
		}
	} else {
		rg.rit.Text = ""
		rg.rit.TextBgColor = ui.ColorDefault
		rg.rit.Bg = ui.ColorDefault
	}
	ui.Render(rg.rit)
}

func (rg *radioGui) setXit(powerOn bool, xit int32) {
	if powerOn {
		rg.xit.Text = fmt.Sprintf("%v Hz", xit)
		if rg.state.Vfo.Xit != 0 {
			rg.xit.TextBgColor = ui.ColorRed
			rg.xit.Bg = ui.ColorRed
		} else {
			rg.xit.TextBgColor = ui.ColorDefault
			rg.xit.Bg = ui.ColorDefault
		}
	} else {
		rg.xit.Text = ""
		rg.xit.TextBgColor = ui.ColorDefault
		rg.xit.Bg = ui.ColorDefault
	}
	ui.Render(rg.xit)
}

func (rg *radioGui) setAntenna(powerOn bool, ant int32) {
	if powerOn {
		rg.antenna.Text = fmt.Sprintf("%v", ant)
	} else {
		rg.antenna.Text = fmt.Sprintf("")
	}
	ui.Render(rg.antenna)
}

func (rg *radioGui) setTuningStep(powerOn bool, ts int32) {
	if powerOn {
		rg.tuningStep.Text = fmt.Sprintf("%d Hz", ts)
	} else {
		rg.tuningStep.Text = fmt.Sprintf("")
	}
	ui.Render(rg.tuningStep)
}

func (rg *radioGui) setPtt(powerOn bool, ptt bool) {
	if powerOn && ptt {
		rg.ptt.Bg = ui.ColorRed
	} else {
		rg.ptt.Bg = ui.ColorDefault
	}
	ui.Render(rg.ptt)
}

func (rg *radioGui) setPowerOn(hasPowerOn bool, powerOn bool) {
	rg.powerOn.Text = ""
	if powerOn {
		rg.powerOn.Bg = ui.ColorGreen
	} else {
		if hasPowerOn {
			rg.powerOn.Bg = ui.ColorDefault
		} else {
			rg.powerOn.Bg = ui.ColorDefault
			rg.powerOn.Text = "n/a"
		}
	}
	ui.Render(rg.powerOn)
}

func (rg *radioGui) setSplitVfo(enabled bool, txVfo string) {
	if enabled {
		if txVfo != "" {
			rg.split.Text = txVfo
			rg.split.TextBgColor = ui.ColorGreen
			rg.split.Bg = ui.ColorGreen
		} else {
			rg.split.Text = "n/a"
			rg.split.TextBgColor = ui.ColorDefault
			rg.split.Bg = ui.ColorDefault
		}
	} else {
		rg.split.Text = ""
		rg.split.TextBgColor = ui.ColorDefault
		rg.split.Bg = ui.ColorDefault
	}
	ui.Render(rg.split)
}

func (rg *radioGui) setTxFrequency(enabled bool, txFreq float64) {
	if txFreq > 0 {
		rg.txFrequency.Text = fmt.Sprintf("%.2f kHz", txFreq/1000)
	} else {
		if enabled {
			// rig does can not supply txFrequency
			rg.txFrequency.Text = "n/a"
		} else {
			rg.txFrequency.Text = ""
		}
	}
	ui.Render(rg.txFrequency)
}

func (rg *radioGui) setTxMode(enabled bool, mode string) {
	if mode != "" {
		rg.txMode.Text = mode
	} else {
		if enabled {
			// rig can not supply txMode
			rg.txMode.Text = "n/a"
		} else {
			rg.txMode.Text = ""
		}
	}
	ui.Render(rg.txMode)
}

func (rg *radioGui) setTxFilter(enabled bool, pbWidth int32) {
	if pbWidth > 0 {
		rg.txFilter.Text = fmt.Sprintf("%v Hz", pbWidth)
	} else {
		if enabled {
			// rig can not supply txFilter
			rg.txFilter.Text = "n/a"
		} else {
			rg.txFilter.Text = ""
		}
	}
	ui.Render(rg.txFilter)
}

func (rg *radioGui) setPowerMeter(ptt bool, value float32) {
	if ptt && value > 0 {
		rg.powerMeter.Percent = int(value * 100)
		rg.powerMeter.Label = fmt.Sprintf("%v Watt", value)
	} else {
		rg.powerMeter.Label = ""
		rg.powerMeter.Percent = 0
	}
	ui.Render(rg.powerMeter)
}

func (rg *radioGui) setSwrMeter(ptt bool, value float32) {
	if ptt && value > 0 {
		rg.swrMeter.Percent = int(value * 100)
		rg.swrMeter.Label = fmt.Sprintf("1:%.f2", value)
	} else {
		rg.swrMeter.Label = ""
		rg.swrMeter.Percent = 0
	}
	ui.Render(rg.swrMeter)
}

func (rg *radioGui) setSMeter(ptt bool, value float32) {
	if !ptt {
		if value < 0 {
			s := int((59 - value*-1) / 6)
			rg.sMeter.Label = fmt.Sprintf("S%v", s)
			rg.sMeter.Percent = int((59 - value*-1) * 100 / 114)
		} else {
			rg.sMeter.Label = fmt.Sprintf("S9+%vdB", int(value))
			rg.sMeter.Percent = int((value + 59) * 100 / 114)
		}
	} else {
		rg.sMeter.Label = ""
		rg.sMeter.Percent = 0
	}

	ui.Render(rg.sMeter)
}

func (rg *radioGui) setAttenuator(powerOn bool, value float32) {
	if powerOn {
		if value > 0 {
			rg.attenuator.Text = fmt.Sprintf("-%.0f dB", value)
			rg.attenuator.Bg = ui.ColorGreen
			rg.attenuator.TextBgColor = ui.ColorGreen
		} else {
			rg.attenuator.Text = fmt.Sprintf("%.0f dB", value)
			rg.attenuator.Bg = ui.ColorDefault
			rg.attenuator.TextBgColor = ui.ColorDefault
		}
	} else {
		rg.attenuator.Text = ""
		rg.attenuator.Bg = ui.ColorDefault
		rg.attenuator.TextBgColor = ui.ColorDefault
	}
	ui.Render(rg.attenuator)
}

func (rg *radioGui) setPreamp(powerOn bool, value float32) {
	if powerOn {
		if value > 0 {
			rg.preamp.Text = fmt.Sprintf("-%.0f dB", value)
			rg.preamp.Bg = ui.ColorGreen
			rg.preamp.TextBgColor = ui.ColorGreen
		} else {
			rg.preamp.Text = fmt.Sprintf("%.0f dB", value)
			rg.preamp.Bg = ui.ColorDefault
			rg.preamp.TextBgColor = ui.ColorDefault
		}
	} else {
		rg.preamp.Text = ""
		rg.preamp.Bg = ui.ColorDefault
		rg.preamp.TextBgColor = ui.ColorDefault
	}
	ui.Render(rg.preamp)
}

func (rg *radioGui) setFunctions(functions map[string]bool) {
	for i, f := range rg.functionsData {
		if _, ok := functions[f.Label]; ok {
			rg.functionsData[i].Set = functions[f.Label]
		}
	}
	rg.functions.Items = SprintFunctions(rg.functionsData)
	ui.Render(rg.functions)

}

func (rg *radioGui) setLevels(levels map[string]float32) {

	for i, el := range rg.levelsData {
		for name, value := range levels {
			if el.Label == name {
				rg.levelsData[i].Value = value
			}
		}
	}
	rg.levels.Items = SprintLevels(rg.levelsData)
	ui.Render(rg.levels)
}

func (rg *radioGui) setParameters(parameters map[string]float32) {

	for i, el := range rg.parametersData {
		for name, value := range parameters {
			if el.Label == name {
				rg.parametersData[i].Value = value
			}
		}
	}
	rg.parameters.Items = SprintLevels(rg.parametersData)
	ui.Render(rg.parameters)
}

// updateLatency updates the Latency chart (2 way ping)
func (rg *radioGui) updateLatency(ev ui.Event) {
	latency := ev.Data.(int64) / 1000000 // milli seconds
	if len(rg.latency.Lines[0].Data) > 20 {
		rg.latency.Lines[0].Data = rg.latency.Lines[0].Data[2:]
	}
	rg.latency.Lines[0].Data = append(rg.latency.Lines[0].Data, int(latency))
	rg.latency.Lines[0].Title = fmt.Sprintf("%dms", latency)
	ui.Render(rg.latency)
}

// updateRadioStatus handle the events in case the radio
// goes offline or becomes online
func (rg *radioGui) updateRadioStatus(ev ui.Event) {
	if ev.Data.(bool) {
		//we should update the entire GUI
		// rg.frequency.Text = utils.FormatFreq(rg.internalFreq)
		rg.radioOnline = true
		rg.updateGUI()
	} else {
		// this is a hack to remove artifacts from the
		// log widget when the canvas shrinks after
		// reinitalization
		rg.log.Height = 10
		ui.Render(rg.log)
		//reinit canvas
		rg.init()
	}
	ui.Render(ui.Body)
}

func (rg *radioGui) syncFrequency(ev ui.Event) {

	if !rg.radioOnline {
		return
	}

	if rg.caps.HasPowerstat && !rg.state.RadioOn {
		return
	}

	if time.Since(rg.lastFreqChange) > time.Millisecond*300 {
		rg.internalFreq = rg.state.Vfo.Frequency
		rg.frequency.Text = utils.FormatFreq(rg.internalFreq)
	}

	ui.Render(rg.frequency)
}

func guiLoop(caps sbRadio.Capabilities, evPS *pubsub.PubSub) {

	rg := &radioGui{}
	rg.init()

	ui.Handle("/radio/caps", rg.updateCaps)
	ui.Handle("/radio/state", rg.updateState)
	ui.Handle("/log/msg", rg.addLogEntry)
	ui.Handle("/network/latency", rg.updateLatency)
	ui.Handle("/radio/status", rg.updateRadioStatus)
	ui.Handle("/timer/1s", rg.syncFrequency)

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		if rg.radioOnline {
			if (rg.caps.HasPowerstat && rg.state.RadioOn) || !rg.caps.HasPowerstat {
				rg.internalFreq += float64(rg.state.Vfo.TuningStep)
				freq := rg.internalFreq / 1000
				cmd := []string{"set_freq", fmt.Sprintf("%.2f", freq)}
				evPS.Pub(cmd, events.CliInput)
				rg.frequency.Text = utils.FormatFreq(rg.internalFreq)
				ui.Render(rg.frequency)
				rg.lastFreqChange = time.Now()
			}
		}
	})

	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		if rg.radioOnline {
			if (rg.radioOnline && rg.state.RadioOn) || !rg.caps.HasPowerstat {
				rg.internalFreq -= float64(rg.state.Vfo.TuningStep)
				freq := rg.internalFreq / 1000
				cmd := []string{"set_freq", fmt.Sprintf("%.2f", freq)}
				evPS.Pub(cmd, events.CliInput)
				rg.frequency.Text = utils.FormatFreq(rg.internalFreq)
				ui.Render(rg.frequency)
				rg.lastFreqChange = time.Now()
			}
		}
	})

	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
		evPS.Pub(true, events.Shutdown)
	})

	ui.Handle("/input/kbd", func(ev ui.Event) {
		evData := ev.Data.(EvtInput)
		if evData.KeyStr == "<enter>" && len(rg.cli.Text()) > 0 {
			cmd := strings.Split(rg.cli.Text(), " ")
			evPS.Pub(cmd, events.CliInput)
			rg.cli.Clear()
			ui.Render(ui.Body)
		}
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Loop()
}

type GuiFunction struct {
	Label   string
	Set     bool
	SetOnly bool
}

func SprintFunctions(fs []GuiFunction) []string {
	s := make([]string, 0, len(fs))
	for _, el := range fs {
		item := el.Label
		for i := len(item); i < 8; i++ {
			item = item + " "
		}
		if el.Set {
			item = item + "[X]"
		} else if el.SetOnly {
			item = item + "[SetOnly]"
		} else {
			item = item + "[ ]"
		}
		s = append(s, item)
	}
	return s
}

func SprintLevels(levels []GuiLevel) []string {
	s := make([]string, 0, len(levels))

	for _, level := range levels {
		item := level.Label
		// add some spacing
		for i := len(item); i < 13; i++ {
			item = item + " "
		}
		if level.SetOnly {
			item = item + fmt.Sprint("SetOnly")
		} else {
			intr, frac := math.Modf(float64(level.Value))
			if frac > 0 {
				item = item + fmt.Sprintf("%.2f", level.Value)
			} else {
				item = item + fmt.Sprintf("%.0f", intr)
			}
			s = append(s, item)
		}
	}
	return s
}

type GuiLevel struct {
	Label   string
	Value   float32
	SetOnly bool
}
