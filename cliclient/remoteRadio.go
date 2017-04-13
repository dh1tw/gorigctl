package cliClient

import (
	"fmt"
	"html/template"
	"log"
	"reflect"
	"sync"

	"os"

	"github.com/cskr/pubsub"
	"github.com/dh1tw/gorigctl/comms"
	"github.com/dh1tw/gorigctl/events"
	sbRadio "github.com/dh1tw/gorigctl/sb_radio"
	sbStatus "github.com/dh1tw/gorigctl/sb_status"
	"github.com/dh1tw/gorigctl/utils"
	"github.com/spf13/viper"
)

type RemoteRadioSettings struct {
	CatResponseCh   chan []byte
	RadioStatusCh   chan []byte
	CatRequestTopic string
	ToWireCh        chan comms.IOMsg
	CapabilitiesCh  chan []byte
	WaitGroup       *sync.WaitGroup
	Events          *pubsub.PubSub
}

type remoteRadio struct {
	state           sbRadio.State
	newState        sbRadio.SetState
	caps            sbRadio.Capabilities
	settings        RemoteRadioSettings
	cliCmds         []cliCmd
	printRigUpdates bool
	userID          string
	radioOnline     bool
}

type cliCmd struct {
	Cmd         func(r *remoteRadio, args []string)
	Name        string
	Shortcut    string
	Parameters  string
	Description string
	Example     string
}

func HandleRemoteRadio(rs RemoteRadioSettings) {
	defer rs.WaitGroup.Done()

	shutdownCh := rs.Events.Sub(events.Shutdown)
	cliInputCh := rs.Events.Sub(events.CliInput)

	r := remoteRadio{}
	r.state.Vfo = &sbRadio.Vfo{}
	r.state.Vfo.Functions = make(map[string]bool)
	r.state.Vfo.Levels = make(map[string]float32)
	r.state.Vfo.Parameters = make(map[string]float32)
	r.state.Vfo.Split = &sbRadio.Split{}

	r.settings = rs

	r.cliCmds = make([]cliCmd, 0, 30)
	r.populateCliCmds()

	if viper.IsSet("general.user_id") {
		r.userID = viper.GetString("general.user_id")
	} else {
		r.userID = "unknown_" + utils.RandStringRunes(5)
	}

	// rs.Events.Pub(true, events.ForwardCat)

	fmt.Println("Rig command: ")

	for {
		select {
		case msg := <-rs.CapabilitiesCh:
			r.deserializeCaps(msg)
			// r.PrintCapabilities()
		case msg := <-rs.CatResponseCh:
			r.deserializeCatResponse(msg)
			// r.PrintState()
		case msg := <-rs.RadioStatusCh:
			r.deserializeRadioStatus(msg)
		case msg := <-cliInputCh:
			r.parseCli(msg.([]string))
		case <-shutdownCh:
			log.Println("Disconnecting from Radio")
			return
		}
	}
}

func (r *remoteRadio) deserializeRadioStatus(data []byte) error {

	rStatus := sbStatus.Status{}
	if err := rStatus.Unmarshal(data); err != nil {
		return err
	}

	if r.radioOnline != rStatus.GetOnline() {
		r.radioOnline = rStatus.GetOnline()
		fmt.Println("Update Radio Online:", r.radioOnline)
	}

	return nil
}

func (r *remoteRadio) sendCatRequest(req sbRadio.SetState) error {
	data, err := req.Marshal()
	if err != nil {
		return err
	}

	msg := comms.IOMsg{}
	msg.Data = data
	msg.Topic = r.settings.CatRequestTopic
	msg.Retain = false
	msg.Qos = 0

	r.settings.ToWireCh <- msg

	return nil
}

var stateTmpl = template.Must(template.New("").Parse(
	`
Current Vfo: {{.CurrentVfo}}
  Frequency: {{.Vfo.Frequency}}Hz
  Mode: {{.Vfo.Mode}}
  PBWidth: {{.Vfo.PbWidth}}
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
  Functions: {{range $f := .Vfo.Functions}}{{$f}} {{end}}
  Levels: {{range $name, $val := .Vfo.Levels}}
    {{$name}}: {{$val}} {{end}}
  Parameters: {{range $name, $val := .Vfo.Parameters}}
    {{$name}}: {{$val}} {{end}}
Radio On: {{.RadioOn}}
Ptt: {{.Ptt}}
Update Rate: {{.PollingInterval}}

`,
))

var levelsTmpl = template.Must(template.New("").Parse(
	`
Levels: {{range $name, $val := .}}
    {{$name}}: {{$val}} {{end}}
`,
))

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

func (r *remoteRadio) PrintCapabilities() {
	err := capsTmpl.Execute(os.Stdout, r.caps)
	if err != nil {
		fmt.Println(err)
	}
}

