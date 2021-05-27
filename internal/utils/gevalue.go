package utils

func GetValue(val *string) string {
	if val != nil {
		return *val
	}
	return ""
}
func GetGuidValue(val string) string {
	if val != "" {
		return val
	}
	return GetGuidDefaultValue()
}
func GetGuidDefaultValue() string {
	return "00000000-0000-0000-0000-000000000000"
}

func GetValueInt(val *int) int {
	if val != nil {
		return *val
	}
	return 0
}

func SliceStringsContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SliceIntsContains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
