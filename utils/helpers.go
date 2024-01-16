package utils

var ErrorCodes = map[string]map[string]string{}
var supportedLang = []string{"en", "fr"}

// CheckKeyInSlice function return true if the passed key is present in the slice, else return false
func CheckKeyInSlice(strArray []string, key string) bool {
	if strArray == nil {
		return false
	}
	for _, val := range strArray {
		if val == key {
			return true
		}
	}
	return false
}

func GetError(msg string, lang string) string {
	if !CheckKeyInSlice(supportedLang, lang) {
		lang = "en"
	}
	return ErrorCodes[lang][msg]
}
