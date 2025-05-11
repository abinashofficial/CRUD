package app

import (
	"crud/handlers"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

// CORS middleware to set the headers
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from specific frontend origins
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace with your production client URL

		// Allow specific methods (GET, POST, OPTIONS, etc.)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allow specific headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func runServer(envPort string, h handlers.Store) {
	// Ensure envPort has a valid default value
	if envPort == "" {
		envPort = "8080"
	}

	// Create a new router
	r := mux.NewRouter()

	// Public routes
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

	// WebSocket route
	r.HandleFunc("/ws", h.FieldsHandler.HandleConnections).Methods(http.MethodGet)


   r.HandleFunc("/events",h.FieldsHandler.SSEHandler)
	// Wrap the router with CORS middleware
	http.Handle("/", enableCORS(r))

	// Start the server
	fmt.Printf("Server listening on port %s...\n", envPort)
	log.Fatal(http.ListenAndServe(":"+envPort, nil))
}
