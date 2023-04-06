package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jayfaust3/auth.api/pkg/handlers"
)

func main() {
	log.Print("starting server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	log.Printf("listening on port %s", port)
	router := mux.NewRouter()

	router.HandleFunc("/api/auth/token", handlers.GetToken).Methods(http.MethodGet)

	http.ListenAndServe(":"+port, router)
}
