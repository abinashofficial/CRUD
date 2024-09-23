package app

import (
	"crud/handlers"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	handler "github.com/gorilla/handlers"


)

func runServer(envPort string, h handlers.Store) {
	r := mux.NewRouter()


	r.HandleFunc("/public/create-all", h.FieldsHandler.CreateAll).Methods(http.MethodPost)
	r.HandleFunc("/public/create", h.FieldsHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/public/get/{student-info}", h.FieldsHandler.Get).Methods(http.MethodGet)
	r.HandleFunc("/public/get-all", h.FieldsHandler.GetAll).Methods(http.MethodGet)

	r.HandleFunc("/public/update/{student-info}", h.FieldsHandler.Update).Methods(http.MethodPut)
	r.HandleFunc("/public/delete/{student-info}", h.FieldsHandler.Delete).Methods(http.MethodDelete)


	r.HandleFunc("/public/signin", h.FieldsHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/public/signup", h.FieldsHandler.Signup).Methods(http.MethodPost)
	r.HandleFunc("/public/recovery", h.FieldsHandler.PasswordChange).Methods(http.MethodPut)


corsMiddleware := handler.CORS(

	handler.AllowedOrigins([]string{"*"}), // Allowing all origin as of now

	handler.AllowedHeaders([]string{
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
		"tenant",
		"dealer",
		"tap-api-token",
		"gzip-compress",
		"task",
		"x-tap-accesskey",
		"x-tap-secretkey",
		"access_token",
		"application",
	}),

	handler.AllowedMethods([]string{
		"POST",
		"GET",
		"DELETE",
		"PUT",
		"PATCH",
		"OPTIONS"}),

		handler.AllowCredentials(),
	)
 	fmt.Printf("Server listening on port %d...\n", envPort)

    log.Fatal(http.ListenAndServe(":8080", corsMiddleware(r)))
}
