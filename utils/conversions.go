package utils

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
