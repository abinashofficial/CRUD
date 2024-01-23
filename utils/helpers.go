package utils

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func checkValueInMapOfSlice(key string, Data []map[string]string) bool {
	for _, value := range Data {
		if key == value[key] {
			return true
		}
	}
	return false
}

func GetURLParam(r *http.Request, paramName string) (string, error) {

	params := mux.Vars(r)

	param, ok := params[paramName]

	if !ok {
		return "", fmt.Errorf("%s %q", "url parameter not found", paramName)
	}

	return param, nil
}

func GetQueryParam(r *http.Request, paramName string) (string, error) {

	params := r.URL.Query()

	param := params.Get(paramName)

	if param == "" {
		return "", fmt.Errorf("%s %q", "query parameter not found", paramName)
	}

	return param, nil
}
