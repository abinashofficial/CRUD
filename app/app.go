package app

import (
	"test1Project/handlers"
	"test1Project/handlers/fields"
)

var h handlers.Store

func setupHandlers() {
	h = handlers.Store{
		FieldsHandler: fields.New(),
	}
}
func Start() {
	envPort := "8080"
	setupHandlers()
	runServer(envPort, h)
}
