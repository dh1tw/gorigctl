package radio

import (
	hl "github.com/dh1tw/goHamlib"
	sbRadio "github.com/dh1tw/gorigctl/sb_radio"
)

func (r *radio) serializeState() (msg []byte, err error) {

	msg, err = r.state.Marshal()

	return msg, err
}

func (r *radio) serializeCaps() (msg []byte, err error) {

	caps := sbRadio.Capabilities{}
	caps.Vfos = r.rig.Caps.Vfos
	caps.Modes = r.rig.Caps.Modes
	caps.VfoOps = r.rig.Caps.Operations
	caps.GetFunctions = r.rig.Caps.GetFunctions
	caps.SetFunctions = r.rig.Caps.SetFunctions
	caps.GetLevels = hlValuesToPbValues(r.rig.Caps.SetLevels)
	caps.SetLevels = hlValuesToPbValues(r.rig.Caps.SetLevels)
	caps.GetParameters = hlValuesToPbValues(r.rig.Caps.GetParameters)
	caps.SetParameters = hlValuesToPbValues(r.rig.Caps.SetParameters)
	caps.MaxRit = int32(r.rig.Caps.MaxRit)
	caps.MaxXit = int32(r.rig.Caps.MaxXit)
	caps.MaxIfShift = int32(r.rig.Caps.MaxIfShift)
	caps.Filters = hlMapToPbMap(r.rig.Caps.Filters)
	caps.TuningSteps = hlMapToPbMap(r.rig.Caps.TuningSteps)
	caps.Preamps = intListToint32List(r.rig.Caps.Preamps)
	caps.Attenuators = intListToint32List(r.rig.Caps.Attenuators)
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
	status, ok := hl.RigStatusName[r.rig.Caps.Status]
	if ok {
		caps.Status = status
	}
	msg, err = caps.Marshal()

	return msg, err
}

func hlMapToPbMap(hlMap map[string][]int) map[string]*sbRadio.Int32List {

	pbMap := make(map[string]*sbRadio.Int32List)

	for k, v := range hlMap {
		mv := sbRadio.Int32List{}
		mv.Value = intListToint32List(v)
		pbMap[k] = &mv
	}

	return pbMap
}

func intListToint32List(intList []int) []int32 {

	int32List := make([]int32, 0, len(intList))

	for _, i := range intList {
		var v int32
		v = int32(i)
		int32List = append(int32List, v)
	}

	return int32List
}

func hlValuesToPbValues(hlValues hl.Values) []*sbRadio.Value {

	pbValues := make([]*sbRadio.Value, 0, len(hlValues))

	for _, hlValue := range hlValues {
		var v sbRadio.Value
		v.Name = hlValue.Name
		v.Max = hlValue.Max
		v.Min = hlValue.Min
		v.Step = hlValue.Step
		pbValues = append(pbValues, &v)
	}

	return pbValues
}
