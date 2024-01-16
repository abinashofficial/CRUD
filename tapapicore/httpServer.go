package tapapicore

import (
	"net/http"
	"net/http/pprof"

	"crud/log"
	"crud/tapcontext"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmgorilla"
)

var (
	AppName   string
	GitCommit string
	GitBranch string
)

// Start - http servers
func (s *server) Start(ctx tapcontext.TContext) {
	allowedOrigins := handlers.AllowedOrigins([]string{"*"}) // Allowing all origin as of now

	allowedHeaders := handlers.AllowedHeaders([]string{
		"Accept",
		"Content-Type",
		"contentType",
		"Content-Length",
		"Accept-Encoding",
		"Client-Security-Token",
		"X-CSRF-Token",
		"X-Auth-Token",
		"processData",
		"Authorization",
		"Access-Control-Request-Headers",
		"Access-Control-Request-Method",
		"Connection",
		"Host",
		"Origin",
		"User-Agent",
		"Referer",
		"Cache-Control",
		"X-header",
		"X-Requested-With",
		"timezone",
		"locale",
		"email",
		"tap-api-token",
		"gzip-compress",
		"task",
		"x-tap-accesskey",
		"x-tap-secretkey",
		"access_token",
		"application",
	})

	allowedMethods := handlers.AllowedMethods([]string{
		"POST",
		"GET",
		"DELETE",
		"PUT",
		"PATCH",
		"OPTIONS"})

	allowCredential := handlers.AllowCredentials()

	serverHandler := handlers.CORS(
		allowedHeaders,
		allowedMethods,
		allowedOrigins,
		allowCredential)(
		context.ClearHandler(
			s.newRouter(s.subRoute),
		),
	)
	log.GenericInfo(ctx, "Starting Tap/IOT Server",
		log.FieldsMap{
			"Port":     s.port,
			"SubRoute": s.subRoute,
			"App":      AppName,
			"branch":   GitBranch,
			"commit":   GitCommit,
		})

	err := http.ListenAndServe(":"+s.port, serverHandler)
	if err != nil {
		log.GenericError(ctx, errors.New("failed to start server"),
			log.FieldsMap{
				"Port":     s.port,
				"SubRoute": s.subRoute,
				"App":      AppName,
				"branch":   GitBranch,
				"commit":   GitCommit,
			})

		return
	}
}

// NewRouter provides a mux Router.
// Handles all incoming request who matches registered routes against the request.
func (s *server) newRouter(subRoute string) *mux.Router {
	muxRouter := mux.NewRouter().StrictSlash(true)
	muxRouter.HandleFunc(subRoute+"/debug/pprof", pprof.Index)
	muxRouter.HandleFunc(subRoute+"/debug/pprof/cmdline", pprof.Cmdline)
	muxRouter.HandleFunc(subRoute+"/debug/pprof/profile", pprof.Profile)
	muxRouter.HandleFunc(subRoute+"/debug/pprof/symbol", pprof.Symbol)
	muxRouter.HandleFunc(subRoute+"/debug/pprof/trace", pprof.Trace)
	muxRouter.Handle(subRoute+"/debug/pprof/goroutine", pprof.Handler("goroutine"))
	muxRouter.Handle(subRoute+"/debug/pprof/heap", pprof.Handler("heap"))
	muxRouter.Handle(subRoute+"/debug/pprof/thread/create", pprof.Handler("threadcreate"))
	muxRouter.Handle(subRoute+"/debug/pprof/block", pprof.Handler("block"))
	muxRouter.Use(SetTraceID, apmgorilla.Middleware())
	for _, r := range s.routes {
		muxRouter.HandleFunc(subRoute+r.Pattern, r.HandlerFunc).Methods(r.Method)
	}

	return muxRouter
}

// useMiddleware applies chains of middleware (ie: log, contextWrapper, validateAuth) handler into incoming request
// For example, logging middleware might write the incoming request details to a log
// Note - It applies in reverse order
func useMiddleware(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}
