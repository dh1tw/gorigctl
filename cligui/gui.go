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
	ifShift              *ui.Par
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
	rg.functionsData = make([]GuiFunction, 0, 32)

	rg.levels = ui.NewList()
	rg.levels.BorderLabel = "Levels"
	rg.levels.Height = 10
	rg.levelsData = make([]GuiLevel, 0, 32)

	rg.parameters = ui.NewList()
	rg.parameters.Items = []string{""}
	rg.parameters.BorderLabel = "Parameters"
	rg.parameters.Height = 10

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

	rg.ifShift = ui.NewPar("")
	rg.ifShift.Height = 3
	rg.ifShift.BorderLabel = "IF Shift"

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
			ui.NewCol(1, 0, rg.ifShift)),
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

	// this event could still thrown by a retained
	// caps message
	// if !rg.radioOnline {
	// 	return
	// }

	rg.caps = ev.Data.(sbRadio.Capabilities)

	// update GUI Layout
	rg.operations.Height = 2 + len(rg.caps.VfoOps)
	rg.levels.Height = 2 + len(rg.caps.GetLevels)
	rg.functions.Height = 2 + len(rg.caps.GetFunctions)
	rg.parameters.Height = 2 + len(rg.caps.GetParameters)
	rg.log.Height = rg.calcLogWindowHeight()

	// update widgets
	rg.info.Items[0] = rg.caps.MfgName + " " + rg.caps.ModelName
	rg.info.Items[1] = rg.caps.Version + " " + rg.caps.Status

	for _, funcName := range rg.caps.GetFunctions {
		fData := GuiFunction{Label: funcName}
		rg.functionsData = append(rg.functionsData, fData)
	}
	rg.functions.Items = SprintFunctions(rg.functionsData)

	for _, level := range rg.caps.GetLevels {
		lData := GuiLevel{Label: level.Name}
		rg.levelsData = append(rg.levelsData, lData)
	}
	rg.levels.Items = SprintLevels(rg.levelsData)

	rg.operations.Items = rg.caps.VfoOps

	if !rg.caps.HasPowerstat {
		rg.powerOn.Text = "n/a"
	}

	rg.drawState()
	ui.Clear()
	ui.Body.Align()
	ui.Render(ui.Body)
}

func (rg *radioGui) addLogEntry(ev ui.Event) {
	msg := ev.Data.(string)
	rg.log.Items = append(rg.log.Items, msg)
	ui.Render(rg.log)
}

func (rg *radioGui) updateState(ev ui.Event) {

	// this event could still thrown by a retained
	// state message
	// if !rg.radioOnline {
	// 	return
	// }

	rg.state = ev.Data.(sbRadio.State)
	rg.drawState()

}

