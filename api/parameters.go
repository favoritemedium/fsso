package api

import (
	"encoding/json"
	"net/http"
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

// AreString returns true if all of the keys specified refer to string values.
func (p Parameters) AreString(keys ...string) bool {
	for _, k := range keys {
		if _, ok := p[k].(string); !ok {
			return false
		}
	}
	return true
}

// AreInt returns true if all of the keys specified refer to numbers.
func (p Parameters) AreNumeric(keys ...string) bool {
	for _, k := range keys {
		// The json library by default returns all numbers as float64.  If we want
		// an int, we'll have to convert it.
		if _, ok := p[k].(float64); !ok {
			return false
		}
	}
	return true
}

// AreBool returns true if all of the keys specified refer to boolean values.
func (p Parameters) AreBool(keys ...string) bool {
	for _, k := range keys {
		if _, ok := p[k].(bool); !ok {
			return false
		}
	}
	return true
}

// AreNull returns true if all of the keys specified refer to null values.
func (p Parameters) AreNull(keys ...string) bool {
	for _, k := range keys {
		if p[k] != nil {
			return false
		}
	}
	return true
}

// ParseParameters parses the request parameters into a map.
//
// For POST, PUT, and PATCH requests, getParams attempts to json-decode the
// request body, and returns an error if the json is malformed.
//
// For all other requests (e.g. GET, HEAD, OPTIONS), getParams parses the query
// parameters.  In this case, all values are strings, and there are no error
// conditions.
func ParseParameters(r *http.Request) (Parameters, error) {

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
