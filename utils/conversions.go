package utils

import sbRadio "github.com/dh1tw/gorigctl/sb_radio"
import hl "github.com/dh1tw/goHamlib"

// Btoi Bool to Int
func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Itob Int to Bool
func Itob(i int) bool {
	if i == 1 {
		return true
	}

	return false
}

func HlMapToPbMap(hlMap map[string][]int) map[string]*sbRadio.Int32List {

	pbMap := make(map[string]*sbRadio.Int32List)

	for k, v := range hlMap {
		mv := sbRadio.Int32List{}
		mv.Value = IntListToint32List(v)
		pbMap[k] = &mv
	}

	return pbMap
}

func IntListToint32List(intList []int) []int32 {

	int32List := make([]int32, 0, len(intList))

	for _, i := range intList {
		var v int32
		v = int32(i)
		int32List = append(int32List, v)
	}

	return int32List
}

func HlValuesToPbValues(hlValues hl.Values) []*sbRadio.Value {

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
