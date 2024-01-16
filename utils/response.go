package utils

import (
	"bytes"
	"crud/log"
	"crud/tapcontext"
	"encoding/json"
	"net/http"
)

type ErrResponse struct { //to-do: it must be renamed to a generic response struct
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

// ReturnResponse forms the http response in json format
func ReturnResponse(w http.ResponseWriter, statusCode int, status interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	en := json.NewEncoder(w)
	_ = en.Encode(status)
}

// ErrorResponse returns generic error response
func ErrorResponse(ctx tapcontext.TContext, w http.ResponseWriter, responseErrorMessage string, statusCode int, logError error, fields log.FieldsMap) {
	w.Header().Set("Content-Type", "application/json")
	var buf = new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	if logError != nil {
		log.GenericError(ctx, logError, fields)
	}
	_ = encoder.Encode(ErrResponse{Message: responseErrorMessage})
	w.WriteHeader(statusCode)
	_, _ = w.Write(buf.Bytes())
}
