package utils

import "fmt"

func FormatFreq(freq float64) string {
	initial := fmt.Sprintf("%07.0f", freq/10)
	var freqStr string = ""

	// above > 100GHz
	if freq >= 100000000000 {
		freqStr = initial[:3] + "." + initial[3:6] + "." + initial[6:9] + "." + initial[9:]
		// 10 ... 100GHz
	} else if freq < 100000000000 && freq >= 10000000000 {
		freqStr = initial[:2] + "." + initial[2:5] + "." + initial[5:8] + "." + initial[8:]
		// 1 ... 10GHz
	} else if freq < 10000000000 && freq >= 1000000000 {
		freqStr = initial[:1] + "." + initial[1:4] + "." + initial[4:7] + "." + initial[7:]
		// 100 ... 1000 MHz
	} else if freq < 1000000000 && freq >= 100000000 {
		freqStr = initial[:3] + "." + initial[3:6] + "." + initial[6:]
		// 10 ... 100MHz
	} else if freq < 100000000 && freq >= 10000000 {
		freqStr = initial[:2] + "." + initial[2:5] + "." + initial[5:]
		// 1 ... 10 MHz
	} else if freq < 10000000 && freq >= 1000000 {
		freqStr = initial[1:2] + "." + initial[2:5] + "." + initial[5:]
		// < 1 MHz
	} else if freq < 1000000 {
		freqStr = initial[2:5] + "." + initial[5:]
	}

	return freqStr
}
