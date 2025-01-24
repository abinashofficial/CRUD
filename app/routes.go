package app

import (
	"crud/handlers"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	// handler "github.com/gorilla/handlers"


)

// CORS middleware to set the headers
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from the frontend (e.g., http://localhost:3000)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow specific methods (GET, POST, OPTIONS, etc.)
		w.Header().Set("Access-Control-Allow-Methods", "*")

		// Allow specific headers
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
	r.HandleFunc("/ws", h.FieldsHandler.HandleConnections)
	// http.HandleFunc("/ws", h.FieldsHandler.HandleConnections)





	






	// Wrap the router with CORS middleware
	http.Handle("/", enableCORS(r))
 	fmt.Printf("Server listening on port %d...\n", envPort)

    log.Fatal(http.ListenAndServe(":8080", nil))
}
