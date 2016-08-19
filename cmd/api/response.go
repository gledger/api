package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

func created(
	e func(context.Context, http.ResponseWriter, interface{}) error,
) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(c context.Context, w http.ResponseWriter, r interface{}) error {
		w.WriteHeader(http.StatusCreated)
		return e(c, w, r)
	}
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Add("Content-Type", "application/vnd.api+json")
	return json.NewEncoder(w).Encode(response)
}
