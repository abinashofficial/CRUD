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
	r.HandleFunc("/public/updateprofile", h.FieldsHandler.UpdateProfile).Methods(http.MethodPut)
	r.HandleFunc("/public/send-otp-email", h.FieldsHandler.SendOTPEmailHandler).Methods(http.MethodPost)
	r.HandleFunc("/public/verify-otp", h.FieldsHandler.VerifyOTPHandler).Methods(http.MethodPost)
	r.HandleFunc("/public/get-user", h.FieldsHandler.GetUser).Methods(http.MethodPost)
	r.HandleFunc("/public/send-otp-mobile-no", h.FieldsHandler.SendOTPMobHandler).Methods(http.MethodPost)








    // Allow all origins with the CORS middleware
    corsMiddleware := handler.CORS(
        handler.AllowedOrigins([]string{"*"}), // Allow all origins
        handler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}), // Allowed methods
        handler.AllowedHeaders([]string{"Content-Type", "Authorization"}), // Allowed headers
    )
 	fmt.Printf("Server listening on port %d...\n", envPort)

    log.Fatal(http.ListenAndServe(":8080", corsMiddleware(r)))
}
