package tapapicore

import "net/http"

// AddNoAuthRoutes - route without any Auth like for Twilio Callback Url
func (s *server) AddNoAuthRoutes(methodName string, methodType string, mRoute string, handlerFunc http.HandlerFunc) {
	r := route{
		Name:        methodName,
		Method:      methodType,
		Pattern:     mRoute,
		HandlerFunc: useMiddleware(handlerFunc, recovery, enableCompression, logRequest, createContext)}
	s.routes = append(s.routes, r)
}

// AddBasicRoute is to create route without role validation
func (s *server) AddBasicRoute(methodName, methodType, mRoute string, m map[string][]string, handlerFunc http.HandlerFunc) {
	r := route{
		Name:                   methodName,
		Method:                 methodType,
		Pattern:                mRoute,
		ResourcesPermissionMap: m,
		HandlerFunc:            useMiddleware(handlerFunc, recovery, enableCompression, enableCorsMiddleware, enableLogging, createContext),
	}
	s.routes = append(s.routes, r)
}

// AddRouteForApplication is to create routes with role validation by getting application from header
func (s *server) AddRouteForApplication(methodName, methodType, mRoute string, handlerFunc http.HandlerFunc) {
	r := route{
		Name:        methodName,
		Method:      methodType,
		Pattern:     mRoute,
		HandlerFunc: useMiddleware(enableUserValidationForApplication(handlerFunc), recovery, enableCompression, enableCorsMiddleware, enableLogging, createContext),
	}
	s.routes = append(s.routes, r)
}
