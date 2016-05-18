package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/favoritemedium/fsso/sso"
)

// Type ErrorResponse represents an error as returned to the caller.
type ErrorResponse struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

var (
	ErrInvalidJson      = ErrorResponse{400, "format", "Invalid JSON."}
	ErrBadParameters    = ErrorResponse{400, "parameters", "Invalid Parameters."}
	ErrMethodNotAllowed = ErrorResponse{405, "method", "Method not allowed."}
	ErrUnknown          = ErrorResponse{500, "unknown", "Unknown error."}
)

// wrap adds json encoding/decoding and authentication to an endpoint handler.
func wrap(handler func(*http.Request, *sso.Member, Parameters) (interface{}, error)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)

		dataIn, err := ParseParameters(r)
		if err != nil {
			w.WriteHeader(ErrInvalidJson.Status)
			enc.Encode(&ErrInvalidJson)
			return
		}

		// m := sso.Member{}
		dataOut, err := handler(r, nil, dataIn)
		if err != nil {
			if xerr, ok := err.(ErrorResponse); ok {
				// If our error is an instance of ErrorResponse, that means that the
				// handler generated it and we should pass it on to the caller
				w.WriteHeader(xerr.Status)
				enc.Encode(&xerr)
			} else {
				// We have an unexpected error (such as database failure).
				// Log it so that we can debug.
				log.Println(err)
				w.WriteHeader(ErrUnknown.Status)
				enc.Encode(&ErrUnknown)
			}
			return
		}

		enc.Encode(&dataOut)
	}
}

// InitApi adds handlers for all the API endpoints.
// prefix should probably be "/api/auth/".
func Initialize(prefix string) {
	http.HandleFunc(prefix+"signin", wrap(doSignin))
	http.HandleFunc(prefix+"connect", wrap(doConnect))
	http.HandleFunc(prefix+"signout", wrap(notImplemented))
	http.HandleFunc(prefix+"email/check", wrap(notImplemented))
	http.HandleFunc(prefix+"email/verify", wrap(notImplemented))
	http.HandleFunc(prefix+"new", wrap(notImplemented))
	http.HandleFunc(prefix+"password", wrap(notImplemented))
	http.HandleFunc(prefix+"list", wrap(notImplemented))
	http.HandleFunc(prefix+"clear", wrap(notImplemented))
	http.HandleFunc(prefix+"delete", wrap(notImplemented))
	http.HandleFunc(prefix+"add", wrap(notImplemented))
	http.HandleFunc(prefix+"accounts", wrap(notImplemented))
	http.HandleFunc(prefix+"primary", wrap(notImplemented))
	http.HandleFunc(prefix+"remove", wrap(notImplemented))
}

// doSignin handles the /signin endpoint.
func doSignin(r *http.Request, m *sso.Member, p Parameters) (interface{}, error) {

	if r.Method != "POST" {
		return nil, ErrMethodNotAllowed
	}

	if p.HasExactly("email", "password") && p.AreString("email", "password") {
		return sso.SigninEmail(p["email"].(string), p["password"].(string))
	}

	if p.HasExactly("rtoken") && p.AreString("rtoken") {
		return sso.SigninRefresh(p["rtoken"].(string))
	}

	if p.HasExactly("provider", "id_token") && p.AreString("provider", "id_token") {
		return sso.SigninSocial(p["provider"].(string), p["id_token"].(string))
	}

	return nil, ErrBadParameters
}

// doConnect handles the /connect endpoint.
func doConnect(r *http.Request, m *sso.Member, p Parameters) (interface{}, error) {

	if r.Method != "POST" {
		return nil, ErrMethodNotAllowed
	}

	if p.HasExactly("email", "password") && p.AreString("email", "password") {
		return sso.ConnectEmail(p["email"].(string), p["password"].(string))
	}

	if p.HasExactly("provider", "id_token") && p.AreString("provider", "id_token") {
		return sso.ConnectSocial(p["provider"].(string), p["id_token"].(string))
	}

	return nil, ErrBadParameters
}

// notImplemented is a placeholder for an endpoint that is not implemented.
// For testing purposes it returns the input data and the currently signed-in member.
func notImplemented(r *http.Request, member *sso.Member, params Parameters) (interface{}, error) {
	resp := make(map[string]interface{})
	resp[r.Method] = params
	resp["member"] = member
	return resp, nil
}