func (r *remoteRadio) PrintLevels() {
	err := levelsTmpl.Execute(os.Stdout, r.state.Vfo.Levels)
	if err != nil {
		fmt.Println(err)
	}
}

func (r *remoteRadio) PrintState() {
	err := stateTmpl.Execute(os.Stdout, r.state)
	if err != nil {
		fmt.Println(err)
	}
}

func (r *remoteRadio) deserializeCaps(msg []byte) error {

	caps := sbRadio.Capabilities{}
	err := caps.Unmarshal(msg)
	if err != nil {
		return err
	}

	r.caps = caps

	return nil
}

func (r *remoteRadio) deserializeCatResponse(msg []byte) error {

	ns := sbRadio.State{}
	err := ns.Unmarshal(msg)
	if err != nil {
		return err
	}

	if ns.CurrentVfo != r.state.CurrentVfo {
		r.state.CurrentVfo = ns.CurrentVfo
		if r.printRigUpdates {
			fmt.Println("Updated Current Vfo:", r.state.CurrentVfo)
		}
	}

	if ns.Vfo != nil {

		if ns.Vfo.GetFrequency() != r.state.Vfo.Frequency {
			r.state.Vfo.Frequency = ns.Vfo.GetFrequency()
			if r.printRigUpdates {
				fmt.Printf("Updated Frequency: %.3fkHz\n", r.state.Vfo.Frequency/1000)
			}
		}

		if ns.Vfo.GetMode() != r.state.Vfo.Mode {
			r.state.Vfo.Mode = ns.Vfo.GetMode()
			if r.printRigUpdates {
				fmt.Println("Updated Mode:", r.state.Vfo.Mode)
			}
		}

		if ns.Vfo.GetPbWidth() != r.state.Vfo.PbWidth {
			r.state.Vfo.PbWidth = ns.Vfo.GetPbWidth()
			if r.printRigUpdates {
				fmt.Printf("Updated Filter: %dHz\n", r.state.Vfo.PbWidth)
			}
		}

		if ns.Vfo.GetAnt() != r.state.Vfo.Ant {
			r.state.Vfo.Ant = ns.Vfo.GetAnt()
			if r.printRigUpdates {
				fmt.Println("Updated Antenna:", r.state.Vfo.Ant)
			}
		}

		if ns.Vfo.GetRit() != r.state.Vfo.Rit {
			r.state.Vfo.Rit = ns.Vfo.GetRit()
			if r.printRigUpdates {
				fmt.Printf("Updated Rit: %dHz\n", r.state.Vfo.Rit)
			}
		}

		if ns.Vfo.GetXit() != r.state.Vfo.Xit {
			r.state.Vfo.Xit = ns.Vfo.GetXit()
			if r.printRigUpdates {
				fmt.Printf("Updated Xit: %dHz\n", r.state.Vfo.Xit)
			}
		}

		if ns.Vfo.GetSplit() != nil {
			if !reflect.DeepEqual(ns.Vfo.GetSplit(), r.state.Vfo.Split) {
				if err := r.updateSplit(ns.Vfo.Split); err != nil {
					log.Println(err)
				}
			}
		}

		if ns.Vfo.GetTuningStep() != r.state.Vfo.TuningStep {
			r.state.Vfo.TuningStep = ns.Vfo.GetTuningStep()
			if r.printRigUpdates {
				fmt.Printf("Updated Tuning Step: %dHz\n", r.state.Vfo.TuningStep)
			}
		}

		if ns.Vfo.Functions != nil {
			if !reflect.DeepEqual(ns.Vfo.Functions, r.state.Vfo.Functions) {
				if err := r.updateFunctions(ns.Vfo.GetFunctions()); err != nil {
					log.Println(err)
				}
			}
		}

		if ns.Vfo.Levels != nil {
			if !reflect.DeepEqual(ns.Vfo.Levels, r.state.Vfo.Levels) {
				if err := r.updateLevels(ns.Vfo.GetLevels()); err != nil {
					log.Println(err)
				}
			}
		}

		if ns.Vfo.Parameters != nil {
			if !reflect.DeepEqual(ns.Vfo.Parameters, r.state.Vfo.Parameters) {
				if err := r.updateParams(ns.Vfo.GetParameters()); err != nil {
					log.Println(err)
				}
			}
		}

	}

	if ns.GetRadioOn() != r.state.RadioOn {
		r.state.RadioOn = ns.GetRadioOn()
		if r.printRigUpdates {
			fmt.Println("Updated Radio Power On:", r.state.RadioOn)
		}
	}

	if ns.GetPtt() != r.state.Ptt {
		r.state.Ptt = ns.GetPtt()
		if r.printRigUpdates {
			fmt.Println("Updated PTT On:", r.state.Ptt)
		}
	}

	if ns.GetPollingInterval() != r.state.PollingInterval {
		r.state.PollingInterval = ns.GetPollingInterval()
		if r.printRigUpdates {
			fmt.Printf("Updated rig polling interval: %dms\n", r.state.PollingInterval)
		}
	}

	return nil
}

