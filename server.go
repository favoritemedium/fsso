package main

import (
  "net/http"
  "log"
  "github.com/favoritemedium/fsso/api"
)

// Run the API server
func main() {
  api.Initialize("/api/auth/")
  err := http.ListenAndServe(":8000", nil)
  log.Fatal(err)
}
