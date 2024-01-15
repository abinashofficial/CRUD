package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"test1Project/handlers"
)

func runServer(envPort string, h handlers.Store) {
	r := mux.NewRouter()
	r.HandleFunc("/public/create", h.FieldsHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/public/get", h.FieldsHandler.Get).Methods(http.MethodGet)
	r.HandleFunc("/public/update", h.FieldsHandler.Update).Methods(http.MethodPut)
	r.HandleFunc("/public/delete", h.FieldsHandler.Delete).Methods(http.MethodDelete)

	fmt.Printf("Server listening on port %d...\n", envPort)
	log.Fatal(http.ListenAndServe(":"+envPort, r))
}
