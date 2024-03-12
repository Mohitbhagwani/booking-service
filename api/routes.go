package api

import (
    "github.com/gorilla/mux"
    "booking-service/api/handlers"
	"booking-service/auth"
	"net/http"
)


func SetupRoutes(r *mux.Router) {
    r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
    // r.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
    // r.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")
    // r.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
    // r.HandleFunc("/users", handlers.GetAllUsers).Methods("GET")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	
r.Handle("/users/{id}", auth.ValidateTokenMiddleware(http.HandlerFunc(handlers.UpdateUser))).Methods("PUT")
r.Handle("/users/{id}", auth.ValidateTokenMiddleware(http.HandlerFunc(handlers.DeleteUser))).Methods("DELETE")
r.Handle("/users/{id}", auth.ValidateTokenMiddleware(http.HandlerFunc(handlers.GetUser))).Methods("GET")
r.Handle("/users", auth.ValidateTokenMiddleware(http.HandlerFunc(handlers.GetAllUsers))).Methods("GET")

	
}
