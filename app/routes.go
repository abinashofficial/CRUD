package app

import (
	"crud/handlers"
	"crud/tapapicore"
	"crud/tapcontext"
	"net/http"
)

const serviceRoute = "/tapapi/helpdesk/v1"

func runServer(envPort string, h handlers.Store, ctx tapcontext.TContext) {
	s := tapapicore.NewTapServer(envPort, serviceRoute)
	s.AddRouteForApplication("Create Info", http.MethodPost, "/public/create",
		h.FieldsHandler.Create)

	s.Start(ctx)
}