func (rg *radioGui) drawState() {
	if !rg.frequencyInitialized {
		rg.internalFreq = rg.state.Vfo.Frequency
		rg.frequencyInitialized = true
		rg.frequency.Text = utils.FormatFreq(rg.internalFreq)
	}

	rg.mode.Text = rg.state.Vfo.Mode
	rg.filter.Text = fmt.Sprintf("%v Hz", rg.state.Vfo.PbWidth)
	rg.vfo.Text = rg.state.CurrentVfo
	rg.rit.Text = fmt.Sprintf("%v Hz", rg.state.Vfo.Rit)
	if rg.state.Vfo.Rit != 0 {
		rg.rit.TextBgColor = ui.ColorGreen
		rg.rit.Bg = ui.ColorGreen
	} else {
		rg.rit.TextBgColor = ui.ColorDefault
		rg.rit.Bg = ui.ColorDefault
	}

	rg.xit.Text = fmt.Sprintf("%v Hz", rg.state.Vfo.Xit)
	if rg.state.Vfo.Xit != 0 {
		rg.xit.TextBgColor = ui.ColorRed
		rg.xit.Bg = ui.ColorRed
	} else {
		rg.xit.TextBgColor = ui.ColorDefault
		rg.xit.Bg = ui.ColorDefault
	}

	rg.antenna.Text = fmt.Sprintf("%v", rg.state.Vfo.Ant)

	// rg.attenuator = fmt.Sprintf("%v dB", rg.state.Vfo.Attenuator)
	// rg.preamp = fmt.Sprintf("%v dB", rg.state.Vfo.Preamp)

	if rg.state.Ptt {
		rg.ptt.Bg = ui.ColorRed
	} else {
		rg.ptt.Bg = ui.ColorDefault
	}
	if rg.state.RadioOn {
		rg.powerOn.Bg = ui.ColorGreen
	} else {
		rg.powerOn.Bg = ui.ColorDefault
	}
	if rg.state.Vfo.Split.Enabled {
		if rg.state.Vfo.Split.Vfo != "" {
			rg.split.Text = rg.state.Vfo.Split.Vfo
		} else {
			rg.split.Text = "n/a"
		}
		rg.split.TextBgColor = ui.ColorGreen
		rg.split.Bg = ui.ColorGreen
		if rg.state.Vfo.Split.Frequency > 0 {
			rg.txFrequency.Text = fmt.Sprintf("%.2f kHz", rg.state.Vfo.Split.Frequency/1000)
		} else {
			rg.txFrequency.Text = "n/a"
		}
		if rg.state.Vfo.Split.Mode != "" {
			rg.txMode.Text = rg.state.Vfo.Split.Mode
		} else {
			rg.txMode.Text = "n/a"
		}
		if rg.state.Vfo.Split.PbWidth > 0 {
			rg.txFilter.Text = fmt.Sprintf("%v Hz", rg.state.Vfo.Split.PbWidth)
		} else {
			rg.txFilter.Text = "n/a"
		}
	} else {
		rg.split.Bg = ui.ColorDefault
		rg.split.TextBgColor = ui.ColorDefault
		rg.split.Text = ""
		rg.txFrequency.Text = ""
		rg.txMode.Text = ""
		rg.txFilter.Text = ""
	}
	if rg.state.Ptt {
		rg.sMeter.Percent = 0
		rg.sMeter.Label = ""
		if pValue, ok := rg.state.Vfo.Levels["METER"]; ok {
			rg.powerMeter.Label = fmt.Sprintf("%vW", pValue)
		}
		if swrValue, ok := rg.state.Vfo.Levels["SWR"]; ok {
			rg.swrMeter.Label = fmt.Sprintf("1:%.2f", swrValue)
		}
	} else {
		rg.powerMeter.Percent = 0
		rg.powerMeter.Label = ""
		rg.powerMeter.Percent = 0
		rg.swrMeter.Label = ""
		if sValue, ok := rg.state.Vfo.Levels["STRENGTH"]; ok {
			if sValue < 0 {
				s := int((59 - sValue*-1) / 6)
				rg.sMeter.Label = fmt.Sprintf("S%v", s)
				rg.sMeter.Percent = int((59 - sValue*-1) * 100 / 114)
			} else {
				rg.sMeter.Label = fmt.Sprintf("S9+%vdB", int(sValue))
				rg.sMeter.Percent = int((sValue + 59) * 100 / 114)
			}
		}
	}

	if attValue, ok := rg.state.Vfo.Levels["ATT"]; ok {
		if attValue > 0 {
			rg.attenuator.Text = fmt.Sprintf("-%.0f dB", attValue)
			rg.attenuator.Bg = ui.ColorGreen
			rg.attenuator.TextBgColor = ui.ColorGreen
		} else {
			rg.attenuator.Text = fmt.Sprintf("%.0f dB", attValue)
			rg.attenuator.Bg = ui.ColorDefault
			rg.attenuator.TextBgColor = ui.ColorDefault
		}
	} else {
		rg.attenuator.Text = "n/a"
		rg.attenuator.Bg = ui.ColorDefault
		rg.attenuator.TextBgColor = ui.ColorDefault
	}

	if preampValue, ok := rg.state.Vfo.Levels["PREAMP"]; ok {
		rg.preamp.Text = fmt.Sprintf("%.0f dB", preampValue)
		if preampValue > 0 {
			rg.preamp.Bg = ui.ColorGreen
			rg.preamp.TextBgColor = ui.ColorGreen
		} else {
			rg.preamp.Bg = ui.ColorDefault
			rg.preamp.TextBgColor = ui.ColorDefault
		}
	} else {
		rg.preamp.Text = "n/a"
		rg.preamp.Bg = ui.ColorDefault
		rg.preamp.TextBgColor = ui.ColorDefault
	}

	for i, el := range rg.levelsData {
		for name, value := range rg.state.Vfo.Levels {
			if el.Label == name {
				rg.levelsData[i].Value = value
			}
		}
	}
	rg.levels.Items = SprintLevels(rg.levelsData)

	for i, el := range rg.functionsData {
		found := false
		for _, funcName := range rg.state.Vfo.Functions {
			if el.Label == funcName {
				rg.functionsData[i].Set = true
				found = true
			}
			if !found {
				rg.functionsData[i].Set = false
			}
		}
	}
	rg.functions.Items = SprintFunctions(rg.functionsData)

	ui.Render(ui.Body)
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
		rg.frequency.Text = utils.FormatFreq(rg.internalFreq)
		rg.radioOnline = true
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
	if rg.radioOnline {
		// if rg.state.RadioOn {
		if time.Since(rg.lastFreqChange) > time.Millisecond*300 {
			rg.internalFreq = rg.state.Vfo.Frequency
			rg.frequency.Text = utils.FormatFreq(rg.internalFreq)
		}
		// } else {
		// 	rg.frequency.Text = "RADIO OFF"
		// }
	} else {
		rg.frequency.Text = "RADIO OFFLINE"
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

	// ui.Handle("/sys/kbd/<up>", func(ui.Event) {
	// 	if serverOnline && rg.state.RadioOn {
	// 		intFreq += float64(rg.state.Vfo.TuningStep)
	// 		freq := intFreq / 1000
	// 		cmd := []string{"set_freq", fmt.Sprintf("%.2f", freq)}
	// 		evPS.Pub(cmd, events.CliInput)
	// 		frequencyWidget.Text = utils.FormatFreq(intFreq)
	// 		ui.Render(frequencyWidget)
	// 		lastFreqChange = time.Now()
	// 	}
	// })

	// ui.Handle("/sys/kbd/<down>", func(ui.Event) {
	// 	if serverOnline && state.RadioOn {
	// 		intFreq -= float64(state.Vfo.TuningStep)
	// 		freq := intFreq / 1000
	// 		cmd := []string{"set_freq", fmt.Sprintf("%.2f", freq)}
	// 		evPS.Pub(cmd, events.CliInput)
	// 		frequencyWidget.Text = utils.FormatFreq(intFreq)
	// 		ui.Render(frequencyWidget)
	// 		lastFreqChange = time.Now()
	// 	}
	// })

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
	Label string
	Set   bool
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
		} else {
			item = item + "[ ]"
		}
		s = append(s, item)
	}
	return s
}

func SprintLevels(lv []GuiLevel) []string {
	s := make([]string, 0, len(lv))
	for _, el := range lv {
		item := el.Label
		for i := len(item); i < 13; i++ {
			item = item + " "
		}
		intr, frac := math.Modf(float64(el.Value))
		if frac > 0 {
			item = item + fmt.Sprintf("%.2f", el.Value)
		} else {
			item = item + fmt.Sprintf("%.0f", intr)
		}
		s = append(s, item)
	}
	return s
}

type GuiLevel struct {
	Label string
	Value float32
}
