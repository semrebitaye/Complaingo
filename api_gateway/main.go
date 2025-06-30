package main

import (
	"api_gateway/internal/gateway"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// global rate limiting
	r.Use(gateway.RateLimitMiddleware)

	// public routes no token required
	r.PathPrefix("/register").HandlerFunc(gateway.ForwardTo("http://localhost:8090")).Methods("POST")
	r.PathPrefix("/login").HandlerFunc(gateway.ForwardTo("http://localhost:8090")).Methods("POST")

	// protected routes that require auth token
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(gateway.AuthMiddleware)
	protected.PathPrefix("/").HandlerFunc(gateway.ForwardTo("http://localhost:8090"))

	log.Println("API Gateway listeninig on port :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
