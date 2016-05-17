package api

import (
  "net/http"
  "encoding/json"
  "log"
  "strings"

  "github.com/favoritemedium/fsso/sso"
)

// Type Parameters is used for input parameters parsed either from the query
// URL (GET/HEAD/OPTIONS) or from the request body (POST/PUT/PATCH).
type Parameters map[string]interface{}

// HasAll returns true if all of the specified keys are present in the map.
func (p Parameters) HasAll(keys ...string) bool {
  for _, k := range keys {
      if _, ok := p[k]; !ok {
        return false
      }
  }
  return true
}

// HasOther returns true if the map has any keys that are not in the list.
func (p Parameters) HasOther(keys ...string) bool {
  has := make(map[string]struct{})
  for k := range p {
    has[k] = struct{}{}
  }
  for _, k := range keys {
    delete(has, k)
  }
  return len(has) > 0
}

// HasExactly returns true if all of the keys specified and no others are in the map.
func (p Parameters) HasExactly(keys ...string) bool {
  has := make(map[string]struct{})
  for k := range p {
    has[k] = struct{}{}
  }
  for _, k := range keys {
    if _, ok := has[k]; !ok {
      return false
    }
    delete(has, k)
  }
  return len(has) == 0
}


// Type ErrorResponse represents an error as returned to the caller.
type ErrorResponse struct {
  Code string      `json:"code"`
  Message string   `json:"message"`
}

func (e ErrorResponse) Error() string {
  return e.Message
}

var (
  apiErrJson = ErrorResponse{"format", "Invalid JSON."}
  apiErrParameters = ErrorResponse{"parameters", "Invalid Parameters."}
  apiErrUnknown = ErrorResponse{"unknown", "Unknown error."}
)


// getParams parses the request parameters into a map.
//
// For POST, PUT, and PATCH requests, getParams attempts to json-decode the
// request body, and returns apiErrJson if the json is malformed.
//
// For all other requests (e.g. GET, HEAD, OPTIONS), getParams parses the query
// parameters.  In this case, all values are strings, and there are no error
// conditions.
func getParams(r *http.Request) (Parameters, error) {

  params := make(Parameters)

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

// wrap adds json encoding/decoding and authentication to an endpoint handler.
func wrap(handler (func(*http.Request, *sso.Member, Parameters) (interface{}, error))) func(http.ResponseWriter, *http.Request) {

  return func(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json")
    enc := json.NewEncoder(w)

    dataIn, err := getParams(r)
    if err != nil {
      enc.Encode(&apiErrJson)
      return
    }

    // m := sso.Member{}
    dataOut, err := handler(r, nil, dataIn)
    if err != nil {
      if xerr, ok := err.(ErrorResponse); ok {
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

// InitApi adds handlers for all the API endpoints.
// prefix should probably be "/api/auth/".
func Initialize(prefix string) {
  http.HandleFunc(prefix + "auth", wrap(doAuth))
  http.HandleFunc(prefix + "token", wrap(doToken))
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
}

// doAuth handles the /auth endpoint.
func doAuth(r *http.Request, m *sso.Member, p Parameters) (interface{}, error) {

  if p.HasExactly("email", "password") {
    return sso.AuthEmail(p["email"].(string), p["password"].(string))
  }

  if p.HasExactly("rtoken") {
    return sso.AuthRefresh(p["rtoken"].(string))
  }

  if p.HasExactly("provider", "id_token") {
    return sso.AuthSocial(p["provider"].(string), p["id_token"].(string))
  }

  return nil, apiErrParameters
}


// doToken handles the /token endpoint.
func doToken(r *http.Request, m *sso.Member, p Parameters) (interface{}, error) {

  if p.HasExactly("email", "password") {
    return sso.TokenEmail(p["email"].(string), p["password"].(string))
  }

  if p.HasExactly("provider", "id_token") {
    return sso.TokenSocial(p["provider"].(string), p["id_token"].(string))
  }

  return nil, apiErrParameters
}

// notImplemented is a placeholder for an endpoint that is not implemented.
// For testing purposes it returns the input data and the currently signed-in member.
func notImplemented(r *http.Request, member *sso.Member, params Parameters) (interface{}, error) {
  resp := make(map[string]interface{})
  resp[strings.ToLower(r.Method)] = params
  resp["member"] = member
  return resp, nil
}
