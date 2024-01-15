package utils

func checkValueInMapOfSlice(key string, Data []map[string]string) bool {
	for _, value := range Data {
		if key == value[key] {
			return true
		}
	}
	return false
}
