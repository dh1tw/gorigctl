package localradio

import (
	"log"

	hl "github.com/dh1tw/goHamlib"
)

type LocalRadio struct {
	rig hl.Rig
	log *log.Logger
	vfo int
}

func NewLocalRadio(rigModel, debugLevel int, port hl.Port, log *log.Logger) (*LocalRadio, error) {
	lr := LocalRadio{}
	lr.rig = hl.Rig{}
	lr.log = log
	lr.vfo = hl.RIG_VFO_CURR

	lr.rig.SetDebugLevel(debugLevel)

	if err := lr.rig.Init(rigModel); err != nil {
		return nil, err
	}

	if err := lr.rig.SetPort(port); err != nil {
		return nil, err
	}

	if err := lr.rig.Open(); err != nil {
		return nil, err
	}

	vfo, err := lr.rig.GetVfo()
	if err != nil {
		lr.log.Println(err)
	} else {
		lr.vfo = vfo
	}

	return &lr, nil
}
