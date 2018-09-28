package trackers

import (
	"strconv"
)

// stringSliceContains checks if a slice contiains a string.
func stringSliceContains(slice []string, str string) bool {
	var sliceStr string

	for _, sliceStr = range slice {
		if str == sliceStr {
			return true
		}
	}

	return false
}

// stringSliceRemove removes a string from a slice.
func stringSliceRemove(slice []string, remove string) []string {
	var str string
	var cleanSlice []string

	for _, str = range slice {
		if remove == str {
			continue
		}
		cleanSlice = append(cleanSlice, str)
	}

	return cleanSlice
}

// intSliceContains checks if a integer exists in a slice of integers.
func intSliceContains(slice []int, i int) bool {
	var j int

	for _, j = range slice {
		if i == j {
			return true
		}
	}

	return false
}

// intSliceEq compares if both slices are equal.
func intSliceEq(sliceA []int, sliceB []int) bool {
	var i int

	if (sliceA == nil) != (sliceB == nil) {
		return false
	}
	if len(sliceA) != len(sliceB) {
		return false
	}

	for i = range sliceA {
		if sliceA[i] != sliceB[i] {
			return false
		}
	}

	return true
}

// StringSliceToInt converts a slice of strings to a slice of integers.
func StringSliceToInt(slice []string) ([]int, error) {
	var sliceInt []int
	var s string
	var i int
	var err error

	for _, s = range slice {
		if i, err = strconv.Atoi(s); err != nil {
			return nil, err
		}
		sliceInt = append(sliceInt, i)
	}

	return sliceInt, nil
}
