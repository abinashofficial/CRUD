package tapapicore

import (
	"crud/tapcontext"
	"net/http"
)

type validationResponse struct {
	Response string
	Message  string
	PermMap  map[string][]string
}

type route struct {
	Name                   string
	Method                 string
	Pattern                string
	Modules                []string
	ResourcesPermissionMap interface{}
	HandlerFunc            http.HandlerFunc
}

/*
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

*/

type requestParams struct {
	Endpoint string
	Method   string
	Token    string
	Tenant   string
	Dealer   string
	Body     interface{}
}

const (
	validateEndpoint = "/taprole/private/profile/validate"
)

const (
	//requestID of the parent request
	requestID = "requestId"

	//userEmail of the parent request
	userEmail = "email"

	//tapApiToken for authentication
	tapApiToken = "tap-api-token"
	//application for authentication
	application = "application"
	//user language preference
	locale = "locale"

	integrationsAuthToken = "integrations_auth_token"
)

type server struct {
	port     string
	subRoute string
	routes   []route
}

type TapServer interface {
	Start(ctx tapcontext.TContext)
	AddNoAuthRoutes(methodName string, methodType string, mRoute string, handlerFunc http.HandlerFunc)
	AddBasicRoute(methodName, methodType, mRoute string, m map[string][]string, handlerFunc http.HandlerFunc)
	AddRouteForApplication(methodName, methodType, mRoute string, handlerFunc http.HandlerFunc)
}

// NewTapServer returns a new server instance.
func NewTapServer(port, subRoute string) TapServer {
	return &server{
		port:     port,
		subRoute: subRoute,
	}
}
