package main

import (
  "net/http"
  "encoding/json"
  "log"
  "strings"
)

const prefix = "/api/auth/"

type Member struct {
  Id int64         `json:"-"`
  Email string     `json:"email"`
  ShortName string `json:"shortname"`
  FullName string  `json:"fullname"`
}

type ErrorResponse struct {
  Code string      `json:"code"`
  Message string   `json:"message"`
}

func (e *ErrorResponse) Error() string {
  return e.Message
}

var (
  apiErrJson = ErrorResponse{"format", "Invalid JSON."}
  apiErrMissingParam = ErrorResponse{"missing", "Requred parameter missing."}
  apiErrExtraneousParam = ErrorResponse{"extra", "Extraneous parameter provided."}
  apiErrUnknown = ErrorResponse{"unknown", "Unknown error."}
)

// notImplemented is a placeholder for an enpoint that is not implemented.
// Returns the input data and the current member; for testing.
func notImplemented(r *http.Request, member *Member, params map[string]interface{}) (interface{}, error) {
  resp := make(map[string]interface{})
  resp[strings.ToLower(r.Method)] = params
  resp["member"] = member
  return resp, nil
}

// getParams parses the request parameters into a map.
//
// For POST, PUT, and PATCH requests, getParams attempts to json-decode the
// request body, and returns an error if the json is malformed.
//
// For all other requests (e.g. GET, HEAD, OPTIONS), getParams parses the query
// parameters.  In this case, all values are strings, and there are no error
// conditions.
func getParams(r *http.Request) (map[string]interface{}, error) {

  params := make(map[string]interface{})

  if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {

    dec := json.NewDecoder(r.Body)
    if err := dec.Decode(&params); err != nil {
      return nil, err
    }

  } else {

    for k, v := range r.URL.Query() {
      params[k] = v[0]
    }
  }

  return params, nil
}


func wrap(handler (func(*http.Request, *Member, map[string]interface{}) (interface{}, error))) func(http.ResponseWriter, *http.Request) {

  return func(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json")
    enc := json.NewEncoder(w)

    dataIn, err := getParams(r)
    if err != nil {
      enc.Encode(&apiErrJson)
      return
    }

    // m := Member{}
    dataOut, err := handler(r, nil, dataIn)
    if err != nil {
      if xerr, ok := err.(*ErrorResponse); ok {
        enc.Encode(&xerr)
      } else {
        log.Println(err)
        enc.Encode(&apiErrUnknown)
      }
      return
    }

    enc.Encode(&dataOut)
  }
}

func main() {

  http.HandleFunc(prefix + "auth", wrap(notImplemented))
  http.HandleFunc(prefix + "token", wrap(notImplemented))
  http.HandleFunc(prefix + "signout", wrap(notImplemented))
  http.HandleFunc(prefix + "email", wrap(notImplemented))
  http.HandleFunc(prefix + "verify", wrap(notImplemented))
  http.HandleFunc(prefix + "new", wrap(notImplemented))
  http.HandleFunc(prefix + "password", wrap(notImplemented))
  http.HandleFunc(prefix + "list", wrap(notImplemented))
  http.HandleFunc(prefix + "clear", wrap(notImplemented))
  http.HandleFunc(prefix + "delete", wrap(notImplemented))
  http.HandleFunc(prefix + "add", wrap(notImplemented))
  http.HandleFunc(prefix + "accounts", wrap(notImplemented))
  http.HandleFunc(prefix + "primary", wrap(notImplemented))
  http.HandleFunc(prefix + "remove", wrap(notImplemented))

  err := http.ListenAndServe(":8000", nil)
  log.Fatal(err)
}