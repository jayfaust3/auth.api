package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jayfaust3/auth.api/pkg/handlers"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Received a request")
}

func main() {
	log.Print("starting server...")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

	router := mux.NewRouter()

	router.HandleFunc("/auth/token", handlers.getToken).Methods(http.MethodGet)
}
