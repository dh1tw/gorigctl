package server

import (
	hl "github.com/dh1tw/goHamlib"
	sbRadio "github.com/dh1tw/gorigctl/sb_radio"
	"github.com/dh1tw/gorigctl/utils"
)

func (r *localRadio) serializeState() (msg []byte, err error) {

	msg, err = r.state.Marshal()

	return msg, err
}

func (r *localRadio) serializeCaps() (msg []byte, err error) {

	caps := sbRadio.Capabilities{}
	caps.Vfos = r.rig.Caps.Vfos
	caps.Modes = r.rig.Caps.Modes
	caps.VfoOps = r.rig.Caps.Operations
	caps.GetFunctions = r.rig.Caps.GetFunctions
	caps.SetFunctions = r.rig.Caps.SetFunctions
	caps.GetLevels = utils.HlValuesToPbValues(r.rig.Caps.GetLevels)
	caps.SetLevels = utils.HlValuesToPbValues(r.rig.Caps.SetLevels)
	caps.GetParameters = utils.HlValuesToPbValues(r.rig.Caps.GetParameters)
	caps.SetParameters = utils.HlValuesToPbValues(r.rig.Caps.SetParameters)
	caps.MaxRit = int32(r.rig.Caps.MaxRit)
	caps.MaxXit = int32(r.rig.Caps.MaxXit)
	caps.MaxIfShift = int32(r.rig.Caps.MaxIfShift)
	caps.Filters = utils.HlMapToPbMap(r.rig.Caps.Filters)
	caps.TuningSteps = utils.HlMapToPbMap(r.rig.Caps.TuningSteps)
	caps.Preamps = utils.IntListToint32List(r.rig.Caps.Preamps)
	caps.Attenuators = utils.IntListToint32List(r.rig.Caps.Attenuators)
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
	caps.HasAnt = r.rig.Caps.HasGetAnt
	status, ok := hl.RigStatusName[r.rig.Caps.Status]
	if ok {
		caps.Status = status
	}
	msg, err = caps.Marshal()

	return msg, err
}
