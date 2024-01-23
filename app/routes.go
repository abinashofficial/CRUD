package app

import (
	"crud/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func runServer(envPort string, h handlers.Store) {
	r := mux.NewRouter()
	r.HandleFunc("/public/create", h.FieldsHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/public/get/{student-info}", h.FieldsHandler.Get).Methods(http.MethodGet)
	r.HandleFunc("/public/get-all", h.FieldsHandler.GetAll).Methods(http.MethodGet)

	r.HandleFunc("/public/update", h.FieldsHandler.Update).Methods(http.MethodPut)
	r.HandleFunc("/public/delete", h.FieldsHandler.Delete).Methods(http.MethodDelete)

	fmt.Printf("Server listening on port %d...\n", envPort)
	log.Fatal(http.ListenAndServe(":"+envPort, r))
}