func (r *remoteRadio) updateSplit(newSplit *sbRadio.Split) error {

	if newSplit.GetEnabled() != r.state.Vfo.Split.Enabled {
		r.state.Vfo.Split.Enabled = newSplit.GetEnabled()
		if r.printRigUpdates {
			fmt.Println("Updated Split Enabled:", r.state.Vfo.Split.Enabled)
		}
	}

	if newSplit.GetFrequency() != r.state.Vfo.Split.Frequency {
		r.state.Vfo.Split.Frequency = newSplit.GetFrequency()
		if r.printRigUpdates {
			fmt.Printf("Updated TX (Split) Frequency: %.3fkHz\n", r.state.Vfo.Split.Frequency/1000)
		}
	}

	if newSplit.GetVfo() != r.state.Vfo.Split.Vfo {
		r.state.Vfo.Split.Vfo = newSplit.GetVfo()
		if r.printRigUpdates {
			fmt.Println("Updated TX (Split) Vfo:", r.state.Vfo.Split.Vfo)
		}
	}

	if newSplit.GetMode() != r.state.Vfo.Split.Mode {
		r.state.Vfo.Split.Mode = newSplit.GetMode()
		if r.printRigUpdates {
			fmt.Println("Updated TX (Split) Mode:", r.state.Vfo.Split.Mode)
		}
	}

	if newSplit.GetPbWidth() != r.state.Vfo.Split.PbWidth {

		r.state.Vfo.Split.PbWidth = newSplit.GetPbWidth()
		if r.printRigUpdates {
			fmt.Printf("Split PbWidth: %dHz\n", r.state.Vfo.Split.PbWidth)
		}
	}

	return nil
}

func (r *remoteRadio) updateFunctions(newFuncs map[string]bool) error {

	r.state.Vfo.Functions = newFuncs

	// vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	// functions to be enabled
	// diff := utils.SliceDiff(newFuncs, r.state.Vfo.Functions)
	// for _, f := range diff {
	// 	funcValue, ok := hl.FuncValue[f]
	// 	if !ok {
	// 		return errors.New("unknown function")
	// 	}
	// 	// err := r.rig.SetFunc(vfo, funcValue, true)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// }

	// // functions to be disabled
	// diff = utils.SliceDiff(r.state.Vfo.Functions, newFuncs)
	// for _, f := range diff {
	// 	funcValue, ok := hl.FuncValue[f]
	// 	if !ok {
	// 		return errors.New("unknown function")
	// 	}
	// 	// err := r.rig.SetFunc(vfo, funcValue, false)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// }

	return nil
}

func (r *remoteRadio) updateLevels(newLevels map[string]float32) error {

	r.state.Vfo.Levels = newLevels

	// vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	// for k, v := range newLevels {
	// 	levelValue, ok := hl.LevelValue[k]
	// 	if !ok {
	// 		return errors.New("unknown Level")
	// 	}
	// 	if _, ok := r.state.Vfo.Levels[k]; !ok {
	// 		return errors.New("unsupported Level for this rig")
	// 	}

	// if r.state.Vfo.Levels[k] != v {
	// 	err := r.rig.SetLevel(vfo, levelValue, v)
	// 	if err != nil {
	// 		return nil
	// 	}
	// }
	// }

	return nil
}

func (r *remoteRadio) updateParams(newParams map[string]float32) error {

	r.state.Vfo.Parameters = newParams

	// vfo, _ := hl.VfoValue[r.state.CurrentVfo]

	// for k, v := range newParams {
	// 	paramValue, ok := hl.ParmValue[k]
	// 	if !ok {
	// 		return errors.New("unknown Parameter")
	// 	}
	// if _, ok := r.state.Vfo.Parameters[k]; !ok {
	// 	return errors.New("unsupported Parameter for this rig")
	// }
	// if r.state.Vfo.Levels[k] != v {
	// 	err := r.rig.SetLevel(vfo, paramValue, v)
	// 	if err != nil {
	// 		return nil
	// 	}
	// }
	// }

	return nil
}

func (r *remoteRadio) initSetState() sbRadio.SetState {
	request := sbRadio.SetState{}

	request.CurrentVfo = r.state.CurrentVfo
	request.Vfo = &sbRadio.Vfo{}
	request.Vfo.Split = &sbRadio.Split{}
	request.Md = &sbRadio.MetaData{}
	request.UserId = r.userID

	return request
}
