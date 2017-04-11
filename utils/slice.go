package utils

func SliceDiff(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		// String not found. We add it to return slice
		if !found {
			diff = append(diff, s1)
		}
	}
	return diff
}

func StringMapDiff(map1 map[string]float64, map2 map[string]float64) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for key, _ := range map1 {
		_, found := map2[key]
		// String not found. We add it to return slice
		if !found {
			diff = append(diff, key)
		}
	}
	return diff
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Int32InSlice(a int32, list []int32) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func RemoveStringFromSlice(el string, slice []string) []string {
	for i, v := range slice {
		if v == el {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
