package main

import (
	"log"
	"net/http"

	"github.com/Arinji2/websockets/routes"
)

func main() {
	apiRouter := routes.Router()
	http.Handle("/api/", http.StripPrefix("/api", apiRouter))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
