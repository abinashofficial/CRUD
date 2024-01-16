package taputils

import (
	"net/http"
)

const (
	RequestID   = "requestId"
	EmailString = "email"
)

var Client HTTPClient

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
