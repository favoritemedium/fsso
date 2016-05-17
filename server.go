package main

import (
	"github.com/favoritemedium/fsso/api"
	"log"
	"net/http"
)

// Run the API server
func main() {
	api.Initialize("/api/auth/")
	err := http.ListenAndServe(":8000", nil)
	log.Fatal(err)
}
